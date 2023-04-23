package ggeoip

import (
	"context"
	"fmt"
	"gitlab.ftsview.com/fotoable-go/ggeoip"
	"testing"
	"time"
)

func TestGetIpData(t *testing.T) {
	parentCxt := context.TODO()
	ctx, cancel := context.WithCancel(parentCxt)
	util := ggeoip.NewGeoIpUtil(parentCxt, &ggeoip.GeoIpConfig{
		Ctx:                ctx,
		FileName:           "ggeoip.mmdb",
		Path:               "E:\\go\\src\\game_slots_vsn\\data\\",
		Scope:              1,
		UpdateIntervalHour: 5,
		Fun:                isSuccess,
	})
	util.GeoIPInit()
	time.Sleep(time.Second * 10)
	fmt.Println(ggeoip.GetCountryAndCityByIP("1.202.246.19"))
	cancel()
	time.Sleep(time.Second * 10)
}

func TestLoadLocalFile(t *testing.T) {
	parentCxt := context.TODO()
	config := &ggeoip.GeoIpConfig{
		Ctx:                parentCxt,
		FileName:           "ggeoip.mmdb",
		Path:               "E:\\go\\src\\game_slots_vsn\\data\\",
		Scope:              1,
		UpdateIntervalHour: 5,
		Fun:                isSuccess,
	}
	ggeoip.LoadLocalFile(parentCxt, config, true)
	fmt.Println(ggeoip.GetCountryAndCityByIP("1.202.246.19"))
}
