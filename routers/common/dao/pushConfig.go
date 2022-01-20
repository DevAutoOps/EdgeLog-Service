package dao

import (
	"edgelog/routers/common/notice"
	"encoding/json"
)

type PushConfig struct {
}

func (PushConfig) GetSMTPConfig() (smtp notice.SMTP, err error) {
	config, err := (&Config{}).GetSystemConfig("smtp_push_config")
	if err != nil {
		return
	}
	err = json.Unmarshal([]byte(config), &smtp)
	return
}

func (PushConfig) GetWeChatConfig() (wechat notice.WeChat, err error) {
	config, err := (&Config{}).GetSystemConfig("wechat_push_config")
	if err != nil {
		return
	}
	err = json.Unmarshal([]byte(config), &wechat)
	return
}

func (PushConfig) GetDingTalkConfig() (ding notice.DingTalk, err error) {
	config, err := (&Config{}).GetSystemConfig("ding_push_config")
	if err != nil {
		return
	}
	err = json.Unmarshal([]byte(config), &ding)
	return
}
