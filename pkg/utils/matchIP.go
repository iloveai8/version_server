package utils

import (
	"game_slots_vsn/internal/service"
)

//	func MatchIP(IP string) bool {
//		for _, SIP := range consts.IPList {
//			if IP == SIP {
//				return true
//			}
//		}
//		return false
//	}

func IsGmIP(ip string) bool {
	ipService := service.GmIPService{}
	return ipService.IsGmIP(ip)
}
