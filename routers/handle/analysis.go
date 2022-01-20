package handle

import (
	"context"
	"edgelog/app/global/variable"
	"edgelog/app/table"
	"edgelog/routers/common"
	"edgelog/routers/common/cache"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func HandleAnalysis(group *gin.RouterGroup) {
	group.Use(common.TokenMid)
	group.GET("/cacheSize", cacheSize)
	// group.GET("/nodeOption", nodeOption)
	group.GET("/chart", chart)
	group.GET("/list", list)
	group.GET("/detailed", detailed)
}

const timeTemplate = "2006-01-02 15:04:05"

// @Summary  host option
// @Tags  expand / nginx / analysis
// @Router /api/v1/module/nginx/analysis/nodeOption [get]
// func nodeOption(c *gin.Context) {
// 	nodes := make([]model.Node, 0)
// 	err := variable.GormDb.Model(&table.Node{}).
// 		Where("type = ? ", 1).Find(&nodes).Error
// 	if err != nil {
// 		Error(c, err)
// 		return
// 	}
// 	result := make([]interface{}, len(nodes))
// 	for i, node := range nodes {
// 		result[i] = struct {
// 			Id   int
// 			Name string
// 			Ip   string
// 		}{
// 			int(node.ID), node.Name, node.Ip,
// 		}
// 	}
// 	Ok(c, result)
// }

func cacheSize(c *gin.Context) {
	common.Ok(c, cache.DebugSize())
}

// @Summary  analysis chart
// @Tags  analysis
// @Param startTime query string false " Start time "
// @Param endTime query string false " End time "
// @Router /api/v1/analysis/chart [get]
func chart(c *gin.Context) {
	startTimeStr := c.DefaultQuery("startTime", "")
	endTimeStr := c.DefaultQuery("endTime", "")
	var stime, etime time.Time
	var err error
	if startTimeStr != "" {
		stime, err = time.Parse(timeTemplate, startTimeStr)
		if err != nil {
			common.Error(c, errors.New("startTime error"))
			return
		}
	}
	if endTimeStr != "" {
		etime, err = time.Parse(timeTemplate, endTimeStr)
		if err != nil {
			common.Error(c, errors.New("endTime error"))
			return
		}
	}
	data := cache.ConditionalFilter2("", "", "", stime, etime)
	//ip analysis
	ipMap := make(map[string]map[string]int)
	ipAnalysis := make(map[string]int)
	//traffic
	trafficAnalysis := make(map[string]map[string]int)
	//status code analysis
	statusAnalysis := make(map[string]map[string]int)
	//resource analysis
	resourceAnalysis := make(map[string]int)

	for _, v := range data {
		if len(v.CreatedAt) < 10 {
			continue
		}
		dateStr := v.CreatedAt[0:10]
		//ip
		if _, ok := ipMap[dateStr]; !ok {
			ipMap[dateStr] = make(map[string]int)
		}
		ipMap[dateStr][v.RemoteAddr]++
		ipAnalysis[dateStr] = len(ipMap[dateStr])
		//traffic
		if _, ok := trafficAnalysis[dateStr]; !ok {
			trafficAnalysis[dateStr] = make(map[string]int)
		}
		// trafficAnalysis[dateStr]["visitor"] = ipAnalysis[dateStr]
		trafficAnalysis[dateStr]["visitor"]++
		trafficAnalysis[dateStr]["size"] += int(v.Size)
		//status
		if _, ok := statusAnalysis[dateStr]; !ok {
			statusAnalysis[dateStr] = make(map[string]int)
		}
		if len(v.Status) >= 1 {
			status := v.Status
			if strings.HasPrefix(status, "1") {
				statusAnalysis[dateStr]["1XX"]++
			} else if strings.HasPrefix(status, "2") {
				statusAnalysis[dateStr]["2XX"]++
			} else if strings.HasPrefix(status, "3") {
				statusAnalysis[dateStr]["3XX"]++
			} else if strings.HasPrefix(status, "4") {
				statusAnalysis[dateStr]["4XX"]++
			} else if strings.HasPrefix(status, "5") {
				statusAnalysis[dateStr]["5XX"]++
			}
		}
		// statusAnalysis[dateStr][v.Status]++
		//resource
		resourceAnalysis[dateStr]++
	}
	//time > ip,traffic,status,resource > chart
	result := make(map[string]interface{})
	result["ip"] = ipAnalysis
	result["traffic"] = trafficAnalysis
	result["status"] = statusAnalysis
	result["resource"] = resourceAnalysis
	common.Ok(c, result)
}

// @Summary  analysis list
// @Tags  analysis
// @Param status query string false " status"
// @Param reqUrl query string false " reqUrl"
// @Param clientIp query string false " clientIp"
// @Param startTime query string false " Start time "
// @Param endTime query string false " End time "
// @Param timeout query string false " Interface timeout "
// @Router /api/v1/analysis/list [get]
func list(c *gin.Context) {
	result := make([]interface{}, 0)
	resCh := make(chan interface{}, 1)
	done := make(chan struct{})
	go func() {
		template := table.Template{}
		err := variable.GormDb.Model(&table.Template{}).
			Where("id = ?", variable.Node.TemplateId).
			First(&template).Error
		if err != nil {
			done <- struct{}{}
			return
		}
		status := c.DefaultQuery("status", "")
		reqUrl := c.DefaultQuery("reqUrl", "")
		clientIp := c.DefaultQuery("clientIp", "")
		startTimeStr := c.DefaultQuery("startTime", "")
		endTimeStr := c.DefaultQuery("endTime", "")
		dataModel, err := cache.GetLatestNodeLog(10000, 0, template,
			status, reqUrl, clientIp, startTimeStr, endTimeStr)
		if err != nil {
			done <- struct{}{}
			return
		}
		resCh <- struct {
			Name        string
			Ip          string
			Time        string
			Host        string
			Method      string
			RemoteAddr  string
			Request     string
			RequestTime float64
			Status      string
		}{
			variable.Node.Name,
			variable.Node.Ip,
			dataModel.CreatedAt,
			dataModel.Host,
			dataModel.Method,
			dataModel.RemoteAddr,
			dataModel.Method,
			dataModel.RequestTime,
			dataModel.Status,
		}
		done <- struct{}{}
	}()
	timeout, err := strconv.Atoi(c.DefaultQuery("timeout", "30"))
	if err != nil {
		timeout = 30
	}
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	for {
		select {
		case r := <-resCh:
			result = append(result, r)
		case <-done:
			common.Ok(c, result)
			return
		case <-ctx.Done():
			c.JSON(http.StatusOK, gin.H{
				"code":    200,
				"timeout": true,
				"data":    result,
			})
			return
		}
	}
}

// @Summary  analysis detailed
// @Tags  analysis
// @Param size query string false " size"
// @Param status query string false " status"
// @Param reqUrl query string false " reqUrl"
// @Param clientIp query string false " clientIp"
// @Param startTime query string false " Start time "
// @Param endTime query string false " End time "
// @Router /api/v1/analysis/detailed [get]
func detailed(c *gin.Context) {
	size, err := strconv.Atoi(c.DefaultQuery("size", "20"))
	if err != nil {
		common.Error(c, err)
		return
	}
	status := c.DefaultQuery("status", "")
	reqUrl := c.DefaultQuery("reqUrl", "")
	clientIp := c.DefaultQuery("clientIp", "")
	startTimeStr := c.DefaultQuery("startTime", "")
	endTimeStr := c.DefaultQuery("endTime", "")
	template := table.Template{}
	err = variable.GormDb.Model(&table.Template{}).
		Where("id = ?", variable.Node.TemplateId).
		First(&template).Error
	if err != nil {
		common.Error(c, err)
		return
	}
	result := make([]cache.DataModel, 0)
	err = cache.GetLatestNodeLogList(10000, 0, template,
		status, reqUrl, clientIp, startTimeStr, endTimeStr, size, &result)
	if err != nil {
		common.Error(c, err)
		return
	}
	common.Ok(c, result)
}
