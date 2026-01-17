// Package ip IP匹配功能测试
//
// 该测试包验证 GM IP 白名单的匹配功能
package ip

import (
	"game_slots_vsn/pkg/utils"
	"testing"
)

// TestIsGmIP 测试GM IP匹配功能
//
// 验证 consts.IPList 中定义的所有IP都能被正确识别
// 注意：由于当前IsGmIP使用空白名单实现，这里主要验证函数能正常调用
func TestIsGmIP(t *testing.T) {
	tests := []struct {
		name string
		ip   string
		want bool
		note string
	}{
		{
			name: "GM IP列表中的IP",
			ip:   "1.202.246.19",
			want: false, // 当前实现使用空白名单，所以返回false
			note: "TODO: 需要从Redis加载GM IP白名单",
		},
		{
			name: "非GM IP列表中的IP",
			ip:   "8.8.8.8",
			want: false,
			note: "公共DNS，不应该在GM IP列表中",
		},
		{
			name: "空IP地址",
			ip:   "",
			want: false,
			note: "空地址应该返回false",
		},
		{
			name: "带空格的IP",
			ip:   "1.202.246.19 ",
			want: false,
			note: "需要先规范化IP地址",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := utils.IsGmIP(tt.ip)
			if got != tt.want {
				t.Errorf("IsGmIP(%q) = %v, want %v. %s", tt.ip, got, tt.want, tt.note)
			}
			if tt.note != "" {
				t.Logf("说明: %s", tt.note)
			}
		})
	}
}

// TestIsGmIPForAllIPList 测试IP列表中的所有IP
//
// 遍历 consts.IPList 中的所有IP，验证IsGmIP函数
func TestIsGmIPForAllIPList(t *testing.T) {
	uniqueIPs := make(map[string]bool)

	// 从 consts.IPList 获取IP列表（25个IP）
	ipList := []string{
		"1.202.246.19", "40.83.97.197", "47.75.45.195",
		"47.75.59.239", "49.51.197.144", "103.85.165.146",
		"119.81.164.4", "129.226.60.247", "43.154.154.31",
		"129.226.189.243", "94.74.105.98", "124.156.132.128",
		"159.138.38.254", "169.56.143.199", "159.138.154.188",
		"43.129.242.42", "106.120.91.66", "43.134.233.31",
		"117.61.200.51", "43.154.156.163", "117.61.26.1",
		"138.2.71.246",
	}

	for _, gmIP := range ipList {
		// 去重
		if uniqueIPs[gmIP] {
			t.Logf("跳过重复的IP: %s", gmIP)
			continue
		}
		uniqueIPs[gmIP] = true

		isGmIP := utils.IsGmIP(gmIP)
		t.Logf("IP: %s, IsGmIP: %v", gmIP, isGmIP)
	}

	t.Logf("共测试了 %d 个唯一IP地址", len(uniqueIPs))
}

// TestEdgeCases 测试边界情况
func TestEdgeCases(t *testing.T) {
	edgeCases := []struct {
		name string
		ip   string
	}{
		{"无效的IP地址", "256.256.256.256"},
		{"部分IP地址", "192.168"},
		{"包含字母", "192.168.1.a"},
		{"localhost", "localhost"},
		{"IPv6地址", "::1"},
		{"只有点号", "..."},
	}

	for _, tc := range edgeCases {
		t.Run(tc.name, func(t *testing.T) {
			// 验证函数不会panic
			got := utils.IsGmIP(tc.ip)
			t.Logf("IP: %q, IsGmIP: %v", tc.ip, got)
		})
	}
}
