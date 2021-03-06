package ip

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
	"strings"
)

const (
	HEADER_LENGTH   = 4
	BODYLINE_LENGTH = 20
)

const (
	AREA_COUNTRY  = "country"
	AREA_PROVINCE = "province"
	AREA_CITY     = "city"
	AREA_ZONE     = "zone"
	AREA_LOCATION = "location"
	AREA_OPERATOR = "operator"
)

//Obtain IP segment information
type IpRow struct {
	Start    uint32
	End      uint32
	Country  uint16
	Province uint16
	City     uint16
	Zone     uint16
	Location uint16
	Operator uint16
}

type IpStore struct {
	bodyLength   uint32
	metaLength   uint32
	headerBuffer []byte
	bodyBuffer   []byte
	metaBuffer   []byte
	IpTable      []IpRow //IP information table increases by range
	metaTable    map[string][]string
}

func NewIpStore(filename string) *IpStore {
	store := IpStore{headerBuffer: make([]byte, HEADER_LENGTH), metaTable: make(map[string][]string)}
	store.parseStore(filename)
	return &store
}

//Get IP location information
func (i *IpStore) GetGeoByIp(ipSearch string) (location map[string]string, err error) {
	row, err := i.searchIpRow(ipSearch)
	if err != nil {
		return location, err
	}
	location, err = i.parseIpGeo(row)
	return location, err
}

//Get the area code of IP
func (i *IpStore) GetGeocodeByIp(ipSearch string) (uint64, error) {
	row, err := i.searchIpRow(ipSearch)
	if err != nil {
		return 0, err
	}
	areacode := i.getGeocodeByRow(row)
	codeUint64, err := strconv.ParseUint(areacode, 10, 64)
	if err != nil {
		return 0, err
	}
	return codeUint64, nil

}

func (i *IpStore) GetGeoByGeocode(areacode uint64) map[string]string {
	result := map[string]string{}
	result[AREA_OPERATOR] = i.metaTable[AREA_OPERATOR][areacode%100]
	areacode /= 100
	result[AREA_LOCATION] = i.metaTable[AREA_LOCATION][areacode%100]
	areacode /= 100
	result[AREA_ZONE] = i.metaTable[AREA_ZONE][areacode%10000]
	areacode /= 10000
	result[AREA_CITY] = i.metaTable[AREA_CITY][areacode%10000]
	areacode /= 10000
	result[AREA_PROVINCE] = i.metaTable[AREA_PROVINCE][areacode%10000]
	areacode /= 10000
	result[AREA_COUNTRY] = i.metaTable[AREA_COUNTRY][areacode%10000]
	return result
}

//Get the IP zone information list
func (i *IpStore) GetMetaTable() map[string][]string {
	return i.metaTable
}

//Get the information of the IP segment where the IP is located
func (i *IpStore) searchIpRow(ipSearch string) (row IpRow, err error) {
	search := uint32(IP2Num(ipSearch))
	//fmt.Println(search)
	var start uint32 = 0
	var end uint32 = uint32(len(i.IpTable) - 1)
	var offset uint32 = 0
	for start <= end {
		mid := uint32(math.Floor(float64((end - start) / 2)))
		offset = start + mid
		IpRow := i.IpTable[offset]
		//fmt.Println(IpRow)
		if search >= IpRow.Start {
			if search <= IpRow.End {
				return IpRow, nil
			} else {
				start = offset + 1
				continue
			}
		} else {
			end = offset - 1
			continue
		}
	}
	return row, errors.New("fail to find")
}

func (i *IpStore) parseStore(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		panic("error opening file: %v\n" + err.Error())
	}
	defer file.Close()
	//fmt.Println("open file: ", filename)
	var buf [HEADER_LENGTH]byte

	if _, err := file.Read(buf[0:4]); err != nil {
		panic("error read header" + err.Error())
	}

	i.bodyLength = binary.BigEndian.Uint32(buf[0:4])
	//fmt.Println("body length is: ", i.bodyLength)
	if _, err := file.Read(buf[0:4]); err != nil {
		panic("error read header" + err.Error())
	}
	i.metaLength = binary.BigEndian.Uint32(buf[0:4])
	//fmt.Println("meta length is: ", i.metaLength)
	if err := i.paseBody(file); err != nil {
		panic("parse body  failed:" + err.Error())
	}

	if err := i.parseMeta(file); err != nil {
		panic("pase meta failed" + err.Error())
	}
}

