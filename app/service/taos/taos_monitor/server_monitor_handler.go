package taos_monitor

import (
	"edgelog/app/global/consts"
	"edgelog/app/global/my_errors"
	"edgelog/app/global/variable"
	"edgelog/app/model/commun"
	"edgelog/app/service/taos"
	"edgelog/tools"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"time"
)

func dTypeToTitle(dType int) string {
	result := ""
	switch uint8(dType) {
	case consts.MonitorCpuUsage:
		result = "CPU Utilization rate (%)"
	case consts.MonitorCpuLoad:
		result = "CPU Average load "
	case consts.MonitorMemRate:
		result = "Memory usage (%)"
	case consts.MonitorDiskPartRate:
		result = "Hard disk partition utilization (%)"
	case consts.MonitorNetworkReception:
		result = "Network reception (kb)"
	case consts.MonitorNetworkSending:
		result = "Network sending (kb)"
	}
	return result
}

func granularityToSql(granularity int) string {
	result := ""
	switch granularity {
	case 0:
		result = "1m"
	case 1:
		result = "5m"
	case 2:
		result = "15m"
	case 3:
		result = "1h"
	case 4:
		result = "1d"
	default:
		result = "1m"
	}
	return result
}

func tTypeToSql(tType int, beginTime, endTime string) string {
	result := ""
	switch tType {
	case 0:
		result = " ts >= NOW-1h "
	case 1:
		result = " ts >= NOW-1d "
	case 2:
		result = " ts >= NOW-7d "
	case 3:
		result = fmt.Sprintf(" ts >= '%s' and ts <= '%s'", beginTime, endTime)
	default:
		result = " ts >= NOW-1h "
	}
	return result
}

func convertMonitorValue(oldValue float64, dType uint8) (newValue string) {
	switch dType {
	case consts.MonitorCpuUsage, consts.MonitorCpuBaseUsage, consts.MonitorCpuLoad,
		consts.MonitorDiskReadFlow, consts.MonitorDiskWriteFlow, consts.MonitorDiskPartTotal,
		consts.MonitorDiskPartUsage, consts.MonitorDiskPartRate, consts.MonitorDiskIoWait,
		consts.MonitorDiskIoService:
		newValue = tools.FormatFloat(float64(oldValue)/100, 1)
	case consts.MonitorMemUsed, consts.MonitorNetworkReception, consts.MonitorNetworkSending, consts.MonitorMemRate:
		newValue = tools.FormatFloat(oldValue, 1)
	case consts.MonitorDiskIoBusy:
		newValue = tools.FormatFloat(float64(oldValue)/10000, 1)
	case consts.MonitorNginxConcurrency:
		newValue = fmt.Sprintf("%f", oldValue)
	}
	return
}

func GetServerMonitorData(dType, tType int, beginTime, endTime string, granularity int) (commun.MonitorChart2, error) {
	title := dTypeToTitle(dType)
	chart := commun.MonitorChart2{Title: title, X: make([]string, 0), Y: make([]string, 0)}
	timeCond := tTypeToSql(tType, beginTime, endTime)
	interval := granularityToSql(granularity)
	sql := ""
	sql = fmt.Sprintf("select AVG(value) as value from %s.%s where type=%d and %s INTERVAL(%s)",
		taos.DataBase, taos.MonitorTableName, dType, timeCond, interval)
	rows, err := variable.TaosDb.Query(sql)
	if err != nil {
		fmt.Printf("%s error: %s\n", sql, err.Error())
		variable.ZapLog.Error("[Taos GetServerMonitorData] "+sql+" error:", zap.Error(err))
		return chart, err
	}
	defer rows.Close()

	for rows.Next() {
		var ts string
		var value float64
		err = rows.Scan(&ts, &value)
		if err != nil {
			fmt.Printf("scan table %s data error: %s\n", taos.MonitorTableName, err.Error())
			variable.ZapLog.Error("[Taos GetServerMonitorData] scan table "+taos.MonitorTableName+" data error:", zap.Error(err))
			continue
		}
		thisTime, err := time.ParseInLocation("2006-01-02 15:04:05", ts, time.Local)
		if err != nil {
			continue
		}
		timePeriod := thisTime.Format("2006-01-02 15:04")
		chart.X = append(chart.X, timePeriod)
		currentType := uint8(dType)
		newValue := convertMonitorValue(value, currentType)
		chart.Y = append(chart.Y, newValue)
	}
	return chart, nil
}

