package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

// Load 从指定路径加载配置文件
//
// 从指定路径加载YAML配置文件并解析为Config结构
//
// 参数:
//   configPath: 配置文件路径
//
// 返回:
//   *Config: 配置对象
//   error: 文件读取失败或解析失败时返回错误
//
// 示例:
//   cfg, err := config.Load("./conf/dev.yaml")
//   if err != nil {
//       log.Fatal(err)
//   }
func Load(configPath string) (*Config, error) {
	// 设置配置文件
	v := viper.New()

	// 设置配置文件路径
	v.SetConfigFile(configPath)

	// 读取配置文件
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// 解析配置
	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// 设置全局配置
	GlobalConfig = &cfg

	return &cfg, nil
}

// LoadFromEnv 从环境变量加载配置路径
//
// 从CONFIG_PATH环境变量读取配置文件路径并加载配置
// 如果环境变量未设置，则使用默认路径./conf/dev.yaml
//
// 返回:
//   *Config: 配置对象
//   error: 加载失败时返回错误
//
// 示例:
//   os.Setenv("CONFIG_PATH", "./conf/pro.yaml")
//   cfg, err := config.LoadFromEnv()
func LoadFromEnv() (*Config, error) {
	// 从环境变量获取配置文件路径
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		// 默认配置文件路径
		configPath = "./conf/dev.yaml"
	}

	return Load(configPath)
}

// LoadByEnv 根据环境加载配置
//
// 根据环境名称加载对应的配置文件（dev/pre/pro）
// 如果env为空，则从APP_ENV环境变量读取，默认为"dev"
//
// 参数:
//   env: 环境名称（dev/pre/pro），为空则使用环境变量APP_ENV
//
// 返回:
//   *Config: 配置对象
//   error: 加载失败时返回错误
//
// 示例:
//   cfg, err := config.LoadByEnv("dev")
func LoadByEnv(env string) (*Config, error) {
	if env == "" {
		env = os.Getenv("APP_ENV")
		if env == "" {
			env = "dev"
		}
	}

	// 标准化环境名称
	env = strings.ToLower(env)

	configPath := fmt.Sprintf("./conf/%s.yaml", env)
	return Load(configPath)
}

// LoadOrDefault 加载配置，如果失败则返回默认配置
//
// 尝试从指定路径加载配置，如果失败则返回默认配置
//
// 参数:
//   configPath: 配置文件路径
//
// 返回:
//   *Config: 配置对象（加载成功或默认配置）
//   error: 总是返回nil
//
// 示例:
//   cfg, err := config.LoadOrDefault("./conf/dev.yaml")
//   // 不会返回错误，至少返回默认配置
func LoadOrDefault(configPath string) (*Config, error) {
	cfg, err := Load(configPath)
	if err != nil {
		// 返回默认配置
		cfg = DefaultConfig()
		GlobalConfig = cfg
		return cfg, nil
	}
	return cfg, nil
}

// DefaultConfig 返回默认配置
//
// 返回一个预配置的默认配置对象，适用于开发和测试
//
// 返回:
//   *Config: 默认配置对象
//
// 默认配置:
// - RunMode: debug
// - Port: 9091
// - Redis: localhost:6379, DB=0
// - 日志: 文件和控制台都启用
//
// 示例:
//   cfg := config.DefaultConfig()
func DefaultConfig() *Config {
	return &Config{
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
			URL:      "https://geoip.ftsview.com/GeoLite2-Country.mmdb",
			Duration: 1,
			Scope:    0,
		},
	}
}

// Get 获取全局配置
//
// 获取全局配置实例，如果未初始化则自动加载dev环境配置
// 如果加载失败则返回默认配置
//
// 返回:
//   *Config: 全局配置对象
//
// 示例:
//   cfg := config.Get()
//   addr := cfg.Server.GetAddr()
func Get() *Config {
	if GlobalConfig == nil {
		// 如果全局配置未初始化，尝试从默认路径加载
		cfg, err := LoadByEnv("dev")
		if err != nil {
			return DefaultConfig()
		}
		return cfg
	}
	return GlobalConfig
}

// Reload 重新加载配置
//
// 从指定路径重新加载配置并更新全局配置实例
//
// 参数:
//   configPath: 配置文件路径
//
// 返回:
//   error: 加载失败时返回错误
//
// 示例:
//   if err := config.Reload("./conf/pro.yaml"); err != nil {
//       log.Printf("Reload failed: %v", err)
//   }
func Reload(configPath string) error {
	cfg, err := Load(configPath)
	if err != nil {
		return err
	}
	GlobalConfig = cfg
	return nil
}

// MustLoad 加载配置，失败则panic
//
// 从指定路径加载配置，如果失败则panic
// 适用于main函数中初始化配置的场景
//
// 参数:
//   configPath: 配置文件路径
//
// 返回:
//   *Config: 配置对象
//
// Panic:
// 当配置加载失败时panic
//
// 示例:
//   cfg := config.MustLoad("./conf/dev.yaml")
func MustLoad(configPath string) *Config {
	cfg, err := Load(configPath)
	if err != nil {
		panic(fmt.Sprintf("failed to load config: %v", err))
	}
	return cfg
}

// MustLoadByEnv 根据环境加载配置，失败则panic
//
// 根据环境名称加载配置，如果失败则panic
// 适用于main函数中初始化配置的场景
//
// 参数:
//   env: 环境名称（dev/pre/pro）
//
// 返回:
//   *Config: 配置对象
//
// Panic:
// 当配置加载失败时panic
//
// 示例:
//   cfg := config.MustLoadByEnv("dev")
func MustLoadByEnv(env string) *Config {
	cfg, err := LoadByEnv(env)
	if err != nil {
		panic(fmt.Sprintf("failed to load config for env %s: %v", env, err))
	}
	return cfg
}
