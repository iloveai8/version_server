package module

type GM struct {
	GMEnable bool `json:"gmEnable"`
	Block    bool `json:"block"`
}

func NewGM() *GM {
	return &GM{}
}
