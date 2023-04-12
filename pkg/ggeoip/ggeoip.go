package ggeoip

import (
	"context"
	"fmt"
	"game_slots_vsn/pkg/consts"
	"github.com/spf13/viper"
	"gitlab.ftsview.com/fotoable-go/ggeoip"
	"time"
)

type setting struct {
	Name     string `yaml:"name"`
	Path     string `yaml:"path"`
	Url      string `yaml:"url"`
	Duration int    `yaml:"duration"`
}
type gip struct {
	S      *setting
	ipUtil *ggeoip.GeoIpUtil
}

var Gip = &gip{}

func SetUp() {
	ipSetting := &setting{}
	err := viper.UnmarshalKey(consts.ConfigGgeoip, ipSetting)
	if err != nil {
		panic(err)
	}
	Gip.setup(ipSetting)
}

func (g *gip) setup(ips *setting) {
	fmt.Printf("ggeoip setting:%v\n", *ips)

	ctx := context.TODO()
	util := ggeoip.NewGeoIpUtil(
		ctx,
		&ggeoip.GeoIpConfig{
			Cxt:                ctx,
			GeoIpURL:           ips.Url,
			FileName:           ips.Name,
			Path:               ips.Path,
			UpdateIntervalHour: ips.Duration,
			Fun:                isSuccess,
		},
	)
	g.S = ips
	g.ipUtil = util
	g.ipUtil.GeoIPInit()
}

func isSuccess(is bool) {
	fmt.Println("update ipdb data", time.Now(), is)
}

func (g *gip) GetIP(ip string) string {
	return ggeoip.GetCountryByIP(ip)
}
