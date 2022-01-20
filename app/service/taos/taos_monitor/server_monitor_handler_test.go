package taos_monitor

import (
	"edgelog/app/service/taos"
	"fmt"
	"testing"
)

func TestGetServerMonitorData(t *testing.T) {
	taos.InitTestBaseData()
	chart, err := GetServerMonitorData(15, 3, "2021-10-19 5:00:00", "2021-10-19 23:30:00", 3)
	if err != nil {
		fmt.Println("get data err:", err)
		return
	}
	fmt.Println("data time length:", len(chart.X), " ,value length:", len(chart.Y))
	fmt.Printf("%v\n", chart)
}

func TestGetServerMonitorDataWithMultipoint(t *testing.T) {
	taos.InitTestBaseData()
	chart, err := GetServerMonitorDataWithMultipoint(5, 0, "", "", 0)
	if err != nil {
		fmt.Println("get data err:", err)
		return
	}
	fmt.Println("data time length:", len(chart.X), " ,value length:", len(chart.Y))
	fmt.Printf("%v\n", chart)
}
