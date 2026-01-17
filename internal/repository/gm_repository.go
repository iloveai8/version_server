package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"game_slots_vsn/internal/domain"
	redispkg "game_slots_vsn/internal/repository/redis"
	"game_slots_vsn/pkg/errcode"
)

// gmRepository GM配置仓储实现
//
// 使用Redis存储GM配置数据
type gmRepository struct {
	client redispkg.Client
}

// NewGMRepository 创建GM配置仓储
//
// 参数:
//   client: Redis客户端
//
// 返回:
//   GMRepository: GM配置仓储接口
func NewGMRepository(client redispkg.Client) GMRepository {
	return &gmRepository{
		client: client,
	}
}

// Get 获取GM配置
func (r *gmRepository) Get(ctx context.Context) (*domain.GMConfig, error) {
	key := KeyGMConfig

	// 从Redis获取JSON数据
	data, err := r.client.JSONGet(ctx, key)
	if err != nil {
		// 检查是否是"不存在"错误
		if appErr, ok := errcode.AsAppError(err); ok && appErr.Code == errcode.ErrCodeVersionNotFound {
			return nil, errcode.ErrGMConfigNotFound
		}
		return nil, err
	}

	// 反序列化JSON
	var config domain.GMConfig
	if err := json.Unmarshal([]byte(data), &config); err != nil {
		return nil, errcode.ErrJSONError.WithMessage(fmt.Sprintf("反序列化GM配置失败: %v", err))
	}

	return &config, nil
}

// Save 保存GM配置
func (r *gmRepository) Save(ctx context.Context, config *domain.GMConfig) error {
	// 验证GM配置
	if err := config.Validate(); err != nil {
		return err
	}

	key := KeyGMConfig

	// 序列化为JSON
	data, err := json.Marshal(config)
	if err != nil {
		return errcode.ErrJSONError.WithMessage(fmt.Sprintf("序列化GM配置失败: %v", err))
	}

	// 保存到Redis
	if err := r.client.Set(ctx, key, data, 0); err != nil {
		return err
	}

	return nil
}

// Exists 检查GM配置是否存在
func (r *gmRepository) Exists(ctx context.Context) (bool, error) {
	key := KeyGMConfig
	return r.client.Exists(ctx, key)
}
