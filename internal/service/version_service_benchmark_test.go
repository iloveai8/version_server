// Package service Service层性能测试
//
// 该文件包含Service层的benchmark测试
package service

import (
	"context"
	"testing"

	"github.com/alicebob/miniredis/v2"

	redispkg "game_slots_vsn/internal/repository/redis"
	"game_slots_vsn/internal/repository"
	"game_slots_vsn/internal/domain"
)

// BenchmarkVersionService_GetVersion 性能测试：获取版本配置
func BenchmarkVersionService_GetVersion(b *testing.B) {
	// 创建miniredis
	mr := miniredis.RunT(b)
	defer mr.Close()

	// 创建Redis客户端
	redisClient := redispkg.NewRedisClientFromURL(mr.Addr())

	// 创建Repository和Service
	versionRepo := repository.NewVersionRepository(redisClient)
	versionSvc := NewVersionService(versionRepo)

	// 准备测试数据
	ctx := context.Background()
	version := domain.NewVersion("1.0.0")
	_ = versionSvc.CreateVersion(ctx, version, "dev")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = versionSvc.GetVersion(ctx, "1.0.0", "dev")
	}
}

// BenchmarkVersionService_CreateVersion 性能测试：创建版本配置
func BenchmarkVersionService_CreateVersion(b *testing.B) {
	// 创建miniredis
	mr := miniredis.RunT(b)
	defer mr.Close()

	// 创建Redis客户端
	redisClient := redispkg.NewRedisClientFromURL(mr.Addr())

	// 创建Repository和Service
	versionRepo := repository.NewVersionRepository(redisClient)
	versionSvc := NewVersionService(versionRepo)

	ctx := context.Background()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		version := domain.NewVersion("1.0.0")
		_ = versionSvc.CreateVersion(ctx, version, "dev")
		// 每次循环后删除，避免重复
		_ = versionSvc.DeleteVersion(ctx, "1.0.0", "dev")
	}
}

// BenchmarkVersionService_ListVersions 性能测试：列出版本配置
func BenchmarkVersionService_ListVersions(b *testing.B) {
	// 创建miniredis
	mr := miniredis.RunT(b)
	defer mr.Close()

	// 创建Redis客户端
	redisClient := redispkg.NewRedisClientFromURL(mr.Addr())

	// 创建Repository和Service
	versionRepo := repository.NewVersionRepository(redisClient)
	versionSvc := NewVersionService(versionRepo)

	// 准备测试数据
	ctx := context.Background()
	for i := 0; i < 10; i++ {
		version := domain.NewVersion("1.0." + string(rune('0'+i)))
		_ = versionSvc.CreateVersion(ctx, version, "dev")
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = versionSvc.ListVersions(ctx, "dev")
	}
}

// BenchmarkVersionService_UpdateVersion 性能测试：更新版本配置
func BenchmarkVersionService_UpdateVersion(b *testing.B) {
	// 创建miniredis
	mr := miniredis.RunT(b)
	defer mr.Close()

	// 创建Redis客户端
	redisClient := redispkg.NewRedisClientFromURL(mr.Addr())

	// 创建Repository和Service
	versionRepo := repository.NewVersionRepository(redisClient)
	versionSvc := NewVersionService(versionRepo)

	// 准备测试数据
	ctx := context.Background()
	version := domain.NewVersion("1.0.0")
	_ = versionSvc.CreateVersion(ctx, version, "dev")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		updatedVersion := domain.NewVersion("1.0.0")
		updatedVersion.SubServers = map[string]*domain.SubServer{
			"server1": {
				Vsn:    "1.0.0",
				SrvUrl: "http://server1.example.com",
				ResUrl: "http://server1.example.com/res",
				Type:   1,
			},
		}
		_ = versionSvc.UpdateVersion(ctx, updatedVersion, "dev")
	}
}
