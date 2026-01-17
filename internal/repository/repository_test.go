package repository

import (
	"context"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"game_slots_vsn/internal/domain"
	redispkg "game_slots_vsn/internal/repository/redis"
)

// setupTestRedis 设置测试用的Redis
//
// 返回:
//   *miniredis.Miniredis: miniredis实例
//   redispkg.Client: Redis客户端
//   func(): 清理函数
func setupTestRedis(t *testing.T) (*miniredis.Miniredis, redispkg.Client, func()) {
	t.Helper()

	// 创建miniredis实例
	mr := miniredis.RunT(t)

	// 创建Redis客户端
	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	// 包装为我们的Client接口
	redisClient := redispkg.NewRedisClient(client)

	// 返回清理函数
	cleanup := func() {
		client.Close()
		mr.Close()
	}

	return mr, redisClient, cleanup
}

// ==================== VersionRepository 测试 ====================

func TestNewVersionRepository(t *testing.T) {
	_, redisClient, cleanup := setupTestRedis(t)
	defer cleanup()

	repo := NewVersionRepository(redisClient)
	assert.NotNil(t, repo)
}

func TestVersionRepository_Save(t *testing.T) {
	mr, redisClient, cleanup := setupTestRedis(t)
	defer cleanup()

	repo := NewVersionRepository(redisClient)
	ctx := context.Background()

	t.Run("保存有效版本配置", func(t *testing.T) {
		version := domain.NewVersion("1.0.0")
		version.SubServers = map[string]*domain.SubServer{
			"server1": domain.NewSubServer("1.0.0", "http://example.com", "http://example.com/res", 1),
		}

		err := repo.Save(ctx, version, "dev")
		assert.NoError(t, err)

		// 验证Redis中存在数据
		keys := mr.Keys()
		assert.Contains(t, keys, "version:dev:1.0.0")
		assert.Contains(t, keys, "version:dev:list")
	})

	t.Run("保存无效版本配置", func(t *testing.T) {
		_ = mr // 避免未使用变量警告
		invalidVersion := domain.NewVersion("invalid")

		err := repo.Save(ctx, invalidVersion, "dev")
		assert.Error(t, err)
	})

	t.Run("保存包含无效子服务器的版本", func(t *testing.T) {
		_ = mr // 避免未使用变量警告
		version := domain.NewVersion("1.5.0")
		// 添加无效的子服务器
		version.SubServers = map[string]*domain.SubServer{
			"server1": domain.NewSubServer("", "", "", -1),
		}

		err := repo.Save(ctx, version, "dev")
		assert.Error(t, err)
	})

	t.Run("更新已存在的版本", func(t *testing.T) {
		// 先保存一个版本
		version := domain.NewVersion("3.0.0")
		version.SubServers = map[string]*domain.SubServer{
			"server1": domain.NewSubServer("3.0.0", "http://old.com", "http://old.com/res", 1),
		}
		err := repo.Save(ctx, version, "pro")
		require.NoError(t, err)

		// 更新版本
		version.SubServers["server1"].SrvUrl = "http://new.com"
		err = repo.Save(ctx, version, "pro")
		assert.NoError(t, err)

		// 验证更新后的数据
		fetched, err := repo.GetByVersionAndEnv(ctx, "3.0.0", "pro")
		assert.NoError(t, err)
		assert.Equal(t, "http://new.com", fetched.SubServers["server1"].SrvUrl)
	})
}

