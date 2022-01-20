package cache

import (
	"edgelog/app/global/variable"
	"edgelog/app/service/taos/taos_log"
	"edgelog/app/table"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"
)

const timeTemplate = "2006-01-02 15:04:05"

type DataModel struct {
	NodeId               uint
	CreatedAt            string
	Request              string
	Timestamp            string
	HttpHost             string
	Method               string
	RequestTime          float64
	HttpXForwardedFor    string
	Status               string
	Type                 string
	ResponseTime         float64
	UpstreamResponseTime string
	UpstreamHost         string
	Path                 string
	Referer              string
	Host                 string
	RemoteAddr           string
	Size                 int64
	UserAgent            string
}

func init() {
	go syncTaosCache(30 * time.Second)
}

var taosCacheLock *sync.RWMutex = new(sync.RWMutex)
var taosCache = make([]DataModel, 0)
var maxCacheSize = 100000000

func Debug(size int) []DataModel {
	taosCacheLock.RLock()
	defer taosCacheLock.RUnlock()
	result := make([]DataModel, size)
	copy(result, taosCache)
	return result
}

func DebugSize() int {
	taosCacheLock.RLock()
	defer taosCacheLock.RUnlock()
	return len(taosCache)
}

func GetTaosCache() []DataModel {
	taosCacheLock.RLock()
	defer taosCacheLock.RUnlock()
	result := make([]DataModel, len(taosCache))
	copy(result, taosCache)
	return result
}

func syncTaosCache(t time.Duration) {
	endTime := time.Now()
	startTime := time.Now().Add(-time.Hour * 168)
	if datas, err := taos_log.LogSearch(startTime.Format(timeTemplate), endTime.Format(timeTemplate)); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("first time:", time.Since(endTime))
		fmt.Println("first szie:", len(datas))
		appendTaosCache(datas)
	}
	ticker := time.NewTicker(t)
	for {
		<-ticker.C
		endTime = time.Now()
		startTime = time.Now().Add(-t)
		if datas, err := taos_log.LogSearch(startTime.Format(timeTemplate), endTime.Format(timeTemplate)); err != nil {
			fmt.Println(err)
		} else {
			appendTaosCache(datas)
		}
		ticker.Reset(t)
	}
}

func ParseLog(logStr string, template table.Template) map[string]interface{} {
	fields := strings.Split(template.Value, template.Separator)
	result := make(map[string]interface{})
	//template.Format true => json
	if !template.Format {
		var values []string
		if template.Separator == " " {
			isQuotation := false
			values = strings.FieldsFunc(logStr, func(ru rune) bool {
				if ru == '"' {
					isQuotation = !isQuotation
				}
				if ru == '"' {
					return true
				}
				if (ru == ' ' || ru == '\t') && !isQuotation {
					return true
				}
				return false
			})
		} else {
			values = strings.Split(logStr, template.Separator)
		}
		if len(fields) != len(values) {
			return result
		}
		for i, field := range fields {
			result[field] = values[i]
		}
	} else {
		tempMap := make(map[string]interface{})
		err := json.Unmarshal([]byte(logStr), &tempMap)
		if err != nil {
			return result
		}
		for _, field := range fields {
			if v, ok := tempMap[field]; ok {
				result[field] = v
			}
		}
	}
	return result
}

func ParseDataModel(dataMap map[string]interface{}) DataModel {
	result := DataModel{}
	for k, v := range dataMap {
		switch k {
		case "created_at", "time_iso8601":
			if value, ok := v.(string); ok {
				result.CreatedAt = value
			}
		case "remote_addr":
			if value, ok := v.(string); ok {
				result.RemoteAddr = value
			}
		case "method", "request_method":
			if value, ok := v.(string); ok {
				result.Method = value
			}
		case "request", "document_uri":
			if value, ok := v.(string); ok {
				result.Request = value
			}
		case "status":
			if value, ok := v.(string); ok {
				result.Status = value
			}
		case "size", "body_bytes_sent", "bytes_sent":
			if value, ok := v.(int64); ok {
				result.Size = value
			}
			if value, ok := v.(float64); ok {
				result.Size = int64(value)
			}
		case "referer":
			if value, ok := v.(string); ok {
				result.Referer = value
			}
		case "http_host":
			if value, ok := v.(string); ok {
				result.HttpHost = value
				result.Host = value
			}
		case "upstream_response_time":
			if value, ok := v.(string); ok {
				result.ResponseTime, _ = strconv.ParseFloat(value, 64)
			}
		case "response_time":
			if value, ok := v.(float64); ok {
				result.ResponseTime = value
			}
		case "request_time":
			if value, ok := v.(float64); ok {
				result.RequestTime = value
			}
		case "http_x_forwarded_for":
			if value, ok := v.(string); ok {
				result.HttpXForwardedFor = value
			}
		case "user_agent":
			if value, ok := v.(string); ok {
				result.UserAgent = value
			}
		case "upstream_host":
			if value, ok := v.(string); ok {
				result.UpstreamHost = value
			}
		}
	}
	return result
}

