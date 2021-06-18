package structs
//
type Code2Session struct {
	SessionKey string `json:"session_key"`
	OpenID string `json:"openid"`
	ErrCode int32 `json:"errcode,omitempty"`
	ErrMsg string `json:"errmsg,omitempty"`
}

// 登录
type Login struct {
	ID string `json:"id,omitempty"`
	OpenID string `json:"openid"`
	City string `json:"city"`
	Province string `json:"province"`
	Country string `json:"country"`
	AvatarUrl string `json:"avatar_url"`
	Nickname string `json:"nickname"`
	Gender uint32 `json:"gender"`
	Token string `json:"token,omitempty"`
}