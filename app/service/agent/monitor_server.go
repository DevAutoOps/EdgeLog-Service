package agent

import (
	"bufio"
	"bytes"
	"edgelog/app/global/consts"
	"edgelog/app/global/my_errors"
	"edgelog/app/global/variable"
	"edgelog/app/model/agent"
	"edgelog/app/model/commun"
	"edgelog/app/service/taos/taos_log"
	"edgelog/app/service/taos/taos_monitor"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"io"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

const monitorPrefix = "pbn"

type MonitorServer struct {
	listen           net.Listener
	isConnected      bool
	crontab          *cron.Cron
	taskIds          map[uint16]cron.EntryID
	curNodeStatusMap map[int8]int8
	nodeStatusMap    map[int8]int8
	sync.RWMutex
	statusLock *sync.RWMutex
}

func (m *MonitorServer) Start() {
	time.Sleep(10 * time.Second)
	m.isConnected = false
	m.curNodeStatusMap = make(map[int8]int8)
	m.statusLock = &sync.RWMutex{}
	port := variable.ConfigYml.GetInt("HttpServer.ProbesMonitorPort")
	var err error
	m.listen, err = net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		variable.ZapLog.Error("[MonitorServer Start] listen failed, err: ", zap.Error(err))
		return
	}
	variable.ZapLog.Info("[MonitorServer Start] listen ", zap.String("host:port", ":"+strconv.Itoa(port)))
	m.collector()

	for {
		conn, err := m.listen.Accept()
		if err != nil {
			variable.ZapLog.Error("[MonitorServer Start] listen accept failed, err:", zap.Error(err))
			time.Sleep(5 * time.Second)
			continue
		}
		if len(variable.ConfigYml.GetString("Agent.Ip")) > 0 && strings.HasPrefix(conn.RemoteAddr().String(), variable.ConfigYml.GetString("Agent.Ip")) {
			variable.ZapLog.Info(fmt.Sprintf("[MonitorServer Start]  And  %s establish connection", conn.RemoteAddr()))
			go m.process(&conn)
		}
	}
}

func (m *MonitorServer) Stop() {
	_ = m.listen.Close()
	m.isConnected = false
	m.listen = nil
	m.crontab.Stop()
	m.crontab = nil
	delete(m.taskIds, 0)
	m.taskIds = nil
	m.curNodeStatusMap = nil
	m.nodeStatusMap = nil
	m.statusLock = nil
}

func (m *MonitorServer) process(conn *net.Conn) {
	defer (*conn).Close()
	for {
		reader := bufio.NewReader(*conn)
		monitorDataList, err := m.parse(reader, conn)
		if err == io.EOF {
			variable.ZapLog.Error("[MonitorServer process] The agent is disconnected from the server, error: ", zap.Error(err))
			m.isConnected = false
			return
		}
		if err != nil {
			variable.ZapLog.Error("[MonitorServer process] parse data failed, err:", zap.Error(err))
			m.isConnected = false
			return
		}
		if monitorDataList != nil && len(monitorDataList.List) > 0 {
			m.isConnected = true
			m.processData(monitorDataList)
		}
	}
}

