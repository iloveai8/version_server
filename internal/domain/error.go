package domain

import (
	"fmt"

	"game_slots_vsn/pkg/errcode"
)

// DomainError 领域错误类型
//
// 包装了AppError，提供更具体的领域错误信息
// 该类型实现了error接口
type DomainError struct {
	Err     *errcode.AppError // 底层错误
	Context string            // 错误上下文（如操作类型、资源名称等）
}

// Error 实现error接口
//
// 返回格式化的错误消息，包含上下文信息
//
// 返回:
//   string: 错误消息
func (e *DomainError) Error() string {
	if e.Context != "" {
		return fmt.Sprintf("%s: %s", e.Context, e.Err.Message)
	}
	return e.Err.Message
}

// Unwrap 返回底层错误
//
// 支持errors.As和errors.Is
//
// 返回:
//   error: 底层的AppError
func (e *DomainError) Unwrap() error {
	return e.Err
}

// NewDomainError 创建领域错误
//
// 参数:
//   appErr: 底层AppError
//   context: 错误上下文（可选）
//
// 返回:
//   *DomainError: 领域错误实例
func NewDomainError(appErr *errcode.AppError, context string) *DomainError {
	return &DomainError{
		Err:     appErr,
		Context: context,
	}
}

// WrapError 包装错误为领域错误
//
// 如果已经是AppError则直接包装，否则创建内部错误
//
// 参数:
//   err: 待包装的错误
//   context: 错误上下文
//
// 返回:
//   error: 领域错误
func WrapError(err error, context string) error {
	if err == nil {
		return nil
	}

	// 如果是AppError，直接包装
	if appErr, ok := err.(*errcode.AppError); ok {
		return NewDomainError(appErr, context)
	}

	// 如果是DomainError，添加上下文
	if domainErr, ok := err.(*DomainError); ok {
		return &DomainError{
			Err:     domainErr.Err,
			Context: context + ": " + domainErr.Context,
		}
	}

	// 其他错误，转换为内部错误
	return NewDomainError(errcode.ErrInternalError.WithMessage(err.Error()), context)
}

// ==================== 领域特定错误便捷函数 ====================
//
// 以下函数提供了创建常见领域错误的便捷方法

// ErrInvalidVersion 创建无效版本号错误
//
// 参数:
//   version: 无效的版本号
//   reason: 无效原因
//
// 返回:
//   *DomainError: 领域错误
func ErrInvalidVersion(version, reason string) *DomainError {
	return NewDomainError(
		errcode.ErrInvalidParams.WithDetails(map[string]interface{}{
			"field":  "vsn",
			"reason": reason,
			"value":  version,
		}),
		"version validation",
	)
}

// ErrInvalidSubServer 创建无效子服务器错误
//
// 参数:
//   key: 子服务器键
//   reason: 无效原因
//
// 返回:
//   *DomainError: 领域错误
func ErrInvalidSubServer(key, reason string) *DomainError {
	details := map[string]interface{}{
		"field":  "subServer",
		"key":    key,
		"reason": reason,
	}
	return NewDomainError(
		errcode.ErrInvalidParams.WithDetails(details),
		"subServer validation",
	)
}

// ErrSubServerNotFound 创建子服务器不存在错误
//
// 参数:
//   key: 子服务器键
//   available: 可用的子服务器键列表
//
// 返回:
//   *DomainError: 领域错误
func ErrSubServerNotFound(key string, available []string) *DomainError {
	details := map[string]interface{}{
		"key":       key,
		"reason":    "子服务器不存在",
		"available": available,
	}
	return NewDomainError(
		errcode.ErrInvalidParams.WithDetails(details),
		"subServer lookup",
	)
}

// ErrInvalidIPAddress 创建无效IP地址错误
//
// 参数:
//   ip: 无效的IP地址
//
// 返回:
//   *DomainError: 领域错误
func ErrInvalidIPAddress(ip string) *DomainError {
	return NewDomainError(
		errcode.ErrInvalidIPFormat.WithDetails(map[string]interface{}{
			"ip":     ip,
			"reason": "IP地址格式无效",
		}),
		"IP validation",
	)
}

// ErrIPAlreadyExists 创建IP已存在错误
//
// 参数:
//   ip: 已存在的IP地址
//
// 返回:
//   *DomainError: 领域错误
func ErrIPAlreadyExists(ip string) *DomainError {
	return NewDomainError(
		errcode.ErrDuplicateIP.WithDetails(map[string]interface{}{
			"ip":     ip,
			"reason": "IP已存在于白名单中",
		}),
		"IP whitelist add",
	)
}

// ErrIPNotInWhitelist 创建IP不在白名单错误
//
// 参数:
//   ip: 不在白名单中的IP地址
//
// 返回:
//   *DomainError: 领域错误
func ErrIPNotInWhitelist(ip string) *DomainError {
	return NewDomainError(
		errcode.ErrIPNotFound.WithDetails(map[string]interface{}{
			"ip":     ip,
			"reason": "IP不在白名单中",
		}),
		"IP whitelist check",
	)
}

// ErrGMConfigNotFound 创建GM配置不存在错误
//
// 返回:
//   *DomainError: 领域错误
func ErrGMConfigNotFound() *DomainError {
	return NewDomainError(
		errcode.ErrGMConfigNotFound,
		"GM config lookup",
	)
}

// ErrVersionNotFound 创建版本不存在错误
//
// 参数:
//   version: 不存在的版本号
//
// 返回:
//   *DomainError: 领域错误
func ErrVersionNotFound(version string) *DomainError {
	return NewDomainError(
		errcode.NewVersionNotFoundError(version),
		"version lookup",
	)
}

// AsDomainError 转换为DomainError
//
// 尝试将error转换为DomainError类型
//
// 参数:
//   err: 待转换的错误
//
// 返回:
//   *DomainError: DomainError实例（如果不是DomainError则返回nil）
//   bool: 是否转换成功
//
// 示例:
//   if domainErr, ok := AsDomainError(err); ok {
//       log.Printf("Domain error context: %s", domainErr.Context)
//   }
func AsDomainError(err error) (*DomainError, bool) {
	if err == nil {
		return nil, false
	}
	if domainErr, ok := err.(*DomainError); ok {
		return domainErr, true
	}
	return nil, false
}

// IsDomainError 判断是否为DomainError
//
// 参数:
//   err: 待判断的错误
//
// 返回:
//   bool: 如果是DomainError返回true，否则返回false
func IsDomainError(err error) bool {
	_, ok := AsDomainError(err)
	return ok
}
