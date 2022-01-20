package handle

import (
	"edgelog/app/global/variable"
	"edgelog/app/table"
	"edgelog/routers/common"

	"github.com/gin-gonic/gin"
)

func HandleKeyword(group *gin.RouterGroup) {
	group.Use(common.TokenMid)
	group.GET("/all", all)
}

// @Summary  All keywords
// @Tags keyword
// @Router /api/v1/keyword/all [get]
func all(c *gin.Context) {
	result := make([]table.Keyword, 0)
	err := variable.GormDb.Model(table.Keyword{}).
		Find(&result).Error
	if err != nil {
		common.Error(c, err)
		return
	}
	common.Ok(c, result)
}
