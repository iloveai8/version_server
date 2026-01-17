// Package internal 提供应用程序的启动和运行逻辑
//
// 该包实现了：
// - 应用程序启动流程
// - 依赖注入和组件初始化
// - 优雅关闭处理
//
// 设计原则：
// - 分层初始化：按依赖顺序初始化各层组件
// - 错误处理：关键组件初始化失败时退出
// - 优雅关闭：支持SIGTERM/SIGINT信号处理
package internal

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"game_slots_vsn/internal/http"
	redispkg "game_slots_vsn/internal/repository/redis"
	"game_slots_vsn/internal/repository"
	"game_slots_vsn/internal/service"
	"game_slots_vsn/pkg/config"
	applog "game_slots_vsn/pkg/log"
)

// Run 启动服务器
//
// 初始化所有组件并启动HTTP服务器
//
// 启动流程：
// 1. 加载配置
// 2. 初始化日志
// 3. 初始化Redis客户端
// 4. 初始化Repository层
// 5. 初始化Service层
// 6. 初始化HTTP服务器
// 7. 设置优雅关闭
func Run() {
	// 1. 加载配置
	cfg := loadConfig()

	// 2. 初始化日志
	logger := initLogger(cfg)
	defer logger.Sync()

	logger.Info("服务器启动中...",
		zap.String("runMode", cfg.App.RunMode),
		zap.String("ip", cfg.Server.IP),
		zap.Int("port", cfg.Server.Port),
		zap.Strings("redisHosts", cfg.Redis.Hosts),
		zap.Int("redisDB", cfg.Redis.DB),
	)

	// 3. 初始化Redis客户端
	redisClient := initRedisClient(cfg, logger)

	// 测试Redis连接
	if err := testRedisConnection(redisClient, logger); err != nil {
		logger.Fatal("Redis连接失败", zap.Error(err))
	}

	logger.Info("Redis连接成功")

	// 4. 创建Repository层
	versionRepo := repository.NewVersionRepository(redisClient)
	gmRepo := repository.NewGMRepository(redisClient)
	whitelistRepo := repository.NewIPWhitelistRepository(redisClient)

	// 5. 创建Service层
	versionService := service.NewVersionService(versionRepo)
	gmService := service.NewGMService(gmRepo)
	whitelistService := service.NewIPWhitelistService(whitelistRepo)

	// 初始化GM配置
	ctx := context.Background()
	if err := gmService.InitializeGMConfig(ctx); err != nil {
		logger.Warn("初始化GM配置失败", zap.Error(err))
	}

	// 6. 创建HTTP服务器
	server := http.NewServer(
		versionService,
		gmService,
		whitelistService,
		redisClient,
		logger,
		cfg.Server.IP,
		cfg.Server.Port,
	)

	// 7. 设置路由
	server.SetupRoutes()

	// 8. 启动服务器（在goroutine中）
	go func() {
		if err := server.Start(); err != nil {
			logger.Error("服务器启动失败", zap.Error(err))
			os.Exit(1)
		}
	}()

	// 9. 优雅关闭
	setupGracefulShutdown(server, logger)
}

// loadConfig 加载配置
//
// 从环境变量或默认路径加载配置文件
//
// 返回:
//   *config.Config: 配置对象
func loadConfig() *config.Config {
	// 从环境变量获取环境名称
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "dev"
	}

	// 加载配置
	cfg, err := config.LoadByEnv(env)
	if err != nil {
		// 如果加载失败，使用默认配置
		fmt.Printf("配置加载失败，使用默认配置: %v\n", err)
		cfg = config.DefaultConfig()
	}

	return cfg
}

// initLogger 初始化日志
//
// 根据配置初始化日志记录器
//
// 参数:
//   cfg: 配置对象
//
// 返回:
//   *zap.Logger: 日志记录器
func initLogger(cfg *config.Config) *zap.Logger {
	// 使用控制台日志级别初始化
	if err := applog.InitLogger(cfg.Log.ConsoleLevel); err != nil {
		// 如果初始化失败，使用默认级别
		_ = applog.InitLogger("info")
	}

	// 返回zap.Logger（用于兼容现有代码）
	return zap.L()
}

// initRedisClient 初始化Redis客户端
//
// 根据配置创建Redis客户端连接
//
// 参数:
//   cfg: 配置对象
//   logger: 日志记录器
//
// 返回:
//   redispkg.Client: Redis客户端接口
func initRedisClient(cfg *config.Config, logger *zap.Logger) redispkg.Client {
	// 使用第一个Redis节点创建客户端
	// TODO: 支持Redis Sentinel/Cluster模式
	redisAddr := cfg.Redis.Hosts[0]

	logger.Info("初始化Redis客户端",
		zap.String("addr", redisAddr),
		zap.Int("db", cfg.Redis.DB),
	)

	client := redispkg.NewRedisClientFromURL(redisAddr)

	return client
}

// testRedisConnection 测试Redis连接
//
// 尝试PING Redis服务器以验证连接
//
// 参数:
//   client: Redis客户端
//   logger: 日志记录器
//
// 返回:
//   error: 连接失败时返回错误
func testRedisConnection(client redispkg.Client, logger *zap.Logger) error {
	ctx := context.Background()

	logger.Debug("测试Redis连接")

	// 尝试PING Redis
	if err := client.Ping(ctx); err != nil {
		return fmt.Errorf("Redis PING失败: %w", err)
	}

	return nil
}

// setupGracefulShutdown 设置优雅关闭
//
// 监听系统信号并执行优雅关闭
//
// 参数:
//   server: HTTP服务器
//   logger: 日志记录器
func setupGracefulShutdown(server *http.Server, logger *zap.Logger) {
	// 创建信号通道
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP)

	// 等待信号
	sig := <-quit
	logger.Info("接收到关闭信号", zap.String("signal", sig.String()))

	// 创建关闭上下文
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 优雅关闭服务器
	if err := server.Shutdown(ctx); err != nil {
		logger.Error("服务器关闭失败", zap.Error(err))
	}

	logger.Info("服务器已关闭")
}

// Stop 停止服务器
//
// 用于测试或外部停止（预留接口）
func Stop() {
	// TODO: 实现停止逻辑
	// 可以通过context cancellation实现
}