func TestVersionRepository_GetByVersionAndEnv(t *testing.T) {
	_, redisClient, cleanup := setupTestRedis(t)
	defer cleanup()

	repo := NewVersionRepository(redisClient)
	ctx := context.Background()

	t.Run("获取存在的版本", func(t *testing.T) {
		// 先保存一个版本
		version := domain.NewVersion("1.2.3")
		version.SubServers = map[string]*domain.SubServer{
			"server1": domain.NewSubServer("1.2.3", "http://example.com", "http://example.com/res", 1),
		}
		err := repo.Save(ctx, version, "pro")
		require.NoError(t, err)

		// 获取版本
		fetched, err := repo.GetByVersionAndEnv(ctx, "1.2.3", "pro")
		assert.NoError(t, err)
		assert.Equal(t, "1.2.3", fetched.Vsn)
		assert.Len(t, fetched.SubServers, 1)
		assert.Equal(t, "http://example.com", fetched.SubServers["server1"].SrvUrl)
	})

	t.Run("获取不存在的版本", func(t *testing.T) {
		_, err := repo.GetByVersionAndEnv(ctx, "999.0.0", "dev")
		assert.Error(t, err)
	})

	t.Run("获取包含多个子服务器的版本", func(t *testing.T) {
		// 先保存一个包含多个子服务器的版本
		version := domain.NewVersion("2.0.0")
		version.SubServers = map[string]*domain.SubServer{
			"server1": domain.NewSubServer("2.0.0", "http://server1.com", "http://server1.com/res", 1),
			"server2": domain.NewSubServer("2.0.0", "http://server2.com", "http://server2.com/res", 2),
		}
		err := repo.Save(ctx, version, "dev")
		require.NoError(t, err)

		// 获取版本
		fetched, err := repo.GetByVersionAndEnv(ctx, "2.0.0", "dev")
		assert.NoError(t, err)
		assert.Len(t, fetched.SubServers, 2)
	})
}

func TestVersionRepository_Delete(t *testing.T) {
	_, redisClient, cleanup := setupTestRedis(t)
	defer cleanup()

	repo := NewVersionRepository(redisClient)
	ctx := context.Background()

	t.Run("删除存在的版本", func(t *testing.T) {
		// 先保存一个版本
		version := domain.NewVersion("1.0.0")
		err := repo.Save(ctx, version, "dev")
		require.NoError(t, err)

		// 删除版本
		err = repo.Delete(ctx, "1.0.0", "dev")
		assert.NoError(t, err)

		// 验证已删除
		exists, _ := repo.Exists(ctx, "1.0.0", "dev")
		assert.False(t, exists)
	})

	t.Run("删除不存在的版本", func(t *testing.T) {
		err := repo.Delete(ctx, "999.0.0", "dev")
		assert.NoError(t, err) // 删除不存在的版本不应该报错
	})
}

