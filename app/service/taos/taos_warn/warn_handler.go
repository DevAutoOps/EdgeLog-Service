package taos_warn

import (
	"edgelog/app/dao"
	"edgelog/app/global/consts"
	"edgelog/app/global/variable"
	"edgelog/app/model"
	"edgelog/app/model/commun"
	"edgelog/app/service/taos"
	"edgelog/tools"
	"fmt"
	"go.uber.org/zap"
	"math"
	"sort"
	"strconv"
	"time"
)

func GetWarnStatistics(beginTime, endTime string) ([]commun.MonitorChart2, error) {
	result := make([]commun.MonitorChart2, 0)
	cpu, err := (&dao.Config{}).GetSystemConfig(consts.HostCpuThreshold)
	if err != nil {
		cpu = consts.HostDefaultThresholdValueStr
	}
	//cpu := consts.HostDefaultThresholdValueStr
	memory, err := (&dao.Config{}).GetSystemConfig(consts.HostMemoryThreshold)
	if err != nil {
		memory = consts.HostDefaultThresholdValueStr
	}
	//memory := consts.HostDefaultThresholdValueStr
	disk, err := (&dao.Config{}).GetSystemConfig(consts.HostDiskThreshold)
	if err != nil {
		disk = consts.HostDefaultThresholdValueStr
	}
	//disk := consts.HostDefaultThresholdValueStr
	sql := fmt.Sprintf("select count(ts) from %s.%s where ", taos.DataBase, taos.MonitorTableName)
	sqlCpu := sql
	sqlMem := sql
	sqlDisk := sql
	sqlAgentOnline := fmt.Sprintf("select count(ts) from %s.%s where ", taos.DataBase, taos.NodeStatusTableName)
	sqlAgentOffline := sqlAgentOnline
	dateList := make([]string, 0)

	if len(beginTime) > 0 && len(endTime) > 0 {
		sqlCpu = fmt.Sprintf("%s ts >= '%s' and ts <= '%s' and ", sqlCpu, beginTime, endTime)
		sqlMem = fmt.Sprintf("%s ts >= '%s' and ts <= '%s' and ", sqlMem, beginTime, endTime)
		sqlDisk = fmt.Sprintf("%s ts >= '%s' and ts <= '%s' and ", sqlDisk, beginTime, endTime)
		sqlAgentOnline = fmt.Sprintf("%s ts >= '%s' and ts <= '%s' and ", sqlAgentOnline, beginTime, endTime)
		sqlAgentOffline = fmt.Sprintf("%s ts >= '%s' and ts <= '%s' and ", sqlAgentOffline, beginTime, endTime)

		beginTimeObj, err := time.ParseInLocation("2006-01-02 15:04:05", beginTime, time.Local)
		if err != nil {
			return result, err
		}
		endTimeObj, err := time.ParseInLocation("2006-01-02 15:04:05", endTime, time.Local)
		if err != nil {
			return result, err
		}
		timeFormatTpl := "2006-01-02"
		endTimeStr := endTimeObj.Format(timeFormatTpl)
		dateList = append(dateList, beginTimeObj.Format(timeFormatTpl))
		for {
			beginTimeObj = beginTimeObj.AddDate(0, 0, 1)
			dateStr := beginTimeObj.Format(timeFormatTpl)
			dateList = append(dateList, dateStr)
			if dateStr == endTimeStr {
				break
			}
		}
	} else {
		sqlCpu = fmt.Sprintf("%s ts >= NOW-6d and ", sqlCpu)
		sqlMem = fmt.Sprintf("%s ts >= NOW-6d and ", sqlMem)
		sqlDisk = fmt.Sprintf("%s ts >= NOW-6d and ", sqlDisk)
		sqlAgentOnline = fmt.Sprintf("%s ts >= NOW-6d and ", sqlAgentOnline)
		sqlAgentOffline = fmt.Sprintf("%s ts >= NOW-6d and ", sqlAgentOffline)
		currentTime := time.Now()
		for i := 6; i >= 0; i-- {
			oldTime := currentTime.AddDate(0, 0, -i)
			oldTimeStr := oldTime.Format("2006-01-02")
			dateList = append(dateList, oldTimeStr)
		}
	}

	sqlCpu = fmt.Sprintf("%s type= %d and value >= %s*100 INTERVAL(1d)",
		sqlCpu, consts.MonitorCpuUsage, cpu)
	sqlMem = fmt.Sprintf("%s type= %d and value >= %s INTERVAL(1d)",
		sql, consts.MonitorMemRate, memory)
	sqlDisk = fmt.Sprintf("%s type= %d and value >= %s*100 INTERVAL(1d)",
		sql, consts.MonitorDiskPartRate, disk)
	sqlAgentOnline = fmt.Sprintf("%s type=0 and status=1 INTERVAL(1d)",
		sqlAgentOnline)
	sqlAgentOffline = fmt.Sprintf("%s type=0 and status=0 INTERVAL(1d)",
		sqlAgentOffline)

	totalCountMap := make(map[string]int)
	// CPU
	cpuMap := make(map[string]int)
	cpuChart := commun.MonitorChart2{Title: "CPU Early warning statistics", X: make([]string, 0), Y: make([]string, 0)}
	rowsCpu, err := variable.TaosDb.Query(sqlCpu)
	if err != nil {
		fmt.Printf("%s error: %s\n", sqlCpu, err.Error())
		variable.ZapLog.Error("[Taos GetWarnStatistics] "+sqlCpu+"  error:", zap.Error(err))
	} else {
		for rowsCpu.Next() {
			var ts string
			var countValue int
			err = rowsCpu.Scan(&ts, &countValue)
			if err != nil {
				fmt.Printf("scan table %s Cpu data error: %s\n", taos.MonitorTableName, err.Error())
				variable.ZapLog.Error("[Taos GetWarnStatistics] scan table "+taos.MonitorTableName+" Cpu data error:", zap.Error(err))
				continue
			}
			thisTime, err := time.ParseInLocation("2006-01-02 15:04:05", ts, time.Local)
			if err != nil {
				continue
			}
			timePeriod := thisTime.Format("2006-01-02")
			cpuMap[timePeriod] = countValue
		}
	}
	defer rowsCpu.Close()
	for i := 0; i < len(dateList); i++ {
		currentCount := 0
		if currentCountTemp, ok := cpuMap[dateList[i]]; ok {
			currentCount = currentCountTemp
		}
		totalCountMap[dateList[i]] = currentCount
		cpuChart.X = append(cpuChart.X, dateList[i])
		cpuChart.Y = append(cpuChart.Y, fmt.Sprintf("%d", currentCount))
	}

	// MEM
	memMap := make(map[string]int)
	memChart := commun.MonitorChart2{Title: "Memory alert statistics", X: make([]string, 0), Y: make([]string, 0)}
	rowsMem, err := variable.TaosDb.Query(sqlMem)
	if err != nil {
		fmt.Printf("%s error: %s\n", sqlCpu, err.Error())
		variable.ZapLog.Error("[Taos GetWarnStatistics] "+sqlMem+"  error:", zap.Error(err))
	} else {
		for rowsMem.Next() {
			var ts string
			var countValue int
			err = rowsMem.Scan(&ts, &countValue)
			if err != nil {
				fmt.Printf("scan table %s Mem data error: %s\n", taos.MonitorTableName, err.Error())
				variable.ZapLog.Error("[Taos GetWarnStatistics] scan table "+taos.MonitorTableName+" Mem data error:", zap.Error(err))
				continue
			}
			thisTime, err := time.ParseInLocation("2006-01-02 15:04:05", ts, time.Local)
			if err != nil {
				continue
			}
			timePeriod := thisTime.Format("2006-01-02")
			memMap[timePeriod] = countValue
		}
	}
	defer rowsMem.Close()
	for i := 0; i < len(dateList); i++ {
		currentCount := 0
		if currentCountTemp, ok := memMap[dateList[i]]; ok {
			currentCount = currentCountTemp
		}
		totalCountMap[dateList[i]] += currentCount
		memChart.X = append(memChart.X, dateList[i])
		memChart.Y = append(memChart.Y, fmt.Sprintf("%d", currentCount))
	}

	// Disk
	diskMap := make(map[string]int)
	diskChart := commun.MonitorChart2{Title: "Hard disk warning statistics", X: make([]string, 0), Y: make([]string, 0)}
	rowsDisk, err := variable.TaosDb.Query(sqlDisk)
	if err != nil {
		fmt.Printf("%s error: %s\n", sqlDisk, err.Error())
		variable.ZapLog.Error("[Taos GetWarnStatistics] "+sqlDisk+"  error:", zap.Error(err))
	} else {
		for rowsDisk.Next() {
			var ts string
			var countValue int
			err = rowsDisk.Scan(&ts, &countValue)
			if err != nil {
				fmt.Printf("scan table %s Disk data error: %s\n", taos.MonitorTableName, err.Error())
				variable.ZapLog.Error("[Taos GetWarnStatistics] scan table "+taos.MonitorTableName+" Disk data error:", zap.Error(err))
				continue
			}
			thisTime, err := time.ParseInLocation("2006-01-02 15:04:05", ts, time.Local)
			if err != nil {
				continue
			}
			timePeriod := thisTime.Format("2006-01-02")
			diskMap[timePeriod] = countValue
		}
	}
	defer rowsDisk.Close()
	for i := 0; i < len(dateList); i++ {
		currentCount := 0
		if currentCountTemp, ok := diskMap[dateList[i]]; ok {
			currentCount = currentCountTemp
		}
		totalCountMap[dateList[i]] += currentCount
		diskChart.X = append(diskChart.X, dateList[i])
		diskChart.Y = append(diskChart.Y, fmt.Sprintf("%d", currentCount))
	}

	// Offline Agent
	agentOfflineMap := make(map[string]int)
	agentChart := commun.MonitorChart2{Title: "agent Survival warning statistics", X: make([]string, 0), Y: make([]string, 0)}
	rowsAgentOffline, err := variable.TaosDb.Query(sqlAgentOffline)
	if err != nil {
		fmt.Printf("%s error: %s\n", sqlAgentOffline, err.Error())
		variable.ZapLog.Error("[Taos GetWarnStatistics] "+sqlAgentOffline+"  error:", zap.Error(err))
	} else {
		for rowsAgentOffline.Next() {
			var ts string
			var countValue int
			err = rowsAgentOffline.Scan(&ts, &countValue)
			if err != nil {
				fmt.Printf("scan table %s offline agent data error: %s\n", taos.MonitorTableName, err.Error())
				variable.ZapLog.Error("[Taos GetWarnStatistics] scan table "+taos.MonitorTableName+" offline agent data error:", zap.Error(err))
				continue
			}
			thisTime, err := time.ParseInLocation("2006-01-02 15:04:05", ts, time.Local)
			if err != nil {
				continue
			}
			timePeriod := thisTime.Format("2006-01-02")
			agentOfflineMap[timePeriod] = countValue
		}
	}
	defer rowsAgentOffline.Close()
	for i := 0; i < len(dateList); i++ {
		currentCount := 0
		if currentCountTemp, ok := agentOfflineMap[dateList[i]]; ok {
			currentCount = currentCountTemp
		}
		totalCountMap[dateList[i]] += currentCount
		agentChart.X = append(agentChart.X, dateList[i])
		agentChart.Y = append(agentChart.Y, fmt.Sprintf("%d", currentCount))
	}

	// Total
	totalChart := commun.MonitorChart2{Title: "Alert summary statistics", X: make([]string, 0), Y: make([]string, 0)}
	for i := 0; i < len(dateList); i++ {
		currentCount := 0
		if currentCountTemp, ok := totalCountMap[dateList[i]]; ok {
			currentCount = currentCountTemp
		}
		totalChart.X = append(totalChart.X, dateList[i])
		totalChart.Y = append(totalChart.Y, fmt.Sprintf("%d", currentCount))
	}

	result = append(result, totalChart, cpuChart, memChart, diskChart, agentChart)
	return result, nil
}

