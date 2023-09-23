package service

import (
	"game_slots_vsn/internal/service/dao"
)

type GmIPService struct {
}

func (gmIP *GmIPService) GetGmIPList() []string {
	return dao.GetGmIPList()
}

func (gmIP *GmIPService) AddGmIP(ips ...string) bool {
	return dao.AddGmIP(ips)
}

func (gmIP *GmIPService) RemGmIP(ips ...string) bool {
	return dao.RemGmIP(ips)
}

func (gmIP *GmIPService) IsGmIP(ip string) bool {
	return dao.IsGmIP(ip)
}
