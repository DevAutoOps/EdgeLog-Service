package main

import (
	"edgelog/app/global/variable"
	_ "edgelog/bootstrap"
	"edgelog/routers"
)

//  Back end routes can be stored here （ For example, background management system ）
func main() {
	router := routers.InitWebRouter()
	_ = router.Run(variable.ConfigYml.GetString("HttpServer.Web.Port"))
}