func GetWarnList(wType, beginTime, endTime string, pageSize, pageIndex int) (model.WarnInfoList, int, error) {
	result := make(model.WarnInfoList, 0)
	cpu, err := (&dao.Config{}).GetSystemConfig(consts.HostCpuThreshold)
	if err != nil {
		cpu = consts.HostDefaultThresholdValueStr
	}
	//cpu := defaultWarn
	memory, err := (&dao.Config{}).GetSystemConfig(consts.HostMemoryThreshold)
	if err != nil {
		memory = consts.HostDefaultThresholdValueStr
	}
	//memory := defaultWarn
	disk, err := (&dao.Config{}).GetSystemConfig(consts.HostDiskThreshold)
	if err != nil {
		disk = consts.HostDefaultThresholdValueStr
	}
	//disk := defaultWarn

	timeConf := ""
	if len(beginTime) > 0 && len(endTime) > 0 {
		timeConf = fmt.Sprintf(" ts >= '%s' and ts <= '%s' and ", beginTime, endTime)
	} else {
		timeConf = fmt.Sprintf(" ts >= NOW-1d and ")
	}

	sqlMonitor := fmt.Sprintf("(select ts,type,value from %s.%s where %s type = %d and value >= %s*100) UNION ALL (select ts,type,value from %s.%s where %s type = %d and value >= %s) UNION ALL (select ts,type,value from %s.%s where %s type = %d and value >= %s*100)",
		taos.DataBase, taos.MonitorTableName, timeConf, consts.MonitorCpuUsage, cpu,
		taos.DataBase, taos.MonitorTableName, timeConf, consts.MonitorMemRate, memory,
		taos.DataBase, taos.MonitorTableName, timeConf, consts.MonitorDiskPartRate, disk)

	sqlAgent := fmt.Sprintf("select ts,type,status from %s.%s where %s status = 0",
		taos.DataBase, taos.NodeStatusTableName, timeConf)

	monitorDataList := make(model.WarnInfoList, 0)
	agentDataList := make(model.WarnInfoList, 0)
	if len(wType) <= 0 || wType == "0" {
		rowsMonitor, err := variable.TaosDb.Query(sqlMonitor)
		if err != nil {
			fmt.Printf("%s error: %s\n", sqlMonitor, err.Error())
			variable.ZapLog.Error("[Taos GetWarnList] "+sqlMonitor+"  error:", zap.Error(err))
		} else {
			for rowsMonitor.Next() {
				var ts string
				var dType int
				var value int
				err = rowsMonitor.Scan(&ts, &dType, &value)
				if err != nil {
					fmt.Printf("scan table %s data error: %s\n", taos.MonitorTableName, err.Error())
					variable.ZapLog.Error("[Taos GetWarnStatistics] scan table "+taos.MonitorTableName+" data error:", zap.Error(err))
					continue
				}
				thisTime, err := time.ParseInLocation("2006-01-02 15:04:05", ts, time.Local)
				if err != nil {
					continue
				}
				timePeriod := thisTime.Format("2006-01-02 15:04")

				info := ""
				switch uint8(dType) {
				case consts.MonitorCpuUsage:
					info = fmt.Sprintf("CPU Utilization has reached %s%%", tools.FormatFloat(float64(value)/100, 1))
				case consts.MonitorMemRate:
					info = fmt.Sprintf(" Memory usage reached %d%%", value)
				case consts.MonitorDiskPartRate:
					info = fmt.Sprintf(" Disk utilization has reached %s%%", tools.FormatFloat(float64(value)/100, 1))
				}
				monitorDataList = append(monitorDataList, model.WarnInfo{
					NodeName:    variable.Node.Name,
					NodeIp:      variable.Node.Ip,
					Type:        "Host resources",
					MonitorType: dType,
					WarnTime:    timePeriod,
					Ts:          thisTime.Unix(),
					Info:        info,
				})
			}
		}
		defer rowsMonitor.Close()
	}
	if len(wType) <= 0 || wType == "1" {
		rowsAgent, err := variable.TaosDb.Query(sqlAgent)
		if err != nil {
			fmt.Printf("%s error: %s\n", sqlAgent, err.Error())
			variable.ZapLog.Error("[Taos GetWarnList] "+sqlAgent+" error:", zap.Error(err))
		} else {
			for rowsAgent.Next() {
				var ts string
				var dType int
				var status int
				err = rowsAgent.Scan(&ts, &dType, &status)
				if err != nil {
					fmt.Printf("scan table %s data error: %s\n", taos.NodeStatusTableName, err.Error())
					variable.ZapLog.Error("[Taos GetWarnStatistics] scan table "+taos.NodeStatusTableName+" data error:", zap.Error(err))
					continue
				}
				thisTime, err := time.ParseInLocation("2006-01-02 15:04:05", ts, time.Local)
				if err != nil {
					continue
				}
				timePeriod := thisTime.Format("2006-01-02 15:04")

				info := ""
				switch int8(dType) {
				case consts.StatusNode:
					info = fmt.Sprintf("This node is offline ")
				case consts.StatusNginx:
					info = fmt.Sprintf("On this node Nginx Offline ")
				}
				agentDataList = append(agentDataList, model.WarnInfo{
					NodeName:       variable.Node.Name,
					NodeIp:         variable.Node.Ip,
					Type:           "Agent survival",
					NodeStatusType: dType,
					WarnTime:       timePeriod,
					Ts:             thisTime.Unix(),
					Info:           info,
				})
			}
		}
		defer rowsAgent.Close()
	}

	totalDataList := append(monitorDataList, agentDataList...)
	if len(totalDataList) <= 0 {
		return totalDataList, 0, nil
	}

	sort.Sort(totalDataList)
	count := len(totalDataList)
	offset := (pageIndex - 1) * pageSize
	end := offset + pageSize
	if offset >= count {
		return result, count, nil
	}
	if offset < count && end > count {
		end = count
	}
	result = totalDataList[offset:end]

	return result, count, nil
}

