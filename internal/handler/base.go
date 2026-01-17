// Package handler 提供HTTP处理层的实现
//
// 该包实现了：
// - Handler基础功能
// - 请求参数绑定和验证
// - 统一的响应格式
// - 错误处理
//
// 设计原则：
// - 薄层：Handler只负责HTTP相关逻辑，业务逻辑在Service层
// - 参数验证：在Handler层进行HTTP参数验证
// - 标准响应：使用统一的JSON响应格式
// - 错误处理：统一的错误响应和HTTP状态码
package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"game_slots_vsn/pkg/errcode"
)

// ==================== 响应辅助函数 ====================

// RespondWithSuccess 成功响应
//
// 返回标准的成功响应格式
//
// 参数:
//   c: Gin上下文
//   data: 响应数据
func RespondWithSuccess(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"data": data,
	})
}

// RespondWithCreated 创建成功响应
//
// 返回201 Created状态码
//
// 参数:
//   c: Gin上下文
//   data: 响应数据
func RespondWithCreated(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, gin.H{
		"data": data,
	})
}

// RespondWithNoContent 无内容响应
//
// 返回204 No Content状态码（用于删除操作）
//
// 参数:
//   c: Gin上下文
func RespondWithNoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

// RespondWithError 错误响应
//
// 返回标准的错误响应格式，并设置适当的HTTP状态码
//
// 参数:
//   c: Gin上下文
//   err: 错误对象
func RespondWithError(c *gin.Context, err error) {
	// 如果是AppError，使用其HTTP状态码
	if appErr, ok := errcode.AsAppError(err); ok {
		// 设置Retry-After头（如果可重试）
		if appErr.Retryable && appErr.RetryAfter > 0 {
			c.Header("Retry-After", formatIntToString(appErr.RetryAfter))
		}

		c.JSON(appErr.HTTPStatus, gin.H{
			"error": gin.H{
				"code":    appErr.Code,
				"message": appErr.Message,
				"details": appErr.Details,
			},
		})
		return
	}

	// 其他错误，返回内部错误
	c.JSON(http.StatusInternalServerError, gin.H{
		"error": gin.H{
			"code":    errcode.ErrCodeInternalError,
			"message": "服务器内部错误",
		},
	})
}

// ==================== 参数绑定辅助函数 ====================

// ShouldBindJSON JSON参数绑定
//
// 绑定JSON请求体，并在失败时返回错误响应
//
// 参数:
//   c: Gin上下文
//   obj: 目标对象
//
// 返回:
//   bool: 绑定成功返回true，失败返回false
func ShouldBindJSON(c *gin.Context, obj interface{}) bool {
	if err := c.ShouldBindJSON(obj); err != nil {
		RespondWithError(c, errcode.ErrInvalidParams.WithDetails(map[string]interface{}{
			"reason": "JSON格式错误或参数缺失",
		}))
		return false
	}
	return true
}

// ShouldBindQuery Query参数绑定
//
// 绑定Query参数，并在失败时返回错误响应
//
// 参数:
//   c: Gin上下文
//   obj: 目标对象
//
// 返回:
//   bool: 绑定成功返回true，失败返回false
func ShouldBindQuery(c *gin.Context, obj interface{}) bool {
	if err := c.ShouldBindQuery(obj); err != nil {
		RespondWithError(c, errcode.ErrInvalidParams.WithDetails(map[string]interface{}{
			"reason": "Query参数格式错误",
		}))
		return false
	}
	return true
}

// ShouldBindURI URI参数绑定
//
// 绑定URI参数，并在失败时返回错误响应
//
// 参数:
//   c: Gin上下文
//   obj: 目标对象
//
// 返回:
//   bool: 绑定成功返回true，失败返回false
func ShouldBindURI(c *gin.Context, obj interface{}) bool {
	if err := c.ShouldBindUri(obj); err != nil {
		RespondWithError(c, errcode.ErrInvalidParams.WithDetails(map[string]interface{}{
			"reason": "URI参数格式错误",
		}))
		return false
	}
	return true
}

// ==================== 请求参数结构 ====================

// GetVersionRequest 获取版本请求参数
type GetVersionRequest struct {
	Vsn string `form:"vsn" binding:"required"` // 版本号
	Env string `form:"env" binding:"required"` // 环境
}

// CreateVersionRequest 创建版本请求参数
type CreateVersionRequest struct {
	Env string `form:"env" binding:"required"` // 环境
}

// UpdateVersionRequest 更新版本请求参数
type UpdateVersionRequest struct {
	Env string `form:"env" binding:"required"` // 环境
}

// DeleteVersionRequest 删除版本请求参数
type DeleteVersionRequest struct {
	Env string `form:"env" binding:"required"` // 环境
}

// ListVersionsRequest 列出版本请求参数
type ListVersionsRequest struct {
	Env string `form:"env" binding:"required"` // 环境
}

// AddSubServerRequest 添加子服务器请求参数
type AddSubServerRequest struct {
	Env  string `form:"env" binding:"required"`  // 环境
	Key  string `form:"key" binding:"required"`  // 子服务器键
	Vsn  string `json:"vsn" binding:"required"`  // 版本号
	URL  string `json:"srvUrl" binding:"required"` // 服务器URL
	Res  string `json:"resUrl" binding:"required"` // 资源URL
	Type int    `json:"type" binding:"required"`   // 服务器类型
}

// RemoveSubServerRequest 移除子服务器请求参数
type RemoveSubServerRequest struct {
	Env string `form:"env" binding:"required"` // 环境
	Key string `form:"key" binding:"required"` // 子服务器键
}

// ToggleGMRequest 切换GM请求参数
type ToggleGMRequest struct {
	Field string `form:"field" binding:"required,oneof=gmEnable block"` // 切换字段
}

// AddIPRequest 添加IP请求参数
type AddIPRequest struct {
	IP string `json:"ip" binding:"required"` // IP地址
}

// RemoveIPRequest 移除IP请求参数
type RemoveIPRequest struct {
	IP string `json:"ip" binding:"required"` // IP地址
}

// BatchAddIPsRequest 批量添加IP请求参数
type BatchAddIPsRequest struct {
	IPs []string `json:"ips" binding:"required"` // IP列表
}

// BatchRemoveIPsRequest 批量移除IP请求参数
type BatchRemoveIPsRequest struct {
	IPs []string `json:"ips" binding:"required"` // IP列表
}

// ==================== 辅助函数 ====================

// formatIntToString 将整数转换为字符串
//
// 用于Retry-After头的设置
func formatIntToString(n int) string {
	if n < 10 {
		return string('0' + rune(n))
	}
	var result string
	for n > 0 {
		result = string('0'+rune(n%10)) + result
		n /= 10
	}
	return result
}
