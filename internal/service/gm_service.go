package service

import (
	"context"

	"game_slots_vsn/internal/domain"
	"game_slots_vsn/internal/repository"
	"game_slots_vsn/pkg/errcode"
)

// gmService GM配置服务实现
//
// 提供GM配置的业务逻辑
type gmService struct {
	gmRepo repository.GMRepository
}

// NewGMService 创建GM配置服务
//
// 参数:
//   gmRepo: GM配置仓储
//
// 返回:
//   GMService: GM配置服务接口
func NewGMService(gmRepo repository.GMRepository) GMService {
	return &gmService{
		gmRepo: gmRepo,
	}
}

// GetGMConfig 获取GM配置
func (s *gmService) GetGMConfig(ctx context.Context) (*domain.GMConfig, error) {
	// 从仓储获取
	config, err := s.gmRepo.Get(ctx)
	if err != nil {
		return nil, err
	}

	return config, nil
}

// UpdateGMConfig 更新GM配置
func (s *gmService) UpdateGMConfig(ctx context.Context, config *domain.GMConfig) error {
	// 参数验证
	if config == nil {
		return errcode.ErrInvalidParams.WithDetails(map[string]interface{}{
			"reason": "GM配置不能为nil",
		})
	}

	// 验证配置（GMConfig总是有效的，但保持一致性）
	if err := config.Validate(); err != nil {
		return err
	}

	// 保存到仓储
	return s.gmRepo.Save(ctx, config)
}

// ToggleGM 切换GM功能开关
func (s *gmService) ToggleGM(ctx context.Context) (bool, error) {
	// 获取当前配置
	config, err := s.gmRepo.Get(ctx)
	if err != nil {
		// 如果配置不存在，创建默认配置
		if appErr, ok := errcode.AsAppError(err); ok && appErr.Code == errcode.ErrCodeGMConfigNotFound {
			config = domain.NewGMConfig(false, false)
		} else {
			return false, err
		}
	}

	// 切换GM功能
	newState := config.ToggleGM()

	// 保存更新
	if err := s.gmRepo.Save(ctx, config); err != nil {
		return false, err
	}

	return newState, nil
}

// ToggleBlock 切换封锁状态
func (s *gmService) ToggleBlock(ctx context.Context) (bool, error) {
	// 获取当前配置
	config, err := s.gmRepo.Get(ctx)
	if err != nil {
		// 如果配置不存在，创建默认配置
		if appErr, ok := errcode.AsAppError(err); ok && appErr.Code == errcode.ErrCodeGMConfigNotFound {
			config = domain.NewGMConfig(false, false)
		} else {
			return false, err
		}
	}

	// 切换封锁状态
	newState := config.ToggleBlock()

	// 保存更新
	if err := s.gmRepo.Save(ctx, config); err != nil {
		return false, err
	}

	return newState, nil
}

// InitializeGMConfig 初始化GM配置（如果不存在）
func (s *gmService) InitializeGMConfig(ctx context.Context) error {
	// 检查配置是否存在
	exists, err := s.gmRepo.Exists(ctx)
	if err != nil {
		return err
	}

	// 如果已存在，不进行初始化
	if exists {
		return nil
	}

	// 创建默认配置
	defaultConfig := domain.NewGMConfig(false, false)

	// 保存到仓储
	return s.gmRepo.Save(ctx, defaultConfig)
}