func (m *MonitorServer) parse(reader *bufio.Reader, conn *net.Conn) (*commun.MonitorDataList, error) {
	//Read prefix
	prefixByte := make([]byte, len(monitorPrefix))
	_, err := reader.Read(prefixByte)
	if err != nil {
		return nil, err
	}
	if string(prefixByte) != monitorPrefix {
		fmt.Printf("monitor server (%s) parse monitor data prefix unequal\n", (*conn).RemoteAddr())
		return nil, errors.New(my_errors.ErrorMonitorDataPrefixError)
	}

	countByte := make([]byte, 1)
	_, err = reader.Read(countByte)
	if err != nil {
		return nil, err
	}
	countBuff := bytes.NewBuffer(countByte)
	var count uint8
	err = binary.Read(countBuff, binary.LittleEndian, &count)
	if err != nil {
		return nil, err
	}

	m.statusLock.Lock()
	if m.curNodeStatusMap == nil {
		m.curNodeStatusMap = make(map[int8]int8)
	}
	m.curNodeStatusMap[consts.StatusNode] = 1
	m.statusLock.Unlock()
	if count == 0 {
		return &commun.MonitorDataList{HasData: false}, nil
	} else {
		lengthByte := make([]byte, 4)
		_, err = reader.Read(lengthByte)
		if err != nil {
			variable.ZapLog.Error("[MonitorServer parse] parse data length failed, err:", zap.Error(err))
			return &commun.MonitorDataList{HasData: false}, nil
		}
		lengthBuff := bytes.NewBuffer(lengthByte)
		var length uint32
		err = binary.Read(lengthBuff, binary.LittleEndian, &length)
		if err != nil {
			variable.ZapLog.Error("[MonitorServer parse] parse data length failed, err:", zap.Error(err))
			return &commun.MonitorDataList{HasData: false}, nil
		}

		if uint32(reader.Buffered()) < length {
			variable.ZapLog.Error("[MonitorServer parse] reader buffered size is smaller than length")
			return &commun.MonitorDataList{HasData: false}, nil
		}

		dataList := &commun.MonitorDataList{HasData: true}
		var i uint8 = 0
		for i = 0; i < count; i++ {
			typeByte := make([]byte, 1)
			_, err = reader.Read(typeByte)
			if err != nil {
				variable.ZapLog.Error("[MonitorServer parse] parse data type failed, err:", zap.Error(err))
				break
			}
			typeBuff := bytes.NewBuffer(typeByte)
			var dType uint8
			err = binary.Read(typeBuff, binary.LittleEndian, &dType)
			if err != nil {
				variable.ZapLog.Error("[MonitorServer parse] parse data type failed, err:", zap.Error(err))
				break
			}

			valueByte := make([]byte, 4)
			_, err = reader.Read(valueByte)
			if err != nil {
				variable.ZapLog.Error("[MonitorServer parse] parse data value failed, err:", zap.Error(err))
				break
			}
			valueBuff := bytes.NewBuffer(valueByte)
			var intValue int32
			err = binary.Read(valueBuff, binary.LittleEndian, &intValue)
			if err != nil {
				variable.ZapLog.Error("[MonitorServer parse] parse data value failed, err:", zap.Error(err))
				break
			}
			data := commun.MonitorData{Type: dType, Value: intValue}
			if dataList.List == nil {
				dataList.List = make([]commun.MonitorData, 0)
			}

			if data.Type > 127 {
				m.statusLock.Lock()
				if m.curNodeStatusMap == nil {
					m.curNodeStatusMap = make(map[int8]int8)
				}
				switch data.Type - 127 {
				case uint8(consts.StatusNginx):
					m.curNodeStatusMap[consts.StatusNginx] = int8(data.Value)
				}
				m.statusLock.Unlock()
				continue
			}

			dataList.List = append(dataList.List, data)
		}
		if dataList.List == nil || len(dataList.List) <= 0 {
			dataList.HasData = false
		}
		return dataList, nil
	}
}

func (m *MonitorServer) processData(monitorDataList *commun.MonitorDataList) {
	if !monitorDataList.HasData || monitorDataList.List == nil || len(monitorDataList.List) <= 0 || variable.TaosDb == nil {
		return
	}
	taos_monitor.MonitorDataAdd(monitorDataList)
}

func (m *MonitorServer) collector() {
	if m.crontab == nil {
		m.crontab = cron.New(cron.WithSeconds())
	}
	spec := fmt.Sprintf("*/%d * * * * ?", consts.CollectingMonitorInterval)
	taskId, err := m.crontab.AddFunc(spec, m.collectingData)
	if err != nil {
		variable.ZapLog.Error("[MonitorServer collector] add task error: ", zap.Error(err))
		return
	}
	if m.taskIds == nil {
		m.taskIds = make(map[uint16]cron.EntryID)
	}
	m.taskIds[0] = taskId
	m.crontab.Start()
}