func TestVersionRepository_Exists(t *testing.T) {
	_, redisClient, cleanup := setupTestRedis(t)
	defer cleanup()

	repo := NewVersionRepository(redisClient)
	ctx := context.Background()

	t.Run("检查存在的版本", func(t *testing.T) {
		version := domain.NewVersion("1.0.0")
		err := repo.Save(ctx, version, "dev")
		require.NoError(t, err)

		exists, err := repo.Exists(ctx, "1.0.0", "dev")
		assert.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("检查不存在的版本", func(t *testing.T) {
		exists, err := repo.Exists(ctx, "999.0.0", "dev")
		assert.NoError(t, err)
		assert.False(t, exists)
	})
}

func TestVersionRepository_List(t *testing.T) {
	_, redisClient, cleanup := setupTestRedis(t)
	defer cleanup()

	repo := NewVersionRepository(redisClient)
	ctx := context.Background()

	t.Run("获取版本列表", func(t *testing.T) {
		// 保存多个版本
		version1 := domain.NewVersion("1.0.0")
		version2 := domain.NewVersion("1.2.0")
		version3 := domain.NewVersion("2.0.0")

		err := repo.Save(ctx, version1, "dev")
		require.NoError(t, err)
		err = repo.Save(ctx, version2, "dev")
		require.NoError(t, err)
		err = repo.Save(ctx, version3, "dev")
		require.NoError(t, err)

		// 获取列表
		versions, err := repo.List(ctx, "dev")
		assert.NoError(t, err)
		assert.Len(t, versions, 3)

		// 验证版本号
		vsns := make([]string, 0, len(versions))
		for _, v := range versions {
			vsns = append(vsns, v.Vsn)
		}
		assert.Contains(t, vsns, "1.0.0")
		assert.Contains(t, vsns, "1.2.0")
		assert.Contains(t, vsns, "2.0.0")
	})

	t.Run("获取空列表", func(t *testing.T) {
		versions, err := repo.List(ctx, "pro")
		assert.NoError(t, err)
		assert.Len(t, versions, 0)
	})

	t.Run("获取包含子服务器的版本列表", func(t *testing.T) {
		// 保存版本并设置子服务器
		version := domain.NewVersion("1.8.0")
		version.SubServers = map[string]*domain.SubServer{
			"server1": domain.NewSubServer("1.8.0", "http://test.com", "http://test.com/res", 1),
		}
		err := repo.Save(ctx, version, "dev")
		require.NoError(t, err)

		// 获取列表
		versions, err := repo.List(ctx, "dev")
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(versions), 1)

		// 验证子服务器也被正确加载
		for _, v := range versions {
			if v.Vsn == "1.8.0" {
				assert.Len(t, v.SubServers, 1)
			}
		}
	})
}

// ==================== GMRepository 测试 ====================

func TestNewGMRepository(t *testing.T) {
	_, redisClient, cleanup := setupTestRedis(t)
	defer cleanup()

	repo := NewGMRepository(redisClient)
	assert.NotNil(t, repo)
}

func TestGMRepository_Save(t *testing.T) {
	mr, redisClient, cleanup := setupTestRedis(t)
	defer cleanup()

	repo := NewGMRepository(redisClient)
	ctx := context.Background()

	t.Run("保存GM配置", func(t *testing.T) {
		config := domain.NewGMConfig(true, false)

		err := repo.Save(ctx, config)
		assert.NoError(t, err)

		// 验证Redis中存在数据
		keys := mr.Keys()
		assert.Contains(t, keys, "gm:config")
	})
}

func TestGMRepository_Get(t *testing.T) {
	mr, redisClient, cleanup := setupTestRedis(t)
	defer cleanup()

	repo := NewGMRepository(redisClient)
	ctx := context.Background()

	t.Run("获取存在的GM配置", func(t *testing.T) {
		// 先保存配置
		config := domain.NewGMConfig(true, false)
		err := repo.Save(ctx, config)
		require.NoError(t, err)

		// 获取配置
		fetched, err := repo.Get(ctx)
		assert.NoError(t, err)
		assert.True(t, fetched.GMEnable)
		assert.False(t, fetched.Block)
	})

	t.Run("获取不存在的GM配置", func(t *testing.T) {
		// 清空Redis
		mr.FlushAll()

		_, err := repo.Get(ctx)
		assert.Error(t, err)
	})

	t.Run("获取不同状态的GM配置", func(t *testing.T) {
		// 测试GM关闭，Block开启
		config := domain.NewGMConfig(false, true)
		err := repo.Save(ctx, config)
		require.NoError(t, err)

		fetched, err := repo.Get(ctx)
		assert.NoError(t, err)
		assert.False(t, fetched.GMEnable)
		assert.True(t, fetched.Block)
	})
}

func TestGMRepository_Exists(t *testing.T) {
	mr, redisClient, cleanup := setupTestRedis(t)
	defer cleanup()

	repo := NewGMRepository(redisClient)
	ctx := context.Background()

	t.Run("检查存在的GM配置", func(t *testing.T) {
		config := domain.NewGMConfig(true, false)
		err := repo.Save(ctx, config)
		require.NoError(t, err)

		exists, err := repo.Exists(ctx)
		assert.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("检查不存在的GM配置", func(t *testing.T) {
		mr.FlushAll()

		exists, err := repo.Exists(ctx)
		assert.NoError(t, err)
		assert.False(t, exists)
	})
}

// ==================== IPWhitelistRepository 测试 ====================

func TestNewIPWhitelistRepository(t *testing.T) {
	_, redisClient, cleanup := setupTestRedis(t)
	defer cleanup()

	repo := NewIPWhitelistRepository(redisClient)
	assert.NotNil(t, repo)
}

func TestIPWhitelistRepository_Save(t *testing.T) {
	mr, redisClient, cleanup := setupTestRedis(t)
	defer cleanup()

	repo := NewIPWhitelistRepository(redisClient)
	ctx := context.Background()

	t.Run("保存IP白名单", func(t *testing.T) {
		whitelist := domain.NewIPWhitelist()
		err := whitelist.AddIP("192.168.1.1")
		require.NoError(t, err)
		err = whitelist.AddIP("10.0.0.0/24")
		require.NoError(t, err)

		err = repo.Save(ctx, whitelist)
		assert.NoError(t, err)

		// 验证Redis中存在数据
		keys := mr.Keys()
		assert.Contains(t, keys, "ip:whitelist")
	})

	t.Run("保存空IP白名单", func(t *testing.T) {
		_ = mr // 避免未使用变量警告
		whitelist := domain.NewIPWhitelist()

		err := repo.Save(ctx, whitelist)
		assert.NoError(t, err)
	})

	t.Run("保存包含多个IP的白名单", func(t *testing.T) {
		_ = mr // 避免未使用变量警告
		whitelist := domain.NewIPWhitelist()
		ips := []string{"192.168.1.1", "192.168.1.2", "10.0.0.1", "172.16.0.0/16"}
		for _, ip := range ips {
			err := whitelist.AddIP(ip)
			require.NoError(t, err)
		}

		err := repo.Save(ctx, whitelist)
		assert.NoError(t, err)

		// 验证保存后的数据
		fetched, err := repo.Get(ctx)
		assert.NoError(t, err)
		assert.Equal(t, 4, fetched.Count())
	})

	t.Run("更新IP白名单", func(t *testing.T) {
		_ = mr // 避免未使用变量警告
		whitelist := domain.NewIPWhitelist()
		err := whitelist.AddIP("192.168.1.1")
		require.NoError(t, err)
		err = repo.Save(ctx, whitelist)
		require.NoError(t, err)

		// 更新白名单，添加新的IP
		err = whitelist.AddIP("192.168.1.2")
		require.NoError(t, err)
		err = repo.Save(ctx, whitelist)
		assert.NoError(t, err)

		// 验证更新后的数据
		fetched, err := repo.Get(ctx)
		assert.NoError(t, err)
		assert.Equal(t, 2, fetched.Count())
	})
}

func TestIPWhitelistRepository_Get(t *testing.T) {
	mr, redisClient, cleanup := setupTestRedis(t)
	defer cleanup()

	repo := NewIPWhitelistRepository(redisClient)
	ctx := context.Background()

	t.Run("获取存在的IP白名单", func(t *testing.T) {
		// 先保存白名单
		whitelist := domain.NewIPWhitelist()
		err := whitelist.AddIP("192.168.1.1")
		require.NoError(t, err)
		err = whitelist.AddIP("10.0.0.0/24")
		require.NoError(t, err)
		err = repo.Save(ctx, whitelist)
		require.NoError(t, err)

		// 获取白名单
		fetched, err := repo.Get(ctx)
		assert.NoError(t, err)
		assert.Equal(t, 2, fetched.Count())
		assert.True(t, fetched.Contains("192.168.1.1"))
		assert.True(t, fetched.Contains("10.0.0.0/24"))
	})

	t.Run("获取不存在的IP白名单", func(t *testing.T) {
		mr.FlushAll()

		// 应该返回空的白名单而不是错误
		whitelist, err := repo.Get(ctx)
		assert.NoError(t, err)
		assert.NotNil(t, whitelist)
		assert.True(t, whitelist.IsEmpty())
	})
}

func TestIPWhitelistRepository_Exists(t *testing.T) {
	mr, redisClient, cleanup := setupTestRedis(t)
	defer cleanup()

	repo := NewIPWhitelistRepository(redisClient)
	ctx := context.Background()

	t.Run("检查存在的IP白名单", func(t *testing.T) {
		whitelist := domain.NewIPWhitelist()
		err := whitelist.AddIP("192.168.1.1")
		require.NoError(t, err)
		err = repo.Save(ctx, whitelist)
		require.NoError(t, err)

		exists, err := repo.Exists(ctx)
		assert.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("检查不存在的IP白名单", func(t *testing.T) {
		mr.FlushAll()

		exists, err := repo.Exists(ctx)
		assert.NoError(t, err)
		assert.False(t, exists)
	})
}

func TestIPWhitelistRepository_Contains(t *testing.T) {
	mr, redisClient, cleanup := setupTestRedis(t)
	defer cleanup()

	repo := NewIPWhitelistRepository(redisClient)
	ctx := context.Background()

	t.Run("检查IP在白名单中", func(t *testing.T) {
		// 先保存白名单
		whitelist := domain.NewIPWhitelist()
		err := whitelist.AddIP("192.168.1.1")
		require.NoError(t, err)
		err = whitelist.AddIP("10.0.0.0/24")
		require.NoError(t, err)
		err = repo.Save(ctx, whitelist)
		require.NoError(t, err)

		// 检查精确匹配
		contains, err := repo.Contains(ctx, "192.168.1.1")
		assert.NoError(t, err)
		assert.True(t, contains)

		// 检查CIDR匹配
		contains, err = repo.Contains(ctx, "10.0.0.1")
		assert.NoError(t, err)
		assert.True(t, contains)

		// 检查不在白名单中的IP
		contains, err = repo.Contains(ctx, "172.16.0.1")
		assert.NoError(t, err)
		assert.False(t, contains)
	})

	t.Run("检查空白名单", func(t *testing.T) {
		mr.FlushAll()

		contains, err := repo.Contains(ctx, "192.168.1.1")
		assert.NoError(t, err)
		assert.False(t, contains)
	})

	t.Run("检查多个CIDR网段", func(t *testing.T) {
		_ = mr // 避免未使用变量警告
		whitelist := domain.NewIPWhitelist()
		err := whitelist.AddIP("10.0.0.0/24")
		require.NoError(t, err)
		err = whitelist.AddIP("172.16.0.0/16")
		require.NoError(t, err)
		err = whitelist.AddIP("192.168.1.0/24")
		require.NoError(t, err)
		err = repo.Save(ctx, whitelist)
		require.NoError(t, err)

		// 测试各个网段
		testCases := []struct {
			ip       string
			expected bool
		}{
			{"10.0.0.1", true},
			{"10.0.0.255", true},
			{"172.16.0.1", true},
			{"172.16.255.255", true},
			{"192.168.1.1", true},
			{"192.168.1.255", true},
			{"8.8.8.8", false},
			{"172.15.0.1", false},
		}

		for _, tc := range testCases {
			contains, err := repo.Contains(ctx, tc.ip)
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, contains, "IP %s should have contains=%v", tc.ip, tc.expected)
		}
	})
}

// ==================== 辅助函数测试 ====================

func TestBuildVersionKey(t *testing.T) {
	tests := []struct {
		name     string
		vsn      string
		env      string
		expected string
	}{
		{"dev环境1.0.0版本", "1.0.0", "dev", "version:dev:1.0.0"},
		{"pro环境2.3.4版本", "2.3.4", "pro", "version:pro:2.3.4"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := BuildVersionKey(tt.vsn, tt.env)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestBuildVersionListKey(t *testing.T) {
	tests := []struct {
		name     string
		env      string
		expected string
	}{
		{"dev环境", "dev", "version:dev:list"},
		{"pro环境", "pro", "version:pro:list"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := BuildVersionListKey(tt.env)
			assert.Equal(t, tt.expected, result)
		})
	}
}
