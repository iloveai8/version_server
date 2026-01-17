package domain

import (
	"fmt"
	"net"
	"strings"
	"sync"

	"game_slots_vsn/pkg/errcode"
	"game_slots_vsn/pkg/validator"
)

// IPWhitelist IP白名单领域模型
//
// 表示允许访问的IP地址列表，用于GM后台等需要IP限制的场景
//
// 使用示例：
//   whitelist := domain.NewIPWhitelist()
//   if err := whitelist.AddIP("192.168.1.1"); err != nil {
//       return err
//   }
//   if whitelist.Contains("192.168.1.1") {
//       // 允许访问
//   }
type IPWhitelist struct {
	mu  sync.RWMutex
	IPs map[string]bool `json:"ips"` // IP地址集合（使用map保证唯一性）
}

// NewIPWhitelist 创建新的IP白名单
//
// 返回:
//   *IPWhitelist: IP白名单对象
func NewIPWhitelist() *IPWhitelist {
	return &IPWhitelist{
		IPs: make(map[string]bool),
	}
}

// NewIPWhitelistFromSlice 从IP列表创建IP白名单
//
// 支持单个IP地址和CIDR网段格式
//
// 参数:
//   ips: IP地址或CIDR列表
//
// 返回:
//   *IPWhitelist: IP白名单对象
//   error: 如果有IP格式无效，返回错误
func NewIPWhitelistFromSlice(ips []string) (*IPWhitelist, error) {
	w := &IPWhitelist{
		IPs: make(map[string]bool),
	}

	if err := w.AddIPs(ips); err != nil {
		return nil, err
	}

	return w, nil
}

// Validate 验证IP白名单配置
//
// 验证所有IP地址的格式
//
// 返回:
//   error: 如果有IP格式无效，返回错误
func (i *IPWhitelist) Validate() error {
	i.mu.RLock()
	defer i.mu.RUnlock()

	// IP白名单总是有效的，即使是空的
	return nil
}

// IsValid 快速检查IP白名单是否有效
//
// 返回:
//   bool: 总是返回true
func (i *IPWhitelist) IsValid() bool {
	return true
}

// AddIP 添加IP到白名单
//
// 支持单个IP地址和CIDR网段格式
//
// 参数:
//   ip: IP地址或CIDR字符串
//
// 返回:
//   error: IP格式无效或已存在时返回错误
func (i *IPWhitelist) AddIP(ip string) error {
	// 标准化IP地址
	normalizedIP := strings.TrimSpace(ip)

	// 验证IP或CIDR格式
	if err := validator.ValidateIPOrCIDR(normalizedIP); err != nil {
		return errcode.ErrInvalidIPFormat.WithDetails(map[string]interface{}{
			"ip":     normalizedIP,
			"reason": "IP地址或CIDR格式无效",
		})
	}

	i.mu.Lock()
	defer i.mu.Unlock()

	// 检查是否已存在
	if _, exists := i.IPs[normalizedIP]; exists {
		return errcode.ErrDuplicateIP.WithDetails(map[string]interface{}{
			"ip":     normalizedIP,
			"reason": "IP已存在于白名单中",
		})
	}

	// 添加IP
	i.IPs[normalizedIP] = true
	return nil
}

// RemoveIP 从白名单移除IP
//
// 参数:
//   ip: IP地址
//
// 返回:
//   error: IP不存在时返回错误
func (i *IPWhitelist) RemoveIP(ip string) error {
	normalizedIP := strings.TrimSpace(ip)

	i.mu.Lock()
	defer i.mu.Unlock()

	if _, exists := i.IPs[normalizedIP]; !exists {
		return errcode.ErrIPNotFound.WithDetails(map[string]interface{}{
			"ip":     normalizedIP,
			"reason": "IP不在白名单中",
		})
	}

	delete(i.IPs, normalizedIP)
	return nil
}

// Contains 检查IP是否在白名单中
//
// 参数:
//   ip: IP地址
//
// 返回:
//   bool: IP在白名单中返回true，否则返回false
func (i *IPWhitelist) Contains(ip string) bool {
	normalizedIP := strings.TrimSpace(ip)

	i.mu.RLock()
	defer i.mu.RUnlock()

	return i.IPs[normalizedIP]
}

// ContainsIPNet 检查IP是否在白名单的IP网络范围内
//
// 支持CIDR格式的IP网段匹配
//
// 参数:
//   ip: IP地址
//
// 返回:
//   bool: IP在白名单或某个IP网段范围内返回true
func (i *IPWhitelist) ContainsIPNet(ip string) bool {
	normalizedIP := strings.TrimSpace(ip)
	checkIP := net.ParseIP(normalizedIP)
	if checkIP == nil {
		return false
	}

	i.mu.RLock()
	defer i.mu.RUnlock()

	// 检查精确匹配
	if i.IPs[normalizedIP] {
		return true
	}

	// 检查CIDR网段匹配
	for whitelistIP := range i.IPs {
		if strings.Contains(whitelistIP, "/") {
			_, ipNet, err := net.ParseCIDR(whitelistIP)
			if err != nil {
				continue
			}
			if ipNet.Contains(checkIP) {
				return true
			}
		}
	}

	return false
}

