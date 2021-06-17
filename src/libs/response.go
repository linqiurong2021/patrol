package libs

// Response 返回数据
type Response struct {
	// 返回码
	Code string `json:"code"`
	// 返回说明
	Msg string `json:"msg"`
	// 数据
	Data interface{} `json:"data,omitempty";`
}

// NewResponse 返回数据
func NewResponse(code string, msg string, data interface{}) Response {
	return Response{
		Code: code,
		Msg:  msg,
		Data: data,
	}
}

// Success 成功返回
func Success(msg string, data interface{}) Response {

	return NewResponse(RequestOk, msg, data)
}

// CreateFailure 新增失败
func CreateFailure(msg string, data interface{}) Response {
	return NewResponse(OperateFailure, "create failure", "")
}

// BadRequest 失败返回
func BadRequest(msg string, data interface{}) Response {
	return NewResponse(RequestErr, msg, data)
}

// Forbidden 无权限
func Forbidden(msg string, data interface{}) Response {
	return NewResponse(AccessDenied, msg, data)
}

// Unauthorized 未授权
func Unauthorized(msg string, data interface{}) Response {
	return NewResponse(UnAuth, msg, data)
}

// ServerError 系统内部错误
func ServerError(msg string, data interface{}) Response {
	return NewResponse(ServerErr, msg, data)
}

// ValidFailure 校验
func ValidFailure(data interface{}) Response {
	return NewResponse(ParamInvalidate, "valid failure! please see the field data for details.", data)
}
