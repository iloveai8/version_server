package domain

import (
	"fmt"

	"game_slots_vsn/pkg/errcode"
)

// GMConfig GM配置领域模型
//
// 表示GM（Game Master）后台的开关配置，包括GM功能开关和封锁状态
//
// 使用示例：
//   config := &domain.GMConfig{
//       GMEnable: true,
//       Block:    false,
//   }
//   if err := config.Validate(); err != nil {
//       return err
//   }
type GMConfig struct {
	GMEnable bool `json:"gmEnable"` // GM功能开关
	Block    bool `json:"block"`    // 封锁状态
}

// NewGMConfig 创建新的GM配置
//
// 参数:
//   gmEnable: GM功能是否启用
//   block: 是否封锁
//
// 返回:
//   *GMConfig: GM配置对象
func NewGMConfig(gmEnable, block bool) *GMConfig {
	return &GMConfig{
		GMEnable: gmEnable,
		Block:    block,
	}
}

// Validate 验证GM配置
//
// GM配置总是有效的，因为所有字段都是布尔类型且有默认值
// 此方法为了保持与其他领域模型的一致性而提供
//
// 返回:
//   error: 总是返回nil
func (g *GMConfig) Validate() error {
	// GM配置没有需要验证的字段，所有字段都是布尔类型
	return nil
}

// IsValid 快速检查GM配置是否有效
//
// 返回:
//   bool: 总是返回true
func (g *GMConfig) IsValid() bool {
	return true
}

// ToggleGM 切换GM功能开关
//
// 返回:
//   bool: 切换后的GM功能状态
func (g *GMConfig) ToggleGM() bool {
	g.GMEnable = !g.GMEnable
	return g.GMEnable
}

// ToggleBlock 切换封锁状态
//
// 返回:
//   bool: 切换后的封锁状态
func (g *GMConfig) ToggleBlock() bool {
	g.Block = !g.Block
	return g.Block
}

// SetGMEnable 设置GM功能开关
//
// 参数:
//   enable: GM功能状态
func (g *GMConfig) SetGMEnable(enable bool) {
	g.GMEnable = enable
}

// SetBlock 设置封锁状态
//
// 参数:
//   block: 封锁状态
func (g *GMConfig) SetBlock(block bool) {
	g.Block = block
}

// IsGMEnabled 检查GM功能是否启用
//
// 返回:
//   bool: GM功能启用返回true，否则返回false
func (g *GMConfig) IsGMEnabled() bool {
	return g.GMEnable
}

// IsBlocked 检查是否处于封锁状态
//
// 返回:
//   bool: 封锁状态返回true，否则返回false
func (g *GMConfig) IsBlocked() bool {
	return g.Block
}

// GetStatus 获取GM配置状态摘要
//
// 返回:
//   map[string]interface{}: 包含GM功能状态和封锁状态
func (g *GMConfig) GetStatus() map[string]interface{} {
	return map[string]interface{}{
		"gmEnable": g.GMEnable,
		"block":    g.Block,
	}
}

// String 返回GM配置的字符串表示
//
// 返回:
//   string: GM配置的描述
func (g *GMConfig) String() string {
	return fmt.Sprintf("GMConfig{GMEnable: %v, Block: %v}", g.GMEnable, g.Block)
}

// Clone 克隆GM配置
//
// 创建当前GM配置的深拷贝
//
// 返回:
//   *GMConfig: GM配置的副本
func (g *GMConfig) Clone() *GMConfig {
	return &GMConfig{
		GMEnable: g.GMEnable,
		Block:    g.Block,
	}
}

// Merge 合并另一个GM配置
//
// 当另一个配置不为nil时，使用其值覆盖当前配置
//
// 参数:
//   other: 要合并的GM配置
//
// 返回:
//   error: other为nil时返回错误
func (g *GMConfig) Merge(other *GMConfig) error {
	if other == nil {
		return errcode.ErrInvalidParams.WithDetails(map[string]interface{}{
			"reason": "不能合并nil配置",
		})
	}

	g.GMEnable = other.GMEnable
	g.Block = other.Block

	return nil
}