// AddIPs 批量添加IP到白名单
//
// 参数:
//   ips: IP地址列表
//
// 返回:
//   error: 如果有任何IP格式无效，返回错误（已添加的IP不会回滚）
func (i *IPWhitelist) AddIPs(ips []string) error {
	var lastErr error
	successCount := 0

	for _, ip := range ips {
		if err := i.AddIP(ip); err != nil {
			lastErr = err
		} else {
			successCount++
		}
	}

	// 如果有错误，返回最后一个错误
	if lastErr != nil && successCount == 0 {
		return lastErr
	}

	return nil
}

// RemoveIPs 批量从白名单移除IP
//
// 参数:
//   ips: IP地址列表
//
// 返回:
//   error: 如果有任何IP不存在，返回错误（已移除的IP不会回滚）
func (i *IPWhitelist) RemoveIPs(ips []string) error {
	var lastErr error
	successCount := 0

	for _, ip := range ips {
		if err := i.RemoveIP(ip); err != nil {
			lastErr = err
		} else {
			successCount++
		}
	}

	// 如果有错误且没有成功移除任何IP，返回最后一个错误
	if lastErr != nil && successCount == 0 {
		return lastErr
	}

	return nil
}

// GetIPs 获取所有IP地址列表
//
// 返回:
//   []string: IP地址列表
func (i *IPWhitelist) GetIPs() []string {
	i.mu.RLock()
	defer i.mu.RUnlock()

	ips := make([]string, 0, len(i.IPs))
	for ip := range i.IPs {
		ips = append(ips, ip)
	}
	return ips
}

// Count 获取白名单中IP数量
//
// 返回:
//   int: IP数量
func (i *IPWhitelist) Count() int {
	i.mu.RLock()
	defer i.mu.RUnlock()

	return len(i.IPs)
}

// IsEmpty 检查白名单是否为空
//
// 返回:
//   bool: 白名单为空返回true，否则返回false
func (i *IPWhitelist) IsEmpty() bool {
	return i.Count() == 0
}

// Clear 清空白名单
func (i *IPWhitelist) Clear() {
	i.mu.Lock()
	defer i.mu.Unlock()

	i.IPs = make(map[string]bool)
}

// Clone 克隆IP白名单
//
// 创建当前IP白名单的深拷贝
//
// 返回:
//   *IPWhitelist: IP白名单的副本
func (i *IPWhitelist) Clone() *IPWhitelist {
	i.mu.RLock()
	defer i.mu.RUnlock()

	newWhitelist := NewIPWhitelist()
	for ip := range i.IPs {
		newWhitelist.IPs[ip] = true
	}

	return newWhitelist
}

// ToSlice 转换为IP地址切片
//
// 返回:
//   []string: IP地址列表
func (i *IPWhitelist) ToSlice() []string {
	return i.GetIPs()
}

// String 返回IP白名单的字符串表示
//
// 返回:
//   string: IP白名单的描述
func (i *IPWhitelist) String() string {
	ips := i.GetIPs()
	return fmt.Sprintf("IPWhitelist{Count: %d, IPs: [%s]}", len(ips), strings.Join(ips, ", "))
}

// Merge 合并另一个IP白名单
//
// 将另一个白名单的所有IP添加到当前白名单
//
// 参数:
//   other: 要合并的IP白名单
//
// 返回:
//   error: other为nil时返回错误
func (i *IPWhitelist) Merge(other *IPWhitelist) error {
	if other == nil {
		return errcode.ErrInvalidParams.WithDetails(map[string]interface{}{
			"reason": "不能合并nil白名单",
		})
	}

	other.mu.RLock()
	ips := make([]string, 0, len(other.IPs))
	for ip := range other.IPs {
		ips = append(ips, ip)
	}
	other.mu.RUnlock()

	return i.AddIPs(ips)
}

// ValidateIPs 验证IP地址列表的格式
//
// 支持单个IP地址和CIDR网段格式
//
// 参数:
//   ips: IP地址或CIDR列表
//
// 返回:
//   error: 如果有IP格式无效，返回包含所有错误信息的错误
func (i *IPWhitelist) ValidateIPs(ips []string) error {
	errors := make([]string, 0)

	for _, ip := range ips {
		normalizedIP := strings.TrimSpace(ip)
		if err := validator.ValidateIPOrCIDR(normalizedIP); err != nil {
			errors = append(errors, fmt.Sprintf("IP '%s': %v", normalizedIP, err))
		}
	}

	if len(errors) > 0 {
		return errcode.ErrInvalidIPFormat.WithDetails(map[string]interface{}{
			"reason": "部分IP格式无效",
			"errors": errors,
		})
	}

	return nil
}
