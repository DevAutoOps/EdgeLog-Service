package tools

import (
	"strconv"
	"time"
)

func UintToStr(u uint) string {
	return strconv.FormatUint(uint64(u), 10)
}

func StrToUint(s string) uint {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return uint(i)
}

func StringToInt(e string) (int, error) {
	return strconv.Atoi(e)
}

func StringToInt64(e string) (int64, error) {
	return strconv.ParseInt(e, 10, 64)
}

func GetCurrentTimeStr() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func GetCurrentTime() time.Time {
	return time.Now()
}

// NumericalCheck  Can strings be converted to numbers
func NumericalCheck(value string) (float64, bool) {
	if len(value) <= 0 {
		return 0, false
	}
	resultValue, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, false
	}
	return resultValue, true
}
