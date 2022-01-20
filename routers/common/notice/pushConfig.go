package notice

type PushConfig struct {
	SMTP     SMTP
	WeChat   WeChat
	DingTalk DingTalk
}

type SMTP struct {
	Addr         string
	Port         int
	SSL          bool
	User         string
	Pass         string
	Topic        string
	ReceiveEmail string
}

type WeChat struct {
	EID           string
	EAppID        int
	EAppVoucher   string
	ReceiveUserID string
}

type DingTalk struct {
	AgentID    string
	AppKey     string
	AppSecret  string
	SessionID  string
	HookAddr   string
	HookSecret string
	HookTitle  string
	Receive    string
}
