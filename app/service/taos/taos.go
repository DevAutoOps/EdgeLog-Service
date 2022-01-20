package taos

import (
	"database/sql"
	"edgelog/app/global/variable"
	"fmt"
	_ "github.com/taosdata/driver-go/taosSql"
	"go.uber.org/zap"
	"sync"
)

const driverName = "taosSql"

var (
	DataBase            = ""
	LogTableName        = "nginx_log"
	LogLock             = &sync.RWMutex{}
	MonitorTableName    = "monitor"
	MonitorLock         = &sync.RWMutex{}
	NodeStatusTableName = "node_status"
	NodeStatusLock      = &sync.RWMutex{}
)

func InitTaos() {
	if variable.TaosDb == nil {
		host := variable.ConfigYml.GetString("Taos.Host")
		port := variable.ConfigYml.GetInt("Taos.Port")
		user := variable.ConfigYml.GetString("Taos.User")
		pass := variable.ConfigYml.GetString("Taos.Pass")
		dataBase := variable.ConfigYml.GetString("Taos.DataBase")
		DataBase = dataBase

		url := fmt.Sprintf("%s:%s@tcp(%s:%v)/", user, pass, host, port)
		db, err := sql.Open(driverName, url)
		if err != nil {
			fmt.Printf("Open taos database error: %s\n", err.Error())
			variable.ZapLog.Error("[Taos] Open taos database error:", zap.Error(err))
		} else {
			variable.TaosDb = db
			createDatabase(dataBase)
			createTable(dataBase)
		}
	}
}

func createDatabase(database string) {
	sqlStr := fmt.Sprintf("create database if not exists %s", database)
	if variable.TaosDb != nil {
		_, err := variable.TaosDb.Exec(sqlStr)
		if err != nil {
			fmt.Printf("create database error: %s\n", err.Error())
			variable.ZapLog.Error("[Taos] Create database error:", zap.Error(err))
		}
	}
}

func createTable(database string) {
	sqlStr := fmt.Sprintf("create table if not exists %s.%s (ts timestamp,log BINARY(1000))", database, LogTableName)
	sqlStr2 := fmt.Sprintf("create table if not exists %s.%s (ts timestamp,type TINYINT,value int)", database, MonitorTableName)
	sqlStr3 := fmt.Sprintf("create table if not exists %s.%s (ts timestamp,type TINYINT,status TINYINT)", database, NodeStatusTableName)
	if variable.TaosDb != nil {
		_, err := variable.TaosDb.Exec(sqlStr)
		_, err = variable.TaosDb.Exec(sqlStr2)
		_, err = variable.TaosDb.Exec(sqlStr3)
		if err != nil {
			fmt.Printf("create table error: %s\n", err.Error())
			variable.ZapLog.Error("[Taos] Create table error:", zap.Error(err))
		}
	}
}

func InitTestBaseData() {
	if variable.TaosDb == nil {
		host := "42.192.173.120"
		port := 6030
		user := "root"
		pass := "taosdata"
		dataBase := "edgelog_demo"
		DataBase = dataBase

		url := fmt.Sprintf("%s:%s@tcp(%s:%v)/", user, pass, host, port)
		db, err := sql.Open(driverName, url)
		if err != nil {
			fmt.Printf("Open taos database error: %s\n", err.Error())
		} else {
			variable.TaosDb = db
		}
	}
}
