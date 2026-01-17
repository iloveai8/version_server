// Package service 提供业务逻辑层的实现
//
// 该包实现了：
// - Service接口定义
// - VersionService: 版本配置业务逻辑
// - GMService: GM配置业务逻辑
// - IPWhitelistService: IP白名单业务逻辑
//
// 设计原则：
// - 业务逻辑封装：核心业务规则在Service层实现
// - 协调层：协调多个Repository完成复杂业务
// - 事务处理：跨多个Repository的操作协调
// - 可测试性：使用mock Repository进行单元测试
package service

import (
	"context"

	"game_slots_vsn/internal/domain"
)

// ==================== Service接口定义 ====================

// VersionService 版本配置服务接口
//
// 定义版本配置的业务操作
type VersionService interface {
	// GetVersion 获取版本配置
	//
	// 参数:
	//   ctx: 上下文
	//   vsn: 版本号
	//   env: 环境（dev或pro）
	//
	// 返回:
	//   *domain.Version: 版本配置对象
	//   error: 版本不存在或发生错误时返回错误
	GetVersion(ctx context.Context, vsn, env string) (*domain.Version, error)

	// CreateVersion 创建新版本
	//
	// 参数:
	//   ctx: 上下文
	//   version: 版本配置对象
	//   env: 环境（dev或pro）
	//
	// 返回:
	//   error: 创建失败时返回错误
	CreateVersion(ctx context.Context, version *domain.Version, env string) error

	// UpdateVersion 更新版本配置
	//
	// 参数:
	//   ctx: 上下文
	//   version: 版本配置对象
	//   env: 环境（dev或pro）
	//
	// 返回:
	//   error: 更新失败时返回错误
	UpdateVersion(ctx context.Context, version *domain.Version, env string) error

	// DeleteVersion 删除版本
	//
	// 参数:
	//   ctx: 上下文
	//   vsn: 版本号
	//   env: 环境（dev或pro）
	//
	// 返回:
	//   error: 删除失败时返回错误
	DeleteVersion(ctx context.Context, vsn, env string) error

	// ListVersions 获取所有版本列表
	//
	// 参数:
	//   ctx: 上下文
	//   env: 环境（dev或pro）
	//
	// 返回:
	//   []*domain.Version: 版本配置列表
	//   error: 查询失败时返回错误
	ListVersions(ctx context.Context, env string) ([]*domain.Version, error)

	// AddSubServer 添加子服务器到版本
	//
	// 参数:
	//   ctx: 上下文
	//   vsn: 版本号
	//   env: 环境（dev或pro）
	//   key: 子服务器键
	//   sub: 子服务器配置
	//
	// 返回:
	//   error: 添加失败时返回错误
	AddSubServer(ctx context.Context, vsn, env, key string, sub *domain.SubServer) error

	// RemoveSubServer 从版本移除子服务器
	//
	// 参数:
	//   ctx: 上下文
	//   vsn: 版本号
	//   env: 环境（dev或pro）
	//   key: 子服务器键
	//
	// 返回:
	//   error: 移除失败时返回错误
	RemoveSubServer(ctx context.Context, vsn, env, key string) error
}

// GMService GM配置服务接口
//
// 定义GM配置的业务操作
type GMService interface {
	// GetGMConfig 获取GM配置
	//
	// 参数:
	//   ctx: 上下文
	//
	// 返回:
	//   *domain.GMConfig: GM配置对象
	//   error: 配置不存在或发生错误时返回错误
	GetGMConfig(ctx context.Context) (*domain.GMConfig, error)

	// UpdateGMConfig 更新GM配置
	//
	// 参数:
	//   ctx: 上下文
	//   config: GM配置对象
	//
	// 返回:
	//   error: 更新失败时返回错误
	UpdateGMConfig(ctx context.Context, config *domain.GMConfig) error

	// ToggleGM 切换GM功能开关
	//
	// 参数:
	//   ctx: 上下文
	//
	// 返回:
	//   bool: 切换后的GM功能状态
	//   error: 操作失败时返回错误
	ToggleGM(ctx context.Context) (bool, error)

	// ToggleBlock 切换封锁状态
	//
	// 参数:
	//   ctx: 上下文
	//
	// 返回:
	//   bool: 切换后的封锁状态
	//   error: 操作失败时返回错误
	ToggleBlock(ctx context.Context) (bool, error)

	// InitializeGMConfig 初始化GM配置（如果不存在）
	//
	// 参数:
	//   ctx: 上下文
	//
	// 返回:
	//   error: 初始化失败时返回错误
	InitializeGMConfig(ctx context.Context) error
}

// IPWhitelistService IP白名单服务接口
//
// 定义IP白名单的业务操作
type IPWhitelistService interface {
	// GetWhitelist 获取IP白名单
	//
	// 参数:
	//   ctx: 上下文
	//
	// 返回:
	//   *domain.IPWhitelist: IP白名单对象
	//   error: 查询失败时返回错误
	GetWhitelist(ctx context.Context) (*domain.IPWhitelist, error)

	// UpdateWhitelist 更新IP白名单
	//
	// 参数:
	//   ctx: 上下文
	//   whitelist: IP白名单对象
	//
	// 返回:
	//   error: 更新失败时返回错误
	UpdateWhitelist(ctx context.Context, whitelist *domain.IPWhitelist) error

	// AddIP 添加IP到白名单
	//
	// 参数:
	//   ctx: 上下文
	//   ip: IP地址或CIDR
	//
	// 返回:
	//   error: 添加失败时返回错误
	AddIP(ctx context.Context, ip string) error

	// RemoveIP 从白名单移除IP
	//
	// 参数:
	//   ctx: 上下文
	//   ip: IP地址
	//
	// 返回:
	//   error: 移除失败时返回错误
	RemoveIP(ctx context.Context, ip string) error

	// CheckIP 检查IP是否在白名单中
	//
	// 参数:
	//   ctx: 上下文
	//   ip: IP地址
	//
	// 返回:
	//   bool: 在白名单中返回true
	//   error: 检查失败时返回错误
	CheckIP(ctx context.Context, ip string) (bool, error)

	// BatchAddIPs 批量添加IP
	//
	// 参数:
	//   ctx: 上下文
	//   ips: IP地址列表
	//
	// 返回:
	//   error: 添加失败时返回错误
	BatchAddIPs(ctx context.Context, ips []string) error

	// BatchRemoveIPs 批量移除IP
	//
	// 参数:
	//   ctx: 上下文
	//   ips: IP地址列表
	//
	// 返回:
	//   error: 移除失败时返回错误
	BatchRemoveIPs(ctx context.Context, ips []string) error

	// ClearWhitelist 清空白名单
	//
	// 参数:
	//   ctx: 上下文
	//
	// 返回:
	//   error: 清空失败时返回错误
	ClearWhitelist(ctx context.Context) error
}
