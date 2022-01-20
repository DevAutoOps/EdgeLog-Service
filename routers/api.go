package routers

import (
	"edgelog/app/global/consts"
	"edgelog/app/global/variable"
	"edgelog/app/http/middleware/cors"
	validatorFactory "edgelog/app/http/validator/core/factory"
	_ "edgelog/docs"
	"edgelog/routers/handle"
	"io"
	"net/http"
	"os"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

//  This route is mainly used to set foreground routes such as portal websites 

func InitApiRouter() *gin.Engine {
	var router *gin.Engine
	//  Non debug mode （ Production mode ）  Log write to log file 
	if variable.ConfigYml.GetBool("AppDebug") == false {
		//1. Write log to log file 
		gin.DisableConsoleColor()
		f, _ := os.Create(variable.BasePath + variable.ConfigYml.GetString("Logs.GinLogName"))
		gin.DefaultWriter = io.MultiWriter(f)
		// 2. If yes nginx Acting as agent ， Basically not required gin Framework logging access log ， Open the following line of code ， Mask the above three lines of code ， Performance improvement  5%
		//gin.SetMode(gin.ReleaseMode)

		router = gin.Default()
	} else {
		//  Debug mode ， open  pprof  package ， It is convenient to analyze program performance in development stage 
		router = gin.Default()
		pprof.Register(router)
	}

	// Set cross domain according to configuration 
	if variable.ConfigYml.GetBool("HttpServer.AllowCrossDomain") {
		router.Use(cors.Next())
	}

	router.GET("/", func(context *gin.Context) {
		context.String(http.StatusOK, "Api  Module interface  hello word！")
	})

	//  Document interface access URL
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Processing static resources （ Not recommended gin The framework handles static resources ）
	router.Static("/public", "./public") //   Define the mapping relationship between static resource routing and actual directory 

	//   Create a portal class interface routing group 
	vApi := router.Group("/api/v1/")
	{
		//  Simulate a home page routing 
		home := vApi.Group("home/")
		{
			//  Second parameter description ：
			// 1. It is a form parameter validator function code snippet ， This function parses from the container ， The whole code segment is slightly complex ， But for users ， You only need to know the usage ， It's easy to use ， Look down  ↓↓↓
			// 2. Write a verifier for this interface ， position ：app/http/validator/api/home/news.go
			// 3. Register the above validators in the container ：app/http/validator/common/register_validator/api_register_validator.go  18  Key at behavior registration （consts.ValidatorPrefix + "HomeNews"）。 Then you can use this key to get from the container 
			home.GET("news", validatorFactory.Create(consts.ValidatorPrefix+"HomeNews"))
		}
		handle.HandleUser(vApi.Group("user"))
		handle.HandleMonitor(vApi.Group("monitor"))
		handle.HandleWarn(vApi.Group("warn"))
		handle.HandleTemplate(vApi.Group("template"))
		handle.HandleNode(vApi.Group("node"))
		handle.HandleAnalysis(vApi.Group("analysis"))
		handle.HandleDownload(vApi.Group("download"))
		handle.HandleKeyword(vApi.Group("keyword"))
		handle.HandleBigScreen(vApi.Group("bigscreen"))
	}
	return router
}
