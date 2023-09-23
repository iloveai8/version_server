package ip

import (
	"fmt"
	"game_slots_vsn/internal/service"
	"game_slots_vsn/pkg/consts"
	"game_slots_vsn/pkg/redis"
	"game_slots_vsn/pkg/utils"
	"github.com/spf13/viper"
	"testing"
)

func init() {
	viper.AddConfigPath("../../conf")
	viper.SetConfigType("yaml")
	viper.SetConfigName("pre")
	viper.AutomaticEnv() // read in environment variables that match
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
	redis.SetUp()
}

func TestAddIP(t *testing.T) {
	ipService := service.GmIPService{}
	for _, gmIP := range consts.IPList {
		add := ipService.AddGmIP(gmIP)
		fmt.Println("ip：", gmIP, " add:", add)
	}
}

func TestGetIPList(t *testing.T) {
	ipService := service.GmIPService{}
	ipList := ipService.GetGmIPList()
	for i, gmIP := range ipList {
		fmt.Println("i:", i, " gmIP", gmIP)
	}
}

func TestRemIP(t *testing.T) {
	ipService := service.GmIPService{}
	ipList := ipService.GetGmIPList()
	for i, gmIP := range ipList {
		rem := ipService.RemGmIP(gmIP)
		fmt.Println("i:", i, " ip", gmIP, " rem:", rem)
	}
}

func TestMatchIP(t *testing.T) {
	for i, gmIP := range consts.IPList {
		isGmIP := utils.IsGmIP(gmIP)
		if !isGmIP {
			fmt.Println("i:", i, " gmIP", gmIP, " isGmIP:", isGmIP)
		}
	}
	gmIP := ""
	isGmIP := utils.IsGmIP(gmIP)
	fmt.Println(" gmIP", gmIP, " isGmIP:", isGmIP)

	gmIP = "1.202.246.19 "
	isGmIP = utils.IsGmIP(gmIP)
	fmt.Println(" gmIP", gmIP, " isGmIP:", isGmIP)

	gmIP = "1.202.246.19"
	isGmIP = utils.IsGmIP(gmIP)
	fmt.Println(" gmIP", gmIP, " isGmIP:", isGmIP)

}
