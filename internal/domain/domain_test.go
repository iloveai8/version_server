package domain

import (
	"testing"

	"game_slots_vsn/pkg/errcode"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ==================== Version 测试 ====================

func TestNewVersion(t *testing.T) {
	v := NewVersion("1.2.3")

	assert.Equal(t, "1.2.3", v.Vsn)
	assert.NotNil(t, v.SubServers)
	assert.Empty(t, v.SubServers)
}

func TestVersion_Validate(t *testing.T) {
	tests := []struct {
		name    string
		version *Version
		wantErr bool
		errCode string
	}{
		{
			name: "有效版本号",
			version: &Version{
				Vsn:        "1.2.3",
				SubServers: map[string]*SubServer{},
			},
			wantErr: false,
		},
		{
			name: "有效版本号带Build",
			version: &Version{
				Vsn:        "1.2.3.4",
				SubServers: map[string]*SubServer{},
			},
			wantErr: false,
		},
		{
			name: "空版本号",
			version: &Version{
				Vsn:        "",
				SubServers: map[string]*SubServer{},
			},
			wantErr: true,
			errCode: "INVALID_PARAMS",
		},
		{
			name: "空格版本号",
			version: &Version{
				Vsn:        "   ",
				SubServers: map[string]*SubServer{},
			},
			wantErr: true,
			errCode: "INVALID_PARAMS",
		},
		{
			name: "无效版本号格式",
			version: &Version{
				Vsn:        "invalid",
				SubServers: map[string]*SubServer{},
			},
			wantErr: true,
			errCode: "INVALID_PARAMS",
		},
		{
			name: "版本号含子服务器-有效",
			version: &Version{
				Vsn: "1.0.0",
				SubServers: map[string]*SubServer{
					"server1": {
						Vsn:    "1.0.0",
						SrvUrl: "http://example.com",
						ResUrl: "http://example.com/res",
						Type:   1,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "版本号含子服务器-子服务器无效",
			version: &Version{
				Vsn: "1.0.0",
				SubServers: map[string]*SubServer{
					"server1": {
						Vsn:    "",
						SrvUrl: "invalid-url",
						ResUrl: "",
						Type:   -1,
					},
				},
			},
			wantErr: true,
			errCode: "INVALID_PARAMS",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.version.Validate()
			if tt.wantErr {
				require.Error(t, err)
				if tt.errCode != "" {
					appErr, ok := errcode.AsAppError(err)
					require.True(t, ok)
					assert.Equal(t, tt.errCode, appErr.Code)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestVersion_IsValid(t *testing.T) {
	validVersion := NewVersion("1.2.3")
	assert.True(t, validVersion.IsValid())

	invalidVersion := NewVersion("invalid")
	assert.False(t, invalidVersion.IsValid())
}

func TestVersion_AddSubServer(t *testing.T) {
	v := NewVersion("1.0.0")

	t.Run("添加有效子服务器", func(t *testing.T) {
		sub := NewSubServer("1.0.0", "http://example.com", "http://example.com/res", 1)
		err := v.AddSubServer("server1", sub)
		assert.NoError(t, err)
		assert.True(t, v.HasSubServer("server1"))
	})

	t.Run("添加无效子服务器", func(t *testing.T) {
		invalidSub := NewSubServer("", "", "", -1)
		err := v.AddSubServer("server2", invalidSub)
		assert.Error(t, err)
		assert.False(t, v.HasSubServer("server2"))
	})
}

func TestVersion_GetSubServer(t *testing.T) {
	v := NewVersion("1.0.0")
	sub := NewSubServer("1.0.0", "http://example.com", "http://example.com/res", 1)
	v.SubServers["server1"] = sub

	t.Run("获取存在的子服务器", func(t *testing.T) {
		got, err := v.GetSubServer("server1")
		assert.NoError(t, err)
		assert.Equal(t, sub, got)
	})

	t.Run("获取不存在的子服务器", func(t *testing.T) {
		_, err := v.GetSubServer("nonexistent")
		require.Error(t, err)
		appErr, ok := errcode.AsAppError(err)
		require.True(t, ok)
		assert.Equal(t, "INVALID_PARAMS", appErr.Code)
	})
}

func TestVersion_RemoveSubServer(t *testing.T) {
	v := NewVersion("1.0.0")
	sub := NewSubServer("1.0.0", "http://example.com", "http://example.com/res", 1)
	v.SubServers["server1"] = sub

	t.Run("移除存在的子服务器", func(t *testing.T) {
		v.RemoveSubServer("server1")
		assert.False(t, v.HasSubServer("server1"))
	})

	t.Run("移除不存在的子服务器不报错", func(t *testing.T) {
		v.RemoveSubServer("nonexistent")
		assert.NotPanics(t, func() {
			v.RemoveSubServer("nonexistent")
		})
	})
}

func TestVersion_GetSubServerKeys(t *testing.T) {
	v := NewVersion("1.0.0")
	v.SubServers["server1"] = &SubServer{}
	v.SubServers["server2"] = &SubServer{}

	keys := v.GetSubServerKeys()
	assert.Len(t, keys, 2)
	assert.Contains(t, keys, "server1")
	assert.Contains(t, keys, "server2")
}

func TestVersion_HasSubServer(t *testing.T) {
	v := NewVersion("1.0.0")
	v.SubServers["server1"] = &SubServer{}

	assert.True(t, v.HasSubServer("server1"))
	assert.False(t, v.HasSubServer("server2"))
}

// ==================== SubServer 测试 ====================

func TestNewSubServer(t *testing.T) {
	sub := NewSubServer("1.0.0", "http://example.com", "http://example.com/res", 1)

	assert.Equal(t, "1.0.0", sub.Vsn)
	assert.Equal(t, "http://example.com", sub.SrvUrl)
	assert.Equal(t, "http://example.com/res", sub.ResUrl)
	assert.Equal(t, 1, sub.Type)
}

func TestSubServer_Validate(t *testing.T) {
	tests := []struct {
		name    string
		sub     *SubServer
		wantErr bool
	}{
		{
			name: "有效子服务器",
			sub: NewSubServer(
				"1.0.0",
				"http://example.com",
				"http://example.com/res",
				1,
			),
			wantErr: false,
		},
		{
			name: "空版本号",
			sub: NewSubServer(
				"",
				"http://example.com",
				"http://example.com/res",
				1,
			),
			wantErr: true,
		},
		{
			name: "无效版本号",
			sub: NewSubServer(
				"invalid",
				"http://example.com",
				"http://example.com/res",
				1,
			),
			wantErr: true,
		},
		{
			name: "空服务器URL",
			sub: NewSubServer(
				"1.0.0",
				"",
				"http://example.com/res",
				1,
			),
			wantErr: true,
		},
		{
			name: "无效服务器URL",
			sub: NewSubServer(
				"1.0.0",
				"invalid-url",
				"http://example.com/res",
				1,
			),
			wantErr: true,
		},
		{
			name: "空资源URL",
			sub: NewSubServer(
				"1.0.0",
				"http://example.com",
				"",
				1,
			),
			wantErr: true,
		},
		{
			name: "负数类型",
			sub: NewSubServer(
				"1.0.0",
				"http://example.com",
				"http://example.com/res",
				-1,
			),
			wantErr: true,
		},
		{
			name: "零类型-有效",
			sub: NewSubServer(
				"1.0.0",
				"http://example.com",
				"http://example.com/res",
				0,
			),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.sub.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// ==================== GMConfig 测试 ====================

func TestNewGMConfig(t *testing.T) {
	config := NewGMConfig(true, false)

	assert.True(t, config.GMEnable)
	assert.False(t, config.Block)
}

func TestGMConfig_Validate(t *testing.T) {
	config := NewGMConfig(true, false)

	// GM配置总是有效的
	assert.NoError(t, config.Validate())
}

func TestGMConfig_IsValid(t *testing.T) {
	config := NewGMConfig(true, false)
	assert.True(t, config.IsValid())
}

func TestGMConfig_ToggleGM(t *testing.T) {
	config := NewGMConfig(true, false)

	result := config.ToggleGM()
	assert.False(t, result)
	assert.False(t, config.GMEnable)

	result = config.ToggleGM()
	assert.True(t, result)
	assert.True(t, config.GMEnable)
}

func TestGMConfig_ToggleBlock(t *testing.T) {
	config := NewGMConfig(true, false)

	result := config.ToggleBlock()
	assert.True(t, result)
	assert.True(t, config.Block)

	result = config.ToggleBlock()
	assert.False(t, result)
	assert.False(t, config.Block)
}

func TestGMConfig_SetGMEnable(t *testing.T) {
	config := NewGMConfig(true, false)

	config.SetGMEnable(false)
	assert.False(t, config.GMEnable)

	config.SetGMEnable(true)
	assert.True(t, config.GMEnable)
}

func TestGMConfig_SetBlock(t *testing.T) {
	config := NewGMConfig(true, false)

	config.SetBlock(true)
	assert.True(t, config.Block)

	config.SetBlock(false)
	assert.False(t, config.Block)
}

func TestGMConfig_IsGMEnabled(t *testing.T) {
	config := NewGMConfig(true, false)
	assert.True(t, config.IsGMEnabled())

	config.GMEnable = false
	assert.False(t, config.IsGMEnabled())
}

func TestGMConfig_IsBlocked(t *testing.T) {
	config := NewGMConfig(true, false)
	assert.False(t, config.IsBlocked())

	config.Block = true
	assert.True(t, config.IsBlocked())
}

func TestGMConfig_GetStatus(t *testing.T) {
	config := NewGMConfig(true, false)

	status := config.GetStatus()
	assert.Equal(t, true, status["gmEnable"])
	assert.Equal(t, false, status["block"])
}

func TestGMConfig_String(t *testing.T) {
	config := NewGMConfig(true, false)
	str := config.String()

	assert.Contains(t, str, "GMEnable")
	assert.Contains(t, str, "Block")
}

func TestGMConfig_Clone(t *testing.T) {
	config := NewGMConfig(true, false)
	clone := config.Clone()

	assert.Equal(t, config, clone)

	// 修改克隆不应影响原对象
	clone.GMEnable = false
	assert.True(t, config.GMEnable)
	assert.False(t, clone.GMEnable)
}

func TestGMConfig_Merge(t *testing.T) {
	t.Run("合并有效配置", func(t *testing.T) {
		config := NewGMConfig(true, false)
		other := NewGMConfig(false, true)

		err := config.Merge(other)
		assert.NoError(t, err)
		assert.False(t, config.GMEnable)
		assert.True(t, config.Block)
	})

	t.Run("合并nil配置", func(t *testing.T) {
		config := NewGMConfig(true, false)
		err := config.Merge(nil)
		assert.Error(t, err)
	})
}

// ==================== IPWhitelist 测试 ====================

func TestNewIPWhitelist(t *testing.T) {
	whitelist := NewIPWhitelist()

	assert.NotNil(t, whitelist)
	assert.Empty(t, whitelist.IPs)
}

func TestNewIPWhitelistFromSlice(t *testing.T) {
	t.Run("有效IP列表", func(t *testing.T) {
		ips := []string{"192.168.1.1", "10.0.0.1"}
		whitelist, err := NewIPWhitelistFromSlice(ips)

		assert.NoError(t, err)
		assert.Len(t, whitelist.IPs, 2)
		assert.True(t, whitelist.Contains("192.168.1.1"))
		assert.True(t, whitelist.Contains("10.0.0.1"))
	})

	t.Run("无效IP列表", func(t *testing.T) {
		ips := []string{"invalid-ip"}
		_, err := NewIPWhitelistFromSlice(ips)

		assert.Error(t, err)
	})
}

func TestIPWhitelist_Validate(t *testing.T) {
	whitelist := NewIPWhitelist()
	assert.NoError(t, whitelist.Validate())
}

func TestIPWhitelist_IsValid(t *testing.T) {
	whitelist := NewIPWhitelist()
	assert.True(t, whitelist.IsValid())
}

func TestIPWhitelist_AddIP(t *testing.T) {
	whitelist := NewIPWhitelist()

	t.Run("添加有效IP", func(t *testing.T) {
		err := whitelist.AddIP("192.168.1.1")
		assert.NoError(t, err)
		assert.True(t, whitelist.Contains("192.168.1.1"))
	})

	t.Run("添加带空格的IP", func(t *testing.T) {
		err := whitelist.AddIP("  192.168.1.2  ")
		assert.NoError(t, err)
		assert.True(t, whitelist.Contains("192.168.1.2"))
	})

	t.Run("添加无效IP", func(t *testing.T) {
		err := whitelist.AddIP("invalid-ip")
		assert.Error(t, err)
	})

	t.Run("添加重复IP", func(t *testing.T) {
		whitelist.AddIP("192.168.1.3")
		err := whitelist.AddIP("192.168.1.3")
		assert.Error(t, err)
	})
}

func TestIPWhitelist_RemoveIP(t *testing.T) {
	whitelist := NewIPWhitelist()
	whitelist.AddIP("192.168.1.1")

	t.Run("移除存在的IP", func(t *testing.T) {
		err := whitelist.RemoveIP("192.168.1.1")
		assert.NoError(t, err)
		assert.False(t, whitelist.Contains("192.168.1.1"))
	})

	t.Run("移除不存在的IP", func(t *testing.T) {
		err := whitelist.RemoveIP("10.0.0.1")
		assert.Error(t, err)
	})

	t.Run("移除带空格的IP", func(t *testing.T) {
		whitelist.AddIP("192.168.1.2")
		err := whitelist.RemoveIP("  192.168.1.2  ")
		assert.NoError(t, err)
	})
}

func TestIPWhitelist_Contains(t *testing.T) {
	whitelist := NewIPWhitelist()
	whitelist.AddIP("192.168.1.1")

	t.Run("包含IP", func(t *testing.T) {
		assert.True(t, whitelist.Contains("192.168.1.1"))
	})

	t.Run("不包含IP", func(t *testing.T) {
		assert.False(t, whitelist.Contains("10.0.0.1"))
	})

	t.Run("带空格的IP", func(t *testing.T) {
		assert.True(t, whitelist.Contains("  192.168.1.1  "))
	})
}

func TestIPWhitelist_ContainsIPNet(t *testing.T) {
	whitelist := NewIPWhitelist()
	whitelist.AddIP("192.168.1.1")
	whitelist.AddIP("10.0.0.0/24")

	t.Run("精确匹配", func(t *testing.T) {
		assert.True(t, whitelist.ContainsIPNet("192.168.1.1"))
	})

	t.Run("CIDR匹配", func(t *testing.T) {
		assert.True(t, whitelist.ContainsIPNet("10.0.0.1"))
		assert.True(t, whitelist.ContainsIPNet("10.0.0.255"))
	})

	t.Run("不匹配", func(t *testing.T) {
		assert.False(t, whitelist.ContainsIPNet("172.16.0.1"))
	})

	t.Run("无效IP", func(t *testing.T) {
		assert.False(t, whitelist.ContainsIPNet("invalid"))
	})
}

func TestIPWhitelist_AddIPs(t *testing.T) {
	whitelist := NewIPWhitelist()

	t.Run("批量添加有效IP", func(t *testing.T) {
		ips := []string{"192.168.1.1", "192.168.1.2", "10.0.0.1"}
		err := whitelist.AddIPs(ips)
		assert.NoError(t, err)
		assert.Equal(t, 3, whitelist.Count())
	})

	t.Run("部分无效IP", func(t *testing.T) {
		// 创建新的whitelist以避免IP冲突
		newWhitelist := NewIPWhitelist()
		ips := []string{"192.168.1.3", "invalid", "10.0.0.2"}
		err := newWhitelist.AddIPs(ips)
		// 应该成功添加有效的IP（虽然有一个无效IP）
		assert.NoError(t, err)
	})

	t.Run("全部无效IP", func(t *testing.T) {
		ips := []string{"invalid1", "invalid2"}
		err := whitelist.AddIPs(ips)
		assert.Error(t, err)
	})
}

func TestIPWhitelist_RemoveIPs(t *testing.T) {
	whitelist := NewIPWhitelist()
	whitelist.AddIP("192.168.1.1")
	whitelist.AddIP("192.168.1.2")
	whitelist.AddIP("10.0.0.1")

	t.Run("批量移除存在的IP", func(t *testing.T) {
		ips := []string{"192.168.1.1", "192.168.1.2"}
		err := whitelist.RemoveIPs(ips)
		assert.NoError(t, err)
		assert.Equal(t, 1, whitelist.Count())
	})

	t.Run("移除部分不存在的IP", func(t *testing.T) {
		ips := []string{"10.0.0.1", "nonexistent"}
		err := whitelist.RemoveIPs(ips)
		// 应该成功移除存在的IP
		assert.NoError(t, err)
		assert.Equal(t, 0, whitelist.Count())
	})

	t.Run("移除全部不存在的IP", func(t *testing.T) {
		whitelist.Clear()
		ips := []string{"192.168.1.1", "192.168.1.2"}
		err := whitelist.RemoveIPs(ips)
		assert.Error(t, err)
	})
}

func TestIPWhitelist_GetIPs(t *testing.T) {
	whitelist := NewIPWhitelist()
	whitelist.AddIP("192.168.1.1")
	whitelist.AddIP("10.0.0.1")

	ips := whitelist.GetIPs()
	assert.Len(t, ips, 2)
	assert.Contains(t, ips, "192.168.1.1")
	assert.Contains(t, ips, "10.0.0.1")
}

func TestIPWhitelist_Count(t *testing.T) {
	whitelist := NewIPWhitelist()

	assert.Equal(t, 0, whitelist.Count())

	whitelist.AddIP("192.168.1.1")
	assert.Equal(t, 1, whitelist.Count())

	whitelist.AddIP("10.0.0.1")
	assert.Equal(t, 2, whitelist.Count())
}

func TestIPWhitelist_IsEmpty(t *testing.T) {
	whitelist := NewIPWhitelist()

	assert.True(t, whitelist.IsEmpty())

	whitelist.AddIP("192.168.1.1")
	assert.False(t, whitelist.IsEmpty())
}

func TestIPWhitelist_Clear(t *testing.T) {
	whitelist := NewIPWhitelist()
	whitelist.AddIP("192.168.1.1")
	whitelist.AddIP("10.0.0.1")

	whitelist.Clear()
	assert.True(t, whitelist.IsEmpty())
	assert.Equal(t, 0, whitelist.Count())
}

func TestIPWhitelist_Clone(t *testing.T) {
	whitelist := NewIPWhitelist()
	whitelist.AddIP("192.168.1.1")
	whitelist.AddIP("10.0.0.1")

	clone := whitelist.Clone()
	assert.Equal(t, whitelist.Count(), clone.Count())
	assert.True(t, clone.Contains("192.168.1.1"))

	// 修改克隆不应影响原对象
	clone.AddIP("172.16.0.1")
	assert.Equal(t, 2, whitelist.Count())
	assert.Equal(t, 3, clone.Count())
}

func TestIPWhitelist_ToSlice(t *testing.T) {
	whitelist := NewIPWhitelist()
	whitelist.AddIP("192.168.1.1")
	whitelist.AddIP("10.0.0.1")

	slice := whitelist.ToSlice()
	assert.Len(t, slice, 2)
	assert.Contains(t, slice, "192.168.1.1")
	assert.Contains(t, slice, "10.0.0.1")
}

func TestIPWhitelist_String(t *testing.T) {
	whitelist := NewIPWhitelist()
	whitelist.AddIP("192.168.1.1")
	whitelist.AddIP("10.0.0.1")

	str := whitelist.String()
	assert.Contains(t, str, "IPWhitelist")
	assert.Contains(t, str, "Count: 2")
}

func TestIPWhitelist_Merge(t *testing.T) {
	whitelist1 := NewIPWhitelist()
	whitelist1.AddIP("192.168.1.1")

	whitelist2 := NewIPWhitelist()
	whitelist2.AddIP("10.0.0.1")

	t.Run("合并有效白名单", func(t *testing.T) {
		err := whitelist1.Merge(whitelist2)
		assert.NoError(t, err)
		assert.Equal(t, 2, whitelist1.Count())
		assert.True(t, whitelist1.Contains("10.0.0.1"))
	})

	t.Run("合并nil白名单", func(t *testing.T) {
		err := whitelist1.Merge(nil)
		assert.Error(t, err)
	})
}

func TestIPWhitelist_ValidateIPs(t *testing.T) {
	whitelist := NewIPWhitelist()

	t.Run("全部有效IP", func(t *testing.T) {
		ips := []string{"192.168.1.1", "10.0.0.1"}
		err := whitelist.ValidateIPs(ips)
		assert.NoError(t, err)
	})

	t.Run("部分无效IP", func(t *testing.T) {
		ips := []string{"192.168.1.1", "invalid", "10.0.0.1"}
		err := whitelist.ValidateIPs(ips)
		assert.Error(t, err)
	})
}

// ==================== DomainError 测试 ====================

func TestNewDomainError(t *testing.T) {
	appErr := errcode.ErrInvalidParams
	domainErr := NewDomainError(appErr, "test context")

	assert.Equal(t, appErr, domainErr.Err)
	assert.Equal(t, "test context", domainErr.Context)
}

func TestDomainError_Error(t *testing.T) {
	t.Run("有上下文", func(t *testing.T) {
		appErr := errcode.ErrInvalidParams
		domainErr := NewDomainError(appErr, "test context")

		errMsg := domainErr.Error()
		assert.Contains(t, errMsg, "test context")
		assert.Contains(t, errMsg, appErr.Message)
	})

	t.Run("无上下文", func(t *testing.T) {
		appErr := errcode.ErrInvalidParams
		domainErr := NewDomainError(appErr, "")

		errMsg := domainErr.Error()
		assert.Equal(t, appErr.Message, errMsg)
	})
}

func TestDomainError_Unwrap(t *testing.T) {
	appErr := errcode.ErrInvalidParams
	domainErr := NewDomainError(appErr, "test context")

	assert.Equal(t, appErr, domainErr.Unwrap())
}

func TestWrapError(t *testing.T) {
	t.Run("包装AppError", func(t *testing.T) {
		appErr := errcode.ErrInvalidParams
		wrapped := WrapError(appErr, "test context")

		domainErr, ok := AsDomainError(wrapped)
		require.True(t, ok)
		assert.Equal(t, appErr, domainErr.Err)
		assert.Equal(t, "test context", domainErr.Context)
	})

	t.Run("包装DomainError", func(t *testing.T) {
		originalErr := NewDomainError(errcode.ErrInvalidParams, "context1")
		wrapped := WrapError(originalErr, "context2")

		domainErr, ok := AsDomainError(wrapped)
		require.True(t, ok)
		assert.Contains(t, domainErr.Context, "context2")
		assert.Contains(t, domainErr.Context, "context1")
	})

	t.Run("包装nil错误", func(t *testing.T) {
		wrapped := WrapError(nil, "context")
		assert.Nil(t, wrapped)
	})

	t.Run("包装普通错误", func(t *testing.T) {
		originalErr := assert.AnError
		wrapped := WrapError(originalErr, "test context")

		domainErr, ok := AsDomainError(wrapped)
		require.True(t, ok)
		assert.Equal(t, errcode.ErrInternalError.Code, domainErr.Err.Code)
	})
}

func TestErrInvalidVersion(t *testing.T) {
	err := ErrInvalidVersion("invalid", "format error")

	assert.True(t, IsDomainError(err))
	domainErr, _ := AsDomainError(err)
	assert.Equal(t, "version validation", domainErr.Context)
	assert.Equal(t, errcode.ErrInvalidParams.Code, domainErr.Err.Code)
}

func TestErrInvalidSubServer(t *testing.T) {
	err := ErrInvalidSubServer("server1", "invalid URL")

	assert.True(t, IsDomainError(err))
	domainErr, _ := AsDomainError(err)
	assert.Equal(t, "subServer validation", domainErr.Context)
}

func TestErrSubServerNotFound(t *testing.T) {
	err := ErrSubServerNotFound("server1", []string{"server2", "server3"})

	assert.True(t, IsDomainError(err))
	domainErr, _ := AsDomainError(err)
	assert.Equal(t, "subServer lookup", domainErr.Context)
	assert.Equal(t, errcode.ErrInvalidParams.Code, domainErr.Err.Code)
}

func TestErrInvalidIPAddress(t *testing.T) {
	err := ErrInvalidIPAddress("invalid-ip")

	assert.True(t, IsDomainError(err))
	domainErr, _ := AsDomainError(err)
	assert.Equal(t, "IP validation", domainErr.Context)
	assert.Equal(t, errcode.ErrInvalidIPFormat.Code, domainErr.Err.Code)
}

func TestErrIPAlreadyExists(t *testing.T) {
	err := ErrIPAlreadyExists("192.168.1.1")

	assert.True(t, IsDomainError(err))
	domainErr, _ := AsDomainError(err)
	assert.Equal(t, "IP whitelist add", domainErr.Context)
	assert.Equal(t, errcode.ErrDuplicateIP.Code, domainErr.Err.Code)
}

func TestErrIPNotInWhitelist(t *testing.T) {
	err := ErrIPNotInWhitelist("192.168.1.1")

	assert.True(t, IsDomainError(err))
	domainErr, _ := AsDomainError(err)
	assert.Equal(t, "IP whitelist check", domainErr.Context)
	assert.Equal(t, errcode.ErrIPNotFound.Code, domainErr.Err.Code)
}

func TestErrGMConfigNotFound(t *testing.T) {
	err := ErrGMConfigNotFound()

	assert.True(t, IsDomainError(err))
	domainErr, _ := AsDomainError(err)
	assert.Equal(t, "GM config lookup", domainErr.Context)
	assert.Equal(t, errcode.ErrGMConfigNotFound.Code, domainErr.Err.Code)
}

func TestErrVersionNotFound(t *testing.T) {
	err := ErrVersionNotFound("1.2.3")

	assert.True(t, IsDomainError(err))
	domainErr, _ := AsDomainError(err)
	assert.Equal(t, "version lookup", domainErr.Context)
	assert.Equal(t, errcode.ErrVersionNotFound.Code, domainErr.Err.Code)
}

func TestAsDomainError(t *testing.T) {
	t.Run("DomainError", func(t *testing.T) {
		err := ErrInvalidVersion("1.0.0", "test")
		domainErr, ok := AsDomainError(err)

		assert.True(t, ok)
		assert.NotNil(t, domainErr)
	})

	t.Run("nil错误", func(t *testing.T) {
		domainErr, ok := AsDomainError(nil)

		assert.False(t, ok)
		assert.Nil(t, domainErr)
	})

	t.Run("普通错误", func(t *testing.T) {
		domainErr, ok := AsDomainError(assert.AnError)

		assert.False(t, ok)
		assert.Nil(t, domainErr)
	})
}

func TestIsDomainError(t *testing.T) {
	t.Run("是DomainError", func(t *testing.T) {
		err := ErrInvalidVersion("1.0.0", "test")
		assert.True(t, IsDomainError(err))
	})

	t.Run("不是DomainError", func(t *testing.T) {
		assert.False(t, IsDomainError(assert.AnError))
		assert.False(t, IsDomainError(nil))
	})
}
