// Package domain 提供业务领域模型的定义
//
// 该包实现了所有核心业务实体：
// - Version: 版本配置领域模型
// - GMConfig: GM配置领域模型
// - IPWhitelist: IP白名单领域模型
//
// 领域模型特点：
// - 纯净性：不依赖任何外部包（除基础设施包）
// - 自验证：每个领域对象都有Validate()方法
// - 业务逻辑封装：业务规则在领域模型内部实现
package domain

import (
	"fmt"
	"strings"

	"game_slots_vsn/pkg/errcode"
	"game_slots_vsn/pkg/validator"
)

// Version 版本配置领域模型
//
// 表示一个游戏版本的完整配置，包含版本号和所有子服务器信息
//
// 使用示例：
//   v := &domain.Version{
//       Vsn: "1.2.3",
//       SubServers: make(map[string]*domain.SubServer),
//   }
//   if err := v.Validate(); err != nil {
//       return err
//   }
type Version struct {
	Vsn        string                 `json:"vsn"`        // 版本号 X.Y.Z或X.Y.Z.N
	SubServers map[string]*SubServer  `json:"subServers"` // 子服务器信息映射
}

// SubServer 子服务器配置
//
// 表示一个子服务器的配置信息
type SubServer struct {
	Vsn    string `json:"vsn"`    // 版本号
	SrvUrl string `json:"srvUrl"` // 服务器URL
	ResUrl string `json:"resUrl"` // 资源URL
	Type   int    `json:"type"`   // 服务器类型
}

// NewVersion 创建新的版本配置
//
// 参数:
//   vsn: 版本号字符串
//
// 返回:
//   *Version: 版本配置对象
func NewVersion(vsn string) *Version {
	return &Version{
		Vsn:        vsn,
		SubServers: make(map[string]*SubServer),
	}
}

// NewSubServer 创建新的子服务器配置
//
// 参数:
//   vsn: 版本号
//   srvUrl: 服务器URL
//   resUrl: 资源URL
//   type: 服务器类型
//
// 返回:
//   *SubServer: 子服务器配置对象
func NewSubServer(vsn, srvUrl, resUrl string, serverType int) *SubServer {
	return &SubServer{
		Vsn:    vsn,
		SrvUrl: srvUrl,
		ResUrl: resUrl,
		Type:   serverType,
	}
}

// Validate 验证版本配置
//
// 验证版本号格式和所有子服务器配置的有效性
//
// 返回:
//   error: 配置无效时返回错误，格式正确返回nil
func (v *Version) Validate() error {
	// 验证版本号非空
	if strings.TrimSpace(v.Vsn) == "" {
		return errcode.ErrInvalidParams.WithDetails(map[string]interface{}{
			"field":     "vsn",
			"reason":    "版本号不能为空",
			"vsn":       v.Vsn,
			"rule":      "X.Y.Z 或 X.Y.Z.N",
			"examples":  []string{"1.0.0", "1.2.3.4"},
		})
	}

	// 验证版本号格式
	if err := validator.ValidateVersion(v.Vsn); err != nil {
		return errcode.ErrInvalidParams.WithDetails(map[string]interface{}{
			"field":     "vsn",
			"reason":    "版本号格式无效",
			"vsn":       v.Vsn,
			"rule":      "X.Y.Z 或 X.Y.Z.N",
			"examples":  []string{"1.0.0", "1.2.3.4"},
		})
	}

	// 验证子服务器配置
	for key, sub := range v.SubServers {
		if err := sub.Validate(); err != nil {
			// 如果是AppError，添加子服务器信息到详情
			if appErr, ok := errcode.AsAppError(err); ok {
				details := make(map[string]interface{})
				if appErr.Details != nil {
					for k, v := range appErr.Details {
						details[k] = v
					}
				}
				details["subServerKey"] = key
				return appErr.WithDetails(details)
			}
			// 如果不是AppError，包装为AppError
			return errcode.ErrInvalidParams.WithDetails(map[string]interface{}{
				"subServerKey": key,
				"reason":       err.Error(),
			})
		}
	}

	return nil
}

// IsValid 快速检查版本配置是否有效
//
// 返回:
//   bool: 配置有效返回true，否则返回false
func (v *Version) IsValid() bool {
	return v.Validate() == nil
}

// AddSubServer 添加子服务器
//
// 参数:
//   key: 子服务器键
//   sub: 子服务器配置
//
// 返回:
//   error: 子服务器配置无效时返回错误
func (v *Version) AddSubServer(key string, sub *SubServer) error {
	if err := sub.Validate(); err != nil {
		return fmt.Errorf("invalid subServer: %w", err)
	}

	v.SubServers[key] = sub
	return nil
}

// GetSubServer 获取子服务器
//
// 参数:
//   key: 子服务器键
//
// 返回:
//   *SubServer: 子服务器配置
//   error: 子服务器不存在时返回错误
func (v *Version) GetSubServer(key string) (*SubServer, error) {
	sub, ok := v.SubServers[key]
	if !ok {
		return nil, errcode.ErrInvalidParams.WithDetails(map[string]interface{}{
			"field":     "subServer",
			"key":       key,
			"reason":    "子服务器不存在",
			"available": v.GetSubServerKeys(),
		})
	}

	return sub, nil
}

// RemoveSubServer 移除子服务器
//
// 参数:
//   key: 子服务器键
func (v *Version) RemoveSubServer(key string) {
	delete(v.SubServers, key)
}

// GetSubServerKeys 获取所有子服务器键
//
// 返回:
//   []string: 子服务器键列表
func (v *Version) GetSubServerKeys() []string {
	keys := make([]string, 0, len(v.SubServers))
	for key := range v.SubServers {
		keys = append(keys, key)
	}
	return keys
}

// HasSubServer 检查是否存在指定子服务器
//
// 参数:
//   key: 子服务器键
//
// 返回:
//   bool: 存在返回true，否则返回false
func (v *Version) HasSubServer(key string) bool {
	_, ok := v.SubServers[key]
	return ok
}

// Validate 验证子服务器配置
//
// 验证子服务器的所有必需字段
//
// 返回:
//   error: 配置无效时返回错误，格式正确返回nil
func (s *SubServer) Validate() error {
	errors := make(map[string]string)

	// 验证版本号
	if strings.TrimSpace(s.Vsn) == "" {
		errors["vsn"] = "版本号不能为空"
	} else if err := validator.ValidateVersion(s.Vsn); err != nil {
		errors["vsn"] = "版本号格式无效"
	}

	// 验证服务器URL
	if strings.TrimSpace(s.SrvUrl) == "" {
		errors["srvUrl"] = "服务器URL不能为空"
	} else if err := validator.ValidateURL(s.SrvUrl); err != nil {
		errors["srvUrl"] = "服务器URL格式无效"
	}

	// 验证资源URL
	if strings.TrimSpace(s.ResUrl) == "" {
		errors["resUrl"] = "资源URL不能为空"
	} else if err := validator.ValidateURL(s.ResUrl); err != nil {
		errors["resUrl"] = "资源URL格式无效"
	}

	// 验证服务器类型
	if s.Type < 0 {
		errors["type"] = "服务器类型不能为负数"
	}

	if len(errors) > 0 {
		return errcode.ErrInvalidParams.WithDetails(map[string]interface{}{
			"field":  "subServer",
			"errors": errors,
		})
	}

	return nil
}
