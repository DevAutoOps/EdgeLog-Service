package handle

import (
	"edgelog/app/global/consts"
	"edgelog/app/global/variable"
	"edgelog/app/model"
	"edgelog/app/service/taos/taos_monitor"
	"edgelog/app/utils/response"
	"edgelog/routers/common"
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
)

func HandleMonitor(group *gin.RouterGroup) {
	group.Use(common.TokenMid)
	group.GET("/nodeMonitor", nodeMonitor)
	group.GET("/nodeStatus", getNodeStatus)
}

// @Summary  Query agent status
// @Tags  Monitor
// @Param chartType query string true "chart type, 0 Cpu usage, 1 Cpu load, 2 Mem rate, 3 Disk part rate, 4 Network recv, 5 Network send"
// @Param timeType query string true "time type, 0 real time, 1 24hour, 2 7days, 3 custom"
// @Param granular query string true "time granular, 0 1minute, 1 5minute, 2 15minute, 3 1hour, 4 1day"
// @Param beginTime query string false "Used when time type = 3 is customized"
// @Param endTime query string false "Used when time type = 3 is customized"
// @Success 200 {object} commun.MonitorChart2 "{"code": 200, "data": [...]}"
// @Router /api/v1/monitor/nodeMonitor [GET]
func nodeMonitor(context *gin.Context) {
	chartTypeStr := context.Query("chartType")
	timeTypeStr := context.Query("timeType")
	granularStr := context.Query("granular")
	beginTime := context.DefaultQuery("beginTime", "")
	endTime := context.DefaultQuery("endTime", "")

	if len(chartTypeStr) < 0 {
		response.Error(context, errors.New("missing chart type"))
		return
	}
	chartType, err := strconv.Atoi(chartTypeStr)
	if err != nil {
		response.Error(context, errors.New("chart type not a number"))
		return
	}
	if len(timeTypeStr) < 0 {
		response.Error(context, errors.New("missing time type"))
		return
	}
	timeType, err := strconv.Atoi(timeTypeStr)
	if err != nil {
		response.Error(context, errors.New("time type not a number"))
		return
	}
	if len(granularStr) < 0 {
		response.Error(context, errors.New("missing granular"))
		return
	}
	granular, err := strconv.Atoi(granularStr)
	if err != nil {
		response.Error(context, errors.New("granular not a number"))
		return
	}
	if timeType == 3 && (len(beginTime) <= 0 || len(endTime) <= 0) {
		response.Error(context, errors.New("need begin time and end time"))
		return
	}

	dataType := 0
	switch chartType {
	case 0:
		dataType = int(consts.MonitorCpuUsage)
	case 1:
		dataType = int(consts.MonitorCpuLoad)
	case 2:
		dataType = int(consts.MonitorMemRate)
	case 3:
		dataType = int(consts.MonitorDiskPartRate)
	case 4:
		dataType = int(consts.MonitorNetworkReception)
	case 5:
		dataType = int(consts.MonitorNetworkSending)
	}

	if chartType == 2 || chartType == 3 {
		chart, err := taos_monitor.GetServerMonitorDataWithMultipoint(dataType, timeType, beginTime, endTime, granular)
		if err != nil {
			response.Error(context, err)
			return
		}
		response.Success(context, consts.CurdStatusOkMsg, chart)
	} else {
		chart, err := taos_monitor.GetServerMonitorData(dataType, timeType, beginTime, endTime, granular)
		if err != nil {
			response.Error(context, err)
			return
		}
		response.Success(context, consts.CurdStatusOkMsg, chart)
	}
}

// @Summary  Query agent status
// @Tags  Monitor
// @Success 200 {object} []model.Node "{"code": 200, "data": [...]}"
// @Router /api/v1/monitor/nodeStatus [get]
func getNodeStatus(context *gin.Context) {
	response.Success(context, consts.CurdStatusOkMsg, []model.Node{*variable.Node})
}
