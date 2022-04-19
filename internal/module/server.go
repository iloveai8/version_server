package module

type Server struct {
	Vsn    string `json:"vsn"`
	SrvUrl string `json:"srvUrl"`
	ResUrl string `json:"resUrl"`
	Enable bool   `json:"enable"`
	Type   int8   `json:"type"`
}

func NewServer() *Server {
	return &Server{}
}
