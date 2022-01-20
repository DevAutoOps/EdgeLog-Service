package agent

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"edgelog/app/global/consts"
	"edgelog/app/global/my_errors"
	"edgelog/app/global/variable"
	"edgelog/app/model/commun"
	"edgelog/app/model/proto/log"
	"edgelog/app/service/taos/taos_log"
	"encoding/binary"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"io"
	"net"
	"runtime/debug"
	"strconv"
	"strings"
	"time"
)

const (
	dataLengthUpperLimit uint32 = 2000
	dataLengthLowerLimit uint32 = 3
	logPrefix                   = "taq"
)

type LogServer struct {
	Listen         net.Listener
	dataBufferChan chan *commun.DataBufferClient
}

func (l *LogServer) Start() {
	time.Sleep(10 * time.Second)
	port := variable.ConfigYml.GetInt("HttpServer.ProbesLogPort")
	if port <= 0 || port >= 65536 {
		variable.ZapLog.Error(consts.PortCrossBorderMsg)
		return
	}
	var err error
	l.Listen, err = net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		fmt.Println("log server listen failed, err:", err)
		variable.ZapLog.Error("[LogServer] listen failed, err:", zap.Error(err))
		return
	}
	l.dataBufferChan = make(chan *commun.DataBufferClient, 1)
	go l.receive()
	variable.ZapLog.Info("[LogServer] listen ", zap.String("host:port", ":"+strconv.Itoa(port)))
	for {
		conn, err := l.Listen.Accept() //Establish connection
		if err != nil {
			fmt.Println("log server accept failed, err:", err)
			variable.ZapLog.Error("[LogServer] accept failed, err:", zap.Error(err))
			continue
		}
		if len(variable.ConfigYml.GetString("Agent.Ip")) > 0 && strings.HasPrefix(conn.RemoteAddr().String(), variable.ConfigYml.GetString("Agent.Ip")) {
			fmt.Printf("log server  And  %s Establish connection\n", conn.RemoteAddr())
			variable.ZapLog.Info("[LogServer] Establish connection:", zap.String("host:port", conn.RemoteAddr().String()))
			go l.process(conn)
		}
	}
}

func (l *LogServer) process(conn net.Conn) {
	defer conn.Close()
	dataProperties := &commun.ReceiveDataProperties{
		IsBegin:            true,
		IsFinish:           false,
		TotalLength:        0,
		UnReadFinishLength: 0,
		DataBuffer:         new(bytes.Buffer),
	}
	for {
		reader := bufio.NewReader(conn)
		err := l.parse(&conn, reader, dataProperties)
		if err != nil && err.Error() == my_errors.ErrorCommunicationEnd {
			continue
		}
		if err == io.EOF {
			fmt.Printf("log server read from client failed, err:%s\n", err)
			variable.ZapLog.Error(fmt.Sprintf("[LogServer] read from client failed, err:"), zap.Error(err))
			return
		}
		if err != nil {
			fmt.Printf("log server parse (%s) msg failed, err:%s\n", conn.RemoteAddr(), err)
			variable.ZapLog.Error(fmt.Sprintf("[LogServer] parse (%s) msg failed, err:", conn.RemoteAddr()), zap.Error(err))
			return
		}
	}
}

