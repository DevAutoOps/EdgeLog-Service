package handle

import (
	"bytes"
	"edgelog/app/global/variable"
	"edgelog/routers/common"
	"edgelog/routers/common/cache"
	"fmt"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
)

func HandleDownload(group *gin.RouterGroup) {
	group.GET("/download_agent", download_agent)
	group.GET("/download_log", download_log)
}

// @Summary  window agent download
// @Tags  download
// @Router /api/v1/download/download_agent [get]
func download_agent(c *gin.Context) {
	c.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", "agent_windows.zip"))
	c.Writer.Header().Add("Content-Type", "application/octet-stream")
	c.File(filepath.Join(variable.BasePath, "public", "install", "agent_windows.zip"))
}

// @Summary  log download
// @Tags  download
// @Param size query string false " size"
// @Router /api/v1/download/download_log [get]
func download_log(c *gin.Context) {
	size, err := strconv.Atoi(c.DefaultQuery("size", "50000"))
	if err != nil {
		common.Error(c, err)
		return
	}
	logs := make([]string, 0)
	err = cache.GetLatestNodeLogList2(size, 0, size, &logs)
	if err != nil {
		common.Error(c, err)
		return
	}
	buf := make([][]byte, len(logs))
	for i, log := range logs {
		buf[i] = []byte(log)
	}
	c.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", "log.log"))
	c.Writer.Header().Add("Content-Type", "application/octet-stream")
	c.Writer.Write(bytes.Join(buf, []byte("\n")))
}
