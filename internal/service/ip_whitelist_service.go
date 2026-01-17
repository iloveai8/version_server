package service

import (
	"context"

	"game_slots_vsn/internal/domain"
	"game_slots_vsn/internal/repository"
	"game_slots_vsn/pkg/errcode"
)

// ipWhitelistService IP白名单服务实现
//
// 提供IP白名单的业务逻辑
type ipWhitelistService struct {
	whitelistRepo repository.IPWhitelistRepository
}

// NewIPWhitelistService 创建IP白名单服务
//
// 参数:
//   whitelistRepo: IP白名单仓储
//
// 返回:
//   IPWhitelistService: IP白名单服务接口
func NewIPWhitelistService(whitelistRepo repository.IPWhitelistRepository) IPWhitelistService {
	return &ipWhitelistService{
		whitelistRepo: whitelistRepo,
	}
}

// GetWhitelist 获取IP白名单
func (s *ipWhitelistService) GetWhitelist(ctx context.Context) (*domain.IPWhitelist, error) {
	// 从仓储获取
	whitelist, err := s.whitelistRepo.Get(ctx)
	if err != nil {
		return nil, err
	}

	return whitelist, nil
}

// UpdateWhitelist 更新IP白名单
func (s *ipWhitelistService) UpdateWhitelist(ctx context.Context, whitelist *domain.IPWhitelist) error {
	// 参数验证
	if whitelist == nil {
		return errcode.ErrInvalidParams.WithDetails(map[string]interface{}{
			"reason": "IP白名单不能为nil",
		})
	}

	// 验证白名单
	if err := whitelist.Validate(); err != nil {
		return err
	}

	// 保存到仓储
	return s.whitelistRepo.Save(ctx, whitelist)
}

// AddIP 添加IP到白名单
func (s *ipWhitelistService) AddIP(ctx context.Context, ip string) error {
	// 参数验证
	if ip == "" {
		return errcode.ErrInvalidParams.WithDetails(map[string]interface{}{
			"field": "ip",
			"reason": "IP地址不能为空",
		})
	}

	// 获取当前白名单
	whitelist, err := s.whitelistRepo.Get(ctx)
	if err != nil {
		return err
	}

	// 添加IP
	if err := whitelist.AddIP(ip); err != nil {
		return err
	}

	// 保存更新
	return s.whitelistRepo.Save(ctx, whitelist)
}

// RemoveIP 从白名单移除IP
func (s *ipWhitelistService) RemoveIP(ctx context.Context, ip string) error {
	// 参数验证
	if ip == "" {
		return errcode.ErrInvalidParams.WithDetails(map[string]interface{}{
			"field": "ip",
			"reason": "IP地址不能为空",
		})
	}

	// 获取当前白名单
	whitelist, err := s.whitelistRepo.Get(ctx)
	if err != nil {
		return err
	}

	// 移除IP
	if err := whitelist.RemoveIP(ip); err != nil {
		// 如果IP不存在，不返回错误（幂等操作）
		if appErr, ok := errcode.AsAppError(err); ok && appErr.Code == errcode.ErrCodeIPNotFound {
			return nil
		}
		return err
	}

	// 保存更新
	return s.whitelistRepo.Save(ctx, whitelist)
}

// CheckIP 检查IP是否在白名单中
func (s *ipWhitelistService) CheckIP(ctx context.Context, ip string) (bool, error) {
	// 参数验证
	if ip == "" {
		return false, errcode.ErrInvalidParams.WithDetails(map[string]interface{}{
			"field": "ip",
			"reason": "IP地址不能为空",
		})
	}

	// 通过仓储检查
	return s.whitelistRepo.Contains(ctx, ip)
}

// BatchAddIPs 批量添加IP
func (s *ipWhitelistService) BatchAddIPs(ctx context.Context, ips []string) error {
	// 参数验证
	if len(ips) == 0 {
		return errcode.ErrInvalidParams.WithDetails(map[string]interface{}{
			"reason": "IP列表不能为空",
		})
	}

	// 获取当前白名单
	whitelist, err := s.whitelistRepo.Get(ctx)
	if err != nil {
		return err
	}

	// 批量添加IP
	if err := whitelist.AddIPs(ips); err != nil {
		return err
	}

	// 保存更新
	return s.whitelistRepo.Save(ctx, whitelist)
}

// BatchRemoveIPs 批量移除IP
func (s *ipWhitelistService) BatchRemoveIPs(ctx context.Context, ips []string) error {
	// 参数验证
	if len(ips) == 0 {
		return errcode.ErrInvalidParams.WithDetails(map[string]interface{}{
			"reason": "IP列表不能为空",
		})
	}

	// 获取当前白名单
	whitelist, err := s.whitelistRepo.Get(ctx)
	if err != nil {
		return err
	}

	// 批量移除IP（忽略不存在的IP）
	_ = whitelist.RemoveIPs(ips)

	// 保存更新
	return s.whitelistRepo.Save(ctx, whitelist)
}

// ClearWhitelist 清空白名单
func (s *ipWhitelistService) ClearWhitelist(ctx context.Context) error {
	// 创建空白名单
	emptyWhitelist := domain.NewIPWhitelist()

	// 保存（覆盖原有白名单）
	return s.whitelistRepo.Save(ctx, emptyWhitelist)
}
