package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServerConfig_GetAddr(t *testing.T) {
	tests := []struct {
		name     string
		ip       string
		port     int
		expected string
	}{
		{
			name:     "有IP地址",
			ip:       "127.0.0.1",
			port:     9091,
			expected: "127.0.0.1:9091",
		},
		{
			name:     "无IP地址",
			ip:       "",
			port:     9091,
			expected: ":9091",
		},
		{
			name:     "IPv6地址",
			ip:       "::1",
			port:     9091,
			expected: "::1:9091",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := ServerConfig{
				IP:   tt.ip,
				Port: tt.port,
			}
			assert.Equal(t, tt.expected, cfg.GetAddr())
		})
	}
}

func TestAppConfig_IsDebug(t *testing.T) {
	tests := []struct {
		name     string
		runMode  string
		expected bool
	}{
		{
			name:     "调试模式",
			runMode:  "debug",
			expected: true,
		},
		{
			name:     "生产模式",
			runMode:  "release",
			expected: false,
		},
		{
			name:     "测试模式",
			runMode:  "test",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := AppConfig{RunMode: tt.runMode}
			assert.Equal(t, tt.expected, cfg.IsDebug())
		})
	}
}

func TestAppConfig_GetRunMode(t *testing.T) {
	tests := []struct {
		name     string
		runMode  string
		expected string
	}{
		{
			name:     "调试模式",
			runMode:  "debug",
			expected: "debug",
		},
		{
			name:     "调试模式（大写）",
			runMode:  "DEBUG",
			expected: "debug",
		},
		{
			name:     "生产模式",
			runMode:  "release",
			expected: "release",
		},
		{
			name:     "测试模式",
			runMode:  "test",
			expected: "debug",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := AppConfig{RunMode: tt.runMode}
			assert.Equal(t, tt.expected, cfg.GetRunMode())
		})
	}
}

func TestLogConfig_GetLogFilePath(t *testing.T) {
	tests := []struct {
		name        string
		fileEnable  bool
		fileName    string
		expectError bool
	}{
		{
			name:        "文件日志未启用",
			fileEnable:  false,
			fileName:    "./logs/app.log",
			expectError: false,
		},
		{
			name:        "文件名为空",
			fileEnable:  true,
			fileName:    "",
			expectError: true,
		},
		{
			name:        "相对路径",
			fileEnable:  true,
			fileName:    "./logs/app.log",
			expectError: false,
		},
		{
			name:        "绝对路径",
			fileEnable:  true,
			fileName:    "/var/log/app.log",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := LogConfig{
				FileEnable: tt.fileEnable,
				FileName:   tt.fileName,
			}

			path, err := cfg.GetLogFilePath()
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				if tt.fileEnable && tt.fileName != "" {
					assert.NotEmpty(t, path)
					if filepath.IsAbs(tt.fileName) {
						assert.Equal(t, tt.fileName, path)
					}
				}
			}
		})
	}
}

func TestLogConfig_EnsureLogDir(t *testing.T) {
	// 创建临时目录
	tempDir := os.TempDir()
	testLogPath := filepath.Join(tempDir, "test_logs", "app.log")

	cfg := LogConfig{
		FileEnable: true,
		FileName:   testLogPath,
	}

	// 确保日志目录
	err := cfg.EnsureLogDir()
	assert.NoError(t, err)

	// 验证目录已创建
	logDir := filepath.Dir(testLogPath)
	info, err := os.Stat(logDir)
	assert.NoError(t, err)
	assert.True(t, info.IsDir())

	// 清理
	os.RemoveAll(logDir)
}

