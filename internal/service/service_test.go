package service

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"game_slots_vsn/internal/domain"
	"game_slots_vsn/pkg/errcode"
)

// ==================== Mock Repository ====================

// MockVersionRepository 版本配置仓储的Mock实现
type MockVersionRepository struct {
	mock.Mock
}

func (m *MockVersionRepository) GetByVersionAndEnv(ctx context.Context, vsn, env string) (*domain.Version, error) {
	args := m.Called(ctx, vsn, env)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Version), args.Error(1)
}

func (m *MockVersionRepository) Save(ctx context.Context, version *domain.Version, env string) error {
	args := m.Called(ctx, version, env)
	return args.Error(0)
}

func (m *MockVersionRepository) Delete(ctx context.Context, vsn, env string) error {
	args := m.Called(ctx, vsn, env)
	return args.Error(0)
}

func (m *MockVersionRepository) Exists(ctx context.Context, vsn, env string) (bool, error) {
	args := m.Called(ctx, vsn, env)
	return args.Bool(0), args.Error(1)
}

func (m *MockVersionRepository) List(ctx context.Context, env string) ([]*domain.Version, error) {
	args := m.Called(ctx, env)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Version), args.Error(1)
}

// MockGMRepository GM配置仓储的Mock实现
type MockGMRepository struct {
	mock.Mock
}

func (m *MockGMRepository) Get(ctx context.Context) (*domain.GMConfig, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.GMConfig), args.Error(1)
}

func (m *MockGMRepository) Save(ctx context.Context, config *domain.GMConfig) error {
	args := m.Called(ctx, config)
	return args.Error(0)
}

func (m *MockGMRepository) Exists(ctx context.Context) (bool, error) {
	args := m.Called(ctx)
	return args.Bool(0), args.Error(1)
}

// MockIPWhitelistRepository IP白名单仓储的Mock实现
type MockIPWhitelistRepository struct {
	mock.Mock
}

func (m *MockIPWhitelistRepository) Get(ctx context.Context) (*domain.IPWhitelist, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.IPWhitelist), args.Error(1)
}

func (m *MockIPWhitelistRepository) Save(ctx context.Context, whitelist *domain.IPWhitelist) error {
	args := m.Called(ctx, whitelist)
	return args.Error(0)
}

func (m *MockIPWhitelistRepository) Exists(ctx context.Context) (bool, error) {
	args := m.Called(ctx)
	return args.Bool(0), args.Error(1)
}

func (m *MockIPWhitelistRepository) Contains(ctx context.Context, ip string) (bool, error) {
	args := m.Called(ctx, ip)
	return args.Bool(0), args.Error(1)
}

// ==================== VersionService 测试 ====================

func TestNewVersionService(t *testing.T) {
	mockRepo := new(MockVersionRepository)
	service := NewVersionService(mockRepo)

	assert.NotNil(t, service)
}

