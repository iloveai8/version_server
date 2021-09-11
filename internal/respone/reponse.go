package respone

var (
	OK           = Error{200, "ok"}
	ParamsError  = Error{201, "参数错误"}
)

type Error struct {
	code    uint16
	message string
}

type Response struct {
	Code    uint16      `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
}

// 成功结构
func Success(data interface{}) Response {
	return Response{
		Code:    OK.code,
		Message: OK.message,
		Data:  data ,
	}
}

// 失败
func Fail(error Error, errors interface{}) Response {
	return Response{
		Code:    error.code,
		Message: error.message,
		Errors:  errors,
	}
}
