package table

type Template struct {
	Model
	Name      string `gorm:"uniqueIndex;type:VARCHAR(50) NOT NULL default '' comment ' Template name '"`
	Value     string `gorm:"type:VARCHAR(1000) NOT NULL default '' comment ' Template content '"`
	Format    bool   `gorm:"type:TINYINT(1) UNSIGNED NOT NULL default '0' comment ' Log format ，0  default ，1 JSON'"`
	Type      int    `gorm:"type:TINYINT(3) UNSIGNED NOT NULL default '0' comment ' Log type ，0 Nginx，1 Nginx plus，2 Apache，3 Tomcat'"`
	Separator string `gorm:"type:VARCHAR(50) NOT NULL default '' comment ' Separator '"`
}

func (Template) TableName() string {
	return "template"
}
