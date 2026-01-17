package service

import (
	"context"

	"game_slots_vsn/internal/domain"
	"game_slots_vsn/internal/repository"
	"game_slots_vsn/pkg/errcode"
	"game_slots_vsn/pkg/validator"
)

// versionService 版本配置服务实现
//
// 提供版本配置的业务逻辑
type versionService struct {
	versionRepo repository.VersionRepository
}

// NewVersionService 创建版本配置服务
//
// 参数:
//   versionRepo: 版本配置仓储
//
// 返回:
//   VersionService: 版本配置服务接口
func NewVersionService(versionRepo repository.VersionRepository) VersionService {
	return &versionService{
		versionRepo: versionRepo,
	}
}

// GetVersion 获取版本配置
func (s *versionService) GetVersion(ctx context.Context, vsn, env string) (*domain.Version, error) {
	// 参数验证
	if vsn == "" {
		return nil, errcode.ErrInvalidParams.WithDetails(map[string]interface{}{
			"field": "vsn",
			"reason": "版本号不能为空",
		})
	}
	if env == "" {
		return nil, errcode.ErrInvalidEnv
	}

	// 验证版本号格式
	if err := validator.ValidateVersion(vsn); err != nil {
		return nil, errcode.ErrInvalidParams.WithDetails(map[string]interface{}{
			"field": "vsn",
			"reason": "版本号格式无效",
		})
	}

	// 从仓储获取
	version, err := s.versionRepo.GetByVersionAndEnv(ctx, vsn, env)
	if err != nil {
		return nil, err
	}

	return version, nil
}

// CreateVersion 创建新版本
func (s *versionService) CreateVersion(ctx context.Context, version *domain.Version, env string) error {
	// 参数验证
	if version == nil {
		return errcode.ErrInvalidParams.WithDetails(map[string]interface{}{
			"reason": "版本配置不能为nil",
		})
	}
	if env == "" {
		return errcode.ErrInvalidEnv
	}

	// 验证版本配置
	if err := version.Validate(); err != nil {
		return err
	}

	// 检查版本是否已存在
	exists, err := s.versionRepo.Exists(ctx, version.Vsn, env)
	if err != nil {
		return err
	}
	if exists {
		return errcode.ErrInvalidParams.WithDetails(map[string]interface{}{
			"reason":  "版本已存在",
			"version": version.Vsn,
			"env":     env,
		})
	}

	// 保存到仓储
	return s.versionRepo.Save(ctx, version, env)
}

// UpdateVersion 更新版本配置
func (s *versionService) UpdateVersion(ctx context.Context, version *domain.Version, env string) error {
	// 参数验证
	if version == nil {
		return errcode.ErrInvalidParams.WithDetails(map[string]interface{}{
			"reason": "版本配置不能为nil",
		})
	}
	if env == "" {
		return errcode.ErrInvalidEnv
	}

	// 验证版本配置
	if err := version.Validate(); err != nil {
		return err
	}

	// 检查版本是否存在
	exists, err := s.versionRepo.Exists(ctx, version.Vsn, env)
	if err != nil {
		return err
	}
	if !exists {
		return errcode.NewVersionNotFoundError(version.Vsn)
	}

	// 保存到仓储
	return s.versionRepo.Save(ctx, version, env)
}

// DeleteVersion 删除版本
func (s *versionService) DeleteVersion(ctx context.Context, vsn, env string) error {
	// 参数验证
	if vsn == "" {
		return errcode.ErrInvalidParams.WithDetails(map[string]interface{}{
			"field": "vsn",
			"reason": "版本号不能为空",
		})
	}
	if env == "" {
		return errcode.ErrInvalidEnv
	}

	// 检查版本是否存在
	exists, err := s.versionRepo.Exists(ctx, vsn, env)
	if err != nil {
		return err
	}
	if !exists {
		return errcode.NewVersionNotFoundError(vsn)
	}

	// 从仓储删除
	return s.versionRepo.Delete(ctx, vsn, env)
}

// ListVersions 获取所有版本列表
func (s *versionService) ListVersions(ctx context.Context, env string) ([]*domain.Version, error) {
	// 参数验证
	if env == "" {
		return nil, errcode.ErrInvalidEnv
	}

	// 从仓储获取列表
	versions, err := s.versionRepo.List(ctx, env)
	if err != nil {
		return nil, err
	}

	return versions, nil
}

// AddSubServer 添加子服务器到版本
func (s *versionService) AddSubServer(ctx context.Context, vsn, env, key string, sub *domain.SubServer) error {
	// 参数验证
	if vsn == "" {
		return errcode.ErrInvalidParams.WithDetails(map[string]interface{}{
			"field": "vsn",
			"reason": "版本号不能为空",
		})
	}
	if env == "" {
		return errcode.ErrInvalidEnv
	}
	if key == "" {
		return errcode.ErrInvalidParams.WithDetails(map[string]interface{}{
			"field": "key",
			"reason": "子服务器键不能为空",
		})
	}
	if sub == nil {
		return errcode.ErrInvalidParams.WithDetails(map[string]interface{}{
			"reason": "子服务器配置不能为nil",
		})
	}

	// 验证子服务器配置
	if err := sub.Validate(); err != nil {
		return err
	}

	// 获取版本
	version, err := s.versionRepo.GetByVersionAndEnv(ctx, vsn, env)
	if err != nil {
		return err
	}

	// 添加子服务器
	if err := version.AddSubServer(key, sub); err != nil {
		return err
	}

	// 保存更新
	return s.versionRepo.Save(ctx, version, env)
}

// RemoveSubServer 从版本移除子服务器
func (s *versionService) RemoveSubServer(ctx context.Context, vsn, env, key string) error {
	// 参数验证
	if vsn == "" {
		return errcode.ErrInvalidParams.WithDetails(map[string]interface{}{
			"field": "vsn",
			"reason": "版本号不能为空",
		})
	}
	if env == "" {
		return errcode.ErrInvalidEnv
	}
	if key == "" {
		return errcode.ErrInvalidParams.WithDetails(map[string]interface{}{
			"field": "key",
			"reason": "子服务器键不能为空",
		})
	}

	// 获取版本
	version, err := s.versionRepo.GetByVersionAndEnv(ctx, vsn, env)
	if err != nil {
		return err
	}

	// 检查子服务器是否存在
	if !version.HasSubServer(key) {
		return domain.ErrSubServerNotFound(key, version.GetSubServerKeys())
	}

	// 移除子服务器
	version.RemoveSubServer(key)

	// 保存更新
	return s.versionRepo.Save(ctx, version, env)
}
