package mod

type Vsn struct {
	Vsn    string `json:"vsn"`
	SrvUrl string `json:"server_url"`
	ResUrl string `json:"resource_url"`
}

func NewVsn() * Vsn {
	return &Vsn{}
}
