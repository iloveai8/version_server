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
	Name     string           `yaml:"name"`
	Path     string           `yaml:"path"`
	Url      string           `yaml:"url"`
	Duration int              `yaml:"duration"`
	Scope    ggeoip.ScopeType `yaml:"scope"`
}
type gip struct {
	S   *setting
	ctx context.Context
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
	ctx := g.ctx
	if ctx != nil {
		ctx1, cancel := context.WithCancel(ctx)
		ctx = ctx1
		cancel()
	} else {
		ctx = context.TODO()
	}
	geoIpConfig := &ggeoip.GeoIpConfig{
		Ctx:                ctx,
		GeoIpURL:           ips.Url,
		FileName:           ips.Name,
		Path:               ips.Path,
		UpdateIntervalHour: ips.Duration,
		Scope:              ips.Scope,
		Fun:                isSuccess,
	}
	ggeoip.LoadLocalFile(ctx, geoIpConfig, true)
	g.S = ips
	g.ctx = ctx
}

func isSuccess(is bool) {
	fmt.Println("update ipdb data", time.Now(), is)
}

func (g *gip) GetCountryByIP(ip string) string {
	return ggeoip.GetCountryByIP(ip)
}

func (g *gip) GetCountryAndCityByIP(ip string) (string, map[string]string) {
	return ggeoip.GetCountryAndCityByIP(ip)
}
