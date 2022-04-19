package rsp

var (
	OK                 = Error{200, "ok"}
	ParamsError        = Error{201, "params error"}
	GMInfoNotFound     = Error{202, "gm info not found"}
	ServerInfoNotFound = Error{203, "server info not found"}
	DataError          = Error{204, "data error"}
)

type Error struct {
	code    uint16
	message string
}

type Rsp struct {
	Code    uint16      `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
}

// 成功结构
func Success(data interface{}) Rsp {
	return Rsp{
		Code:    OK.code,
		Message: OK.message,
		Data:    data,
	}
}

// 失败
func Fail(error Error, errors interface{}) Rsp {
	return Rsp{
		Code:    error.code,
		Message: error.message,
		Errors:  errors,
	}
}
