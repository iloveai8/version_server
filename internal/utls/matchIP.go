package utls

import (
	`game_slots_vsn/internal/consts`
)

func MatchIp(IP string) bool {
	for _, SIP := range consts.IPList {
		if IP == SIP {
			return true
		}
	}
	return false
}
