package table

type Config struct {
	Model
	Item  string `gorm:"uniqueIndex;type:VARCHAR(255) NOT NULL default '' comment ' Configuration item '"`
	Value string `gorm:"type:VARCHAR(1000) NOT NULL default '' comment ' Configuration value '"`
}

func (Config) TableName() string {
	return "config"
}