func TestVersionService_GetVersion(t *testing.T) {
	mockRepo := new(MockVersionRepository)
	service := NewVersionService(mockRepo)
	ctx := context.Background()

	t.Run("获取存在的版本", func(t *testing.T) {
		version := domain.NewVersion("1.0.0")
		mockRepo.On("GetByVersionAndEnv", ctx, "1.0.0", "dev").Return(version, nil)

		result, err := service.GetVersion(ctx, "1.0.0", "dev")
		assert.NoError(t, err)
		assert.Equal(t, version, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("获取不存在的版本", func(t *testing.T) {
		mockRepo := new(MockVersionRepository)
		service := NewVersionService(mockRepo)

		notFoundErr := errcode.NewVersionNotFoundError("999.0.0")
		mockRepo.On("GetByVersionAndEnv", ctx, "999.0.0", "dev").Return(nil, notFoundErr)

		_, err := service.GetVersion(ctx, "999.0.0", "dev")
		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("参数验证-空版本号", func(t *testing.T) {
		mockRepo := new(MockVersionRepository)
		service := NewVersionService(mockRepo)

		_, err := service.GetVersion(ctx, "", "dev")
		assert.Error(t, err)
	})

	t.Run("参数验证-空环境", func(t *testing.T) {
		mockRepo := new(MockVersionRepository)
		service := NewVersionService(mockRepo)

		_, err := service.GetVersion(ctx, "1.0.0", "")
		assert.Error(t, err)
	})
}

func TestVersionService_CreateVersion(t *testing.T) {
	mockRepo := new(MockVersionRepository)
	service := NewVersionService(mockRepo)
	ctx := context.Background()

	t.Run("创建新版本", func(t *testing.T) {
		version := domain.NewVersion("1.0.0")
		mockRepo.On("Exists", ctx, "1.0.0", "dev").Return(false, nil)
		mockRepo.On("Save", ctx, version, "dev").Return(nil)

		err := service.CreateVersion(ctx, version, "dev")
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("创建已存在的版本", func(t *testing.T) {
		mockRepo := new(MockVersionRepository)
		service := NewVersionService(mockRepo)

		version := domain.NewVersion("1.0.0")
		mockRepo.On("Exists", ctx, "1.0.0", "dev").Return(true, nil)

		err := service.CreateVersion(ctx, version, "dev")
		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("创建无效版本", func(t *testing.T) {
		mockRepo := new(MockVersionRepository)
		service := NewVersionService(mockRepo)

		invalidVersion := domain.NewVersion("invalid")
		err := service.CreateVersion(ctx, invalidVersion, "dev")
		assert.Error(t, err)
	})

	t.Run("参数验证-nil版本", func(t *testing.T) {
		mockRepo := new(MockVersionRepository)
		service := NewVersionService(mockRepo)

		err := service.CreateVersion(ctx, nil, "dev")
		assert.Error(t, err)
	})
}

func TestVersionService_UpdateVersion(t *testing.T) {
	mockRepo := new(MockVersionRepository)
	service := NewVersionService(mockRepo)
	ctx := context.Background()

	t.Run("更新存在的版本", func(t *testing.T) {
		version := domain.NewVersion("1.0.0")
		mockRepo.On("Exists", ctx, "1.0.0", "dev").Return(true, nil)
		mockRepo.On("Save", ctx, version, "dev").Return(nil)

		err := service.UpdateVersion(ctx, version, "dev")
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("更新不存在的版本", func(t *testing.T) {
		mockRepo := new(MockVersionRepository)
		service := NewVersionService(mockRepo)

		version := domain.NewVersion("999.0.0")
		mockRepo.On("Exists", ctx, "999.0.0", "dev").Return(false, nil)

		err := service.UpdateVersion(ctx, version, "dev")
		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestVersionService_DeleteVersion(t *testing.T) {
	mockRepo := new(MockVersionRepository)
	service := NewVersionService(mockRepo)
	ctx := context.Background()

	t.Run("删除存在的版本", func(t *testing.T) {
		mockRepo.On("Exists", ctx, "1.0.0", "dev").Return(true, nil)
		mockRepo.On("Delete", ctx, "1.0.0", "dev").Return(nil)

		err := service.DeleteVersion(ctx, "1.0.0", "dev")
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("删除不存在的版本", func(t *testing.T) {
		mockRepo := new(MockVersionRepository)
		service := NewVersionService(mockRepo)

		mockRepo.On("Exists", ctx, "999.0.0", "dev").Return(false, nil)

		err := service.DeleteVersion(ctx, "999.0.0", "dev")
		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestVersionService_ListVersions(t *testing.T) {
	mockRepo := new(MockVersionRepository)
	service := NewVersionService(mockRepo)
	ctx := context.Background()

	t.Run("获取版本列表", func(t *testing.T) {
		versions := []*domain.Version{
			domain.NewVersion("1.0.0"),
			domain.NewVersion("1.2.0"),
		}
		mockRepo.On("List", ctx, "dev").Return(versions, nil)

		result, err := service.ListVersions(ctx, "dev")
		assert.NoError(t, err)
		assert.Len(t, result, 2)
		mockRepo.AssertExpectations(t)
	})

	t.Run("参数验证-空环境", func(t *testing.T) {
		mockRepo := new(MockVersionRepository)
		service := NewVersionService(mockRepo)

		_, err := service.ListVersions(ctx, "")
		assert.Error(t, err)
	})
}

func TestVersionService_AddSubServer(t *testing.T) {
	mockRepo := new(MockVersionRepository)
	service := NewVersionService(mockRepo)
	ctx := context.Background()

	t.Run("添加子服务器", func(t *testing.T) {
		version := domain.NewVersion("1.0.0")
		sub := domain.NewSubServer("1.0.0", "http://example.com", "http://example.com/res", 1)

		mockRepo.On("GetByVersionAndEnv", ctx, "1.0.0", "dev").Return(version, nil)
		mockRepo.On("Save", ctx, version, "dev").Return(nil)

		err := service.AddSubServer(ctx, "1.0.0", "dev", "server1", sub)
		assert.NoError(t, err)
		assert.True(t, version.HasSubServer("server1"))
		mockRepo.AssertExpectations(t)
	})

	t.Run("添加无效子服务器", func(t *testing.T) {
		mockRepo := new(MockVersionRepository)
		service := NewVersionService(mockRepo)

		invalidSub := domain.NewSubServer("", "", "", -1)

		err := service.AddSubServer(ctx, "1.0.0", "dev", "server1", invalidSub)
		assert.Error(t, err)
	})
}

func TestVersionService_RemoveSubServer(t *testing.T) {
	mockRepo := new(MockVersionRepository)
	service := NewVersionService(mockRepo)
	ctx := context.Background()

	t.Run("移除存在的子服务器", func(t *testing.T) {
		version := domain.NewVersion("1.0.0")
		sub := domain.NewSubServer("1.0.0", "http://example.com", "http://example.com/res", 1)
		version.AddSubServer("server1", sub)

		mockRepo.On("GetByVersionAndEnv", ctx, "1.0.0", "dev").Return(version, nil)
		mockRepo.On("Save", ctx, version, "dev").Return(nil)

		err := service.RemoveSubServer(ctx, "1.0.0", "dev", "server1")
		assert.NoError(t, err)
		assert.False(t, version.HasSubServer("server1"))
		mockRepo.AssertExpectations(t)
	})

	t.Run("移除不存在的子服务器", func(t *testing.T) {
		mockRepo := new(MockVersionRepository)
		service := NewVersionService(mockRepo)

		version := domain.NewVersion("1.0.0")
		mockRepo.On("GetByVersionAndEnv", ctx, "1.0.0", "dev").Return(version, nil)

		err := service.RemoveSubServer(ctx, "1.0.0", "dev", "nonexistent")
		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

// ==================== GMService 测试 ====================

func TestNewGMService(t *testing.T) {
	mockRepo := new(MockGMRepository)
	service := NewGMService(mockRepo)

	assert.NotNil(t, service)
}

func TestGMService_GetGMConfig(t *testing.T) {
	mockRepo := new(MockGMRepository)
	service := NewGMService(mockRepo)
	ctx := context.Background()

	t.Run("获取GM配置", func(t *testing.T) {
		config := domain.NewGMConfig(true, false)
		mockRepo.On("Get", ctx).Return(config, nil)

		result, err := service.GetGMConfig(ctx)
		assert.NoError(t, err)
		assert.Equal(t, config, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("获取不存在的GM配置", func(t *testing.T) {
		mockRepo := new(MockGMRepository)
		service := NewGMService(mockRepo)

		mockRepo.On("Get", ctx).Return(nil, errcode.ErrGMConfigNotFound)

		_, err := service.GetGMConfig(ctx)
		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestGMService_UpdateGMConfig(t *testing.T) {
	mockRepo := new(MockGMRepository)
	service := NewGMService(mockRepo)
	ctx := context.Background()

	t.Run("更新GM配置", func(t *testing.T) {
		config := domain.NewGMConfig(false, true)
		mockRepo.On("Save", ctx, config).Return(nil)

		err := service.UpdateGMConfig(ctx, config)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("更新nil配置", func(t *testing.T) {
		mockRepo := new(MockGMRepository)
		service := NewGMService(mockRepo)

		err := service.UpdateGMConfig(ctx, nil)
		assert.Error(t, err)
	})
}

func TestGMService_ToggleGM(t *testing.T) {
	mockRepo := new(MockGMRepository)
	service := NewGMService(mockRepo)
	ctx := context.Background()

	t.Run("切换GM功能-从关闭到开启", func(t *testing.T) {
		config := domain.NewGMConfig(false, false)
		mockRepo.On("Get", ctx).Return(config, nil)
		mockRepo.On("Save", ctx, mock.AnythingOfType("*domain.GMConfig")).Return(nil)

		state, err := service.ToggleGM(ctx)
		assert.NoError(t, err)
		assert.True(t, state)
		mockRepo.AssertExpectations(t)
	})

	t.Run("切换GM功能-配置不存在时创建", func(t *testing.T) {
		mockRepo := new(MockGMRepository)
		service := NewGMService(mockRepo)

		mockRepo.On("Get", ctx).Return(nil, errcode.ErrGMConfigNotFound)
		mockRepo.On("Save", ctx, mock.AnythingOfType("*domain.GMConfig")).Return(nil)

		state, err := service.ToggleGM(ctx)
		assert.NoError(t, err)
		assert.True(t, state) // 默认是false，切换后变成true
		mockRepo.AssertExpectations(t)
	})
}

func TestGMService_ToggleBlock(t *testing.T) {
	mockRepo := new(MockGMRepository)
	service := NewGMService(mockRepo)
	ctx := context.Background()

	t.Run("切换封锁状态", func(t *testing.T) {
		config := domain.NewGMConfig(true, false)
		mockRepo.On("Get", ctx).Return(config, nil)
		mockRepo.On("Save", ctx, mock.AnythingOfType("*domain.GMConfig")).Return(nil)

		state, err := service.ToggleBlock(ctx)
		assert.NoError(t, err)
		assert.True(t, state)
		mockRepo.AssertExpectations(t)
	})
}

func TestGMService_InitializeGMConfig(t *testing.T) {
	mockRepo := new(MockGMRepository)
	service := NewGMService(mockRepo)
	ctx := context.Background()

	t.Run("初始化不存在的配置", func(t *testing.T) {
		mockRepo.On("Exists", ctx).Return(false, nil)
		mockRepo.On("Save", ctx, mock.AnythingOfType("*domain.GMConfig")).Return(nil)

		err := service.InitializeGMConfig(ctx)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("配置已存在-不初始化", func(t *testing.T) {
		mockRepo := new(MockGMRepository)
		service := NewGMService(mockRepo)

		mockRepo.On("Exists", ctx).Return(true, nil)

		err := service.InitializeGMConfig(ctx)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
}

// ==================== IPWhitelistService 测试 ====================

func TestNewIPWhitelistService(t *testing.T) {
	mockRepo := new(MockIPWhitelistRepository)
	service := NewIPWhitelistService(mockRepo)

	assert.NotNil(t, service)
}

func TestIPWhitelistService_GetWhitelist(t *testing.T) {
	mockRepo := new(MockIPWhitelistRepository)
	service := NewIPWhitelistService(mockRepo)
	ctx := context.Background()

	t.Run("获取IP白名单", func(t *testing.T) {
		whitelist := domain.NewIPWhitelist()
		mockRepo.On("Get", ctx).Return(whitelist, nil)

		result, err := service.GetWhitelist(ctx)
		assert.NoError(t, err)
		assert.Equal(t, whitelist, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestIPWhitelistService_UpdateWhitelist(t *testing.T) {
	mockRepo := new(MockIPWhitelistRepository)
	service := NewIPWhitelistService(mockRepo)
	ctx := context.Background()

	t.Run("更新IP白名单", func(t *testing.T) {
		whitelist := domain.NewIPWhitelist()
		whitelist.AddIP("192.168.1.1")

		mockRepo.On("Save", ctx, whitelist).Return(nil)

		err := service.UpdateWhitelist(ctx, whitelist)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("更新nil白名单", func(t *testing.T) {
		mockRepo := new(MockIPWhitelistRepository)
		service := NewIPWhitelistService(mockRepo)

		err := service.UpdateWhitelist(ctx, nil)
		assert.Error(t, err)
	})

	t.Run("更新空白名单", func(t *testing.T) {
		mockRepo := new(MockIPWhitelistRepository)
		service := NewIPWhitelistService(mockRepo)

		whitelist := domain.NewIPWhitelist()
		mockRepo.On("Save", ctx, whitelist).Return(nil)

		err := service.UpdateWhitelist(ctx, whitelist)
		assert.NoError(t, err) // 空白名单是有效的
		mockRepo.AssertExpectations(t)
	})
}

func TestIPWhitelistService_AddIP(t *testing.T) {
	mockRepo := new(MockIPWhitelistRepository)
	service := NewIPWhitelistService(mockRepo)
	ctx := context.Background()

	t.Run("添加IP", func(t *testing.T) {
		whitelist := domain.NewIPWhitelist()
		mockRepo.On("Get", ctx).Return(whitelist, nil)
		mockRepo.On("Save", ctx, mock.AnythingOfType("*domain.IPWhitelist")).Return(nil)

		err := service.AddIP(ctx, "192.168.1.1")
		assert.NoError(t, err)
		assert.Equal(t, 1, whitelist.Count())
		mockRepo.AssertExpectations(t)
	})

	t.Run("添加无效IP", func(t *testing.T) {
		mockRepo := new(MockIPWhitelistRepository)
		service := NewIPWhitelistService(mockRepo)

		whitelist := domain.NewIPWhitelist()
		mockRepo.On("Get", ctx).Return(whitelist, nil)

		err := service.AddIP(ctx, "invalid-ip")
		assert.Error(t, err)
	})

	t.Run("添加空IP", func(t *testing.T) {
		mockRepo := new(MockIPWhitelistRepository)
		service := NewIPWhitelistService(mockRepo)

		err := service.AddIP(ctx, "")
		assert.Error(t, err)
	})
}

func TestIPWhitelistService_RemoveIP(t *testing.T) {
	mockRepo := new(MockIPWhitelistRepository)
	service := NewIPWhitelistService(mockRepo)
	ctx := context.Background()

	t.Run("移除存在的IP", func(t *testing.T) {
		whitelist := domain.NewIPWhitelist()
		whitelist.AddIP("192.168.1.1")

		mockRepo.On("Get", ctx).Return(whitelist, nil)
		mockRepo.On("Save", ctx, mock.AnythingOfType("*domain.IPWhitelist")).Return(nil)

		err := service.RemoveIP(ctx, "192.168.1.1")
		assert.NoError(t, err)
		assert.Equal(t, 0, whitelist.Count())
		mockRepo.AssertExpectations(t)
	})

	t.Run("移除不存在的IP-幂等操作", func(t *testing.T) {
		mockRepo := new(MockIPWhitelistRepository)
		service := NewIPWhitelistService(mockRepo)

		whitelist := domain.NewIPWhitelist()
		mockRepo.On("Get", ctx).Return(whitelist, nil)
		// IP不存在时，不会调用Save（幂等操作）
		// mockRepo.On("Save", ctx, mock.AnythingOfType("*domain.IPWhitelist")).Return(nil)

		err := service.RemoveIP(ctx, "192.168.1.1")
		assert.NoError(t, err) // 不应该返回错误
		mockRepo.AssertExpectations(t)
	})
}

func TestIPWhitelistService_CheckIP(t *testing.T) {
	mockRepo := new(MockIPWhitelistRepository)
	service := NewIPWhitelistService(mockRepo)
	ctx := context.Background()

	t.Run("检查IP在白名单中", func(t *testing.T) {
		mockRepo.On("Contains", ctx, "192.168.1.1").Return(true, nil)

		contains, err := service.CheckIP(ctx, "192.168.1.1")
		assert.NoError(t, err)
		assert.True(t, contains)
		mockRepo.AssertExpectations(t)
	})

	t.Run("检查IP不在白名单中", func(t *testing.T) {
		mockRepo := new(MockIPWhitelistRepository)
		service := NewIPWhitelistService(mockRepo)

		mockRepo.On("Contains", ctx, "10.0.0.1").Return(false, nil)

		contains, err := service.CheckIP(ctx, "10.0.0.1")
		assert.NoError(t, err)
		assert.False(t, contains)
		mockRepo.AssertExpectations(t)
	})

	t.Run("检查空IP", func(t *testing.T) {
		mockRepo := new(MockIPWhitelistRepository)
		service := NewIPWhitelistService(mockRepo)

		_, err := service.CheckIP(ctx, "")
		assert.Error(t, err)
	})
}

func TestIPWhitelistService_BatchAddIPs(t *testing.T) {
	mockRepo := new(MockIPWhitelistRepository)
	service := NewIPWhitelistService(mockRepo)
	ctx := context.Background()

	t.Run("批量添加IP", func(t *testing.T) {
		ips := []string{"192.168.1.1", "192.168.1.2", "10.0.0.1"}
		whitelist := domain.NewIPWhitelist()

		mockRepo.On("Get", ctx).Return(whitelist, nil)
		mockRepo.On("Save", ctx, mock.AnythingOfType("*domain.IPWhitelist")).Return(nil)

		err := service.BatchAddIPs(ctx, ips)
		assert.NoError(t, err)
		assert.Equal(t, 3, whitelist.Count())
		mockRepo.AssertExpectations(t)
	})

	t.Run("批量添加空列表", func(t *testing.T) {
		mockRepo := new(MockIPWhitelistRepository)
		service := NewIPWhitelistService(mockRepo)

		err := service.BatchAddIPs(ctx, []string{})
		assert.Error(t, err)
	})
}

func TestIPWhitelistService_BatchRemoveIPs(t *testing.T) {
	mockRepo := new(MockIPWhitelistRepository)
	service := NewIPWhitelistService(mockRepo)
	ctx := context.Background()

	t.Run("批量移除IP", func(t *testing.T) {
		ips := []string{"192.168.1.1", "192.168.1.2"}
		whitelist := domain.NewIPWhitelist()
		whitelist.AddIP("192.168.1.1")
		whitelist.AddIP("192.168.1.2")
		whitelist.AddIP("10.0.0.1")

		mockRepo.On("Get", ctx).Return(whitelist, nil)
		mockRepo.On("Save", ctx, mock.AnythingOfType("*domain.IPWhitelist")).Return(nil)

		err := service.BatchRemoveIPs(ctx, ips)
		assert.NoError(t, err)
		assert.Equal(t, 1, whitelist.Count())
		mockRepo.AssertExpectations(t)
	})
}

func TestIPWhitelistService_ClearWhitelist(t *testing.T) {
	mockRepo := new(MockIPWhitelistRepository)
	service := NewIPWhitelistService(mockRepo)
	ctx := context.Background()

	t.Run("清空白名单", func(t *testing.T) {
		whitelist := domain.NewIPWhitelist()
		whitelist.AddIP("192.168.1.1")

		mockRepo.On("Save", ctx, mock.AnythingOfType("*domain.IPWhitelist")).Return(nil)

		err := service.ClearWhitelist(ctx)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
}

// ==================== 辅助测试函数 ====================

func TestVersionService_RepositoryError(t *testing.T) {
	mockRepo := new(MockVersionRepository)
	service := NewVersionService(mockRepo)
	ctx := context.Background()

	t.Run("仓储返回错误", func(t *testing.T) {
		mockRepo.On("GetByVersionAndEnv", ctx, "1.0.0", "dev").Return(nil, errors.New("repository error"))

		_, err := service.GetVersion(ctx, "1.0.0", "dev")
		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}
