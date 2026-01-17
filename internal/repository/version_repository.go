package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"game_slots_vsn/internal/domain"
	redispkg "game_slots_vsn/internal/repository/redis"
	"game_slots_vsn/pkg/errcode"
)

// versionRepository 版本配置仓储实现
//
// 使用Redis存储版本配置数据
type versionRepository struct {
	client redispkg.Client
}

// NewVersionRepository 创建版本配置仓储
//
// 参数:
//   client: Redis客户端
//
// 返回:
//   VersionRepository: 版本配置仓储接口
func NewVersionRepository(client redispkg.Client) VersionRepository {
	return &versionRepository{
		client: client,
	}
}

// GetByVersionAndEnv 根据版本号和环境获取版本配置
func (r *versionRepository) GetByVersionAndEnv(ctx context.Context, vsn, env string) (*domain.Version, error) {
	key := BuildVersionKey(vsn, env)

	// 从Redis获取JSON数据
	data, err := r.client.JSONGet(ctx, key)
	if err != nil {
		// 检查是否是"不存在"错误
		if appErr, ok := errcode.AsAppError(err); ok && appErr.Code == errcode.ErrCodeVersionNotFound {
			return nil, errcode.NewVersionNotFoundError(vsn)
		}
		return nil, err
	}

	// 反序列化JSON
	var version domain.Version
	if err := json.Unmarshal([]byte(data), &version); err != nil {
		return nil, errcode.ErrJSONError.WithMessage(fmt.Sprintf("反序列化版本配置失败: %v", err))
	}

	return &version, nil
}

// Save 保存版本配置
func (r *versionRepository) Save(ctx context.Context, version *domain.Version, env string) error {
	// 验证版本配置
	if err := version.Validate(); err != nil {
		return err
	}

	key := BuildVersionKey(version.Vsn, env)

	// 序列化为JSON
	data, err := json.Marshal(version)
	if err != nil {
		return errcode.ErrJSONError.WithMessage(fmt.Sprintf("序列化版本配置失败: %v", err))
	}

	// 保存到Redis
	if err := r.client.Set(ctx, key, data, 0); err != nil {
		return err
	}

	// 添加到版本列表
	listKey := BuildVersionListKey(env)
	if err := r.client.SAdd(ctx, listKey, version.Vsn); err != nil {
		return err
	}

	return nil
}

// Delete 删除版本配置
func (r *versionRepository) Delete(ctx context.Context, vsn, env string) error {
	key := BuildVersionKey(vsn, env)

	// 从Redis删除
	if err := r.client.Del(ctx, key); err != nil {
		return err
	}

	// 从版本列表移除
	listKey := BuildVersionListKey(env)
	if err := r.client.SRem(ctx, listKey, vsn); err != nil {
		return err
	}

	return nil
}

// Exists 检查版本是否存在
func (r *versionRepository) Exists(ctx context.Context, vsn, env string) (bool, error) {
	key := BuildVersionKey(vsn, env)
	return r.client.Exists(ctx, key)
}

// List 获取所有版本列表
func (r *versionRepository) List(ctx context.Context, env string) ([]*domain.Version, error) {
	listKey := BuildVersionListKey(env)

	// 获取所有版本号
	vsns, err := r.client.SMembers(ctx, listKey)
	if err != nil {
		return nil, err
	}

	// 获取每个版本的详细信息
	versions := make([]*domain.Version, 0, len(vsns))
	for _, vsn := range vsns {
		version, err := r.GetByVersionAndEnv(ctx, vsn, env)
		if err != nil {
			// 跳过获取失败的版本
			continue
		}
		versions = append(versions, version)
	}

	return versions, nil
}
