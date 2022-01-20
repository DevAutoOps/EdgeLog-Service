package init_data

import (
	"edgelog/app/global/variable"
	"edgelog/app/model"
)

func InitData() {
	InitDatabase()
	initNode()
}

func initNode() {
	variable.Node = &model.Node{
		Name:       "Demo host",
		Ip:         "42.192.173.120",
		Port:       9060,
		Account:    "root",
		Password:   "jV@DEN!mBIW$#mB8Nu2MW34",
		IsInit:     true,
		AgentPort:  20210,
		Status:     false,
		Os:         0,
		AppStatus:  "[{\"name\":\"Nginx\",\"status\":0}]",
		Conf:       "/usr/local/nginx/conf/nginx.conf",
		Logs:       "/usr/local/nginx/logs/",
		TemplateId: 1,
	}
}
