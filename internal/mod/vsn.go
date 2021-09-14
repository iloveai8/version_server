package mod

type Vsn struct {
	Vsn    string `json:"vsn"`
	SrvUrl string `json:"srvUrl"`
	ResUrl string `json:"resUrl"`
}

func NewVsn() * Vsn {
	return &Vsn{}
}