func (m *MonitorServer) collectingData() {
	m.nodeStatusMap = make(map[int8]int8)
	if !variable.Node.IsInit {
		m.nodeStatusMap[consts.StatusNode] = -1
	} else {
		if variable.Node.Status {
			m.nodeStatusMap[consts.StatusNode] = 1
		} else {
			m.nodeStatusMap[consts.StatusNode] = 0
		}
	}

	if len(variable.Node.AppStatus) > 0 {
		var appStatusList []agent.NodeAppStatus
		err := json.Unmarshal([]byte(variable.Node.AppStatus), &appStatusList)
		if err != nil {
			m.nodeStatusMap[consts.StatusNginx] = 0
		} else {
			for _, appStatus := range appStatusList {
				if appStatus.Name == "Nginx" {
					m.nodeStatusMap[consts.StatusNginx] = int8(appStatus.Status)
					break
				}
			}
			if _, ok := m.nodeStatusMap[consts.StatusNginx]; !ok {
				m.nodeStatusMap[consts.StatusNginx] = 0
			}
		}
	} else {
		m.nodeStatusMap[consts.StatusNginx] = 0
	}

	currentNodeStatusMap := make(map[int8]int8)
	m.statusLock.Lock()
	if len(m.curNodeStatusMap) > 0 {
		for statusType, statusValue := range m.curNodeStatusMap {
			currentNodeStatusMap[statusType] = statusValue
		}
	}
	m.curNodeStatusMap = make(map[int8]int8)
	m.statusLock.Unlock()

	if len(m.nodeStatusMap) > 0 {
		nodeStatusChangeList := make([]taos_log.NodeStatusModel, 0)
		for statusType, statusValue := range m.nodeStatusMap {
			if currentStatus, ok := currentNodeStatusMap[statusType]; ok {
				if statusValue != currentStatus {
					nodeStatusChangeList = append(nodeStatusChangeList, taos_log.NodeStatusModel{
						Type: statusType, Status: currentStatus,
					})
					// save current status value
					m.nodeStatusMap[statusType] = currentStatus
					m.saveStatusChanged(statusType, currentStatus)
				}
			} else {
				if statusValue > 0 {
					nodeStatusChangeList = append(nodeStatusChangeList, taos_log.NodeStatusModel{
						Type: statusType, Status: 0,
					})
					// save current status value
					m.nodeStatusMap[statusType] = 0
					m.saveStatusChanged(statusType, 0)
				}
			}
		}

		if m.isConnected && len(currentNodeStatusMap) <= 0 {
			if m.nodeStatusMap[consts.StatusNode] == 0 {
				nodeStatusChangeList = append(nodeStatusChangeList, taos_log.NodeStatusModel{
					Type: consts.StatusNode, Status: 1,
				})
				// save current status value
				m.nodeStatusMap[consts.StatusNode] = 1
			}
		}
		if !m.isConnected {
			if m.nodeStatusMap[consts.StatusNode] > 0 {
				nodeStatusChangeList = append(nodeStatusChangeList, taos_log.NodeStatusModel{
					Type: consts.StatusNode, Status: 0,
				})
				// save current status value
				m.nodeStatusMap[consts.StatusNode] = 0
			}
		}

		m.saveStatusChangedForTaos(nodeStatusChangeList)
	}

	m.saveNodeStatusChanged()
}

func (m *MonitorServer) saveStatusChanged(appType int8, status int8) {
	if appType == consts.StatusNode {
		return
	}
	var appStatusList []agent.NodeAppStatus
	switch appType {
	case consts.StatusNginx:
		appStatusList = append(appStatusList, agent.NodeAppStatus{Name: "Nginx", Status: int(status)})
	}
	if len(appStatusList) > 0 {
		appStatusByte, err := json.Marshal(appStatusList)
		if err != nil {
			variable.ZapLog.Error("[MonitorServer saveStatusChanged] convert json error:", zap.Error(err))
			return
		}
		variable.Node.AppStatus = fmt.Sprintf("%s", appStatusByte)
	}
}

func (m *MonitorServer) saveStatusChangedForTaos(nodeStatusChangeList []taos_log.NodeStatusModel) {
	if len(nodeStatusChangeList) > 0 {
		taos_monitor.MonitorStatusAdd(nodeStatusChangeList)
	}
}

func (m *MonitorServer) saveNodeStatusChanged() {
	if m.isConnected {
		variable.Node.Status = true
	} else {
		variable.Node.Status = false
	}
}
