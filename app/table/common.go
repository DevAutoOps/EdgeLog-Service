package table

import "time"

type Model struct {
	ID        uint      `gorm:"primarykey" json:"id" swaggerignore:"true"`
	CreatedAt time.Time `gorm:"default 'DEFAULT CURRENT_TIMESTAMP'" swaggerignore:"true"`
	UpdatedAt time.Time `gorm:"default 'ON UPDATE CURRENT_TIMESTAMP'" swaggerignore:"true"`
}

type ITable interface {
	TableName() string
}
