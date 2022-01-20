package variable

import (
	"database/sql"
	"edgelog/app/global/my_errors"
	"edgelog/app/model"
	"edgelog/app/utils/ip"
	"edgelog/app/utils/snow_flake/snowflake_interf"
	"edgelog/app/utils/yml_config/ymlconfig_interf"
	"edgelog/routers/common/notice"
	"log"
	"os"
	"strings"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

//  All encapsulated global variables support concurrency security ， Please feel free to use 
//  Developer encapsulated global variables ， Please check and confirm concurrency security 

var (
	BasePath           string                  //  Defines the root directory of the project 
	EventDestroyPrefix = "Destroy_"            //   Event prefix that needs to be destroyed when the program exits 
	ConfigKeyPrefix    = "Config_"             //   When configuring file key cache ， Prefix of key 
	DateFormat         = "2006-01-02 15:04:05" //   When configuring file key cache ， Prefix of key 

	//  Global log pointer 
	ZapLog *zap.Logger
	//  Global profile 
	ConfigYml       ymlconfig_interf.YmlConfigInterf //  Global profile pointer 
	ConfigGormv2Yml ymlconfig_interf.YmlConfigInterf //  Global profile pointer 

	//gorm  Database client ， If you operate the database using gorm， Please uncomment the following ， stay  bootstrap>init  file ， It can be used after initialization 
	GormDb *gorm.DB //  overall situation gorm Client connection for 

	// Snowflake algorithm global variable 
	SnowFlake snowflake_interf.InterfaceSnowFlake

	//taos database
	TaosDb *sql.DB

	// warning threshold
	CpuThreshold  = 0
	MemThreshold  = 0
	DiskThreshold = 0

	// notice
	EmailNotice    notice.INotice
	WeChatNotice   notice.INotice
	DingTalkNotice *notice.DingTalkNotice

	// local Node
	Node *model.Node

	//ip database
	IpStore *ip.IpStore
)

func init() {
	// 1. Initialize program root directory 
	if curPath, err := os.Getwd(); err == nil {
		//  Path processing ， Strange path when compatibility unit tester starts 
		if len(os.Args) > 1 && strings.HasPrefix(os.Args[1], "-test") {
			BasePath = strings.Replace(strings.Replace(curPath, `\test`, "", 1), `/test`, "", 1)
		} else {
			BasePath = curPath
		}
	} else {
		log.Fatal(my_errors.ErrorsBasePath)
	}
}
