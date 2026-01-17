// Package repository Repository层性能测试
//
// 该文件包含Repository层的benchmark测试
package repository

import (
	"context"
	"testing"

	"github.com/alicebob/miniredis/v2"

	redispkg "game_slots_vsn/internal/repository/redis"
	"game_slots_vsn/internal/domain"
)

// BenchmarkVersionRepository_Get 性能测试：获取版本配置
func BenchmarkVersionRepository_Get(b *testing.B) {
	// 创建miniredis
	mr := miniredis.RunT(b)
	defer mr.Close()

	// 创建Redis客户端
	redisClient := redispkg.NewRedisClientFromURL(mr.Addr())

	// 创建Repository
	repo := NewVersionRepository(redisClient)

	// 准备测试数据
	ctx := context.Background()
	version := domain.NewVersion("1.0.0")
	version.SubServers = map[string]*domain.SubServer{
		"server1": {
			Vsn:    "1.0.0",
			SrvUrl: "http://server1.example.com",
			ResUrl: "http://server1.example.com/res",
			Type:   1,
		},
	}

	// 保存测试数据
	_ = repo.Save(ctx, version, "dev")

	// 重置计时器
	b.ResetTimer()

	// 执行benchmark
	for i := 0; i < b.N; i++ {
		_, _ = repo.GetByVersionAndEnv(ctx, "1.0.0", "dev")
	}
}

// BenchmarkVersionRepository_Save 性能测试：保存版本配置
func BenchmarkVersionRepository_Save(b *testing.B) {
	// 创建miniredis
	mr := miniredis.RunT(b)
	defer mr.Close()

	// 创建Redis客户端
	redisClient := redispkg.NewRedisClientFromURL(mr.Addr())

	// 创建Repository
	repo := NewVersionRepository(redisClient)

	// 准备测试数据
	ctx := context.Background()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		version := domain.NewVersion("1.0.0")
		version.SubServers = map[string]*domain.SubServer{
			"server1": {
				Vsn:    "1.0.0",
				SrvUrl: "http://server1.example.com",
				ResUrl: "http://server1.example.com/res",
				Type:   1,
			},
		}
		_ = repo.Save(ctx, version, "dev")
	}
}

// BenchmarkVersionRepository_List 性能测试：列出版本配置
func BenchmarkVersionRepository_List(b *testing.B) {
	// 创建miniredis
	mr := miniredis.RunT(b)
	defer mr.Close()

	// 创建Redis客户端
	redisClient := redispkg.NewRedisClientFromURL(mr.Addr())

	// 创建Repository
	repo := NewVersionRepository(redisClient)

	// 准备测试数据
	ctx := context.Background()
	for i := 0; i < 10; i++ {
		version := domain.NewVersion("1.0." + string(rune('0'+i)))
		_ = repo.Save(ctx, version, "dev")
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = repo.List(ctx, "dev")
	}
}

// BenchmarkVersionRepository_Exists 性能测试：检查版本是否存在
func BenchmarkVersionRepository_Exists(b *testing.B) {
	// 创建miniredis
	mr := miniredis.RunT(b)
	defer mr.Close()

	// 创建Redis客户端
	redisClient := redispkg.NewRedisClientFromURL(mr.Addr())

	// 创建Repository
	repo := NewVersionRepository(redisClient)

	// 准备测试数据
	ctx := context.Background()
	version := domain.NewVersion("1.0.0")
	_ = repo.Save(ctx, version, "dev")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = repo.Exists(ctx, "1.0.0", "dev")
	}
}
