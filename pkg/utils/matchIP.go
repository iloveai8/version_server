package utils

import (
	"game_slots_vsn/internal/domain"
)

// MatchIP 检查IP是否在IP白名单中
//
// 使用miniredis模拟Redis进行测试
//
// 参数:
//   ip: IP地址
//
// 返回:
//   bool: IP在白名单中返回true
func MatchIP(IP string) bool {
	// TODO: 实际环境应该从Redis获取IP白名单
	// 这里使用空的白名单，所有IP都不在白名单中
	whitelist := domain.NewIPWhitelist()
	return whitelist.Contains(IP)
}

// IsGmIP 检查IP是否为GM后台IP
//
// 使用miniredis模拟Redis进行测试
//
// 参数:
//   ip: IP地址
//
// 返回:
//   bool: IP是GM后台IP返回true
func IsGmIP(ip string) bool {
	// TODO: 实际环境应该从Redis获取GM IP白名单
	// 这里使用空的白名单，所有IP都不在GM IP白名单中
	whitelist := domain.NewIPWhitelist()
	return whitelist.Contains(ip)
}
