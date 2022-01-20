package dao

import (
	"edgelog/app/global/variable"
	"edgelog/app/table"
)

type Keyword struct {
}

func (k *Keyword) Create(keyword *table.Keyword) (err error) {
	if err = variable.GormDb.Model(&table.Keyword{}).Create(&keyword).Error; err != nil {
		return
	}
	return
}

func (k *Keyword) GetList() (list []table.Keyword, err error) {
	list = make([]table.Keyword, 0)
	err = variable.GormDb.Model(&table.Keyword{}).Find(&list).Error
	return
}
