// Package main Game Slots Version Server 主程序
//
// 该程序提供版本配置管理服务
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	redispkg "game_slots_vsn/internal/repository/redis"
	"game_slots_vsn/internal/repository"
	"game_slots_vsn/internal/service"
	httpserver "game_slots_vsn/internal/http"
)

func main() {
	// 执行主程序
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "错误: %v\n", err)
		os.Exit(1)
	}
}

// run 运行服务器
func run() error {
	// 创建日志记录器
	logger, err := zap.NewProduction()
	if err != nil {
		return fmt.Errorf("创建日志记录器失败: %w", err)
	}
	defer logger.Sync()

	logger.Info("启动 Game Slots Version Server")

	// 创建 Redis 客户端
	redisClient := redispkg.NewRedisClientFromURL("localhost:6379")
	ctx := context.Background()
	if err := redisClient.Ping(ctx); err != nil {
		logger.Warn("Redis 连接失败", zap.Error(err))
		logger.Info("将继续启动，但某些功能可能不可用")
	}

	// 创建 Repository
	versionRepo := repository.NewVersionRepository(redisClient)
	gmRepo := repository.NewGMRepository(redisClient)
	whitelistRepo := repository.NewIPWhitelistRepository(redisClient)

	// 创建 Service
	versionSvc := service.NewVersionService(versionRepo)
	gmSvc := service.NewGMService(gmRepo)
	whitelistSvc := service.NewIPWhitelistService(whitelistRepo)

	// 初始化 GM 配置
	if err := gmSvc.InitializeGMConfig(ctx); err != nil {
		logger.Warn("初始化 GM 配置失败", zap.Error(err))
	}

	// 创建 HTTP 服务器
	server := httpserver.NewServer(
		versionSvc,
		gmSvc,
		whitelistSvc,
		redisClient,
		logger,
		"0.0.0.0",
		8080,
	)

	// 设置路由
	server.SetupRoutes()

	// 启动服务器（在 goroutine 中）
	go func() {
		logger.Info("HTTP 服务器启动",
			zap.String("host", "0.0.0.0"),
			zap.Int("port", 8080),
		)

		if err := server.Start(); err != nil {
			logger.Error("HTTP 服务器启动失败", zap.Error(err))
		}
	}()

	// 等待中断信号以优雅关闭
	return waitForShutdown(server, logger)
}

// waitForShutdown 等待关闭信号并优雅关闭服务器
func waitForShutdown(server *httpserver.Server, logger *zap.Logger) error {
	// 创建信号通道
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// 等待信号
	sig := <-quit
	logger.Info("收到关闭信号", zap.String("signal", sig.String()))

	// 创建关闭上下文
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 优雅关闭服务器
	logger.Info("正在关闭服务器...")
	if err := server.Shutdown(ctx); err != nil {
		logger.Error("服务器关闭失败", zap.Error(err))
		return err
	}

	logger.Info("服务器已关闭")
	return nil
}
