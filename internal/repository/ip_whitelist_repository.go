package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"game_slots_vsn/internal/domain"
	redispkg "game_slots_vsn/internal/repository/redis"
	"game_slots_vsn/pkg/errcode"
)

// ipWhitelistRepository IP白名单仓储实现
//
// 使用Redis存储IP白名单数据
type ipWhitelistRepository struct {
	client redispkg.Client
}

// NewIPWhitelistRepository 创建IP白名单仓储
//
// 参数:
//   client: Redis客户端
//
// 返回:
//   IPWhitelistRepository: IP白名单仓储接口
func NewIPWhitelistRepository(client redispkg.Client) IPWhitelistRepository {
	return &ipWhitelistRepository{
		client: client,
	}
}

// Get 获取IP白名单
func (r *ipWhitelistRepository) Get(ctx context.Context) (*domain.IPWhitelist, error) {
	key := KeyIPWhitelist

	// 从Redis获取JSON数据
	data, err := r.client.JSONGet(ctx, key)
	if err != nil {
		// 检查是否是"不存在"错误
		if appErr, ok := errcode.AsAppError(err); ok && appErr.Code == errcode.ErrCodeVersionNotFound {
			// 返回空的白名单
			return domain.NewIPWhitelist(), nil
		}
		return nil, err
	}

	// 反序列化JSON
	var whitelist domain.IPWhitelist
	if err := json.Unmarshal([]byte(data), &whitelist); err != nil {
		return nil, errcode.ErrJSONError.WithMessage(fmt.Sprintf("反序列化IP白名单失败: %v", err))
	}

	return &whitelist, nil
}

// Save 保存IP白名单
func (r *ipWhitelistRepository) Save(ctx context.Context, whitelist *domain.IPWhitelist) error {
	// 验证IP白名单
	if err := whitelist.Validate(); err != nil {
		return err
	}

	key := KeyIPWhitelist

	// 序列化为JSON
	data, err := json.Marshal(whitelist)
	if err != nil {
		return errcode.ErrJSONError.WithMessage(fmt.Sprintf("序列化IP白名单失败: %v", err))
	}

	// 保存到Redis
	if err := r.client.Set(ctx, key, data, 0); err != nil {
		return err
	}

	return nil
}

// Exists 检查IP白名单是否存在
func (r *ipWhitelistRepository) Exists(ctx context.Context) (bool, error) {
	key := KeyIPWhitelist
	return r.client.Exists(ctx, key)
}

// Contains 检查IP是否在白名单中
func (r *ipWhitelistRepository) Contains(ctx context.Context, ip string) (bool, error) {
	// 获取白名单
	whitelist, err := r.Get(ctx)
	if err != nil {
		return false, err
	}

	// 检查IP是否在白名单中
	return whitelist.ContainsIPNet(ip), nil
}
