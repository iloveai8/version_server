package httputil

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	// XRequestID 请求ID HTTP头
	XRequestID = "X-Request-ID"
)

// RequestID 请求ID中间件
//
// Gin中间件，为每个请求生成或提取唯一的请求ID
// 请求ID会自动设置到响应头中
//
// 行为:
// - 从请求头中读取X-Request-ID
// - 如果请求头中没有，则生成新的UUID
// - 将请求ID设置到上下文和响应头
//
// 示例:
//   router.Use(httputil.RequestID())
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头中获取请求ID
		requestID := c.GetHeader(XRequestID)

		// 如果请求头中没有请求ID，则生成新的
		if requestID == "" {
			requestID = generateRequestID()
		}

		// 设置请求ID到上下文
		c.Set(XRequestID, requestID)

		// 设置请求ID到响应头
		c.Header(XRequestID, requestID)

		c.Next()
	}
}

// GetRequestID 从上下文中获取请求ID
//
// 从Gin上下文中获取当前请求的ID
//
// 参数:
//   c: Gin上下文
//
// 返回:
//   string: 请求ID，如果未设置则返回空字符串
//
// 示例:
//   requestID := httputil.GetRequestID(c)
//   log.Printf("Request ID: %s", requestID)
func GetRequestID(c *gin.Context) string {
	if requestID, exists := c.Get(XRequestID); exists {
		if id, ok := requestID.(string); ok {
			return id
		}
	}
	return ""
}

// generateRequestID 生成请求ID
func generateRequestID() string {
	return uuid.New().String()
}
