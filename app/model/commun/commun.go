package commun

import "bytes"

//Data buffer
type DataBufferClient struct {
	DataBuffer *bytes.Buffer
	Length     uint32
	OriLength  uint32
	RemoteInfo string
}

//Receive data properties
type ReceiveDataProperties struct {
	IsBegin            bool
	IsFinish           bool
	TotalLength        uint32
	UnReadFinishLength int32
	CurrentLength      int32
	OriLength          uint32
	DataBuffer         *bytes.Buffer
}
