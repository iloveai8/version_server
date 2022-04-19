package module

type GM struct {
	GMSrvUrl string `json:"gmSrvUrl"`
	GMResUrl string `json:"gmResUrl"`
	GMEnable bool   `json:"gmEnable"`
}

func NewGM() *GM {
	return &GM{}
}
