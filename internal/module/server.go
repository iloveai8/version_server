package module

type Subs map[string]*SubServer

type Server struct {
	SubServer Subs `json:"subServer"`
}

type SubServer struct {
	Vsn    string `json:"vsn"`
	SrvUrl string `json:"srvUrl"`
	ResUrl string `json:"resUrl"`
	//Enable bool   `json:"enable"`
	Type int `json:"type"`
}

func NewServer() *Server {
	return &Server{
		SubServer: make(map[string]*SubServer),
	}
}

func NewSubServer() *SubServer {
	return &SubServer{}
}
