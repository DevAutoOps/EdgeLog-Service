package main

import (
	"edgelog/app/global/variable"
	_ "edgelog/bootstrap"
	"edgelog/routers"
)

//  Portal portal can be stored here 
func main() {
	router := routers.InitApiRouter()
	_ = router.Run(variable.ConfigYml.GetString("HttpServer.Api.Port"))
}
