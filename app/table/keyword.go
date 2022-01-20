package table

type Keyword struct {
	Model
	Name  string `gorm:"type:VARCHAR(50) NOT NULL default '' comment ' Keyword name '"`
	Value string `gorm:"type:VARCHAR(50) NOT NULL default '' comment ' Corresponding value '"`
	Type  int    `gorm:"type:TINYINT(1) UNSIGNED NOT NULL default '0' comment ' type ，0  keyword , 1  Separator '"`
	Order int    `gorm:"type:INT(11) NOT NULL default '0' comment ' sort '"`
}

func (Keyword) TableName() string {
	return "keyword"
}
