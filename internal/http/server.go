// Package http 提供HTTP服务器和路由配置
//
// 该包实现了：
// - HTTP服务器配置和启动
// - 路由注册
// - 中间件配置
//
// 设计原则：
// - 分离路由配置：对外API和管理API分开
// - 中间件链：日志、恢复、CORS等
// - 优雅关闭：支持服务器优雅关闭
package http

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"

	"game_slots_vsn/internal/handler/admin"
	"game_slots_vsn/internal/handler/api"
	redispkg "game_slots_vsn/internal/repository/redis"
	"game_slots_vsn/internal/service"
	"game_slots_vsn/pkg/metrics"
)

// Server HTTP服务器
//
// 封装了Gin引擎和相关配置
type Server struct {
	engine            *gin.Engine
	server            *http.Server
	versionService    service.VersionService
	gmService         service.GMService
	whitelistService  service.IPWhitelistService
	healthChecker     *HealthChecker
	logger            *zap.Logger
	host              string
	port              int
}

// NewServer 创建HTTP服务器
//
// 参数:
//   versionService: 版本配置服务
//   gmService: GM配置服务
//   whitelistService: IP白名单服务
//   redisClient: Redis客户端
//   logger: 日志记录器
//   host: 监听地址
//   port: 监听端口
//
// 返回:
//   *Server: HTTP服务器实例
func NewServer(
	versionService service.VersionService,
	gmService service.GMService,
	whitelistService service.IPWhitelistService,
	redisClient redispkg.Client,
	logger *zap.Logger,
	host string,
	port int,
) *Server {
	// 创建Gin引擎
	engine := gin.New()

	// 配置中间件
	engine.Use(recoveryMiddleware(logger))
	engine.Use(loggerMiddleware(logger))
	engine.Use(corsMiddleware())
	engine.Use(metrics.PrometheusMiddleware(logger))

	// 创建健康检查器
	healthChecker := NewHealthChecker(redisClient, logger)

	return &Server{
		engine:            engine,
		versionService:    versionService,
		gmService:         gmService,
		whitelistService:  whitelistService,
		healthChecker:     healthChecker,
		logger:            logger,
		host:              host,
		port:              port,
	}
}

// SetupRoutes 设置路由
//
// 注册所有API路由
func (s *Server) SetupRoutes() {
	// 注册健康检查路由
	s.healthChecker.RegisterHealthRoutes(s.engine)

	// 注册Prometheus metrics端点
	s.engine.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// 创建Handler
	versionHandler := api.NewVersionHandler(s.versionService)
	gmHandler := admin.NewGMHandler(s.gmService)
	whitelistHandler := admin.NewIPWhitelistHandler(s.whitelistService)

	// 对外API路由（游戏客户端）
	s.setupAPIRoutes(versionHandler)

	// 管理后台API路由（GM后台）
	s.setupAdminRoutes(gmHandler, whitelistHandler)
}

// setupAPIRoutes 设置对外API路由
func (s *Server) setupAPIRoutes(versionHandler *api.VersionHandler) {
	apiV1 := s.engine.Group("/api/v1")
	{
		// 版本配置相关
		apiV1.GET("/version", versionHandler.GetVersion)
		apiV1.POST("/version", versionHandler.CreateVersion)
		apiV1.PUT("/version/:vsn", versionHandler.UpdateVersion)
		apiV1.DELETE("/version/:vsn", versionHandler.DeleteVersion)
		apiV1.GET("/versions", versionHandler.ListVersions)

		// 子服务器相关
		apiV1.POST("/version/:vsn/subserver", versionHandler.AddSubServer)
		apiV1.DELETE("/version/:vsn/subserver", versionHandler.RemoveSubServer)
	}
}

// setupAdminRoutes 设置管理后台API路由
func (s *Server) setupAdminRoutes(gmHandler *admin.GMHandler, whitelistHandler *admin.IPWhitelistHandler) {
	adminV1 := s.engine.Group("/admin/v1")
	{
		// GM配置相关
		gm := adminV1.Group("/gm")
		{
			gm.GET("/config", gmHandler.GetGMConfig)
			gm.PUT("/config", gmHandler.UpdateGMConfig)
			gm.POST("/toggle", gmHandler.ToggleGM)
			gm.POST("/block/toggle", gmHandler.ToggleBlock)
			gm.POST("/initialize", gmHandler.InitializeGMConfig)
		}

		// IP白名单相关
		whitelist := adminV1.Group("/ip-whitelist")
		{
			whitelist.GET("", whitelistHandler.GetWhitelist)
			whitelist.PUT("", whitelistHandler.UpdateWhitelist)
			whitelist.POST("/ip", whitelistHandler.AddIP)
			whitelist.DELETE("/ip", whitelistHandler.RemoveIP)
			whitelist.POST("/check", whitelistHandler.CheckIP)
			whitelist.POST("/batch", whitelistHandler.BatchAddIPs)
			whitelist.DELETE("/batch", whitelistHandler.BatchRemoveIPs)
			whitelist.POST("/clear", whitelistHandler.ClearWhitelist)
		}
	}
}

// Start 启动HTTP服务器
//
// 开始监听HTTP请求
//
// 返回:
//   error: 启动失败时返回错误
func (s *Server) Start() error {
	// 创建HTTP服务器
	addr := s.host + ":" + httpPortToString(s.port)
	s.server = &http.Server{
		Addr:         addr,
		Handler:      s.engine,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// 记录启动信息
	s.logger.Info("HTTP服务器启动",
		zap.String("host", s.host),
		zap.Int("port", s.port),
	)

	// 启动服务器
	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}

// Shutdown 优雅关闭HTTP服务器
//
// 优雅关闭HTTP服务器，等待现有连接处理完成
//
// 参数:
//   ctx: 上下文
func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info("HTTP服务器正在关闭...")

	// 设置关闭超时时间
	shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// 优雅关闭
	if err := s.server.Shutdown(shutdownCtx); err != nil {
		return err
	}

	s.logger.Info("HTTP服务器已关闭")
	return nil
}

// GetEngine 获取Gin引擎
//
// 用于测试或高级配置
//
// 返回:
//   *gin.Engine: Gin引擎实例
func (s *Server) GetEngine() *gin.Engine {
	return s.engine
}

// ==================== 辅助函数 ====================

// httpPortToString 将端口整数转换为字符串
//
// 参数:
//   n: 端口号
//
// 返回:
//   string: 端口字符串
func httpPortToString(n int) string {
	if n == 0 {
		return "0"
	}
	var result string
	for n > 0 {
		result = string('0'+rune(n%10)) + result
		n /= 10
	}
	return result
}
