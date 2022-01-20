package taos_log

import (
	"edgelog/app/global/my_errors"
	"edgelog/app/global/variable"
	"edgelog/app/model/proto/log"
	"edgelog/app/service/taos"
	"errors"
	"fmt"
	"go.uber.org/zap"
)

func LogAdd(logData *log.LogData) {
	if variable.TaosDb == nil {
		variable.ZapLog.Warn("[Taos] taos database is not initialized")
		return
	}

	taos.LogLock.Lock()
	defer taos.LogLock.Unlock()
	sql := fmt.Sprintf("insert into %s.%s values (now,'%s')",
		taos.DataBase, taos.LogTableName, logData.Log)
	_, err := variable.TaosDb.Exec(sql)
	if err != nil {
		fmt.Printf("insert table %s data, sql: %s, error: %s\n", taos.LogTableName, sql, err.Error())
		variable.ZapLog.Error(fmt.Sprintf("[Taos] Insert table %s data, sql: %s error:", taos.LogTableName, sql), zap.Error(err))
	}
}

//Search log
func LogSearch(beginTime, endTime string) ([]LogModel, error) {
	if (len(beginTime) > 0 && len(endTime) <= 0) || (len(beginTime) <= 0 && len(endTime) > 0) {
		return nil, errors.New(my_errors.ErrorInsufficientParameters)
	}
	sql := fmt.Sprintf("SELECT log FROM %s.%s WHERE ", taos.DataBase, taos.LogTableName)
	if len(beginTime) > 0 {
		sql = fmt.Sprintf("%s ts >= '%s' and ts <= '%s'", sql, beginTime, endTime)
	}
	rows, err := variable.TaosDb.Query(sql)
	if err != nil {
		fmt.Printf("%s error: %s\n", sql, err.Error())
		variable.ZapLog.Error("[Taos LogSearch] "+sql+"  error:", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	logList := make([]LogModel, 0)
	for rows.Next() {
		var data LogModel
		err = rows.Scan(&data.Log)
		if err != nil {
			fmt.Printf("scan table %s data error: %s\n", taos.LogTableName, err.Error())
			variable.ZapLog.Error("[Taos StressServer] scan table "+taos.LogTableName+" data error:", zap.Error(err))
			continue
		}
		logList = append(logList, data)
	}
	return logList, nil
}

// Search the last few logs
func LogLastSearchByLimitAndOffset(limit, offset int, beginTime, endTime string) ([]LogModel, error) {
	if (len(beginTime) > 0 && len(endTime) <= 0) || (len(beginTime) <= 0 && len(endTime) > 0) {
		return nil, errors.New(my_errors.ErrorInsufficientParameters)
	}
	sqlCount := fmt.Sprintf("SELECT COUNT(*) as count FROM %s.%s ", taos.DataBase, taos.LogTableName)
	sql := fmt.Sprintf("SELECT log FROM %s.%s ", taos.DataBase, taos.LogTableName)
	if len(beginTime) > 0 {
		sqlCount = fmt.Sprintf("%s WHERE ts >= '%s' and ts <= '%s'", sqlCount, beginTime, endTime)
		sql = fmt.Sprintf("%s WHERE ts >= '%s' and ts <= '%s'", sql, beginTime, endTime)
	}

	countRows, err := variable.TaosDb.Query(sqlCount)
	if err != nil {
		fmt.Printf("%s error: %s\n", sql, err.Error())
		variable.ZapLog.Error("[Taos LogLastSearchByLimitAndOffset] "+sqlCount+"  error:", zap.Error(err))
		return nil, err
	}
	defer countRows.Close()
	count := 0
	for countRows.Next() {
		var currentCount int
		err = countRows.Scan(&currentCount)
		if err != nil {
			fmt.Printf("scan table %s data error: %s\n", taos.LogTableName, err.Error())
			variable.ZapLog.Error("[Taos LogLastSearchByLimitAndOffset] scan table "+taos.LogTableName+" data error:", zap.Error(err))
			continue
		}
		count = currentCount
	}
	if count <= 0 {
		return nil, errors.New(my_errors.ErrorResultEmpty)
	} else if count < limit {
		limit = count
	}

	sql = fmt.Sprintf("%s ORDER BY ts DESC LIMIT %d OFFSET %d", sql, limit, offset)
	rows, err := variable.TaosDb.Query(sql)
	if err != nil {
		fmt.Printf("%s error: %s\n", sql, err.Error())
		variable.ZapLog.Error("[Taos LogLastSearchByLimitAndOffset] "+sql+"  error:", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	logList := make([]LogModel, 0)
	for rows.Next() {
		var data LogModel
		err = rows.Scan(&data.Log)
		if err != nil {
			fmt.Printf("scan table %s data error: %s\n", taos.LogTableName, err.Error())
			variable.ZapLog.Error("[Taos LogLastSearchByLimitAndOffset] scan table "+taos.LogTableName+" data error:", zap.Error(err))
			continue
		}
		logList = append(logList, data)
	}
	return logList, nil
}
