package app

import (
	"fmt"
	"game_slots_vsn/pkg/consts"
	"github.com/spf13/viper"
)

type setting struct {
	RunMode string `yaml:"runMode"`
}

type app struct {
	S *setting
}

var App = &app{}

func SetUp() {
	as := &setting{}
	err := viper.UnmarshalKey(consts.ConfigApp, as)
	if err != nil {
		panic(err)
	}
	App.setup(as)
}

func (a *app) setup(as *setting) {
	fmt.Printf("app setting:%v\n", *as)
	a.S = as
}

func RunMode() string {
	return App.S.RunMode
}
