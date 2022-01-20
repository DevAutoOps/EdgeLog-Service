package handle

import (
	"edgelog/app/global/consts"
	"edgelog/app/global/variable"
	"edgelog/app/service/taos/taos_warn"
	"edgelog/app/utils/response"
	"edgelog/routers/common"
	"edgelog/routers/common/dao"
	"edgelog/routers/common/notice"
	"edgelog/tools"
	"encoding/json"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func HandleWarn(group *gin.RouterGroup) {
	group.Use(common.TokenMid)
	group.GET("/getPushConfig", getPushConfig)
	group.GET("/getThreshold", getThreshold)
	group.GET("/warnStatistics", warnStatistics)
	group.GET("/warnList", warnList)
	group.POST("/savePushConfig", savePushConfig)
	group.POST("/saveThreshold", saveThreshold)
}

// @Summary getPushConfig
// @Tags Warn
// @Param name formData string true " smtp|wechat|dingtalk "
// @Router /api/v1/warn/getPushConfig [get]
func getPushConfig(c *gin.Context) {
	query, ok := c.GetQuery("name")
	if !ok {
		response.ErrorParam(c, "name")
		return
	}
	var data interface{}
	switch query {
	case "smtp":
		smtp, err := (&dao.PushConfig{}).GetSMTPConfig()
		if err != nil {
			data = "{}"
			// response.Fail(c, consts.CurdSelectFailCode, consts.CurdSelectFailMsg, nil)
			// return
		} else {
			data = smtp
		}
	case "wechat":
		wechat, err := (&dao.PushConfig{}).GetWeChatConfig()
		if err != nil {
			data = "{}"
			// response.Fail(c, consts.CurdSelectFailCode, consts.CurdSelectFailMsg, nil)
			// return
		} else {
			data = wechat
		}
	case "dingtalk":
		ding, err := (&dao.PushConfig{}).GetDingTalkConfig()
		if err != nil {
			data = "{}"
			// response.Fail(c, consts.CurdSelectFailCode, consts.CurdSelectFailMsg, nil)
			// return
		} else {
			data = ding
		}
	default:
		response.ErrorParam(c, "name")
		return
	}
	response.Success(c, consts.Success, data)
}

// @Summary  getThreshold
// @Tags  warn
// @Router /api/v1/warn/getThreshold [get]
func getThreshold(c *gin.Context) {
	cpu, _ := (&dao.Config{}).GetSystemConfig(consts.HostCpuThreshold)
	memory, _ := (&dao.Config{}).GetSystemConfig(consts.HostMemoryThreshold)
	disk, _ := (&dao.Config{}).GetSystemConfig(consts.HostDiskThreshold)
	v1, _ := strconv.Atoi(cpu)
	v2, _ := strconv.Atoi(memory)
	v3, _ := strconv.Atoi(disk)
	response.Success(c, consts.Success, gin.H{
		"CPU":    v1,
		"Memory": v2,
		"Disk":   v3,
	})
}

// @Summary  warnStatistics
// @Tags  warn
// @Param beginTime formData string true " beginTime "
// @Param endTime formData string true " endTime "
// @Router /api/v1/warn/warnStatistics [get]
func warnStatistics(context *gin.Context) {
	beginTime := context.DefaultQuery("beginTime", "")
	endTime := context.DefaultQuery("endTime", "")
	list, err := taos_warn.GetWarnStatistics(beginTime, endTime)
	if err != nil {
		response.Error(context, err)
		return
	}
	response.Success(context, consts.CurdStatusOkMsg, list)
}

// @Summary  warnList
// @Tags  warn
// @Param type formData string true " type "
// @Param beginTime formData string true " beginTime "
// @Param endTime formData string true " endTime "
// @Router /api/v1/warn/warnList [get]
func warnList(context *gin.Context) {
	wType := context.DefaultQuery("type", "")
	beginTime := context.DefaultQuery("beginTime", "")
	endTime := context.DefaultQuery("endTime", "")
	var pageSize = 10
	var pageIndex = 1
	if size := context.Request.FormValue("pageSize"); size != "" {
		pageSize, _ = tools.StringToInt(size)
	}

	if index := context.Request.FormValue("pageNum"); index != "" {
		pageIndex, _ = tools.StringToInt(index)
	}
	list, count, err := taos_warn.GetWarnList(wType, beginTime, endTime, pageSize, pageIndex)
	if err != nil {
		response.Error(context, err)
		return
	}
	response.Success(context, consts.CurdStatusOkMsg, response.PageRes{List: list, Count: count, PageIndex: pageIndex, PageSize: pageSize})
}

// @Summary  savePushConfig
// @Tags  warn
// @Param PushConfig body notice.PushConfig true " object "
// @Router /api/v1/warn/savePushConfig [POST]
func savePushConfig(c *gin.Context) {
	var config notice.PushConfig
	err := c.BindJSON(&config)
	tools.HasError(err, " Data parsing failed ", -1)

	if config.SMTP != (notice.SMTP{}) {
		variable.EmailNotice, _ = notice.CreateEmailNotice(config.SMTP)
		temp, _ := json.Marshal(config.SMTP)
		err = (&dao.Config{}).AddAndSetSystemConfig("smtp_push_config", string(temp))
	}
	if config.WeChat != (notice.WeChat{}) {
		variable.WeChatNotice, _ = notice.CreateWeChatNotice(config.WeChat, 3*time.Second)
		temp, _ := json.Marshal(config.WeChat)
		err = (&dao.Config{}).AddAndSetSystemConfig("wechat_push_config", string(temp))
	}
	if config.DingTalk != (notice.DingTalk{}) {
		variable.DingTalkNotice, _ = notice.NewDingTalkNotice(config.DingTalk, 3*time.Second)
		temp, _ := json.Marshal(config.DingTalk)
		err = (&dao.Config{}).AddAndSetSystemConfig("ding_push_config", string(temp))
	}

	if err != nil {
		response.Fail(c, consts.CurdUpdateFailCode, consts.CurdUpdateFailMsg, "")
		return
	}
	response.Success(c, consts.CurdStatusOkMsg, nil)
}

type Threshold struct {
	CPU    int
	Memory int
	Disk   int
}

// @Summary  saveThreshold
// @Tags  warn
// @Param Threshold body Threshold true " object "
// @Router /api/v1/warn/saveThreshold [POST]
func saveThreshold(c *gin.Context) {
	var config Threshold
	err := c.BindJSON(&config)
	tools.HasError(err, " Data parsing failed ", -1)

	if config.CPU != 0 {
		err = (&dao.Config{}).AddAndSetSystemConfig(consts.HostCpuThreshold, strconv.Itoa(config.CPU))
	}
	if config.Memory != 0 {
		err = (&dao.Config{}).AddAndSetSystemConfig(consts.HostMemoryThreshold, strconv.Itoa(config.Memory))
	}
	if config.Disk != 0 {
		err = (&dao.Config{}).AddAndSetSystemConfig(consts.HostDiskThreshold, strconv.Itoa(config.Disk))
	}

	if err != nil {
		response.Fail(c, consts.CurdUpdateFailCode, consts.CurdUpdateFailMsg, "")
		return
	}
	if config.CPU != 0 {
		variable.CpuThreshold = config.CPU
	}
	if config.Memory != 0 {
		variable.MemThreshold = config.Memory
	}
	if config.Disk != 0 {
		variable.DiskThreshold = config.Disk
	}

	response.Success(c, consts.CurdStatusOkMsg, nil)
}