func TestGeoIPConfig_GetGeoDBFilePath(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		dbName   string
		expected string
	}{
		{
			name:     "路径以/结尾",
			path:     "./data/",
			dbName:   "geoip.mmdb",
			expected: "./data/geoip.mmdb",
		},
		{
			name:     "路径不以/结尾",
			path:     "./data",
			dbName:   "geoip.mmdb",
			expected: "./data/geoip.mmdb",
		},
		{
			name:     "空路径",
			path:     "",
			dbName:   "geoip.mmdb",
			expected: "geoip.mmdb",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := GeoIPConfig{
				Path: tt.path,
				Name: tt.dbName,
			}
			result := cfg.GetGeoDBFilePath()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGeoIPConfig_IsCountryMode(t *testing.T) {
	cfg := GeoIPConfig{Scope: 0}
	assert.True(t, cfg.IsCountryMode())
	assert.False(t, cfg.IsCityMode())
}

func TestGeoIPConfig_IsCityMode(t *testing.T) {
	cfg := GeoIPConfig{Scope: 1}
	assert.True(t, cfg.IsCityMode())
	assert.False(t, cfg.IsCountryMode())
}

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	assert.NotNil(t, cfg)
	assert.Equal(t, "debug", cfg.App.RunMode)
	assert.Equal(t, 9091, cfg.Server.Port)
	assert.NotEmpty(t, cfg.Redis.Hosts)
	assert.Equal(t, 0, cfg.Redis.DB)
	assert.Equal(t, 5, cfg.Redis.PoolSize)
	assert.True(t, cfg.Log.FileEnable)
	assert.True(t, cfg.Log.ConsoleEnable)
	assert.NotEmpty(t, cfg.Log.FileName)
	assert.NotEmpty(t, cfg.GeoIP.Name)
	assert.NotEmpty(t, cfg.GeoIP.Path)
}

func TestLoad(t *testing.T) {
	// 创建测试配置文件
	tempDir := os.TempDir()
	testConfigPath := filepath.Join(tempDir, "test_config.yaml")

	testConfigContent := `
app:
  runMode: "test"

server:
  ip: "127.0.0.1"
  port: 8080

redis:
  masterName: ""
  password: ""
  hosts:
    - "localhost:6379"
  db: 0
  poolSize: 10

log:
  fileEnable: true
  fileName: "./logs/test.log"
  fileLevel: info
  consoleEnable: true
  consoleLevel: debug
  maxSize: 20
  maxBackups: 10
  maxAges: 30
  compress: true
  jsonEnable: true

ggeoip:
  name: "test.mmdb"
  path: "./data/"
  url: "https://example.com/test.mmdb"
  duration: 1
  scope: 0
`

	err := os.WriteFile(testConfigPath, []byte(testConfigContent), 0644)
	assert.NoError(t, err)
	defer os.Remove(testConfigPath)

	// 加载配置
	cfg, err := Load(testConfigPath)
	assert.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.Equal(t, "test", cfg.App.RunMode)
	assert.Equal(t, "127.0.0.1", cfg.Server.IP)
	assert.Equal(t, 8080, cfg.Server.Port)
	assert.Equal(t, 10, cfg.Redis.PoolSize)
}

func TestLoad_FileNotFound(t *testing.T) {
	_, err := Load("/nonexistent/config.yaml")
	assert.Error(t, err)
}

func TestLoadOrDefault(t *testing.T) {
	// 测试加载失败时返回默认配置
	cfg, err := LoadOrDefault("/nonexistent/config.yaml")
	assert.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.Equal(t, "debug", cfg.App.RunMode)
}

func TestGet(t *testing.T) {
	// 重置全局配置
	GlobalConfig = nil

	// 第一次调用应该加载配置或返回默认配置
	cfg := Get()
	assert.NotNil(t, cfg)

	// 第二次调用应该返回相同的实例（因为GlobalConfig已设置）
	cfg2 := Get()
	// 注意：每次Get()可能返回不同的实例，但GlobalConfig应该是同一个
	assert.Equal(t, cfg, cfg2)
}

func TestMustLoad(t *testing.T) {
	// 创建测试配置文件
	tempDir := os.TempDir()
	testConfigPath := filepath.Join(tempDir, "must_load_test.yaml")

	testConfigContent := `
app:
  runMode: "test"
server:
  port: 9091
`

	err := os.WriteFile(testConfigPath, []byte(testConfigContent), 0644)
	assert.NoError(t, err)
	defer os.Remove(testConfigPath)

	// 正常加载
	cfg := MustLoad(testConfigPath)
	assert.NotNil(t, cfg)
	assert.Equal(t, "test", cfg.App.RunMode)
}

func TestMustLoad_Panic(t *testing.T) {
	// 测试加载失败时是否panic
	assert.Panics(t, func() {
		MustLoad("/nonexistent/config.yaml")
	})
}

func TestLoadByEnv(t *testing.T) {
	tests := []struct {
		name        string
		env         string
		expectError bool
	}{
		{
			name:        "开发环境（文件不存在）",
			env:         "dev",
			expectError: true,
		},
		{
			name:        "生产环境（文件不存在）",
			env:         "pro",
			expectError: true,
		},
		{
			name:        "预发布环境（文件不存在）",
			env:         "pre",
			expectError: true,
		},
		{
			name:        "空环境（默认dev，文件不存在）",
			env:         "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := LoadByEnv(tt.env)
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, cfg)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, cfg)
			}
		})
	}
}

func TestConfigStructures(t *testing.T) {
	// 测试配置结构体是否正确初始化
	cfg := &Config{
		App: AppConfig{
			RunMode: "debug",
		},
		Server: ServerConfig{
			IP:   "",
			Port: 9091,
		},
		Redis: RedisConfig{
			MasterName: "",
			Password:   "",
			Hosts:      []string{"localhost:6379"},
			DB:         0,
			PoolSize:   5,
		},
		Log: LogConfig{
			FileEnable:    true,
			FileName:      "./logs/app.log",
			FileLevel:     "info",
			ConsoleEnable: true,
			ConsoleLevel:  "debug",
			MaxSize:       20,
			MaxBackups:    10,
			MaxAges:       30,
			Compress:      true,
			JSONEnable:    true,
		},
		GeoIP: GeoIPConfig{
			Name:     "geoip.mmdb",
			Path:     "./data/",
			URL:      "https://example.com/geoip.mmdb",
			Duration: 1,
			Scope:    0,
		},
	}

	assert.NotNil(t, cfg)
	assert.Equal(t, "debug", cfg.App.RunMode)
	assert.Equal(t, 9091, cfg.Server.Port)
	assert.Len(t, cfg.Redis.Hosts, 1)
	assert.True(t, cfg.Log.FileEnable)
	assert.True(t, cfg.GeoIP.IsCountryMode())
}

