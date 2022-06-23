package module

type GM struct {
	GMEnable bool `json:"gmEnable"`
}

func NewGM() *GM {
	return &GM{}
}
