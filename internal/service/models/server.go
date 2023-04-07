package models

type ServerInfo struct {
	SubServerInfoMap subServerInfoMap `json:"subServer"`
}

type subServerInfoMap map[string]*SubServerInfo

type SubServerInfo struct {
	Vsn    string `json:"vsn"`
	SrvUrl string `json:"srvUrl"`
	ResUrl string `json:"resUrl"`
	Type   int    `json:"type"`
}