func (l *LogServer) parse(conn *net.Conn, reader *bufio.Reader, dataProperties *commun.ReceiveDataProperties) error {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Printf("log server parse (%s) recover, err:%s\n", (*conn).RemoteAddr(), err)
			variable.ZapLog.Error(fmt.Sprintf("[LogServer parse] parse (%s) recover, err: %s\n  recover info: %s", (*conn).RemoteAddr(), err, debug.Stack()))
			l.resetParse(dataProperties)
		}
	}()
	if dataProperties.IsBegin {
		//Read prefix
		prefixByte := make([]byte, len(logPrefix))
		_, err := reader.Read(prefixByte)
		if err != nil {
			return err
		}
		if string(prefixByte) != logPrefix {
			fmt.Printf("log server (%s) parse log data prefix unequal\n", (*conn).RemoteAddr())
			return nil
		}
		//Read message length
		lengthByte := make([]byte, 4)
		_, err = reader.Read(lengthByte)
		if err != nil {
			return err
		}
		lengthBuff := bytes.NewBuffer(lengthByte)
		var length uint32
		err = binary.Read(lengthBuff, binary.LittleEndian, &length)
		if err != nil {
			return err
		}

		if length > dataLengthUpperLimit || length <= dataLengthLowerLimit {
			return nil
		}

		//Read data original length
		oriLengthByte := make([]byte, 4)
		_, err = reader.Read(oriLengthByte)
		if err != nil {
			return err
		}
		oriLengthBuff := bytes.NewBuffer(oriLengthByte)
		var oriLength uint32
		err = binary.Read(oriLengthBuff, binary.LittleEndian, &oriLength)
		if err != nil {
			return err
		}
		dataProperties.TotalLength = length
		dataProperties.UnReadFinishLength = int32(length)
		dataProperties.OriLength = oriLength
		dataProperties.IsBegin = false
	} else {
		//Read next byte
		nextByte := make([]byte, 1)
		_, err := reader.Read(nextByte)
		if err != nil {
			return err
		}
		if dataProperties.DataBuffer == nil {
			dataProperties.DataBuffer = new(bytes.Buffer)
		}
		dataProperties.DataBuffer.Write(nextByte)
		dataProperties.UnReadFinishLength -= 1
	}

	received := int32(reader.Buffered())
	dataProperties.CurrentLength = received
	dataProperties.UnReadFinishLength = dataProperties.UnReadFinishLength - received

	isFinish := false
	if dataProperties.UnReadFinishLength <= 0 {
		isFinish = true
	}

	pack := make([]byte, dataProperties.CurrentLength)
	_, err := reader.Read(pack)
	if err != nil {
		return err
	}
	if dataProperties.DataBuffer == nil {
		dataProperties.DataBuffer = new(bytes.Buffer)
	}
	dataProperties.DataBuffer.Write(pack)

	if isFinish {
		data := &commun.DataBufferClient{
			DataBuffer: bytes.NewBuffer(dataProperties.DataBuffer.Bytes()),
			Length:     dataProperties.TotalLength,
			OriLength:  dataProperties.OriLength,
			RemoteInfo: (*conn).RemoteAddr().String(),
		}
		l.dataBufferChan <- data
		l.resetParse(dataProperties)
		return errors.New(my_errors.ErrorCommunicationEnd)
	}

	return nil
}

func (l *LogServer) resetParse(dataProperties *commun.ReceiveDataProperties) {
	dataProperties.IsBegin = true
	dataProperties.IsFinish = false
	dataProperties.TotalLength = 0
	dataProperties.UnReadFinishLength = 0
	dataProperties.OriLength = 0
	dataProperties.CurrentLength = 0
	dataProperties.DataBuffer.Reset()
}

func (l *LogServer) receive() {
	for {
		dataBufferClient := <-l.dataBufferChan
		if dataBufferClient.DataBuffer.Len() > 0 && dataBufferClient.DataBuffer.Len() == int(dataBufferClient.Length) {
			data := make([]byte, dataBufferClient.DataBuffer.Len())
			n, err := dataBufferClient.DataBuffer.Read(data)
			if n <= 0 || err != nil {
				continue
			}
			gzipBuffer := bytes.NewBuffer(data)
			gzipReader, err := gzip.NewReader(gzipBuffer)
			if err != nil {
				fmt.Printf("log server decode log zip data error: %s\n", err)
				continue
			}
			gzipReader.Close()
			var unGzipBuffer bytes.Buffer
			_, _ = io.Copy(&unGzipBuffer, gzipReader)
			if unGzipBuffer.Len() <= 0 {
				fmt.Printf("log server read log zip data len is zero.\n")
				continue
			}
			logData := &log.LogData{}
			err = proto.Unmarshal(unGzipBuffer.Bytes(), logData)
			if err != nil {
				fmt.Printf("log server decode log data error: %s\n", err)
				continue
			}
			if logData == nil {
				fmt.Sprintf("[CommandServer receive] log data object is empty: %v\n", logData)
				continue
			}
			l.insertDatabase(logData)
		} else {
			time.Sleep(1e8)
		}
	}
}

func (l *LogServer) insertDatabase(logData *log.LogData) {
	taos_log.LogAdd(logData)
}
