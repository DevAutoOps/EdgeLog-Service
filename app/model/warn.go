package model

type WarnInfo struct {
	Id             uint
	NodeName       string `json:"nodeName"`
	NodeIp         string `json:"nodeIp"`
	Type           string `json:"type"`
	MonitorType    int    `json:"-"`
	NodeStatusType int    `json:"-"`
	WarnTime       string `json:"warnTime"`
	Ts             int64  `json:"-"`
	Info           string `json:"info"`
}

type WarnInfoList []WarnInfo

func (n WarnInfoList) Len() int {
	return len(n)
}

func (n WarnInfoList) Less(i, j int) bool {
	return n[i].Ts > n[j].Ts
}

func (n WarnInfoList) Swap(i, j int) {
	n[i], n[j] = n[j], n[i]
}
