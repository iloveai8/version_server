package utils

import (
	"game_slots_vsn/pkg/consts"
)

func MatchIp(IP string) bool {
	for _, SIP := range consts.IPList {
		if IP == SIP {
			return true
		}
	}
	return false
}
