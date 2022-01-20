package bootstrap

import (
	_ "edgelog/app/core/destroy" //  Listener exit signal ， For resource release 
	"edgelog/app/global/my_errors"
	"edgelog/app/global/variable"
	"edgelog/app/http/validator/common/register_validator"
	"edgelog/app/service/agent"
	"edgelog/app/service/init_data"
	"edgelog/app/service/sys_log_hook"
	"edgelog/app/service/taos"
	"edgelog/app/utils/gorm_v2"
	"edgelog/app/utils/ip"
	"edgelog/app/utils/snow_flake"
	"edgelog/app/utils/yml_config"
	"edgelog/app/utils/zap_factory"
	"log"
	"os"
)

//  Check whether the non compiled directory required by the project exists ， Avoid missing directory when compiled. 
func checkRequiredFolders() {
	//1. Check whether the configuration file exists 
	if _, err := os.Stat(variable.BasePath + "/config/config.yml"); err != nil {
		log.Fatal(my_errors.ErrorsConfigYamlNotExists + err.Error())
	}
	if _, err := os.Stat(variable.BasePath + "/config/gorm_v2.yml"); err != nil {
		log.Fatal(my_errors.ErrorsConfigGormNotExists + err.Error())
	}
	//2. inspect public Does the directory exist 
	if _, err := os.Stat(variable.BasePath + "/public/"); err != nil {
		log.Fatal(my_errors.ErrorsPublicNotExists + err.Error())
	}
	//3. inspect storage/logs  Does the directory exist 
	if _, err := os.Stat(variable.BasePath + "/storage/logs/"); err != nil {
		if err := os.MkdirAll("/storage/logs/", os.ModePerm); err != nil {
			log.Fatal(my_errors.ErrorsStorageLogsNotExists + err.Error())
		}
	}
	// 4. Automatically create soft connections 、 Better management of static resources 
	if _, err := os.Stat(variable.BasePath + "/public/storage"); err == nil {
		if err = os.Remove(variable.BasePath + "/public/storage"); err != nil {
			log.Fatal(my_errors.ErrorsSoftLinkDeleteFail + err.Error())
		}
	}
	if err := os.Symlink(variable.BasePath+"/storage/app", variable.BasePath+"/public/storage"); err != nil {
		log.Fatal(my_errors.ErrorsSoftLinkCreateFail + err.Error())
	}
}

func init() {
	// 1.  initialization   Project root path ， See  variable  Constant package ， Correlation path ：app\global\variable\variable.go

	//2. Check the necessary conditions for non compilation such as configuration files and log directories 
	checkRequiredFolders()

	//3. Initialize form parameter validator ， Register in container （Web、Api Shared container ）
	register_validator.WebRegisterValidator()
	register_validator.ApiRegisterValidator()

	// 4. Launch for profile (confgi.yml、gorm_v2.yml) Change monitoring ，  Profile action pointer ， Initialize as global variable 
	variable.ConfigYml = yml_config.CreateYamlFactory()
	variable.ConfigYml.ConfigFileChangeListen()
	// config>gorm_v2.yml  Start file change listening event 
	variable.ConfigGormv2Yml = variable.ConfigYml.Clone("gorm_v2")
	variable.ConfigGormv2Yml.ConfigFileChangeListen()

	// 5. Initialize global log handle ， And load the log hook processing function 
	variable.ZapLog = zap_factory.CreateZapFactory(sys_log_hook.ZapLogHandler)

	// 6. Initialize according to configuration  gorm mysql  overall situation  *gorm.Db
	if variable.ConfigGormv2Yml.GetInt("Gormv2.Mysql.IsInitGolobalGormMysql") == 1 {
		if dbMysql, err := gorm_v2.GetOneMysqlClient(); err != nil {
			log.Fatal(my_errors.ErrorsGormInitFail + err.Error())
		} else {
			variable.GormDb = dbMysql
			init_data.InitData()
		}
	}

	//Initialize ipstore
	variable.IpStore = ip.NewIpStore("storage/ip.dat")

	// 7. Snowflake algorithm global variable 
	variable.SnowFlake = snow_flake.CreateSnowflakeFactory()

	go func() {
		//Initialize Taos database
		taos.InitTaos()
		//Log receiving service
		(&agent.LogServer{}).Start()
	}()
	go (&agent.MonitorServer{}).Start()
}
