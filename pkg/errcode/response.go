package errcode

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ErrorResponse 错误响应结构
//
// 用于生成符合标准的错误响应JSON
type ErrorResponse struct {
	Error *ErrorInfo `json:"error"`
}

// ErrorInfo 错误信息
//
// 包含错误码、消息和可选的详细信息
type ErrorInfo struct {
	Code    string                 `json:"code"`                  // 错误码
	Message string                 `json:"message"`               // 错误消息
	Details map[string]interface{} `json:"details,omitempty"`     // 错误详情
}

// SuccessResponse 成功响应结构
//
// 用于生成符合标准的成功响应JSON
type SuccessResponse struct {
	Data interface{} `json:"data"`  // 响应数据
}

// SetRetryHeaders 设置Retry-After头（RFC 7231标准）
//
// 如果错误可重试，则设置标准的Retry-After HTTP响应头
//
// 参数:
//   c: Gin上下文
//   err: AppError实例
//
// 示例:
//   errcode.SetRetryHeaders(c, appErr)
//   // 如果appErr.Retryable=true且RetryAfter=3，将设置头：Retry-After: 3
func SetRetryHeaders(c *gin.Context, err *AppError) {
	if err.Retryable && err.RetryAfter > 0 {
		c.Header("Retry-After", strconv.Itoa(err.RetryAfter))
	}
}

// RespondWithError 响应错误
//
// 统一的错误响应处理函数，支持AppError和普通error
//
// 参数:
//   c: Gin上下文
//   err: 错误实例（AppError或普通error）
//
// 行为:
// - 如果是AppError：返回对应的状态码和错误信息，设置Retry-After头（如可重试）
// - 如果是普通error：返回500内部错误
// - 如果err为nil：返回500内部错误（未知错误）
//
// 示例:
//   errcode.RespondWithError(c, errcode.ErrVersionNotFound)
//   errcode.RespondWithError(c, errcode.NewParamsError("版本号不能为空"))
func RespondWithError(c *gin.Context, err error) {
	if err == nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: &ErrorInfo{
				Code:    ErrCodeInternalError,
				Message: "未知错误",
			},
		})
		return
	}

	// 判断是否为AppError
	if appErr, ok := err.(*AppError); ok {
		status := appErr.HTTPStatus
		if status == 0 {
			status = http.StatusInternalServerError
		}

		// 设置Retry-After头（标准）
		SetRetryHeaders(c, appErr)

		c.JSON(status, ErrorResponse{
			Error: &ErrorInfo{
				Code:    appErr.Code,
				Message: appErr.Message,
				Details: appErr.Details,
			},
		})
		return
	}

	// 非AppError，返回内部错误
	c.JSON(http.StatusInternalServerError, ErrorResponse{
		Error: &ErrorInfo{
			Code:    ErrCodeInternalError,
			Message: "服务器内部错误",
		},
	})
}

// RespondWithSuccess 响应成功
//
// 返回HTTP 200成功响应
//
// 参数:
//   c: Gin上下文
//   data: 响应数据
//
// 示例:
//   errcode.RespondWithSuccess(c, gin.H{
//       "version": "1.0.0",
//       "url": "http://example.com",
//   })
func RespondWithSuccess(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, SuccessResponse{
		Data: data,
	})
}

// RespondWithSuccessAndStatus 响应成功（指定状态码）
//
// 返回指定HTTP状态码的成功响应
//
// 参数:
//   c: Gin上下文
//   status: HTTP状态码
//   data: 响应数据
//
// 示例:
//   errcode.RespondWithSuccessAndStatus(c, http.StatusAccepted, data)
func RespondWithSuccessAndStatus(c *gin.Context, status int, data interface{}) {
	c.JSON(status, SuccessResponse{
		Data: data,
	})
}

// RespondWithCreated 响应创建成功（201）
//
// 返回HTTP 201 Created响应，用于资源创建成功的场景
//
// 参数:
//   c: Gin上下文
//   data: 创建的资源数据
//
// 示例:
//   errcode.RespondWithCreated(c, createdVersion)
func RespondWithCreated(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, SuccessResponse{
		Data: data,
	})
}

// RespondWithNoContent 响应无内容（204）
//
// 返回HTTP 204 No Content响应，用于删除成功的场景
//
// 参数:
//   c: Gin上下文
//
// 示例:
//   errcode.RespondWithNoContent(c)
func RespondWithNoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

// RespondWithBadRequest 响应参数错误（400）
//
// 快捷方法：返回HTTP 400 Bad Request响应
//
// 参数:
//   c: Gin上下文
//   message: 错误消息
//
// 示例:
//   errcode.RespondWithBadRequest(c, "版本号不能为空")
func RespondWithBadRequest(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, ErrorResponse{
		Error: &ErrorInfo{
			Code:    ErrCodeInvalidParams,
			Message: message,
		},
	})
}

// RespondWithNotFound 响应未找到（404）
//
// 快捷方法：返回HTTP 404 Not Found响应
//
// 参数:
//   c: Gin上下文
//   message: 错误消息
//
// 示例:
//   errcode.RespondWithNotFound(c, "版本1.2.0不存在")
func RespondWithNotFound(c *gin.Context, message string) {
	c.JSON(http.StatusNotFound, ErrorResponse{
		Error: &ErrorInfo{
			Code:    ErrCodeVersionNotFound,
			Message: message,
		},
	})
}

// RespondWithInternalError 响应内部错误（500）
//
// 快捷方法：返回HTTP 500 Internal Server Error响应
//
// 参数:
//   c: Gin上下文
//   message: 错误消息
//
// 示例:
//   errcode.RespondWithInternalError(c, "服务器内部错误")
func RespondWithInternalError(c *gin.Context, message string) {
	c.JSON(http.StatusInternalServerError, ErrorResponse{
		Error: &ErrorInfo{
			Code:    ErrCodeInternalError,
			Message: message,
		},
	})
}

// RespondWithTimeout 响应超时（504）
//
// 快捷方法：返回HTTP 504 Gateway Timeout响应（可重试）
//
// 参数:
//   c: Gin上下文
//   message: 错误消息
//
// 注意:
// 该响应会自动设置Retry-After: 3头
//
// 示例:
//   errcode.RespondWithTimeout(c, "Redis查询超时")
func RespondWithTimeout(c *gin.Context, message string) {
	c.JSON(http.StatusGatewayTimeout, ErrorResponse{
		Error: &ErrorInfo{
			Code:    ErrCodeTimeout,
			Message: message,
		},
	})
}
