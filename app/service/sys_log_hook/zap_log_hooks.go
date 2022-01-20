package sys_log_hook

import (
	"go.uber.org/zap/zapcore"
)

// edgelog  System operation log hook function 
// 1. A single log is a structure format ， This function intercepts each log ， You can perform subsequent processing ， for example ： Push to alicloud log management panel 、ElasticSearch  Log library, etc 

func ZapLogHandler(entry zapcore.Entry) error {

	//  parameter  entry  introduce 
	// entry   Parameter is a single log structure ， The main fields are as follows ：
	//Level       Log level 
	//Time        current time  
	//LoggerName   Log name 
	//Message     Log content 
	//Caller      Call path of each file 
	//Stack       Code call stack 

	// Start a collaborative process here ，hook It will not affect program performance at all ，
	go func(paramEntry zapcore.Entry) {
		//fmt.Println(" edgelog  hook ....， You can continue to process system logs here ....")
		//fmt.Printf("%#+v\n", paramEntry)
	}(entry)
	return nil
}
