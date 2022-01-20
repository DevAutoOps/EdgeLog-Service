package taos_monitor

import (
	"edgelog/app/global/variable"
	"edgelog/app/model/commun"
	"edgelog/app/service/taos"
	"edgelog/app/service/taos/taos_log"
	"fmt"
	"go.uber.org/zap"
	"time"
)

func MonitorDataAdd(monitorDataList *commun.MonitorDataList) {
	if !monitorDataList.HasData || monitorDataList.List == nil || len(monitorDataList.List) <= 0 {
		variable.ZapLog.Warn("[Taos MonitorDataAdd] The data list is empty")
		return
	}
	if variable.TaosDb == nil {
		variable.ZapLog.Warn("[Taos] taos database is not initialized")
		return
	}
	taos.MonitorLock.Lock()
	defer taos.MonitorLock.Unlock()
	values := ""
	timespan := time.Now().UnixNano() / 1e6
	for _, item := range monitorDataList.List {
		values += fmt.Sprintf(" (%d,%d,%d)", timespan, item.Type, item.Value)
		timespan++
	}
	sql := fmt.Sprintf("insert into %s.%s values %s", taos.DataBase, taos.MonitorTableName, values)
	_, err := variable.TaosDb.Exec(sql)
	if err != nil {
		fmt.Printf("insert table %s data error: %s\n", taos.MonitorTableName, err.Error())
		variable.ZapLog.Error(fmt.Sprintf("[Taos] Insert table %s data error:", taos.MonitorTableName), zap.Error(err))
	}
	counter := len(monitorDataList.List) + 1
	time.Sleep(time.Duration(counter) * time.Millisecond)
}

func MonitorStatusAdd(nodeStatusChangeList []taos_log.NodeStatusModel) {
	if nodeStatusChangeList == nil || len(nodeStatusChangeList) <= 0 {
		variable.ZapLog.Warn("[Taos MonitorStatusAdd] The data list is empty")
		return
	}
	if variable.TaosDb == nil {
		variable.ZapLog.Warn("[Taos] taos database is not initialized")
		return
	}
	taos.NodeStatusLock.Lock()
	defer taos.NodeStatusLock.Unlock()

	values := ""
	timespan := time.Now().UnixNano() / 1e6
	for _, item := range nodeStatusChangeList {
		values += fmt.Sprintf(" (%d,%d,%d)", timespan, item.Type, item.Status)
		timespan++
	}
	sql := fmt.Sprintf("insert into %s.%s values %s", taos.DataBase, taos.NodeStatusTableName, values)
	_, err := variable.TaosDb.Exec(sql)
	if err != nil {
		fmt.Printf("insert table %s data error: %s\n", taos.NodeStatusTableName, err.Error())
		variable.ZapLog.Error(fmt.Sprintf("[Taos] Insert table %s data error:", taos.NodeStatusTableName), zap.Error(err))
	}
	counter := len(nodeStatusChangeList) + 1
	time.Sleep(time.Duration(counter) * time.Millisecond)
}