// TestLoadFromEnv 测试从环境变量加载配置
func TestLoadFromEnv(t *testing.T) {
	// 创建测试配置文件
	tempDir := os.TempDir()
	testConfigPath := filepath.Join(tempDir, "env_test_config.yaml")

	testConfigContent := `
app:
  runMode: "test"
server:
  port: 8080
`
	err := os.WriteFile(testConfigPath, []byte(testConfigContent), 0644)
	assert.NoError(t, err)
	defer os.Remove(testConfigPath)

	// 设置CONFIG_PATH环境变量
	os.Setenv("CONFIG_PATH", testConfigPath)
	defer os.Unsetenv("CONFIG_PATH")

	// 从环境变量加载配置
	cfg, err := LoadFromEnv()
	assert.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.Equal(t, "test", cfg.App.RunMode)
	assert.Equal(t, 8080, cfg.Server.Port)
}

// TestLoadFromEnv_DefaultPath 测试使用默认路径
func TestLoadFromEnv_DefaultPath(t *testing.T) {
	// 不设置CONFIG_PATH环境变量，使用默认路径
	os.Unsetenv("CONFIG_PATH")

	// 默认路径./conf/dev.yaml不存在，应该返回错误
	cfg, err := LoadFromEnv()
	assert.Error(t, err)
	assert.Nil(t, cfg)
}

// TestReload 测试重新加载配置
func TestReload(t *testing.T) {
	// 创建两个测试配置文件
	tempDir := os.TempDir()
	configPath1 := filepath.Join(tempDir, "reload_test1.yaml")
	configPath2 := filepath.Join(tempDir, "reload_test2.yaml")

	configContent1 := `
app:
  runMode: "dev"
server:
  port: 9091
`
	configContent2 := `
app:
  runMode: "pro"
server:
  port: 8080
`

	// 写入第一个配置文件
	err := os.WriteFile(configPath1, []byte(configContent1), 0644)
	assert.NoError(t, err)
	defer os.Remove(configPath1)

	// 加载第一个配置
	cfg1, err := Load(configPath1)
	assert.NoError(t, err)
	assert.Equal(t, "dev", cfg1.App.RunMode)
	assert.Equal(t, 9091, cfg1.Server.Port)

	// 写入第二个配置文件
	err = os.WriteFile(configPath2, []byte(configContent2), 0644)
	assert.NoError(t, err)
	defer os.Remove(configPath2)

	// 重新加载配置
	err = Reload(configPath2)
	assert.NoError(t, err)

	// 验证全局配置已更新
	cfg2 := GlobalConfig
	assert.NotNil(t, cfg2)
	assert.Equal(t, "pro", cfg2.App.RunMode)
	assert.Equal(t, 8080, cfg2.Server.Port)
}

// TestReload_FileNotFound 测试重新加载不存在的文件
func TestReload_FileNotFound(t *testing.T) {
	err := Reload("/nonexistent/config.yaml")
	assert.Error(t, err)
}

// TestMustLoadByEnv 测试根据环境加载配置（成功）
func TestMustLoadByEnv(t *testing.T) {
	// 创建测试配置文件
	tempDir := os.TempDir()
	testConfigPath := filepath.Join(tempDir, "must_load_by_env_test.yaml")

	testConfigContent := `
app:
  runMode: "test"
server:
  port: 9999
`
	err := os.WriteFile(testConfigPath, []byte(testConfigContent), 0644)
	assert.NoError(t, err)
	defer os.Remove(testConfigPath)

	// 临时修改LoadByEnv的行为，让它使用我们的测试文件
	// 由于LoadByEnv使用固定路径格式，我们需要创建对应路径的文件
	confDir := filepath.Join(tempDir, "conf")
	err = os.MkdirAll(confDir, 0755)
	assert.NoError(t, err)
	defer os.RemoveAll(confDir)

	devConfigPath := filepath.Join(confDir, "dev.yaml")
	err = os.WriteFile(devConfigPath, []byte(testConfigContent), 0644)
	assert.NoError(t, err)

	// 修改当前工作目录到tempDir
	originalWd, _ := os.Getwd()
	defer os.Chdir(originalWd)
	os.Chdir(tempDir)

	// 正常加载
	cfg := MustLoadByEnv("dev")
	assert.NotNil(t, cfg)
	assert.Equal(t, "test", cfg.App.RunMode)
	assert.Equal(t, 9999, cfg.Server.Port)
}

// TestMustLoadByEnv_Panic 测试根据环境加载配置（失败时panic）
func TestMustLoadByEnv_Panic(t *testing.T) {
	// 在没有配置文件的情况下应该panic
	assert.Panics(t, func() {
		MustLoadByEnv("nonexistent_env")
	})
}
