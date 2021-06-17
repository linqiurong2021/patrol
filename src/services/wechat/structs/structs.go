package structs
//
type Code2Session struct {
	SessionKey string `json:"session_key"`
	OpenID string `json:"openid"`
	ErrCode int32 `json:"errcode,omitempty"`
	ErrMsg string `json:"errmsg,omitempty"`
}
