package notice

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

const (
	WeixinGetTokenUrl       = "https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=%s&corpsecret=%s"
	WeixinMessageSendUrl    = "https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token=%s"
	WeixinMediaUploadUrl    = "https://qyapi.weixin.qq.com/cgi-bin/media/upload?access_token=%s&type=file"
	WeixinDepartmentListUrl = "https://qyapi.weixin.qq.com/cgi-bin/department/list?access_token=%s"
	WeixinMessageSendText   = `{"touser":"%s","msgtype":"text","agentid":%d,"text":{"content":"%s"}}`
	WeixinMessageSendFile   = `{"touser":"%s","msgtype":"file","agentid":%d,"file":{"media_id":"%s"}}`
)

type WeixinNotice struct {
	mu          sync.RWMutex
	accessToken string
	agentId     int
	client      *http.Client
	receive     []string
}

func CreateWeChatNotice(w WeChat, timeOut time.Duration) (INotice, error) {
	tokenApiResult, err := getAccessToken(w.EID, w.EAppVoucher)
	if err != nil {
		return nil, err
	}
	split := strings.Split(w.ReceiveUserID, ",")
	result := &WeixinNotice{
		accessToken: tokenApiResult.AccessToken,
		agentId:     w.EAppID,
		client:      &http.Client{Timeout: timeOut},
		receive:     split,
	}
	go func(_n *WeixinNotice) {
		time.Sleep(time.Duration(tokenApiResult.ExpiresIn) * time.Second)
		_n.mu.Lock()
		r, e := getAccessToken(w.EID, w.EAppVoucher)
		if e != nil {
			fmt.Println(e)
			return
		}
		_n.accessToken = r.AccessToken
		_n.mu.Unlock()
	}(result)
	return result, nil
}

