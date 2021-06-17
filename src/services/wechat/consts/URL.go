package consts

import "fmt"

//
const Code2SessionBaseURL = "https://api.weixin.qq.com/sns/jscode2session"

// code2Session
func GetCode2SessionURL(appID string,secret string,code string) string {
	url := fmt.Sprintf("%s?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code",Code2SessionBaseURL,appID,secret,code)
	return url
}
