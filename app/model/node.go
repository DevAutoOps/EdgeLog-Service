package model

type Node struct {
	Name       string `json:"name"`
	Ip         string `json:"ip"`
	Port       int    `json:"port"`
	Account    string `json:"account"`
	Password   string `json:"password"`
	IsInit     bool   `json:"isInit"`
	AgentPort  int    `json:"agentPort"`
	Status     bool   `json:"status"`
	Os         int    `json:"os"`
	AppStatus  string `json:"appStatus"`
	Conf       string `json:"conf"`
	Logs       string `json:"logs"`
	TemplateId uint   `json:"templateId"`
}
