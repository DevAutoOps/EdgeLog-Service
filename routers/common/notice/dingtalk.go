package notice

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	DingTalkGetTokenUrl     = "https://oapi.dingtalk.com/gettoken?appkey=%s&appsecret=%s"
	DingTalkMessageSendUrl  = "https://oapi.dingtalk.com/topapi/message/corpconversation/asyncsend_v2?access_token=%s"
	DingTalkWebHookUrl      = `%s&timestamp=%s&sign=%s`
	DingTalkWebHookMsgTemp  = `{"msgtype":"text","title":"%s","text":{"content":"%s"}}`
	DingTalkMessageSendText = `{"agent_id": "%s","userid_list": "%s","to_all_user": false,"msg": {"msgtype": "text","text": {"content": "%s"}}}`
)

type DingTalkNotice struct {
	config         DingTalk
	AppAccessToken string
	client         *http.Client
}

func NewDingTalkNotice(dingTalk DingTalk, timeOut time.Duration) (*DingTalkNotice, error) {
	getAccessToken := func(appKey, appSecret string) (string, error) {
		resp, err := http.Get(fmt.Sprintf(DingTalkGetTokenUrl, appKey, appSecret))
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}
		var result getTokenApiResult
		if err = json.Unmarshal(body, &result); err != nil {
			return "", err
		}
		if result.ErrCode != 0 {
			return "", errors.New(result.ErrMsg)
		}
		return result.AccessToken, nil
	}
	accessToken, err := getAccessToken(dingTalk.AppKey, dingTalk.AppSecret)
	if err != nil {
		return nil, err
	}
	return &DingTalkNotice{
		config:         dingTalk,
		AppAccessToken: accessToken,
		client:         &http.Client{Timeout: timeOut},
	}, nil
}

type ApiResult struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

func (s *DingTalkNotice) SendText(content string) (err error) {
	req, err := http.NewRequest("POST",
		fmt.Sprintf(DingTalkMessageSendUrl, s.AppAccessToken), strings.NewReader(fmt.Sprintf(DingTalkMessageSendText, s.config.AgentID, s.config.Receive, content)))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json;charset=utf-8")
	resp, err := s.client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	var result ApiResult
	if err = json.Unmarshal(body, &result); err != nil {
		return
	}
	if result.ErrCode != 0 {
		err = errors.New(result.ErrMsg)
		return
	}
	return
}

func (s *DingTalkNotice) SendTextByWebHook(content string) (result ApiResult, err error) {
	timestamp, sign := getDingTalkSign(s.config.HookSecret)
	req, err := http.NewRequest("POST",
		fmt.Sprintf(DingTalkWebHookUrl, s.config.HookAddr, timestamp, sign), strings.NewReader(fmt.Sprintf(DingTalkWebHookMsgTemp, s.config.HookTitle, content)))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json;charset=utf-8")
	resp, err := s.client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	if err = json.Unmarshal(body, &result); err != nil {
		return
	}
	if result.ErrCode != 0 {
		err = errors.New(result.ErrMsg)
		return
	}
	return
}

func getDingTalkSign(secret string) (timestamp, sign string) {
	timestamp = fmt.Sprint(time.Now().UnixNano() / 1e6)
	stringToSign := timestamp + "\n" + secret
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(stringToSign))
	encodeString := base64.StdEncoding.EncodeToString(h.Sum(nil))
	sign = url.QueryEscape(encodeString)
	return
}