func (i *IpStore) paseBody(file *os.File) error {
	i.bodyBuffer = make([]byte, i.bodyLength)
	if _, err := file.ReadAt(i.bodyBuffer, HEADER_LENGTH+HEADER_LENGTH); err != nil {
		panic("read body error")
	}
	buf := bytes.NewBuffer(i.bodyBuffer)
	var offset uint32 = 0
	for offset < i.bodyLength {
		line := buf.Next(BODYLINE_LENGTH)
		row, err := i.parseBodyLine(line)
		if err != nil {
			return err
		}
		i.IpTable = append(i.IpTable, row)
		offset += BODYLINE_LENGTH
	}
	return nil
}

func (i *IpStore) parseMeta(file *os.File) (err error) {
	i.metaBuffer = make([]byte, i.metaLength)
	if _, err := file.ReadAt(i.metaBuffer, int64(HEADER_LENGTH+HEADER_LENGTH+i.bodyLength)); err != nil {
		panic("read meta error")
	}
	return json.Unmarshal(i.metaBuffer, &i.metaTable)
}

func (i *IpStore) parseIpGeo(row IpRow) (map[string]string, error) {
	geo := make(map[string]string)
	geo[AREA_COUNTRY] = i.metaTable[AREA_COUNTRY][row.Country]
	geo[AREA_PROVINCE] = i.metaTable[AREA_PROVINCE][row.Province]
	geo[AREA_CITY] = i.metaTable[AREA_CITY][row.City]
	geo[AREA_ZONE] = i.metaTable[AREA_ZONE][row.Zone]
	geo[AREA_LOCATION] = i.metaTable[AREA_LOCATION][row.Location]
	geo[AREA_OPERATOR] = i.metaTable[AREA_OPERATOR][row.Operator]
	geo["areacode"] = i.getGeocodeByRow(row)
	return geo, nil

}

func (i *IpStore) getGeocodeByRow(row IpRow) string {
	countryCode := strconv.Itoa(int(row.Country))
	provinceCode := fmt.Sprintf("%04d", row.Province)
	cityCode := fmt.Sprintf("%04d", row.City)
	zoneCode := fmt.Sprintf("%04d", row.Zone)
	provoderCode := fmt.Sprintf("%02d", row.Location)
	OperatorCode := fmt.Sprintf("%02d", row.Operator)
	return countryCode + provinceCode + cityCode + zoneCode + provoderCode + OperatorCode

}

//@TODO Parse by Reflect IpRow
func (i *IpStore) parseBodyLine(buffer []byte) (row IpRow, err error) {
	buf := bytes.NewBuffer(buffer)
	if err = binary.Read(buf, binary.BigEndian, &row.Start); err != nil {
		goto fail
	}
	if err = binary.Read(buf, binary.BigEndian, &row.End); err != nil {
		goto fail
	}
	if err = binary.Read(buf, binary.BigEndian, &row.Country); err != nil {
		goto fail
	}
	if err = binary.Read(buf, binary.BigEndian, &row.Province); err != nil {
		goto fail
	}
	if err = binary.Read(buf, binary.BigEndian, &row.City); err != nil {
		goto fail
	}
	if err = binary.Read(buf, binary.BigEndian, &row.Zone); err != nil {
		goto fail
	}
	if err = binary.Read(buf, binary.BigEndian, &row.Location); err != nil {
		goto fail
	}
	if err = binary.Read(buf, binary.BigEndian, &row.Operator); err != nil {
		goto fail
	}
fail:
	return row, err
}

func IP2Num(requestip string) uint64 {
	//Get the long of the client address
	nowip := strings.Split(requestip, ".")
	if len(nowip) != 4 {
		return 0
	}
	a, _ := strconv.ParseUint(nowip[0], 10, 64)
	b, _ := strconv.ParseUint(nowip[1], 10, 64)
	c, _ := strconv.ParseUint(nowip[2], 10, 64)
	d, _ := strconv.ParseUint(nowip[3], 10, 64)
	ipNum := a<<24 | b<<16 | c<<8 | d
	return ipNum
}