func appendTaosCache(datas []taos_log.LogModel) {
	taosCacheLock.Lock()
	defer taosCacheLock.Unlock()
	template := table.Template{}
	err := variable.GormDb.Model(&table.Template{}).
		Where("id = ?", variable.Node.TemplateId).
		First(&template).Error
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, data := range datas {
		dataMap := ParseLog(data.Log, template)
		item := ParseDataModel(dataMap)
		req := strings.ToLower(item.Method)
		if !strings.HasPrefix(req, "head") &&
			!strings.HasPrefix(req, "get") &&
			!strings.HasPrefix(req, "post") &&
			!strings.HasPrefix(req, "put") &&
			!strings.HasPrefix(req, "delete") {
			continue
		}
		taosCache = append(taosCache, item)
	}
	if len(taosCache) > maxCacheSize {
		taosCache = taosCache[:maxCacheSize]
	}
}

func ConditionalFilter(pname, pip, status string, start, end time.Time) []DataModel {
	result := make([]DataModel, 0)
	taosCacheLock.RLock()
	defer taosCacheLock.RUnlock()
	for _, v := range taosCache {
		if !(start.IsZero() && end.IsZero()) {
			stamp2, err := time.Parse("2006-01-02T15:04:05+08:00", v.CreatedAt)
			if err != nil {
				continue
			}
			if stamp2.Before(start) || stamp2.After(end) {
				continue
			}
		}
		if status != "" && status != v.Status {
			continue
		}
		result = append(result, v)
	}

	return result
}

func ConditionalFilter2(status, reqUrl, clientIp string, start, end time.Time) []DataModel {
	result := make([]DataModel, 0)
	taosCacheLock.RLock()
	defer taosCacheLock.RUnlock()
	for _, v := range taosCache {
		if reqUrl != "" && reqUrl != v.Method {
			continue
		}
		if clientIp != "" && clientIp != v.RemoteAddr {
			continue
		}
		if status != "" && status != v.Status {
			continue
		}
		if !(start.IsZero() && end.IsZero()) {
			stamp, err := time.Parse("2006-01-02T15:04:05+08:00", v.CreatedAt)
			if err != nil {
				continue
			}
			if stamp.Before(start) || stamp.After(end) {
				continue
			}
		}
		result = append(result, v)
	}
	return result
}

func GetLatestNodeLog(limit, offset int, template table.Template,
	status, reqUrl, clientIp, start, end string) (DataModel, error) {
	datas, err := taos_log.LogLastSearchByLimitAndOffset(limit, offset, start, end)
	if err != nil {
		return DataModel{}, err
	}
	if len(datas) == 0 {
		return DataModel{}, errors.New("not found")
	}
	for i := 0; i < len(datas); i++ {
		data := datas[i]
		dataMap := ParseLog(data.Log, template)
		v := ParseDataModel(dataMap)
		if reqUrl != "" && reqUrl != v.Method {
			continue
		}
		if status != "" && status != v.Status {
			continue
		}
		if clientIp != "" && clientIp != v.RemoteAddr {
			continue
		}
		return v, nil
	}
	return GetLatestNodeLog(limit, offset+limit, template, status, reqUrl, clientIp, start, end)
}

func GetLatestNodeLogList(limit, offset int, template table.Template,
	status, reqUrl, clientIp, start, end string, size int, list *[]DataModel) error {
	datas, err := taos_log.LogLastSearchByLimitAndOffset(limit, offset, start, end)
	if err != nil {
		return err
	}
	if len(datas) == 0 {
		return errors.New("not found")
	}
	for i := 0; i < len(datas); i++ {
		data := datas[i]
		dataMap := ParseLog(data.Log, template)
		v := ParseDataModel(dataMap)
		if reqUrl != "" && reqUrl != v.Method {
			continue
		}
		if status != "" && status != v.Status {
			continue
		}
		if clientIp != "" && clientIp != v.RemoteAddr {
			continue
		}
		size--
		*list = append(*list, v)
		if size == 0 {
			return nil
		}
	}
	return GetLatestNodeLogList(limit, offset+limit, template, status, reqUrl, clientIp, start, end, size, list)
}

func GetLatestNodeLogList2(limit, offset, size int, list *[]string) error {
	datas, err := taos_log.LogLastSearchByLimitAndOffset(limit, offset, "", "")
	if err != nil {
		return err
	}
	if len(datas) == 0 {
		return errors.New("not found")
	}
	for i := 0; i < len(datas); i++ {
		size--
		*list = append(*list, datas[i].Log)
		if size == 0 {
			return nil
		}
	}
	return GetLatestNodeLogList2(limit, offset+limit, size, list)
}
