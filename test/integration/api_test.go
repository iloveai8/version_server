// Package integration 提供端到端集成测试
//
// 该测试包实现了：
// - 使用miniredis模拟Redis
// - 完整的HTTP服务器测试
// - 端到端API测试
package integration

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	redispkg "game_slots_vsn/internal/repository/redis"
	"game_slots_vsn/internal/repository"
	"game_slots_vsn/internal/service"
	httpserver "game_slots_vsn/internal/http"
	"game_slots_vsn/internal/domain"
)

// TestServer 测试服务器
type TestServer struct {
	server    *httptest.Server
	mr        *miniredis.Miniredis
	redisCl   redispkg.Client
	logger    *zap.Logger
	versionSvc service.VersionService
	gmSvc     service.GMService
	whitelistSvc service.IPWhitelistService
}

// SetupTestServer 创建测试服务器
//
// 返回:
//   *TestServer: 测试服务器实例
//   func: 清理函数
func SetupTestServer(t *testing.T) *TestServer {
	// 创建miniredis
	mr := miniredis.RunT(t)

	// 创建Redis客户端
	redisClient := redispkg.NewRedisClientFromURL(mr.Addr())

	// 创建logger
	logger := zap.NewNop()

	// 创建Repository
	versionRepo := repository.NewVersionRepository(redisClient)
	gmRepo := repository.NewGMRepository(redisClient)
	whitelistRepo := repository.NewIPWhitelistRepository(redisClient)

	// 创建Service
	versionSvc := service.NewVersionService(versionRepo)
	gmSvc := service.NewGMService(gmRepo)
	whitelistSvc := service.NewIPWhitelistService(whitelistRepo)

	// 初始化GM配置
	ctx := context.Background()
	err := gmSvc.InitializeGMConfig(ctx)
	require.NoError(t, err, "初始化GM配置失败")

	// 创建HTTP服务器
	httpServer := httpserver.NewServer(
		versionSvc,
		gmSvc,
		whitelistSvc,
		redisClient,
		logger,
		"",
		0, // 使用随机端口
	)

	// 设置路由
	httpServer.SetupRoutes()

	// 创建测试服务器
	engine := httpServer.GetEngine()
	testServer := httptest.NewServer(engine)

	return &TestServer{
		server:       testServer,
		mr:           mr,
		redisCl:      redisClient,
		logger:       logger,
		versionSvc:   versionSvc,
		gmSvc:        gmSvc,
		whitelistSvc: whitelistSvc,
	}
}

// Close 关闭测试服务器
func (ts *TestServer) Close() {
	if ts.server != nil {
		ts.server.Close()
	}
	if ts.mr != nil {
		ts.mr.Close()
	}
}

// GetURL 获取测试服务器URL
func (ts *TestServer) GetURL() string {
	return ts.server.URL
}

// ==================== 健康检查测试 ====================

func TestHealthCheck(t *testing.T) {
	ts := SetupTestServer(t)
	defer ts.Close()

	// 测试 /health/live
	t.Run("Liveness", func(t *testing.T) {
		resp, err := http.Get(ts.GetURL() + "/health/live")
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var result map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err)
		assert.Equal(t, "alive", result["status"])
	})

	// 测试 /health/ready
	t.Run("Readiness", func(t *testing.T) {
		resp, err := http.Get(ts.GetURL() + "/health/ready")
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var result map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err)
		assert.Equal(t, "ready", result["status"])
	})

	// 测试 /health
	t.Run("FullHealth", func(t *testing.T) {
		resp, err := http.Get(ts.GetURL() + "/health")
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var result map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err)
		assert.Equal(t, "healthy", result["status"])
		assert.NotEmpty(t, result["timestamp"])
		assert.NotEmpty(t, result["uptime"])
		assert.NotEmpty(t, result["checks"])
	})

	// 测试 /ping
	t.Run("Ping", func(t *testing.T) {
		resp, err := http.Get(ts.GetURL() + "/ping")
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var result map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err)
		assert.Equal(t, "pong", result["message"])
	})
}

// ==================== 版本配置API测试 ====================

func TestVersionAPI(t *testing.T) {
	ts := SetupTestServer(t)
	defer ts.Close()

	ctx := context.Background()

	// 创建测试版本
	version := domain.NewVersion("1.0.0")
	version.SubServers = map[string]*domain.SubServer{
		"server1": {
			Vsn:    "1.0.0",
			SrvUrl: "http://server1.example.com",
			ResUrl: "http://server1.example.com/res",
			Type:   1,
		},
	}

	err := ts.versionSvc.CreateVersion(ctx, version, "dev")
	require.NoError(t, err, "创建版本失败")

	t.Run("GetVersion", func(t *testing.T) {
		resp, err := http.Get(ts.GetURL() + "/api/v1/version?vsn=1.0.0&env=dev")
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var result map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err)

		data := result["data"].(map[string]interface{})
		assert.Equal(t, "1.0.0", data["vsn"])
	})

	t.Run("ListVersions", func(t *testing.T) {
		resp, err := http.Get(ts.GetURL() + "/api/v1/versions?env=dev")
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var result map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err)

		data := result["data"].(map[string]interface{})
		versions := data["versions"].([]interface{})
		count := data["count"].(float64)
		assert.GreaterOrEqual(t, int(count), 1)
		assert.GreaterOrEqual(t, len(versions), 1)
	})

	t.Run("CreateVersion", func(t *testing.T) {
		newVersion := map[string]interface{}{
			"vsn": "2.0.0",
			"subServers": map[string]interface{}{
				"server2": map[string]interface{}{
					"vsn":    "2.0.0",
					"srvUrl": "http://server2.example.com",
					"resUrl": "http://server2.example.com/res",
					"type":   1,
				},
			},
		}
		body, _ := json.Marshal(newVersion)

		resp, err := http.Post(ts.GetURL()+"/api/v1/version?env=dev", "application/json",
			jsonBodyReader(body))
		require.NoError(t, err)
		defer resp.Body.Close()

		// 创建成功返回201
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
	})

	t.Run("DeleteVersion", func(t *testing.T) {
		// 先创建一个临时版本
		tempVersion := domain.NewVersion("3.0.0")
		err := ts.versionSvc.CreateVersion(ctx, tempVersion, "dev")
		require.NoError(t, err)

		// 删除版本
		req, err := http.NewRequest("DELETE", ts.GetURL()+"/api/v1/version/3.0.0?env=dev", nil)
		require.NoError(t, err)

		client := &http.Client{}
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// 删除成功返回204
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})
}

