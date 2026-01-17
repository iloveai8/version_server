// Package metrics 提供Prometheus监控中间件
package metrics

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// PrometheusMiddleware Prometheus监控中间件
//
// 记录HTTP请求指标到Prometheus
//
// 参数:
//   logger: 日志记录器
//
// 返回:
//   gin.HandlerFunc: Gin中间件函数
func PrometheusMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取请求信息
		method := c.Request.Method
		path := c.FullPath()

		// 如果path为空（例如404），使用实际路径
		if path == "" {
			path = c.Request.URL.Path
		}

		// 增加进行中的请求数
		IncHTTPRequestsInProgress(method, path)
		defer DecHTTPRequestsInProgress(method, path)

		// 记录开始时间
		start := time.Now()

		// 处理请求
		c.Next()

		// 记录请求完成后的指标
		duration := time.Since(start)
		statusCode := c.Writer.Status()
		responseSize := c.Writer.Size()

		// 记录指标
		RecordHTTPRequest(method, path, statusCode, duration, responseSize)

		// 如果发生错误，记录日志
		if statusCode >= 400 {
			logger.Debug("HTTP请求错误",
				zap.String("method", method),
				zap.String("path", path),
				zap.Int("status", statusCode),
				zap.Duration("duration", duration),
			)
		}
	}
}

// NormalizePath 标准化路径
//
// 将路径中的变量（如ID）替换为占位符，以便Prometheus聚合
//
// 参数:
//   path: 原始路径
//
// 返回:
//   string: 标准化后的路径
//
// 示例:
//   NormalizePath("/api/v1/version/1.0.0") => "/api/v1/version/:vsn"
func NormalizePath(path string) string {
	// 简单的路径标准化逻辑
	// 可以根据实际需求扩展

	// 版本号路径
	if len(path) > 20 {
		// 检查是否包含版本号模式
		parts := splitPath(path)
		for i, part := range parts {
			// 检查是否为版本号格式（X.Y.Z）
			if isVersionString(part) {
				parts[i] = ":vsn"
			}
			// 检查是否为数字ID
			if isNumeric(part) {
				parts[i] = ":id"
			}
		}
		return joinPath(parts)
	}

	return path
}

// isVersionString 检查字符串是否为版本号格式
//
// 参数:
//   s: 字符串
//
// 返回:
//   bool: 是版本号返回true
func isVersionString(s string) bool {
	dotCount := 0
	for _, c := range s {
		if c == '.' {
			dotCount++
		} else if c < '0' || c > '9' {
			return false
		}
	}
	return dotCount >= 2 && dotCount <= 3
}

// isNumeric 检查字符串是否为纯数字
//
// 参数:
//   s: 字符串
//
// 返回:
//   bool: 是数字返回true
func isNumeric(s string) bool {
	if s == "" {
		return false
	}
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

// splitPath 分割路径
//
// 参数:
//   path: 路径
//
// 返回:
//   []string: 路径片段
func splitPath(path string) []string {
	if path == "" || path == "/" {
		return []string{}
	}

	parts := make([]string, 0)
	start := 0
	for i := 1; i < len(path); i++ {
		if path[i] == '/' {
			parts = append(parts, path[start+1:i])
			start = i
		}
	}
	if start < len(path)-1 {
		parts = append(parts, path[start+1:])
	}

	return parts
}

// joinPath 连接路径
//
// 参数:
//   parts: 路径片段
//
// 返回:
//   string: 完整路径
func joinPath(parts []string) string {
	if len(parts) == 0 {
		return "/"
	}

	result := "/"
	for i, part := range parts {
		if i > 0 {
			result += "/"
		}
		result += part
	}
	return result
}

// GetStatusCodeLabel 获取状态码标签
//
// 将状态码转换为Prometheus友好的标签
//
// 参数:
//   code: HTTP状态码
//
// 返回:
//   string: 状态码标签
func GetStatusCodeLabel(code int) string {
	if code >= 200 && code < 300 {
		return "2xx"
	} else if code >= 300 && code < 400 {
		return "3xx"
	} else if code >= 400 && code < 500 {
		return "4xx"
	} else if code >= 500 {
		return "5xx"
	}
	return strconv.Itoa(code)
}