func timeTypeToSql(timeType int) (timeConf, granularity string, timePoints []string) {
	timePoints = make([]string, 0)
	currentTime := time.Now()
	switch timeType {
	case 1:
		today := currentTime.Format("2006-01-02")
		today = fmt.Sprintf("%s 00:00:00", today)
		timeConf = fmt.Sprintf(" ts >= '%s' ", today)
		granularity = "2h"
		timePoints = getTimePoint2(currentTime)
	case 2:
		timeConf = " ts >= NOW-6d "
		granularity = "1d"
		for i := 6; i >= 0; i-- {
			oldTime := currentTime.AddDate(0, 0, -i)
			oldTimeStr := oldTime.Format("2006-01-02")
			timePoints = append(timePoints, oldTimeStr)
		}
	default:
		timeConf = " ts >= NOW-1h "
		granularity = "5m"
		timePoints = getTimePoint1(currentTime)
	}
	return
}

func getTimePoint1(currentTime time.Time) (timePoints []string) {
	timePoints = make([]string, 0)
	day := currentTime.Format("2006-01-02")
	hour := currentTime.Format("15")
	minute := currentTime.Format("04")
	minuteLastDigit := minute[len(minute)-1:]
	minuteLastDigitNumber, _ := strconv.Atoi(minuteLastDigit)
	if minuteLastDigitNumber >= 5 {
		minuteLastDigitNumber = 5
	} else {
		minuteLastDigitNumber = 0
	}
	startMinute := ""
	if len(minute) == 2 {
		startMinute = fmt.Sprintf("%s%d", minute[:1], minuteLastDigitNumber)
	} else {
		startMinute = fmt.Sprintf("0%d", minuteLastDigitNumber)
	}
	startTimeStr := fmt.Sprintf("%s %s:%s:00", day, hour, startMinute)
	startTime, err := time.ParseInLocation("2006-01-02 15:04:05", startTimeStr, time.Local)
	if err != nil {
		fmt.Println(err)
		return
	}
	for i := 12; i > 0; i-- {
		count := time.Duration(i)
		oldTime := startTime.Add(-5 * time.Minute * count)
		oldTimeStr := oldTime.Format("2006-01-02 15:04")
		timePoints = append(timePoints, oldTimeStr)
	}
	timePoints = append(timePoints, fmt.Sprintf("%s %s:%s", day, hour, startMinute))
	return
}

