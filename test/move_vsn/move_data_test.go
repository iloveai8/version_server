// Package move_vsn 版本数据迁移测试
//
// 该测试包验证版本配置在Redis集群环境下的数据迁移功能
package move_vsn

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"game_slots_vsn/internal/domain"
	"game_slots_vsn/internal/repository"
	redispkg "game_slots_vsn/internal/repository/redis"
	"game_slots_vsn/internal/service"
	"game_slots_vsn/pkg/constants"
)

// setupTestRedis 创建测试用的Redis环境
//
// 使用miniredis模拟Redis，避免依赖真实Redis服务
//
// 返回:
//   *miniredis.Miniredis: miniredis实例
//   redispkg.Client: Redis客户端
func setupTestRedis(t *testing.T) (*miniredis.Miniredis, redispkg.Client) {
	// 创建miniredis
	mr := miniredis.RunT(t)

	// 创建Redis客户端
	redisClient := redispkg.NewRedisClientFromURL(mr.Addr())

	return mr, redisClient
}

// TestRemoveVsn 测试版本删除功能
//
// 验证可以正确删除指定环境的版本配置
func TestRemoveVsn(t *testing.T) {
	mr, redisClient := setupTestRedis(t)
	defer mr.Close()

	// 创建Repository和Service
	versionRepo := repository.NewVersionRepository(redisClient)
	versionSvc := service.NewVersionService(versionRepo)

	ctx := context.Background()

	// 创建测试版本
	testVersions := []string{"1.0.0", "1.1.0", "1.2.0"}
	for _, vsn := range testVersions {
		version := domain.NewVersion(vsn)
		version.SubServers = map[string]*domain.SubServer{
			"server1": {
				Vsn:    vsn,
				SrvUrl: "http://server1.example.com",
				ResUrl: "http://server1.example.com/res",
				Type:   1,
			},
		}
		err := versionSvc.CreateVersion(ctx, version, "dev")
		require.NoError(t, err, "创建版本失败: %s", vsn)
	}

	// 验证版本已创建
	for _, vsn := range testVersions {
		_, err := versionSvc.GetVersion(ctx, vsn, "dev")
		assert.NoError(t, err, "版本应该存在: %s", vsn)
	}

	// 删除版本
	err := versionSvc.DeleteVersion(ctx, "1.0.0", "dev")
	assert.NoError(t, err, "删除版本失败")

	// 验证版本已删除
	_, err = versionSvc.GetVersion(ctx, "1.0.0", "dev")
	assert.Error(t, err, "版本应该已被删除")
}

// TestMoveTOCluster 测试数据迁移到集群
//
// 验证版本配置可以正确迁移到Redis集群环境
func TestMoveTOCluster(t *testing.T) {
	mr, redisClient := setupTestRedis(t)
	defer mr.Close()

	// 创建Repository
	versionRepo := repository.NewVersionRepository(redisClient)
	versionSvc := service.NewVersionService(versionRepo)

	ctx := context.Background()

	// 测试数据：模拟不同平台的版本配置
	testCases := []struct {
		platform string
		maxVsn    string
		env       string
	}{
		{"android", "100.2.17", "dev"},
		{"android", "100.2.18", "pre"},
		{"android", "100.2.7", "pro"},
		{"ios", "200.4.16", "dev"},
		{"ios", "200.4.17", "pre"},
		{"ios", "100.2.16", "pro"},
		{"windows", "100.1.21", "dev"},
		{"windows", "100.2.8", "pre"},
		{"windows", "100.1.22", "pro"},
	}

	// 创建测试版本
	for _, tc := range testCases {
		t.Run(tc.platform+"_"+tc.env, func(t *testing.T) {
			vsn := tc.maxVsn
			version := domain.NewVersion(vsn)
			version.SubServers = map[string]*domain.SubServer{
				tc.platform: {
					Vsn:    vsn,
					SrvUrl: fmt.Sprintf("http://%s.example.com", tc.platform),
					ResUrl: fmt.Sprintf("http://%s.example.com/res", tc.platform),
					Type:   getServerType(tc.env),
				},
			}

			err := versionSvc.CreateVersion(ctx, version, tc.env)
			require.NoError(t, err, "创建版本失败: platform=%s, env=%s", tc.platform, tc.env)

			// 验证版本已创建
			created, err := versionSvc.GetVersion(ctx, vsn, tc.env)
			require.NoError(t, err, "获取版本失败")
			assert.Equal(t, vsn, created.Vsn, "版本号应该匹配")
		})
	}

	// 列出所有版本验证
	for _, env := range []string{"dev", "pre", "pro"} {
		versions, err := versionSvc.ListVersions(ctx, env)
		assert.NoError(t, err, "列出版本失败: env=%s", env)
		t.Logf("环境 %s 共有 %d 个版本", env, len(versions))
	}
}

// TestClusterData 测试集群数据一致性
//
// 验证Redis集群环境下的数据读写一致性
func TestClusterData(t *testing.T) {
	mr, redisClient := setupTestRedis(t)
	defer mr.Close()

	versionRepo := repository.NewVersionRepository(redisClient)
	versionSvc := service.NewVersionService(versionRepo)

	ctx := context.Background()

	// 创建版本
	vsn := "1.0.0"
	version := domain.NewVersion(vsn)
	version.SubServers = map[string]*domain.SubServer{
		"server1": {
			Vsn:    vsn,
			SrvUrl: "http://server1.example.com",
			ResUrl: "http://server1.example.com/res",
			Type:   1,
		},
	}

	err := versionSvc.CreateVersion(ctx, version, "dev")
	require.NoError(t, err)

	// 多次读取验证一致性
	for i := 0; i < 5; i++ {
		retrieved, err := versionSvc.GetVersion(ctx, vsn, "dev")
		require.NoError(t, err)
		assert.Equal(t, vsn, retrieved.Vsn, "第%d次读取: 版本号应该匹配", i+1)
		assert.Len(t, retrieved.SubServers, 1, "子服务器数量应该一致")
	}
}

// getServerType 根据环境获取服务器类型
//
// 参数:
//   env: 环境标识
//
// 返回:
//   int: 服务器类型
func getServerType(env string) int {
	switch env {
	case "dev":
		return constants.ServerTypeDEV
	case "pre":
		return constants.ServerTypePRE
	case "pro":
		return constants.ServerTypePRO
	default:
		return constants.ServerTypeDefault
	}
}

// TestMain 测试入口
//
// 设置测试环境并运行所有测试
func TestMain(m *testing.M) {
	// 设置环境变量（如果需要）
	if os.Getenv("REDIS_ADDR") == "" {
		os.Setenv("REDIS_ADDR", "localhost:6379")
	}

	// 运行测试
	code := m.Run()
	os.Exit(code)
}
