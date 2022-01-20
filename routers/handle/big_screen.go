package handle

import (
	"edgelog/app/global/consts"
	"edgelog/app/global/variable"
	"edgelog/app/service/taos/taos_warn"
	"edgelog/app/utils/response"
	"edgelog/routers/common"
	"edgelog/routers/common/cache"
	"edgelog/tools"

	"github.com/gin-gonic/gin"
)

func HandleBigScreen(group *gin.RouterGroup) {
	group.Use(common.TokenMid)
	group.GET("/1", b1)
	group.GET("/2", b2)
	group.GET("/3", b3)
	group.GET("/4", b4)
	group.GET("/5", b5)
	group.GET("/6", b6)
}

// @Summary  Warn big screen chart
// @Tags  bigscreen
// @Router /api/v1/bigscreen/1 [get]
func b1(c *gin.Context) {
	result := make(map[string]int)
	result["online"] = 5
	result["total"] = 22
	response.Success(c, consts.Success, result)
}

// @Summary  Warn big screen chart
// @Tags  bigscreen
// @Router /api/v1/bigscreen/2 [get]
func b2(c *gin.Context) {
	result := make(map[string]int)
	result["online"] = 10
	result["total"] = 50
	result["nginx"] = 5
	result["tomcat"] = 3
	result["apache"] = 2
	response.Success(c, consts.Success, result)
}

// @Summary  Warn big screen chart
// @Tags  bigscreen
// @Param type formData string true " time type, 0 1hour, 1 1day, 2 7day "
// @Success 200 {object} commun.WarnChart "{"code": 200, "data": [...]}"
// @Router /api/v1/bigscreen/3 [get]
func b3(c *gin.Context) {
	timeTypeStr := c.DefaultQuery("type", "")
	timeType := 0
	if timeTypeStr != "" {
		timeType, _ = tools.StringToInt(timeTypeStr)
	}
	result, err := taos_warn.GetWarnBigScreen(timeType)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, consts.Success, result)
}

func b4(c *gin.Context) {
	response.Success(c, consts.Success, nil)
}

// @Summary  screen 5
// @Tags  bigscreen
// @Router /api/v1/bigscreen/5 [get]
func b5(c *gin.Context) {
	response.Success(c, consts.Success, screen5Func(cache.GetTaosCache()))
}

// @Summary  screen 6
// @Tags  bigscreen
// @Router /api/v1/bigscreen/6 [get]
func b6(c *gin.Context) {
	response.Success(c, consts.Success, screen6Func(cache.GetTaosCache()))
}

func screen5Func(datas []cache.DataModel) interface{} {
	PosMap := make(map[string]int)
	PosIpMap := make(map[string]map[string]int)
	for _, v := range datas {
		pos := " unknown "
		geo, err := variable.IpStore.GetGeoByIp(v.RemoteAddr)
		if err == nil {
			pos = geo["province"]
		}
		PosMap[pos]++
		if _, ok := PosIpMap[pos]; !ok {
			PosIpMap[pos] = make(map[string]int)
		}
		PosIpMap[pos][v.RemoteAddr]++
	}
	result := make([]map[string]interface{}, 0)
	for k, v := range PosMap {
		item := make(map[string]interface{})
		item["pos"] = k
		item["sum"] = v
		res := make([]string, 0)
		for req := range PosIpMap[k] {
			res = append(res, req)
		}
		item["ips"] = PosIpMap[k]
		result = append(result, item)
	}
	return result
}

func screen6Func(datas []cache.DataModel) interface{} {
	timeReqRemoteAddrMap := make(map[string]map[string]map[string]int)
	for _, v := range datas {
		if len(v.CreatedAt) < 10 {
			continue
		}
		time := v.CreatedAt[0:10]
		if _, ok := timeReqRemoteAddrMap[time]; !ok {
			timeReqRemoteAddrMap[time] = make(map[string]map[string]int)
		}
		if _, ok := timeReqRemoteAddrMap[time][v.Request]; !ok {
			timeReqRemoteAddrMap[time][v.Request] = make(map[string]int)
		}
		timeReqRemoteAddrMap[time][v.Request][v.RemoteAddr]++
	}
	result := make(map[string]map[string]map[string]interface{})
	for time, reqRemoteAddrMap := range timeReqRemoteAddrMap {
		result[time] = make(map[string]map[string]interface{})
		for req, m := range reqRemoteAddrMap {
			result[time][req] = make(map[string]interface{})
			result[time][req]["client"] = len(m)
			sum := 0
			for _, v := range m {
				sum += v
			}
			result[time][req]["sum"] = sum
		}
	}
	return result
}
