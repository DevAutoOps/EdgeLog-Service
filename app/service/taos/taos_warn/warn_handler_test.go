package taos_warn

import (
	"edgelog/app/service/taos"
	"fmt"
	"testing"
)

func TestGetWarnStatistics(t *testing.T) {
	taos.InitTestBaseData()
	chart, err := GetWarnStatistics("2021-09-29 00:00:00", "2021-10-08 23:59:59")
	if err != nil {
		fmt.Println("get data err:", err)
		return
	}
	fmt.Printf("%v\n", chart)
}