//{
//	"errcode":0,
//	"errmsg":"",
//	"access_token": "accesstoken000001",
//	"expires_in": 7200
//}
type getTokenApiResult struct {
	ErrCode     int    `json:"errcode"`
	ErrMsg      string `json:"errmsg"`
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

func getAccessToken(corpId, corpSecret string) (result getTokenApiResult, err error) {
	resp, err := http.Get(fmt.Sprintf(WeixinGetTokenUrl, corpId, corpSecret))
	if err != nil {
		return
	}
	defer resp.Body.Close()
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

//{
//	"errcode": 0,
//	"errmsg": "ok",
//	"department": [{
//		"id": 2,
//		"name": " Guangzhou R & D Center ",
//		"name_en": "RDGZ",
//		"parentid": 1,
//		"order": 10
//	}, {
//		"id": 3,
//		"name": " Mailbox product department ",
//		"name_en": "mail",
//		"parentid": 2,
//		"order": 40
//	}]
//}
type departmentListApiResult struct {
	ErrCode    int          `json:"errcode"`
	ErrMsg     string       `json:"errmsg"`
	Department []Department `json:"department"`
}

type Department struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	NameEn   string `json:"name_en"`
	ParentId int    `json:"parentid"`
	Order    int    `json:"order"`
}

func (s *WeixinNotice) getDepartmentList() (result departmentListApiResult, err error) {
	resp, err := http.Get(fmt.Sprintf(WeixinDepartmentListUrl, s.accessToken))
	if err != nil {
		return
	}
	defer resp.Body.Close()
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

//{
//	"errcode" : 0,
//	"errmsg" : "ok",
//	"invaliduser" : "userid1|userid2",
//	"invalidparty" : "partyid1|partyid2",
//	"invalidtag": "tagid1|tagid2"
//}
type massageSendApiResult struct {
	ErrCode      int    `json:"errcode"`
	ErrMsg       string `json:"errmsg"`
	InvalidUser  string `json:"invaliduser"`
	InvalidParty string `json:"invalidparty"`
	InvalidTag   string `json:"invalidtag"`
}

func (s *WeixinNotice) SendText(to []string, msg string) (err error) {
	//{
	//	"touser" : "UserID1|UserID2|UserID3",
	//	"toparty" : "PartyID1|PartyID2",
	//	"totag" : "TagID1 | TagID2",
	//	"msgtype" : "text",
	//	"agentid" : 1,
	//	"text" : {
	//		"content" : " Your express has arrived ， Please bring your work card to the mail center to get it 。\n Check before departure <a href=\"http://work.weixin.qq.com\"> Mail center video live </a>， Smart to avoid queuing 。"
	//	},
	//	"safe":0,
	//	"enable_id_trans": 0,
	//	"enable_duplicate_check": 0,
	//	"duplicate_check_interval": 1800
	//}
	//touser	 no 	 Specifies the member who receives the message ， member ID list （ For multiple recipients ‘|’ separate ， Maximum support 1000 individual ）。
	// exceptional case ： Designated as ”@all”， Send to all members of the enterprise application 
	//toparty	 no 	 Specify the Department that receives the message ， department ID list ， For multiple recipients ‘|’ separate ， Maximum support 100 individual 。
	// When touser by ”@all” Ignore this parameter when 
	//totag	 no 	 Specifies the label of the received message ， label ID list ， For multiple recipients ‘|’ separate ， Maximum support 100 individual 。
	// When touser by ”@all” Ignore this parameter when 
	//msgtype	 yes 	 Message type ， Fixed as ：text
	//agentid	 yes 	 Enterprise application id， integer 。 Enterprise internal development ， It can be viewed on the settings page of the application ； Third party service provider ， Via the interface   Obtain enterprise authorization information   Get the parameter value 
	//content	 yes 	 Message content ， Up to 2048 Bytes ， More than will be truncated （ support id Translation ）
	//safe	 no 	 Indicates whether it is a confidential message ，0 Indicates that it can be shared ，1 Indicates that the content cannot be shared and the watermark is displayed ， Default to 0
	//enable_id_trans	 no 	 Indicates whether it is turned on id Translation ，0 Indicates No ，1 Means yes ， default 0。 Only for third-party applications ， Enterprise self built applications can be ignored 。
	//enable_duplicate_check	 no 	 Indicates whether to enable duplicate message checking ，0 Indicates No ，1 Means yes ， default 0
	//duplicate_check_interval	 no 	 Indicates whether to repeat the message check interval ， default 1800s， Maximum not more than 4 hour 
	s.mu.RLock()
	defer s.mu.RUnlock()
	to = append(to, s.receive...)
	paramsStr := fmt.Sprintf(WeixinMessageSendText, strings.Join(to, "|"), s.agentId, msg)
	req, err := http.NewRequest("POST", fmt.Sprintf(WeixinMessageSendUrl, s.accessToken), bytes.NewBuffer([]byte(paramsStr)))
	if err != nil {
		return
	}
	req.Header.Add("content-type", "application/json")
	return s.Send(req)
}

func (s *WeixinNotice) SendFile(to []string, file *os.File) (err error) {
	//{
	//	"touser" : "UserID1|UserID2|UserID3",
	//	"toparty" : "PartyID1|PartyID2",
	//	"totag" : "TagID1 | TagID2",
	//	"msgtype" : "file",
	//	"agentid" : 1,
	//	"file" : {
	//		"media_id" : "1Yv-zXfHjSjU-7LH-GwtYqDGS-zz6w22KmWAT5COgP7o"
	//	},
	//	"safe":0,
	//	"enable_duplicate_check": 0,
	//	"duplicate_check_interval": 1800
	//}
	//touser	 no 	 member ID list （ Message Receiver  ， For multiple recipients ‘|’ separate ， Maximum support 1000 individual ）。 exceptional case ： Designated as @all， Send it to all members who pay attention to the enterprise application 
	//toparty	 no 	 department ID list ， For multiple recipients ‘|’ separate ， Maximum support 100 individual 。 When touser by @all Ignore this parameter when 
	//totag	 no 	 label ID list ， For multiple recipients ‘|’ separate ， Maximum support 100 individual 。 When touser by @all Ignore this parameter when 
	//msgtype	 yes 	 Message type ， Fixed as ：file
	//agentid	 yes 	 Enterprise application id， integer 。 Enterprise internal development ， It can be viewed on the settings page of the application ； Third party service provider ， Via the interface   Obtain enterprise authorization information   Get the parameter value 
	//media_id	 yes 	 file id， You can call the upload temporary material interface to obtain 
	//safe	 no 	 Indicates whether it is a confidential message ，0 Indicates that it can be shared ，1 Indicates that the content cannot be shared and the watermark is displayed ， Default to 0
	//enable_duplicate_check	 no 	 Indicates whether to enable duplicate message checking ，0 Indicates No ，1 Means yes ， default 0
	//duplicate_check_interval	 no 	 Indicates whether to repeat the message check interval ， default 1800s， Maximum not more than 4 hour 
	s.mu.RLock()
	defer s.mu.RUnlock()
	media, err := s.uploadTempMaterial(file)
	if err != nil {
		return
	}
	to = append(to, s.receive...)
	paramsStr := fmt.Sprintf(WeixinMessageSendFile, strings.Join(to, "|"), s.agentId, media.MediaId)
	req, err := http.NewRequest("POST", fmt.Sprintf(WeixinMessageSendUrl, s.accessToken), bytes.NewBuffer([]byte(paramsStr)))
	if err != nil {
		return
	}
	req.Header.Add("content-type", "application/json")
	return s.Send(req)
}

func (s *WeixinNotice) Send(req *http.Request) (err error) {
	resp, err := s.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	var result massageSendApiResult
	if err = json.Unmarshal(body, &result); err != nil {
		return
	}
	if result.ErrCode != 0 {
		err = errors.New(result.ErrMsg)
		return
	}
	return
}

//{
//	"errcode": 0,
//	"errmsg": ""，
//	"type": "image",
//	"media_id": "1G6nrLmr5EC3MMb_-zK1dDdzmd0p7cNliYu9V5w7o8K0",
//	"created_at": "1380000000"
//}
//type	 Media file type ， There are pictures respectively （image）、 voice （voice）、 video （video）， Ordinary file (file)
//media_id	 Unique identification obtained after media file upload ，3 Valid within days 
//created_at	 Media file upload timestamp 
type mediaUploadApiResult struct {
	ErrCode   int    `json:"errcode"`
	ErrMsg    string `json:"errmsg"`
	Type      string `json:"type"`
	MediaId   string `json:"media_id"`
	CreatedAt string `json:"created_at"`
}

func (s *WeixinNotice) uploadTempMaterial(file *os.File) (result mediaUploadApiResult, err error) {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	fw, err := w.CreateFormFile("media", file.Name())
	if err != nil {
		return
	}
	if _, err = io.Copy(fw, file); err != nil {
		return
	}
	err = w.Close()
	if err != nil {
		return
	}
	req, err := http.NewRequest("POST", fmt.Sprintf(WeixinMediaUploadUrl, s.accessToken), &b)
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", w.FormDataContentType())
	resp, err := s.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("bad status: %s", resp.Status)
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