// ==================== GM配置API测试 ====================

func TestGMAPI(t *testing.T) {
	ts := SetupTestServer(t)
	defer ts.Close()

	t.Run("GetGMConfig", func(t *testing.T) {
		resp, err := http.Get(ts.GetURL() + "/admin/v1/gm/config")
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var result map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err)

		data := result["data"].(map[string]interface{})
		assert.Contains(t, data, "gmEnable")
		assert.Contains(t, data, "block")
	})

	t.Run("ToggleGM", func(t *testing.T) {
		resp, err := http.Post(ts.GetURL()+"/admin/v1/gm/toggle", "application/json", nil)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var result map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err)

		data := result["data"].(map[string]interface{})
		assert.Contains(t, data, "gmEnable")
	})

	t.Run("UpdateGMConfig", func(t *testing.T) {
		updateData := map[string]interface{}{
			"gmEnable": true,
			"block":    false,
		}
		body, _ := json.Marshal(updateData)

		// 使用PUT方法访问 /admin/v1/gm/config
		req, err := http.NewRequest("PUT", ts.GetURL()+"/admin/v1/gm/config",
			jsonBodyReader(body))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}

// ==================== IP白名单API测试 ====================

func TestIPWhitelistAPI(t *testing.T) {
	ts := SetupTestServer(t)
	defer ts.Close()

	t.Run("GetWhitelist", func(t *testing.T) {
		resp, err := http.Get(ts.GetURL() + "/admin/v1/ip-whitelist")
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var result map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err)

		data := result["data"].(map[string]interface{})
		assert.Contains(t, data, "ips")
	})

	t.Run("AddIP", func(t *testing.T) {
		ipData := map[string]interface{}{
			"ip": "192.168.1.100",
		}
		body, _ := json.Marshal(ipData)

		resp, err := http.Post(ts.GetURL()+"/admin/v1/ip-whitelist/ip", "application/json",
			jsonBodyReader(body))
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("CheckIP", func(t *testing.T) {
		checkData := map[string]interface{}{
			"ip": "192.168.1.100",
		}
		body, _ := json.Marshal(checkData)

		resp, err := http.Post(ts.GetURL()+"/admin/v1/ip-whitelist/check", "application/json",
			jsonBodyReader(body))
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var result map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err)

		data := result["data"].(map[string]interface{})
		assert.Equal(t, true, data["contains"])
	})

	t.Run("RemoveIP", func(t *testing.T) {
		// 使用DELETE方法删除IP，参数放在JSON body中
		ipData := map[string]interface{}{
			"ip": "192.168.1.100",
		}
		body, _ := json.Marshal(ipData)

		req, err := http.NewRequest("DELETE",
			ts.GetURL()+"/admin/v1/ip-whitelist/ip", jsonBodyReader(body))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// 删除成功返回204 No Content
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})
}

// ==================== 错误处理测试 ====================

func TestErrorHandling(t *testing.T) {
	ts := SetupTestServer(t)
	defer ts.Close()

	t.Run("VersionNotFound", func(t *testing.T) {
		resp, err := http.Get(ts.GetURL() + "/api/v1/version?vsn=999.0.0&env=dev")
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)

		var result map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err)

		assert.Contains(t, result, "error")
		errData := result["error"].(map[string]interface{})
		assert.Equal(t, "VERSION_NOT_FOUND", errData["code"])
	})

	t.Run("InvalidParams", func(t *testing.T) {
		resp, err := http.Get(ts.GetURL() + "/api/v1/version?vsn=invalid&env=dev")
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var result map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err)

		assert.Contains(t, result, "error")
		errData := result["error"].(map[string]interface{})
		assert.Equal(t, "INVALID_PARAMS", errData["code"])
	})
}

// ==================== 辅助函数 ====================

// jsonBodyReader 创建JSON请求体Reader
func jsonBodyReader(data []byte) *jsonReader {
	return &jsonReader{data: data}
}

type jsonReader struct {
	data []byte
	pos  int
}

func (r *jsonReader) Read(p []byte) (n int, err error) {
	if r.pos >= len(r.data) {
		return 0, io.EOF
	}
	n = copy(p, r.data[r.pos:])
	r.pos += n
	return n, nil
}

// TestMain 测试入口
func TestMain(m *testing.M) {
	// 设置Gin为测试模式
	gin.SetMode(gin.TestMode)

	// 运行测试
	code := m.Run()
	os.Exit(code)
}