func getTimePoint2(currentTime time.Time) (timePoints []string) {
	timePoints = make([]string, 0)
	day := currentTime.Format("2006-01-02")
	hour := currentTime.Format("15")
	startHour, _ := strconv.Atoi(hour)
	if startHour%2 != 0 {
		startHour--
	}
	startTimeStr := fmt.Sprintf("%s %d:00:00", day, startHour)
	startTime, err := time.ParseInLocation("2006-01-02 15:04:05", startTimeStr, time.Local)
	if err != nil {
		fmt.Println(err)
		return
	}

	for i := startHour / 2; i > 0; i-- {
		count := time.Duration(i)
		oldTime := startTime.Add(-2 * time.Hour * count)
		oldTimeStr := oldTime.Format("2006-01-02 15:04")
		timePoints = append(timePoints, oldTimeStr)
	}
	timePoints = append(timePoints, fmt.Sprintf("%s %d:00", day, startHour))
	return
}

func GetWarnBigScreen(timeType int) (commun.WarnChart, error) {
	timeConf, granularity, dateList := timeTypeToSql(timeType)
	result := commun.WarnChart{}
	cpu, err := (&dao.Config{}).GetSystemConfig(consts.HostCpuThreshold)
	if err != nil {
		cpu = consts.HostDefaultThresholdValueStr
	}
	//cpu = "10"
	memory, err := (&dao.Config{}).GetSystemConfig(consts.HostMemoryThreshold)
	if err != nil {
		memory = consts.HostDefaultThresholdValueStr
	}
	//memory = "20"
	disk, err := (&dao.Config{}).GetSystemConfig(consts.HostDiskThreshold)
	if err != nil {
		disk = consts.HostDefaultThresholdValueStr
	}
	//disk = "40"

	sql := fmt.Sprintf("select count(ts) from %s.%s where ", taos.DataBase, taos.MonitorTableName)
	sqlCpu := fmt.Sprintf("%s %s and type= %d and value >= %s*100 INTERVAL(%s)",
		sql, timeConf, consts.MonitorCpuUsage, cpu, granularity)
	sqlMem := fmt.Sprintf("%s %s and type= %d and value >= %s INTERVAL(%s)",
		sql, timeConf, consts.MonitorMemRate, memory, granularity)
	sqlDisk := fmt.Sprintf("%s %s and type= %d and value >= %s*100 INTERVAL(%s)",
		sql, timeConf, consts.MonitorDiskPartRate, disk, granularity)
	sqlAgentOffline := fmt.Sprintf("select count(ts) from %s.%s where %s and type=0 and status=0 INTERVAL(%s)",
		taos.DataBase, taos.NodeStatusTableName, timeConf, granularity)

	totalCountMap := make(map[string]int)
	// CPU
	cpuMap := make(map[string]int)
	cpuChart := commun.WarnPieItem{Name: "CPU"}
	rowsCpu, err := variable.TaosDb.Query(sqlCpu)
	if err != nil {
		fmt.Printf("%s error: %s\n", sqlCpu, err.Error())
		variable.ZapLog.Error("[Taos GetWarnStatistics] "+sqlCpu+"  error:", zap.Error(err))
	} else {
		for rowsCpu.Next() {
			var ts string
			var countValue int
			err = rowsCpu.Scan(&ts, &countValue)
			if err != nil {
				fmt.Printf("scan table %s Cpu data error: %s\n", taos.MonitorTableName, err.Error())
				variable.ZapLog.Error("[Taos GetWarnStatistics] scan table "+taos.MonitorTableName+" Cpu data error:", zap.Error(err))
				continue
			}
			thisTime, err := time.ParseInLocation("2006-01-02 15:04:05", ts, time.Local)
			if err != nil {
				continue
			}
			timePeriod := ""
			if timeType == 2 {
				timePeriod = thisTime.Format("2006-01-02")
			} else {
				timePeriod = thisTime.Format("2006-01-02 15:04")
			}
			cpuMap[timePeriod] = countValue
		}
	}
	defer rowsCpu.Close()
	cpuTotalCount := 0
	for i := 0; i < len(dateList); i++ {
		currentCount := 0
		if currentCountTemp, ok := cpuMap[dateList[i]]; ok {
			currentCount = currentCountTemp
		}
		totalCountMap[dateList[i]] = currentCount
		cpuTotalCount += currentCount
	}
	cpuChart.Count = cpuTotalCount

	// MEM
	memMap := make(map[string]int)
	memChart := commun.WarnPieItem{Name: "mem"}
	rowsMem, err := variable.TaosDb.Query(sqlMem)
	if err != nil {
		fmt.Printf("%s error: %s\n", sqlCpu, err.Error())
		variable.ZapLog.Error("[Taos GetWarnStatistics] "+sqlMem+"  error:", zap.Error(err))
	} else {
		for rowsMem.Next() {
			var ts string
			var countValue int
			err = rowsMem.Scan(&ts, &countValue)
			if err != nil {
				fmt.Printf("scan table %s Mem data error: %s\n", taos.MonitorTableName, err.Error())
				variable.ZapLog.Error("[Taos GetWarnStatistics] scan table "+taos.MonitorTableName+" Mem data error:", zap.Error(err))
				continue
			}
			thisTime, err := time.ParseInLocation("2006-01-02 15:04:05", ts, time.Local)
			if err != nil {
				continue
			}
			timePeriod := ""
			if timeType == 2 {
				timePeriod = thisTime.Format("2006-01-02")
			} else {
				timePeriod = thisTime.Format("2006-01-02 15:04")
			}
			memMap[timePeriod] = countValue
		}
	}
	defer rowsMem.Close()
	memTotalCount := 0
	for i := 0; i < len(dateList); i++ {
		currentCount := 0
		if currentCountTemp, ok := memMap[dateList[i]]; ok {
			currentCount = currentCountTemp
		}
		totalCountMap[dateList[i]] += currentCount
		memTotalCount += currentCount
	}
	memChart.Count = memTotalCount

	// Disk
	diskMap := make(map[string]int)
	diskChart := commun.WarnPieItem{Name: "disk"}
	rowsDisk, err := variable.TaosDb.Query(sqlDisk)
	if err != nil {
		fmt.Printf("%s error: %s\n", sqlDisk, err.Error())
		variable.ZapLog.Error("[Taos GetWarnStatistics] "+sqlDisk+"  error:", zap.Error(err))
	} else {
		for rowsDisk.Next() {
			var ts string
			var countValue int
			err = rowsDisk.Scan(&ts, &countValue)
			if err != nil {
				fmt.Printf("scan table %s Disk data error: %s\n", taos.MonitorTableName, err.Error())
				variable.ZapLog.Error("[Taos GetWarnStatistics] scan table "+taos.MonitorTableName+" Disk data error:", zap.Error(err))
				continue
			}
			thisTime, err := time.ParseInLocation("2006-01-02 15:04:05", ts, time.Local)
			if err != nil {
				continue
			}
			timePeriod := ""
			if timeType == 2 {
				timePeriod = thisTime.Format("2006-01-02")
			} else {
				timePeriod = thisTime.Format("2006-01-02 15:04")
			}
			diskMap[timePeriod] = countValue
		}
	}
	defer rowsDisk.Close()
	diskTotalCount := 0
	for i := 0; i < len(dateList); i++ {
		currentCount := 0
		if currentCountTemp, ok := diskMap[dateList[i]]; ok {
			currentCount = currentCountTemp
		}
		totalCountMap[dateList[i]] += currentCount
		diskTotalCount += currentCount
	}
	diskChart.Count = diskTotalCount

	// Offline Agent
	agentOfflineMap := make(map[string]int)
	agentChart := commun.WarnPieItem{Name: "agent"}
	rowsAgentOffline, err := variable.TaosDb.Query(sqlAgentOffline)
	if err != nil {
		fmt.Printf("%s error: %s\n", sqlAgentOffline, err.Error())
		variable.ZapLog.Error("[Taos GetWarnStatistics] "+sqlAgentOffline+"  error:", zap.Error(err))
	} else {
		for rowsAgentOffline.Next() {
			var ts string
			var countValue int
			err = rowsAgentOffline.Scan(&ts, &countValue)
			if err != nil {
				fmt.Printf("scan table %s offline agent data error: %s\n", taos.MonitorTableName, err.Error())
				variable.ZapLog.Error("[Taos GetWarnStatistics] scan table "+taos.MonitorTableName+" offline agent data error:", zap.Error(err))
				continue
			}
			thisTime, err := time.ParseInLocation("2006-01-02 15:04:05", ts, time.Local)
			if err != nil {
				continue
			}
			timePeriod := ""
			if timeType == 2 {
				timePeriod = thisTime.Format("2006-01-02")
			} else {
				timePeriod = thisTime.Format("2006-01-02 15:04")
			}
			agentOfflineMap[timePeriod] = countValue
		}
	}
	defer rowsAgentOffline.Close()
	agentOfflineTotalCount := 0
	for i := 0; i < len(dateList); i++ {
		currentCount := 0
		if currentCountTemp, ok := agentOfflineMap[dateList[i]]; ok {
			currentCount = currentCountTemp
		}
		totalCountMap[dateList[i]] += currentCount
		agentOfflineTotalCount += currentCount
	}
	agentChart.Count = agentOfflineTotalCount

	// Total
	totalChart := commun.MonitorChart2{X: make([]string, 0), Y: make([]string, 0)}
	for i := 0; i < len(dateList); i++ {
		currentCount := 0
		if currentCountTemp, ok := totalCountMap[dateList[i]]; ok {
			currentCount = currentCountTemp
		}
		totalChart.X = append(totalChart.X, dateList[i])
		totalChart.Y = append(totalChart.Y, fmt.Sprintf("%d", currentCount))
	}
	result.Column = totalChart

	total := cpuChart.Count + memChart.Count + diskChart.Count + agentChart.Count
	if total == 0 {
		cpuChart.Percentage = 0
		memChart.Percentage = 0
		diskChart.Percentage = 0
		agentChart.Percentage = 0
	} else {
		cpuChart.Percentage = math.Round(float64(cpuChart.Count) / float64(total) * 100)
		memChart.Percentage = math.Round(float64(memChart.Count) / float64(total) * 100)
		diskChart.Percentage = math.Round(float64(diskChart.Count) / float64(total) * 100)
		agentChart.Percentage = math.Round(float64(agentChart.Count) / float64(total) * 100)
	}
	result.Pie = append(result.Pie, cpuChart, memChart, diskChart, agentChart)

	return result, nil
}
