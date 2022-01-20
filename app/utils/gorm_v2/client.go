package gorm_v2

import (
	"edgelog/app/global/my_errors"
	"edgelog/app/global/variable"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormLog "gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"
	"strings"
	"time"
)

//  Get a  mysql  client 
func GetOneMysqlClient() (*gorm.DB, error) {
	sqlType := "Mysql"
	readDbIsOpen := variable.ConfigGormv2Yml.GetInt("Gormv2." + sqlType + ".IsOpenReadDb")
	return GetSqlDriver(sqlType, readDbIsOpen)
}

//  Get database driver ,  Can pass options  Connect any number of databases with dynamic parameters 
func GetSqlDriver(sqlType string, readDbIsOpen int, dbConf ...ConfigParams) (*gorm.DB, error) {

	var dbDialector gorm.Dialector
	if val, err := getDbDialector(sqlType, "Write", dbConf...); err != nil {
		variable.ZapLog.Error(my_errors.ErrorsDialectorDbInitFail+sqlType, zap.Error(err))
	} else {
		dbDialector = val
	}
	gormDb, err := gorm.Open(dbDialector, &gorm.Config{
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
		Logger:                 redefineLog(sqlType), // intercept 、 take over  gorm v2  Own log 
	})
	if err != nil {
		//gorm  Database driver initialization failed 
		return nil, err
	}

	//  If read / write separation is enabled ， Configure read database （resource、read、replicas）
	//  Read write separation configuration only 
	if readDbIsOpen == 1 {
		if val, err := getDbDialector(sqlType, "Read", dbConf...); err != nil {
			variable.ZapLog.Error(my_errors.ErrorsDialectorDbInitFail+sqlType, zap.Error(err))
		} else {
			dbDialector = val
		}
		resolverConf := dbresolver.Config{
			Replicas: []gorm.Dialector{dbDialector}, //   read   Operation Library ， Query class 
			Policy:   dbresolver.RandomPolicy{},     // sources/replicas  Load balancing policy applies to 
		}
		err = gormDb.Use(dbresolver.Register(resolverConf).SetConnMaxIdleTime(time.Second * 30).
			SetConnMaxLifetime(variable.ConfigGormv2Yml.GetDuration("Gormv2."+sqlType+".Read.SetConnMaxLifetime") * time.Second).
			SetMaxIdleConns(variable.ConfigGormv2Yml.GetInt("Gormv2." + sqlType + ".Read.SetMaxIdleConns")).
			SetMaxOpenConns(variable.ConfigGormv2Yml.GetInt("Gormv2." + sqlType + ".Read.SetMaxOpenConns")))
		if err != nil {
			return nil, err
		}
	}

	//  Query has no data ， shield  gorm v2  An error will pop up in the bag 
	// https://github.com/go-gorm/gorm/issues/3789   this  issue  The problem reflected is what we have solved this time 
	_ = gormDb.Callback().Query().Before("gorm:query").Register("disable_raise_record_not_found", func(d *gorm.DB) {
		d.Statement.RaiseErrorOnNotFound = false
	})

	//  Set connection pool for primary connection (43 Database driven pointer returned by row )
	if rawDb, err := gormDb.DB(); err != nil {
		return nil, err
	} else {
		rawDb.SetConnMaxIdleTime(time.Second * 30)
		rawDb.SetConnMaxLifetime(variable.ConfigGormv2Yml.GetDuration("Gormv2."+sqlType+".Write.SetConnMaxLifetime") * time.Second)
		rawDb.SetMaxIdleConns(variable.ConfigGormv2Yml.GetInt("Gormv2." + sqlType + ".Write.SetMaxIdleConns"))
		rawDb.SetMaxOpenConns(variable.ConfigGormv2Yml.GetInt("Gormv2." + sqlType + ".Write.SetMaxOpenConns"))
		return gormDb, nil
	}
}

//  Get a database dialect (Dialector), Generally speaking, it is based on different connection parameters ， Get the connection pointer of a specific type of database 
func getDbDialector(sqlType, readWrite string, dbConf ...ConfigParams) (gorm.Dialector, error) {
	var dbDialector gorm.Dialector
	dsn := getDsn(sqlType, readWrite, dbConf...)
	switch strings.ToLower(sqlType) {
	case "mysql":
		dbDialector = mysql.Open(dsn)
	default:
		return nil, errors.New(my_errors.ErrorsDbDriverNotExists + sqlType)
	}
	return dbDialector, nil
}

//   Generate database driver according to configuration parameters  dsn
func getDsn(sqlType, readWrite string, dbConf ...ConfigParams) string {
	Host := variable.ConfigGormv2Yml.GetString("Gormv2." + sqlType + "." + readWrite + ".Host")
	DataBase := variable.ConfigGormv2Yml.GetString("Gormv2." + sqlType + "." + readWrite + ".DataBase")
	Port := variable.ConfigGormv2Yml.GetInt("Gormv2." + sqlType + "." + readWrite + ".Port")
	User := variable.ConfigGormv2Yml.GetString("Gormv2." + sqlType + "." + readWrite + ".User")
	Pass := variable.ConfigGormv2Yml.GetString("Gormv2." + sqlType + "." + readWrite + ".Pass")
	Charset := variable.ConfigGormv2Yml.GetString("Gormv2." + sqlType + "." + readWrite + ".Charset")

	if len(dbConf) > 0 {
		if strings.ToLower(readWrite) == "write" {
			if len(dbConf[0].Write.Host) > 0 {
				Host = dbConf[0].Write.Host
			}
			if len(dbConf[0].Write.DataBase) > 0 {
				DataBase = dbConf[0].Write.DataBase
			}
			if dbConf[0].Write.Port > 0 {
				Port = dbConf[0].Write.Port
			}
			if len(dbConf[0].Write.User) > 0 {
				User = dbConf[0].Write.User
			}
			if len(dbConf[0].Write.Pass) > 0 {
				Pass = dbConf[0].Write.Pass
			}
			if len(dbConf[0].Write.Charset) > 0 {
				Charset = dbConf[0].Write.Charset
			}
		} else {
			if len(dbConf[0].Read.Host) > 0 {
				Host = dbConf[0].Read.Host
			}
			if len(dbConf[0].Read.DataBase) > 0 {
				DataBase = dbConf[0].Read.DataBase
			}
			if dbConf[0].Read.Port > 0 {
				Port = dbConf[0].Read.Port
			}
			if len(dbConf[0].Read.User) > 0 {
				User = dbConf[0].Read.User
			}
			if len(dbConf[0].Read.Pass) > 0 {
				Pass = dbConf[0].Read.Pass
			}
			if len(dbConf[0].Read.Charset) > 0 {
				Charset = dbConf[0].Read.Charset
			}
		}
	}

	switch strings.ToLower(sqlType) {
	case "mysql":
		return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local", User, Pass, Host, Port, DataBase, Charset)
	}
	return ""
}

//  Create custom log module ， yes  gorm  Log interception 、
func redefineLog(sqlType string) gormLog.Interface {
	return createCustomGormLog(sqlType,
		SetInfoStrFormat("[info] %s\n"), SetWarnStrFormat("[warn] %s\n"), SetErrStrFormat("[error] %s\n"),
		SetTraceStrFormat("[traceStr] %s [%.3fms] [rows:%v] %s\n"), SetTracWarnStrFormat("[traceWarn] %s %s [%.3fms] [rows:%v] %s\n"), SetTracErrStrFormat("[traceErr] %s %s [%.3fms] [rows:%v] %s\n"))
}
