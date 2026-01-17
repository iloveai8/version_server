// Package repository 提供数据访问层的抽象和实现
//
// 该包实现了：
// - Repository接口定义
// - Redis客户端封装
// - VersionRepository: 版本配置数据访问
// - GMRepository: GM配置数据访问
// - IPWhitelistRepository: IP白名单数据访问
//
// 设计原则：
// - 接口隔离：定义清晰的Repository接口
// - 依赖倒置：高层模块依赖抽象接口而非具体实现
// - 可测试性：使用miniredis进行集成测试
package repository

import (
	"context"

	"game_slots_vsn/internal/domain"
)

// ==================== Repository接口定义 ====================

// VersionRepository 版本配置仓储接口
//
// 定义版本配置的持久化操作
type VersionRepository interface {
	// GetByVersionAndEnv 根据版本号和环境获取版本配置
	//
	// 参数:
	//   ctx: 上下文
	//   vsn: 版本号
	//   env: 环境（dev或pro）
	//
	// 返回:
	//   *domain.Version: 版本配置对象
	//   error: 版本不存在或发生错误时返回错误
	GetByVersionAndEnv(ctx context.Context, vsn, env string) (*domain.Version, error)

	// Save 保存版本配置
	//
	// 参数:
	//   ctx: 上下文
	//   version: 版本配置对象
	//   env: 环境（dev或pro）
	//
	// 返回:
	//   error: 保存失败时返回错误
	Save(ctx context.Context, version *domain.Version, env string) error

	// Delete 删除版本配置
	//
	// 参数:
	//   ctx: 上下文
	//   vsn: 版本号
	//   env: 环境（dev或pro）
	//
	// 返回:
	//   error: 删除失败时返回错误
	Delete(ctx context.Context, vsn, env string) error

	// Exists 检查版本是否存在
	//
	// 参数:
	//   ctx: 上下文
	//   vsn: 版本号
	//   env: 环境（dev或pro）
	//
	// 返回:
	//   bool: 存在返回true，否则返回false
	//   error: 发生错误时返回错误
	Exists(ctx context.Context, vsn, env string) (bool, error)

	// List 获取所有版本列表
	//
	// 参数:
	//   ctx: 上下文
	//   env: 环境（dev或pro）
	//
	// 返回:
	//   []*domain.Version: 版本配置列表
	//   error: 查询失败时返回错误
	List(ctx context.Context, env string) ([]*domain.Version, error)
}

// GMRepository GM配置仓储接口
//
// 定义GM配置的持久化操作
type GMRepository interface {
	// Get 获取GM配置
	//
	// 参数:
	//   ctx: 上下文
	//
	// 返回:
	//   *domain.GMConfig: GM配置对象
	//   error: 配置不存在或发生错误时返回错误
	Get(ctx context.Context) (*domain.GMConfig, error)

	// Save 保存GM配置
	//
	// 参数:
	//   ctx: 上下文
	//   config: GM配置对象
	//
	// 返回:
	//   error: 保存失败时返回错误
	Save(ctx context.Context, config *domain.GMConfig) error

	// Exists 检查GM配置是否存在
	//
	// 参数:
	//   ctx: 上下文
	//
	// 返回:
	//   bool: 存在返回true，否则返回false
	//   error: 发生错误时返回错误
	Exists(ctx context.Context) (bool, error)
}

// IPWhitelistRepository IP白名单仓储接口
//
// 定义IP白名单的持久化操作
type IPWhitelistRepository interface {
	// Get 获取IP白名单
	//
	// 参数:
	//   ctx: 上下文
	//
	// 返回:
	//   *domain.IPWhitelist: IP白名单对象
	//   error: 查询失败时返回错误
	Get(ctx context.Context) (*domain.IPWhitelist, error)

	// Save 保存IP白名单
	//
	// 参数:
	//   ctx: 上下文
	//   whitelist: IP白名单对象
	//
	// 返回:
	//   error: 保存失败时返回错误
	Save(ctx context.Context, whitelist *domain.IPWhitelist) error

	// Exists 检查IP白名单是否存在
	//
	// 参数:
	//   ctx: 上下文
	//
	// 返回:
	//   bool: 存在返回true，否则返回false
	//   error: 发生错误时返回错误
	Exists(ctx context.Context) (bool, error)

	// Contains 检查IP是否在白名单中
	//
	// 参数:
	//   ctx: 上下文
	//   ip: IP地址
	//
	// 返回:
	//   bool: 在白名单中返回true，否则返回false
	//   error: 查询失败时返回错误
	Contains(ctx context.Context, ip string) (bool, error)
}

// ==================== Redis键定义 ====================

const (
	// Redis键前缀
	KeyPrefixVersion     = "version:"      // 版本配置键前缀
	KeyPrefixGMConfig    = "gm:"           // GM配置键前缀
	KeyPrefixIPWhitelist = "ip_whitelist:" // IP白名单键前缀

	// 特殊键
	KeyGMConfig    = "gm:config"    // GM配置键
	KeyIPWhitelist = "ip:whitelist" // IP白名单键
)

// BuildVersionKey 构建版本配置的Redis键
//
// 参数:
//
//	vsn: 版本号
//	env: 环境（dev或pro）
//
// 返回:
//
//	string: Redis键，格式为 "version:{env}:{vsn}"
func BuildVersionKey(vsn, env string) string {
	return KeyPrefixVersion + env + ":" + vsn
}

// BuildVersionListKey 构建版本列表的Redis键
//
// 参数:
//
//	env: 环境（dev或pro）
//
// 返回:
//
//	string: 版本列表键，格式为 "version:{env}:list"
func BuildVersionListKey(env string) string {
	return KeyPrefixVersion + env + ":list"
}
