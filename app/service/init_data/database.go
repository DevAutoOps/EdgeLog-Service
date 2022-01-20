package init_data

import (
	"edgelog/app/dao"
	"edgelog/app/global/consts"
	"edgelog/app/global/variable"
	"edgelog/app/table"
	"edgelog/tools"
	"go.uber.org/zap"
)

//  Initialize tables and data
func InitDatabase() {
	_ = createTables()
	initData()
}

func createTables() error {
	tables := []table.ITable{
		&table.Config{},
		&table.Keyword{},
		&table.Template{},
	}
	for _, v := range tables {
		err := variable.GormDb.AutoMigrate(v)
		if err != nil {
			return err
		}
	}

	return nil
}

func initData() {
	initKeyword()
	initHostThreshold()
	initTemplate()
}

//  Initialize log format fields
func initKeyword() {
	list, err := (&dao.Keyword{}).GetList()
	if err != nil || len(list) == 0 {
		keyword := []string{
			"args",
			"query_string",
			"arg_name",
			"is_args",
			"uri",
			"document_uri",
			"document_root",
			"host",
			"hostname",
			"https",
			"binary_remote_addr",
			"body_bytes_sent",
			"bytes_sent",
			"connection",
			"connection_requests",
			"content_length",
			"content_type",
			"cookie_name",
			"limit_rate",
			"msec",
			"nginx_version",
			"pid",
			"pipe",
			"proxy_protocol_addr",
			"realpath_root",
			"remote_addr",
			"remote_port",
			"remote_user",
			"request",
			"request_body",
			"request_body_file",
			"request_completion",
			"request_filename",
			"request_length",
			"request_method",
			"request_time",
			"request_uri",
			"scheme",
			"server_addr",
			"server_name",
			"server_port",
			"server_protocol",
			"status",
			"time_iso8601",
			"time_local",
			"cookie_name",
			"http_name",
			"http_cookie",
			"http_host",
			"http_referer",
			"http_user_agent",
			"http_x_forwarded_for",
			"sent_http_name",
			"sent_http_cache_control",
			"sent_http_connection",
			"sent_http_content_type",
			"sent_http_keep_alive",
			"sent_http_last_modified",
			"sent_http_location",
			"sent_http_transfer_encoding",
		}
		defaultKeyword := make([]table.Keyword, len(keyword))
		for i, str := range keyword {
			defaultKeyword[i] = table.Keyword{
				Name:  str,
				Value: str,
				Type:  0,
				Order: i,
			}
		}
		if err = variable.GormDb.Model(&table.Keyword{}).Create(defaultKeyword).Error; err != nil {
			variable.ZapLog.Error("initKeyword err:", zap.Error(err))
		}
	}
}

func initHostThreshold() {
	value := consts.HostDefaultThresholdValueStr
	if !(&dao.Config{}).CheckSystemConfigIsExistence(consts.HostCpuThreshold) {
		_ = (&dao.Config{}).AddAndSetSystemConfig(consts.HostCpuThreshold, value)
		variable.CpuThreshold, _ = tools.StringToInt(value)
	} else {
		cpuThreshold, err := (&dao.Config{}).GetSystemConfig(consts.HostCpuThreshold)
		if err != nil {
			variable.CpuThreshold = consts.HostDefaultThresholdValue
		} else {
			variable.CpuThreshold, err = tools.StringToInt(cpuThreshold)
			if err != nil {
				variable.CpuThreshold = consts.HostDefaultThresholdValue
			}
		}
	}
	if !(&dao.Config{}).CheckSystemConfigIsExistence(consts.HostMemoryThreshold) {
		_ = (&dao.Config{}).AddAndSetSystemConfig(consts.HostMemoryThreshold, value)
		variable.MemThreshold, _ = tools.StringToInt(value)
	} else {
		memThreshold, err := (&dao.Config{}).GetSystemConfig(consts.HostMemoryThreshold)
		if err != nil {
			variable.MemThreshold = consts.HostDefaultThresholdValue
		} else {
			variable.MemThreshold, err = tools.StringToInt(memThreshold)
			if err != nil {
				variable.MemThreshold = consts.HostDefaultThresholdValue
			}
		}
	}
	if !(&dao.Config{}).CheckSystemConfigIsExistence(consts.HostDiskThreshold) {
		_ = (&dao.Config{}).AddAndSetSystemConfig(consts.HostDiskThreshold, value)
		variable.DiskThreshold, _ = tools.StringToInt(value)
	} else {
		diskThreshold, err := (&dao.Config{}).GetSystemConfig(consts.HostDiskThreshold)
		if err != nil {
			variable.DiskThreshold = consts.HostDefaultThresholdValue
		} else {
			variable.DiskThreshold, err = tools.StringToInt(diskThreshold)
			if err != nil {
				variable.DiskThreshold = consts.HostDefaultThresholdValue
			}
		}
	}
}

func initTemplate() {
	template, err := (&dao.Template{}).GetDefault()
	if template.ID == 0 || err != nil {
		defaultTemplate := &table.Template{
			Name:      "Demo Template",
			Value:     "created_at|remote_addr|method|request|status|size|referer|http_host|upstream_response_time|request_time|http_x_forwarded_for|user_agent|upstream_host",
			Format:    true,
			Type:      0,
			Separator: "|",
		}
		defaultTemplate.ID = 1
		err = (&dao.Template{}).Create(defaultTemplate)
		if err != nil {
			variable.ZapLog.Error("initTemplate err:", zap.Error(err))
		}
	}
}
