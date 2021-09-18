package mod

import (
	"go.uber.org/zap/zapcore"
)

type Vsn struct {
	Vsn    string `json:"vsn"`
	SrvUrl string `json:"srvUrl"`
	ResUrl string `json:"resUrl"`
	Enable bool   `json:"enable"`
}

func (v Vsn) MarshalLogObject(zo zapcore.ObjectEncoder) error {
	zo.AddString("vsn", v.Vsn)
	zo.AddString("srvUrl", v.SrvUrl)
	zo.AddString("resUrl", v.ResUrl)
	zo.AddBool("enable", v.Enable)
	return nil
}

func NewVsn() *Vsn {
	return &Vsn{}
}
