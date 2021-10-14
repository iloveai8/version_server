package mod

type Vsn struct {
	Vsn    string `json:"vsn"`
	SrvUrl string `json:"srvUrl"`
	ResUrl string `json:"resUrl"`
	Enable bool   `json:"enable"`
}

func NewVsn() *Vsn {
	return &Vsn{}
}

type GMConf struct {
	GMSrvUrl string `json:"gmSrvUrl"`
	GMResUrl string `json:"gmResUrl"`
	GMEnable bool   `json:"gmEnable"`
}

func NewGmConf() *GMConf {
	return &GMConf{}
}
