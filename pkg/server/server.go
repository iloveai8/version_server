package server

import (
	"fmt"
	"game_slots_vsn/pkg/consts"
	"github.com/spf13/viper"
)

type setting struct {
	IP   string `yaml:"ip"`
	Port int    `yaml:"port"`
}

type server struct {
	S *setting
}

var Server = &server{}

func SetUp() {
	ss := &setting{}
	err := viper.UnmarshalKey(consts.ConfigServer, ss)
	if err != nil {
		panic(err)
	}
	Server.setup(ss)
}

func (s *server) setup(ss *setting) {
	fmt.Printf("server setting:%v\n", *ss)
	s.S = ss
}
