package handle

import (
	"edgelog/app/global/variable"
	"edgelog/app/table"
	"edgelog/routers/common"
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func HandleTemplate(group *gin.RouterGroup) {
	group.Use(common.TokenMid)
	group.GET("/manageList", manageList)
	group.GET("/dropDown", dropDown)
	group.POST("/save", save)
	group.POST("/update/:id", update)
	group.DELETE("/delete/:id", delete)
	group.GET("/info/:id", info)
}

// @Summary  template manageList
// @Tags  template
// @Param pageSize query string false " Entries per page "
// @Param pageNum query string false " Page number "
// @Router /api/v1/template/manageList [get]
func manageList(c *gin.Context) {
	common.StandardManageList(c, func() (*gorm.DB, interface{}, error) {
		db := variable.GormDb.Model(&table.Template{}).Where("type=?", 0)
		qName := c.DefaultQuery("name", "")
		if qName != "" {
			db = db.Where("name like ?", "%"+qName+"%")
		}
		formatStr := c.DefaultQuery("format", "")
		if formatStr != "" {
			qFormat, err := strconv.Atoi(formatStr)
			if err == nil {
				db = db.Where("format = ?", qFormat != 0)
			} else {
				return nil, nil, err
			}
		}
		return db, new([]table.Template), nil
	})
}

// @Summary  template dropDown
// @Tags  template
// @Router /api/v1/template/dropDown [get]
func dropDown(c *gin.Context) {
	result := make([]table.Template, 0)
	err := variable.GormDb.Model(&table.Template{}).Where("type=?", 0).Find(&result).Error
	if err != nil {
		common.Error(c, err)
		return
	}
	common.Ok(c, result)
}

// @Summary  template save
// @Tags  template
// @Param Template body table.Template true " object "
// @Router /api/v1/template/save [post]
func save(c *gin.Context) {
	template := &table.Template{}
	if err := c.ShouldBind(template); err != nil {
		common.Error(c, err)
		return
	}
	if template.Separator != "|" &&
		template.Separator != "#" &&
		template.Separator != "," &&
		template.Separator != " " {
		common.Error(c, errors.New("error separator"))
		return
	}
	if template.Value == "" {
		common.Error(c, errors.New("error value"))
		return
	}
	namecheck := table.Template{}
	variable.GormDb.Model(&table.Template{}).Where("name=?", template.Name).Find(&namecheck)
	if namecheck.ID != 0 {
		common.Error(c, errors.New("error name"))
		return
	}
	template.Type = 0
	err := variable.GormDb.Model(&table.Template{}).Create(template).Error
	if err != nil {
		common.Error(c, err)
		return
	}
	common.Ok(c, template)
}

// @Summary  template update
// @Tags  template
// @Param id path int true " template  ID"
// @Router /api/v1/template/update/:id [post]
func update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		common.Error(c, err)
		return
	}
	model := &table.Template{}
	if err := c.ShouldBind(model); err != nil {
		common.Error(c, err)
		return
	}
	res := variable.GormDb.Model(table.Template{}).
		Where("id=?", id).
		Updates(common.StructToMapFilter(model))
	if res.Error != nil {
		common.Error(c, res.Error)
		return
	}
	if res.RowsAffected == 0 {
		common.Error(c, errors.New("record not found"))
		return
	}
	common.Ok(c, model)
}

// @Summary  template delete
// @Tags  template
// @Param id path int true " template  ID"
// @Router /api/v1/template/delete/:id [delete]
func delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		common.Error(c, err)
		return
	}
	err = variable.GormDb.Where("id=?", id).
		Delete(table.Template{}).Error
	if err != nil {
		common.Error(c, err)
		return
	}
	common.Ok(c, id)
}

// @Summary  template info
// @Tags  template
// @Param id path int true " template  ID"
// @Router /api/v1/template/info/:id [get]
func info(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		common.Error(c, err)
		return
	}
	model := &table.Template{}
	if err := variable.GormDb.Model(table.Template{}).
		Where("id = ?", id).
		First(model).Error; err != nil {
		common.Error(c, err)
		return
	}
	common.Ok(c, model)
}
