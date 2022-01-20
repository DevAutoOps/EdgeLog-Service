package agent

type Server struct {
	Ip          string `yaml:"IP" config:"ServerIP"`
	LogPort     string `yaml:"LogPort" config:"LogPort"`         //Receive log port
	MonitorPort string `yaml:"MonitorPort" config:"MonitorPort"` //Receive monitoring data port
}

type Logs struct {
	Enable              bool   `yaml:"Enable" config:"LogsEnable"`
	CloseStandardOutput bool   `yaml:"CloseStandardOutput" config:"CloseStandardOutput"`
	LogName             string `yaml:"LogName" config:"LogName"`
	TextFormat          string `yaml:"TextFormat" config:"LogsTextFormat"`
	TimePrecision       string `yaml:"TimePrecision" config:"LogsTimePrecision"`
	MaxSize             int    `yaml:"MaxSize" config:"LogsMaxSize"`
	MaxBackups          int    `yaml:"MaxBackups" config:"LogsMaxBackups"`
	MaxAge              int    `yaml:"MaxAge" config:"LogsMaxAge"`
	Compress            bool   `yaml:"Compress" config:"LogsCompress"`
	PrintLevel          uint   `yaml:"PrintLevel" config:"LogsPrintLevel"`
}

type Modules struct {
	LogCollector LogCollector `yaml:"LogCollector" config:"LogCollector"` //Log collector
	Monitor      Monitor      `yaml:"Monitor" config:"Monitor"`           //Server monitoring
}

type LogCollector struct {
	Enable      bool     `yaml:"Enable" config:"LogCollectorEnable"`           //Enable log collection
	Targets     []string `yaml:"Targets" config:"LogCollectorTargets"`         //Target object for collecting logs
	TargetConfs []string `yaml:"TargetConfs" config:"LogCollectorTargetConfs"` //The configuration file corresponding to the target object for collecting logs
}

type Monitor struct {
	Enable     bool     `yaml:"Enable" config:"MonitorEnable"`         // Enable server monitoring
	Freq       uint     `yaml:"Freq" config:"Freq"`                    // Unit: ms, execution frequency
	Check      bool     `yaml:"Check" config:"MonitorCheck"`           // Whether to detect the server hardware information. If false, only the survival status of the server will be detected
	CheckFreq  uint     `yaml:"CheckFreq" config:"CheckFreq"`          // Once every n tests
	Apps       []string `yaml:"Apps" config:"MonitorApps"`             // App to be monitored
	AppProcess []string `yaml:"AppProcess" config:"MonitorAppProcess"` // Corresponding app process
}

type Config struct {
	Server  Server  `yaml:"Server" config:"Server"`   //Server related properties
	Modules Modules `yaml:"Modules" config:"Modules"` //Modular
	Logs    Logs    `yaml:"Logs"`                     //Log correlation
}
