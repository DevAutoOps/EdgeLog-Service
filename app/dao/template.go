package dao

import (
	"edgelog/app/global/variable"
	"edgelog/app/table"
)

type Template struct {
}

func (k *Template) GetDefault() (template table.Template, err error) {
	err = variable.GormDb.Model(&table.Template{}).Find(&template).Where("id=1").Error
	return
}

func (k *Template) Create(template *table.Template) (err error) {
	if err = variable.GormDb.Model(&table.Template{}).Create(&template).Error; err != nil {
		return
	}
	return
}
