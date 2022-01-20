package zap_factory

import (
	"edgelog/app/global/variable"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"time"
)

func CreateZapFactory(entry func(zapcore.Entry) error) *zap.Logger {

	//  Gets the mode the program is in ：   Development and debugging  、  production 
	//variable.ConfigYml := yml_config.CreateYamlFactory()
	appDebug := variable.ConfigYml.GetBool("AppDebug")

	//  Determine the current mode of the program ， Debug mode directly returns a convenient zap Log manager address ， All logs can be printed to the console 
	if appDebug == true {
		if logger, err := zap.NewDevelopment(zap.Hooks(entry)); err == nil {
			return logger
		} else {
			log.Fatal(" establish zap Log package failed ， details ：" + err.Error())
		}
	}

	//  The following is   Non debugging （ production ） Code required for pattern 
	encoderConfig := zap.NewProductionEncoderConfig()

	timePrecision := variable.ConfigYml.GetString("Logs.TimePrecision")
	var recordTimeFormat string
	switch timePrecision {
	case "second":
		recordTimeFormat = "2006-01-02 15:04:05"
	case "millisecond":
		recordTimeFormat = "2006-01-02 15:04:05.000"
	default:
		recordTimeFormat = "2006-01-02 15:04:05"

	}
	encoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format(recordTimeFormat))
	}
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.TimeKey = "created_at" //  generate json Time key field of format log ， Default to  ts, After modification, it is convenient to import logs into  ELK  The server 

	var encoder zapcore.Encoder
	switch variable.ConfigYml.GetString("Logs.TextFormat") {
	case "console":
		encoder = zapcore.NewConsoleEncoder(encoderConfig) //  Normal mode 
	case "json":
		encoder = zapcore.NewJSONEncoder(encoderConfig) // json format 
	default:
		encoder = zapcore.NewConsoleEncoder(encoderConfig) //  Normal mode 
	}

	// Writer 
	fileName := variable.BasePath + variable.ConfigYml.GetString("Logs.edgelogLogName")
	lumberJackLogger := &lumberjack.Logger{
		Filename:   fileName,                                     // Location of log files 
		MaxSize:    variable.ConfigYml.GetInt("Logs.MaxSize"),    // Before cutting ， Maximum size of log file （ with MB Unit ）
		MaxBackups: variable.ConfigYml.GetInt("Logs.MaxBackups"), // Maximum number of old files to keep 
		MaxAge:     variable.ConfigYml.GetInt("Logs.MaxAge"),     // Maximum number of days to keep old files 
		Compress:   variable.ConfigYml.GetBool("Logs.Compress"),  // Compress / Archive old files 
	}
	writer := zapcore.AddSync(lumberJackLogger)
	//  Start initialization zap Log core parameters ，
	// Parameter one ： encoder 
	// Parameter two ： Writer 
	// Parameter three ： Parameter level ，debug Level supports the logging of all subsequent functions ， If it is  fatal  High level ， Then level >=fatal  Before you can write the log 
	zapCore := zapcore.NewCore(encoder, writer, zap.InfoLevel)
	return zap.New(zapCore, zap.AddCaller(), zap.Hooks(entry))
}
