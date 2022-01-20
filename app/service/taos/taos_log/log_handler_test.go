package taos_log

import (
	"edgelog/app/model/proto/log"
	"edgelog/app/service/taos"
	"fmt"
	"testing"
)

func TestLogLastSearchByLimitAndOffset(t *testing.T) {
	taos.InitTestBaseData()
	dataList, err := LogLastSearchByLimitAndOffset(5, 0, "", "")
	if err != nil {
		fmt.Println("search data err:", err)
		return
	}
	fmt.Println("data length:", len(dataList))
	fmt.Printf("%v\n", dataList)
}

func TestLogSearch(t *testing.T) {
	taos.InitTestBaseData()
	dataList, err := LogSearch("2021-09-23 08:00:00", "2021-09-23 16:30:00")
	if err != nil {
		fmt.Println("search data err:", err)
		return
	}
	fmt.Println("data length:", len(dataList))
	fmt.Printf("%v\n", dataList)
}

func TestLogAdd(t *testing.T) {
	logs := make([]*log.LogData, 0)
	logs = append(logs, &log.LogData{Log: `42.192.173.120	-	-	-	2022-01-10T08:15:03+08:00	GET	"/ttlsa-req-status"	"-"	HTTP/1.1	200	2807	0.000	"-"	"curl/7.29.0"	--Src-Ip	42.192.173.120	b44406e41ffe	172.18.0.4	47970	80	"-"	"-"	"-"	V5`})
	logs = append(logs, &log.LogData{Log: `42.192.173.120	-	-	-	2022-01-10T09:15:03+08:00	GET	"/"	"-"	HTTP/1.1	200	61	0.007	"-"	"curl/7.29.0"	--Src-Ip	42.192.173.120	b44406e41ffe	172.18.0.4	47974	80	"42.192.173.120:9023"	"200"	"0.006"	V5`})
	logs = append(logs, &log.LogData{Log: `42.192.173.120	-	-	-	2022-01-10T10:15:04+08:00	GET	"/file1"	"-"	HTTP/1.1	404	196	0.009	"-"	"curl/7.29.0"	--Src-Ip	42.192.173.120	b44406e41ffe	172.18.0.4	47982	80	"42.192.173.120:9023"	"404"	"0.010"	V5`})
	logs = append(logs, &log.LogData{Log: `42.192.173.120	-	-	-	2022-01-10T11:15:04+08:00	GET	"/"	"-"	HTTP/1.1	200	61	0.007	"-"	"curl/7.29.0"	--Src-Ip	42.192.173.120	b44406e41ffe	172.18.0.4	47988	80	"42.192.173.120:9023"	"200"	"0.007"	V5`})
	logs = append(logs, &log.LogData{Log: `42.192.173.120	-	-	-	2022-01-10T12:15:05+08:00	GET	"/file2"	"-"	HTTP/1.1	404	196	0.007	"-"	"curl/7.29.0"	--Src-Ip	42.192.173.120	b44406e41ffe	172.18.0.4	47996	80	"42.192.173.120:9023"	"404"	"0.007"	V5`})
	logs = append(logs, &log.LogData{Log: `42.192.173.120	-	-	-	2022-01-10T13:15:05+08:00	GET	"/"	"-"	HTTP/1.1	200	61	0.006	"-"	"curl/7.29.0"	--Src-Ip	42.192.173.120	b44406e41ffe	172.18.0.4	48002	80	"42.192.173.120:9023"	"200"	"0.007"	V5`})
	logs = append(logs, &log.LogData{Log: `42.192.173.120	-	-	-	2022-01-10T14:15:06+08:00	GET	"/file3"	"-"	HTTP/1.1	404	196	0.007	"-"	"curl/7.29.0"	--Src-Ip	42.192.173.120	b44406e41ffe	172.18.0.4	48012	80	"42.192.173.120:9023"	"404"	"0.007"	V5`})
	logs = append(logs, &log.LogData{Log: `42.192.173.120	-	-	-	2022-01-10T15:15:06+08:00	GET	"/"	"-"	HTTP/1.1	200	61	0.007	"-"	"curl/7.29.0"	--Src-Ip	42.192.173.120	b44406e41ffe	172.18.0.4	48016	80	"42.192.173.120:9023"	"200"	"0.006"	V5`})
	logs = append(logs, &log.LogData{Log: `42.192.173.120	-	-	-	2022-01-10T17:15:07+08:00	GET	"/file4"	"-"	HTTP/1.1	404	196	0.006	"-"	"curl/7.29.0"	--Src-Ip	42.192.173.120	b44406e41ffe	172.18.0.4	48026	80	"42.192.173.120:9023"	"404"	"0.007"	V5`})
	logs = append(logs, &log.LogData{Log: `42.192.173.120	-	-	-	2022-01-10T16:15:07+08:00	GET	"/"	"-"	HTTP/1.1	200	61	0.006	"-"	"curl/7.29.0"	--Src-Ip	42.192.173.120	b44406e41ffe	172.18.0.4	48030	80	"42.192.173.120:9023"	"200"	"0.006"	V5`})

	taos.InitTestBaseData()
	for _, logData := range logs {
		LogAdd(logData)
	}
}
