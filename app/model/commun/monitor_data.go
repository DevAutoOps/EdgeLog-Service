package commun

//Probe monitoring data
type MonitorData struct {
	Type  uint8
	Value int32
	Ts    string
}

//Probe monitoring data list
type MonitorDataList struct {
	HasData bool
	List    []MonitorData
}

//Monitoring chart data
type MonitorChart struct {
	Title string    `json:"title"` //title
	X     []string  `json:"x"`     //X-axis, time data
	Y     []float64 `json:"y"`     //Y-axis, data
}

//Monitoring chart data 2
type MonitorChart2 struct {
	Title string   `json:"title"` //title
	X     []string `json:"x"`     //X-axis, time data
	Y     []string `json:"y"`     //Y-axis, data
}

// chart data struct
type ChartData struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type ChartPointData []ChartData

//Monitoring chart data 3
type MonitorChart3 struct {
	Title string           `json:"title"` //title
	X     []string         `json:"x"`     //X-axis, time data
	Y     []ChartPointData `json:"y"`     //Y-axis, data
}

//Manager monitoring screen
type ManagerDataMap map[string]interface{}

//Single nginx concurrency
type NginxConcurrencyData struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

//Early warning information
type MonitorWarnInfo struct {
	Name string `json:"name"`
	Ip   string `json:"ip"`
	Type string `json:"type"`
	Msg  string `json:"msg"`
	Time string `json:"time"`
}

type MonitorValue struct {
	Value string
	Ts    string
}

type MonitorValueList struct {
	Type    uint8
	HasData bool
	List    []MonitorValue
}

type WarnPieItem struct {
	Name       string  `json:"name"`
	Count      int     `json:"count"`
	Percentage float64 `json:"percentage"`
}

type WarnChart struct {
	Column MonitorChart2 `json:"column"`
	Pie    []WarnPieItem `json:"pie"`
}
