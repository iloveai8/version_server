package e

const (
	SUCCESS            = 200
	ParamError         = 201
	GmInfoNotFound     = 202
	ServerInfoNotFound = 203
	DataError          = 204
)

var MsgFlags = map[int]string{
	SUCCESS:            "ok",
	ParamError:         "params error",
	GmInfoNotFound:     "gm info not found",
	ServerInfoNotFound: "server info not found",
	DataError:          "data error",
}

// GetMsg get error information based on Code
func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}
	return MsgFlags[ParamError]
}
