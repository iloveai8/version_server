// Package http 提供HTTP服务器的健康检查功能
//
// 该文件实现了：
// - 健康检查端点
// - 存活检查（Liveness）- 用于K8s liveness probe
// - 就绪检查（Readiness）- 用于K8s readiness probe
package http

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"game_slots_vsn/internal/repository/redis"
)

// HealthChecker 健康检查器
//
// 提供应用健康状态检查功能
type HealthChecker struct {
	redisClient redis.Client
	logger      *zap.Logger
	// 是否就绪
	ready bool
	mu    sync.RWMutex
}

// NewHealthChecker 创建健康检查器
//
// 参数:
//   redisClient: Redis客户端
//   logger: 日志记录器
//
// 返回:
//   *HealthChecker: 健康检查器实例
func NewHealthChecker(redisClient redis.Client, logger *zap.Logger) *HealthChecker {
	return &HealthChecker{
		redisClient: redisClient,
		logger:      logger,
		ready:       true, // 默认状态为就绪
	}
}

// HealthStatus 健康状态响应
type HealthStatus struct {
	Status    string            `json:"status"`    // overall: healthy/unhealthy
	Timestamp string            `json:"timestamp"` // ISO8601时间戳
	Uptime    string            `json:"uptime"`    // 运行时长
	Checks    map[string]Check  `json:"checks"`    // 各项检查结果
}

// Check 单项检查结果
type Check struct {
	Status  string `json:"status"`  // healthy/unhealthy
	Message string `json:"message,omitempty"` // 详细信息
}

// SetReady 设置就绪状态
//
// 参数:
//   ready: true表示就绪，false表示未就绪
func (h *HealthChecker) SetReady(ready bool) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.ready = ready
}

// IsReady 检查是否就绪
//
// 返回:
//   bool: true表示就绪，false表示未就绪
func (h *HealthChecker) IsReady() bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.ready
}

// startTime 应用启动时间
var startTime = time.Now()

// ==================== 健康检查端点 ====================

// RegisterHealthRoutes 注册健康检查路由
//
// 参数:
//   router: Gin路由引擎
func (h *HealthChecker) RegisterHealthRoutes(router *gin.Engine) {
	// 存活检查 - K8s Liveness Probe
	router.GET("/health/live", h.Liveness)

	// 就绪检查 - K8s Readiness Probe
	router.GET("/health/ready", h.Readiness)

	// 完整健康检查
	router.GET("/health", h.Health)

	// 简化的健康检查（兼容性）
	router.GET("/ping", h.Ping)
}

// Liveness 存活检查
//
// K8s Liveness Probe端点
// 如果应用崩溃，这个端点将无法访问，K8s会重启容器
func (h *HealthChecker) Liveness(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "alive",
	})
}

// Readiness 就绪检查
//
// K8s Readiness Probe端点
// 检查应用是否准备好接收流量
func (h *HealthChecker) Readiness(c *gin.Context) {
	if !h.IsReady() {
		h.logger.Warn("应用未就绪")
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status": "not_ready",
		})
		return
	}

	// 检查Redis连接
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := h.redisClient.Ping(ctx); err != nil {
		h.logger.Error("Redis未就绪", zap.Error(err))
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status": "not_ready",
			"reason": "redis_unavailable",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ready",
	})
}

// Health 完整健康检查
//
// 返回所有组件的健康状态
func (h *HealthChecker) Health(c *gin.Context) {
	checks := make(map[string]Check)
	overallStatus := "healthy"

	// 检查Redis连接
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := h.redisClient.Ping(ctx); err != nil {
		checks["redis"] = Check{
			Status:  "unhealthy",
			Message: "Redis连接失败: " + err.Error(),
		}
		overallStatus = "unhealthy"
	} else {
		checks["redis"] = Check{
			Status:  "healthy",
			Message: "Redis连接正常",
		}
	}

	// 检查应用就绪状态
	if !h.IsReady() {
		checks["app"] = Check{
			Status:  "unhealthy",
			Message: "应用未就绪",
		}
		overallStatus = "unhealthy"
	} else {
		checks["app"] = Check{
			Status:  "healthy",
			Message: "应用运行正常",
		}
	}

	// 计算运行时长
	uptime := time.Since(startTime)

	// 返回HTTP状态码
	httpStatus := http.StatusOK
	if overallStatus == "unhealthy" {
		httpStatus = http.StatusServiceUnavailable
	}

	c.JSON(httpStatus, HealthStatus{
		Status:    overallStatus,
		Timestamp: time.Now().Format(time.RFC3339),
		Uptime:    uptime.String(),
		Checks:    checks,
	})
}

// Ping 简单的ping端点
//
// 用于快速检查服务是否响应
func (h *HealthChecker) Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
		"time":    time.Now().Format(time.RFC3339),
	})
}
