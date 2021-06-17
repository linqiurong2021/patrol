package libs

// Errors
type Errors struct {}

// NewErrors NewErrors
func NewErrors() *Errors {
	return &Errors{}
}

const (
	RequestOk = "100000"
	RequestErr = "100001"
	ParamInvalidate = "100002"
	NoData = "100003"
	MethodInvalidate = "100004"
	AccessDenied = "100006"
	AuthFailure = "100005"
	MethodNotFound = "100007"
	HadLogin = "100008"
	OperateFailure = "100009"
	AuthExpired="100010"
	UnAuth="100011"
	ServerErr = "100012"

	UnknownErr = "999999"
)

var Text = map[string]string{
	RequestOk:  "请求成功",
	RequestErr: "请求失败",
	ParamInvalidate: "参数错误",
	NoData: "无数据",
	MethodInvalidate: "请求方式错误",
	AccessDenied: "无权限",
	AuthFailure: "授权失败",
	MethodNotFound: "请求路径出错",
	HadLogin: "用户已登录",
	OperateFailure: "操作失败",
	AuthExpired: "授权过期",
	UnAuth:"未授权",
	ServerErr: "系统错误",
	UnknownErr: "未知错误",
}

// Text 返回中文信息
func (e *Errors) Text(code string) string {
	str, ok := Text[code]
	if ok {
		return str
	}
	return Text[UnknownErr]
}
//
//func (e *Errors) String()  {
//
//}
