package snow_flake

import (
	"edgelog/app/global/consts"
	"edgelog/app/global/variable"
	"edgelog/app/utils/snow_flake/snowflake_interf"
	"sync"
	"time"
)

//  Create a snowflake algorithm generator ( Build factory )
func CreateSnowflakeFactory() snowflake_interf.InterfaceSnowFlake {
	return &snowflake{
		timestamp: 0,
		machineId: variable.ConfigYml.GetInt64("SnowFlake.SnowFlakeMachineId"),
		sequence:  0,
	}
}

type snowflake struct {
	sync.Mutex
	timestamp int64
	machineId int64
	sequence  int64
}

//  Generate distributed ID
func (s *snowflake) GetId() int64 {
	s.Lock()
	defer func() {
		s.Unlock()
	}()
	now := time.Now().UnixNano() / 1e6
	if s.timestamp == now {
		s.sequence = (s.sequence + 1) & consts.SequenceMask
		if s.sequence == 0 {
			for now <= s.timestamp {
				now = time.Now().UnixNano() / 1e6
			}
		}
	} else {
		s.sequence = 0
	}
	s.timestamp = now
	r := (now-consts.StartTimeStamp)<<consts.TimestampShift | (s.machineId << consts.MachineIdShift) | (s.sequence)
	return r
}
