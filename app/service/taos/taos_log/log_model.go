package taos_log

type LogModel struct {
	Time string
	Log  string
}

type NodeStatusModel struct {
	Time   string
	Type   int8
	Status int8
}