func GetServerMonitorDataWithMultipoint(dType, tType int, beginTime, endTime string, granularity int) (commun.MonitorChart3, error) {
	title := dTypeToTitle(dType)
	chart := commun.MonitorChart3{Title: title, X: make([]string, 0), Y: make([]commun.ChartPointData, 0)}
	if dType != int(consts.MonitorMemRate) && dType != int(consts.MonitorDiskPartRate) {
		return chart, errors.New(my_errors.ErrorResultEmpty)
	}
	timeCond := tTypeToSql(tType, beginTime, endTime)
	interval := granularityToSql(granularity)
	if dType == int(consts.MonitorMemRate) {
		sqlUsed := fmt.Sprintf("select AVG(value) as value from %s.%s where type=%d and %s INTERVAL(%s)",
			taos.DataBase, taos.MonitorTableName, consts.MonitorMemUsed, timeCond, interval)
		timeArray := make([]string, 0)
		memUsedMap := make(map[string]float64)
		rowsUsed, err := variable.TaosDb.Query(sqlUsed)
		if err != nil {
			fmt.Printf("%s error: %s\n", sqlUsed, err.Error())
			variable.ZapLog.Error("[Taos GetWarnStatistics] "+sqlUsed+"  error:", zap.Error(err))
			return chart, err
		} else {
			for rowsUsed.Next() {
				var ts string
				var value float64
				err = rowsUsed.Scan(&ts, &value)
				if err != nil {
					fmt.Printf("scan table %s data error: %s\n", taos.MonitorTableName, err.Error())
					variable.ZapLog.Error("[Taos GetServerMonitorData] scan table "+taos.MonitorTableName+" data error:", zap.Error(err))
					continue
				}
				thisTime, err := time.ParseInLocation("2006-01-02 15:04:05", ts, time.Local)
				if err != nil {
					continue
				}
				timePeriod := thisTime.Format("2006-01-02 15:04")
				memUsedMap[timePeriod] = value
				timeArray = append(timeArray, timePeriod)
			}
		}
		defer rowsUsed.Close()
		if len(memUsedMap) <= 0 || len(timeArray) <= 0 {
			return chart, errors.New(my_errors.ErrorResultEmpty)
		}

		sqlRate := fmt.Sprintf("select AVG(value) as value from %s.%s where type=%d and %s INTERVAL(%s)",
			taos.DataBase, taos.MonitorTableName, consts.MonitorMemRate, timeCond, interval)
		memRateMap := make(map[string]float64)
		rowsRate, err := variable.TaosDb.Query(sqlRate)
		if err != nil {
			fmt.Printf("%s error: %s\n", sqlRate, err.Error())
			variable.ZapLog.Error("[Taos GetWarnStatistics] "+sqlRate+"  error:", zap.Error(err))
			return chart, err
		} else {
			for rowsRate.Next() {
				var ts string
				var value float64
				err = rowsRate.Scan(&ts, &value)
				if err != nil {
					fmt.Printf("scan table %s data error: %s\n", taos.MonitorTableName, err.Error())
					variable.ZapLog.Error("[Taos GetServerMonitorData] scan table "+taos.MonitorTableName+" data error:", zap.Error(err))
					continue
				}
				thisTime, err := time.ParseInLocation("2006-01-02 15:04:05", ts, time.Local)
				if err != nil {
					continue
				}
				timePeriod := thisTime.Format("2006-01-02 15:04")
				memRateMap[timePeriod] = value
			}
		}
		defer rowsRate.Close()
		if len(memRateMap) <= 0 {
			return chart, errors.New(my_errors.ErrorResultEmpty)
		}

		memTotalValue := ""
		for _, timePoint := range timeArray {
			if len(timePoint) <= 0 {
				continue
			}
			memUsed, ok1 := memUsedMap[timePoint]
			memRate, ok2 := memRateMap[timePoint]
			if ok1 && ok2 {
				chart.X = append(chart.X, timePoint)
				multiData := make(commun.ChartPointData, 0)
				if len(memTotalValue) <= 0 {
					memTotal := memUsed / memRate * 100
					memTotalValue = tools.FormatFloat(memTotal, 1) + "M"
				}
				multiData = append(multiData, commun.ChartData{
					Name:  "Memory usage",
					Value: convertMonitorValue(memRate, consts.MonitorMemRate) + "%",
				},
					commun.ChartData{
						Name:  "Memory usage",
						Value: convertMonitorValue(memUsed, consts.MonitorMemUsed) + "M",
					},
					commun.ChartData{
						Name:  "Total memory",
						Value: memTotalValue,
					})
				chart.Y = append(chart.Y, multiData)
			}
		}
	} else {
		sqlUsed := fmt.Sprintf("select AVG(value) as value from %s.%s where type=%d and %s INTERVAL(%s)",
			taos.DataBase, taos.MonitorTableName, consts.MonitorDiskPartUsage, timeCond, interval)
		timeArray := make([]string, 0)
		diskUsedMap := make(map[string]float64)
		rowsUsed, err := variable.TaosDb.Query(sqlUsed)
		if err != nil {
			fmt.Printf("%s error: %s\n", sqlUsed, err.Error())
			variable.ZapLog.Error("[Taos GetWarnStatistics] "+sqlUsed+"  error:", zap.Error(err))
			return chart, err
		} else {
			for rowsUsed.Next() {
				var ts string
				var value float64
				err = rowsUsed.Scan(&ts, &value)
				if err != nil {
					fmt.Printf("scan table %s data error: %s\n", taos.MonitorTableName, err.Error())
					variable.ZapLog.Error("[Taos GetServerMonitorData] scan table "+taos.MonitorTableName+" data error:", zap.Error(err))
					continue
				}
				thisTime, err := time.ParseInLocation("2006-01-02 15:04:05", ts, time.Local)
				if err != nil {
					continue
				}
				timePeriod := thisTime.Format("2006-01-02 15:04")
				diskUsedMap[timePeriod] = value
				timeArray = append(timeArray, timePeriod)
			}
		}
		defer rowsUsed.Close()
		if len(diskUsedMap) <= 0 || len(timeArray) <= 0 {
			return chart, errors.New(my_errors.ErrorResultEmpty)
		}

		sqlRate := fmt.Sprintf("select AVG(value) as value from %s.%s where type=%d and %s INTERVAL(%s)",
			taos.DataBase, taos.MonitorTableName, consts.MonitorDiskPartRate, timeCond, interval)
		diskRateMap := make(map[string]float64)
		rowsRate, err := variable.TaosDb.Query(sqlRate)
		if err != nil {
			fmt.Printf("%s error: %s\n", sqlRate, err.Error())
			variable.ZapLog.Error("[Taos GetWarnStatistics] "+sqlRate+"  error:", zap.Error(err))
			return chart, err
		} else {
			for rowsRate.Next() {
				var ts string
				var value float64
				err = rowsRate.Scan(&ts, &value)
				if err != nil {
					fmt.Printf("scan table %s data error: %s\n", taos.MonitorTableName, err.Error())
					variable.ZapLog.Error("[Taos GetServerMonitorData] scan table "+taos.MonitorTableName+" data error:", zap.Error(err))
					continue
				}
				thisTime, err := time.ParseInLocation("2006-01-02 15:04:05", ts, time.Local)
				if err != nil {
					continue
				}
				timePeriod := thisTime.Format("2006-01-02 15:04")
				diskRateMap[timePeriod] = value
			}
		}
		defer rowsRate.Close()
		if len(diskRateMap) <= 0 {
			return chart, errors.New(my_errors.ErrorResultEmpty)
		}

		sqlTotal := fmt.Sprintf("select value from %s.%s where type=%d and %s order by ts desc limit 1",
			taos.DataBase, taos.MonitorTableName, consts.MonitorDiskPartTotal, timeCond)
		diskTotalValue := ""
		rowsTotal, err := variable.TaosDb.Query(sqlTotal)
		if err != nil {
			fmt.Printf("%s error: %s\n", sqlTotal, err.Error())
			variable.ZapLog.Error("[Taos GetWarnStatistics] "+sqlTotal+"  error:", zap.Error(err))
			return chart, err
		} else {
			for rowsTotal.Next() {
				var value float64
				err = rowsTotal.Scan(&value)
				if err != nil {
					fmt.Printf("scan table %s data error: %s\n", taos.MonitorTableName, err.Error())
					variable.ZapLog.Error("[Taos GetServerMonitorData] scan table "+taos.MonitorTableName+" data error:", zap.Error(err))
					continue
				}
				diskTotalValue = convertMonitorValue(value, consts.MonitorDiskPartTotal) + "G"
			}
		}
		defer rowsTotal.Close()

		for _, timePoint := range timeArray {
			if len(timePoint) <= 0 {
				continue
			}
			diskUsed, ok1 := diskUsedMap[timePoint]
			diskRate, ok2 := diskRateMap[timePoint]
			if ok1 && ok2 {
				chart.X = append(chart.X, timePoint)
				multiData := make(commun.ChartPointData, 0)
				if len(diskTotalValue) <= 0 {
					diskTotal := diskUsed / diskRate * 100
					diskTotalValue = tools.FormatFloat(diskTotal, 1) + "G"
				}
				multiData = append(multiData, commun.ChartData{
					Name:  "Hard disk utilization",
					Value: convertMonitorValue(diskRate, consts.MonitorDiskPartRate) + "%",
				},
					commun.ChartData{
						Name:  "Hard disk usage",
						Value: convertMonitorValue(diskUsed, consts.MonitorDiskPartUsage) + "G",
					},
					commun.ChartData{
						Name:  "Total hard disk",
						Value: diskTotalValue,
					})
				chart.Y = append(chart.Y, multiData)
			}
		}
	}
	return chart, nil
}