func Num2IP(ipnum uint64) string {
	byte1 := ipnum & 0xff
	byte2 := (ipnum & 0xff00)
	byte2 >>= 8
	byte3 := (ipnum & 0xff0000)
	byte3 >>= 16
	byte4 := (ipnum & 0xff000000)
	byte4 >>= 24
	result := strconv.FormatUint(byte4, 10) + "." +
		strconv.FormatUint(byte3, 10) + "." +
		strconv.FormatUint(byte2, 10) + "." +
		strconv.FormatUint(byte1, 10)
	return result
}

type datFile struct {
	err error
	*bytes.Buffer
	headerLength int
	bodyLength   int
	geoMap       map[string]map[string]uint16
	geoSlice     map[string][]string
	operator     map[string]int
	writer       io.Writer
}

func NewDatFile(w io.Writer) *datFile {
	m := map[string]map[string]uint16{
		AREA_COUNTRY:  make(map[string]uint16),
		AREA_PROVINCE: make(map[string]uint16),
		AREA_CITY:     make(map[string]uint16),
		AREA_ZONE:     make(map[string]uint16),
		AREA_LOCATION: make(map[string]uint16),
		AREA_OPERATOR: make(map[string]uint16),
	}
	return &datFile{
		Buffer:   bytes.NewBuffer(nil),
		geoMap:   m,
		geoSlice: make(map[string][]string),
		writer:   bufio.NewWriter(w),
	}
}

//get area code by typ
func (d *datFile) getCode(typ string, area string) uint16 {
	var code uint16
	code, ok := d.geoMap[typ][area]
	if !ok {
		code = uint16(len(d.geoMap[typ]))
		d.geoMap[typ][area] = code
		d.geoSlice[typ] = append(d.geoSlice[typ], area)
	}
	return code
}

//@TODO parse fields by reflect the ip row
func (d *datFile) writeBody(fields []string) error {
	if d.err != nil {
		return d.err
	}
	start, _ := strconv.ParseUint(fields[0], 10, 32)
	end, _ := strconv.ParseUint(fields[1], 10, 32)
	binary.Write(d, binary.BigEndian, uint32(start))
	binary.Write(d, binary.BigEndian, uint32(end))
	binary.Write(d, binary.BigEndian, d.getCode(AREA_COUNTRY, fields[2]))
	binary.Write(d, binary.BigEndian, d.getCode(AREA_PROVINCE, fields[3]))
	binary.Write(d, binary.BigEndian, d.getCode(AREA_CITY, fields[4]))
	binary.Write(d, binary.BigEndian, d.getCode(AREA_ZONE, fields[5]))
	binary.Write(d, binary.BigEndian, d.getCode(AREA_LOCATION, fields[6]))
	binary.Write(d, binary.BigEndian, d.getCode(AREA_OPERATOR, fields[7]))
	return d.err
}

//bodylength|body|metalength|meta
func (d *datFile) writeFile() error {
	if d.err != nil {
		return d.err
	}

	bodyLength := d.Buffer.Len()
	meta, err := json.Marshal(d.geoSlice)
	if err != nil {
		d.err = err
		return d.err
	}
	metaLength := len(meta)

	binary.Write(d.writer, binary.BigEndian, uint32(bodyLength))
	binary.Write(d.writer, binary.BigEndian, uint32(metaLength))
	d.writer.Write(d.Buffer.Bytes())
	d.writer.Write(meta)

	fmt.Println("meta length is: ", metaLength)
	fmt.Println("body length is: ", bodyLength)
	return err
}

func MakeDat(infile, outfile string) error {
	in, err := os.Open(infile)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.OpenFile(outfile, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 755)
	if err != nil {
		return err
	}
	defer out.Close()
	output := NewDatFile(out)
	r := bufio.NewReader(in)
	count := 0
	for {
		count++
		line, err := r.ReadString('\n')
		if err != nil && err != io.EOF {
			return err
		}
		if len(line) != 0 {
			fields := strings.Fields(line)
			if len(fields) != 8 {
				return errors.New("invalid input file invalid line string")
			}
			if err := output.writeBody(fields); err != nil {
				return err
			}
		}
		if err == io.EOF {
			break
		}
	}
	if err := output.writeFile(); err != nil {
		return err
	}
	fmt.Println("amount ip range from ip source: ", count)
	return nil
}
