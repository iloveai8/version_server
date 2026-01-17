// Package config 提供配置管理功能
//
// 该包实现了：
// - 配置文件加载（YAML格式）
// - 多环境配置支持（dev/pre/pro）
// - 全局配置管理
// - 配置验证
// - 日志路径处理
// - GeoIP配置
//
// 使用示例：
//   cfg, err := config.LoadByEnv("dev")
//   if err != nil {
//       log.Fatal(err)
//   }
//   addr := cfg.Server.GetAddr()
package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Config 总配置结构
//
// 包含应用、服务器、Redis、日志和GeoIP的所有配置项
type Config struct {
	App     AppConfig      `yaml:"app" validate:"required"`
	Server  ServerConfig   `yaml:"server" validate:"required"`
	Redis   RedisConfig    `yaml:"redis" validate:"required"`
	Log     LogConfig      `yaml:"log" validate:"required"`
	GeoIP   GeoIPConfig    `yaml:"ggeoip" validate:"required"`
}

// AppConfig 应用配置
//
// 定义应用程序运行模式和基本配置
type AppConfig struct {
	RunMode string `yaml:"runMode" validate:"required,oneof=debug release test"` // 运行模式：debug/release/test
}

// ServerConfig 服务器配置
//
// 定义HTTP服务器监听地址和端口
type ServerConfig struct {
	IP   string `yaml:"ip"`                                                          // 监听IP地址（空表示监听所有网卡）
	Port int    `yaml:"port" validate:"required,min=1,max=65535"`                   // 监听端口
}

// RedisConfig Redis配置
//
// 定义Redis连接参数和连接池设置
type RedisConfig struct {
	MasterName string   `yaml:"masterName"`  // Redis Sentinel主节点名称（可选）
	Password   string   `yaml:"password"`    // Redis密码
	Hosts      []string `yaml:"hosts" validate:"required,min=1"`                   // Redis主机列表（集群模式）
	DB         int      `yaml:"db" validate:"min=0,max=15"`                         // 数据库编号（0-15）
	PoolSize   int      `yaml:"poolSize" validate:"min=1"`                          // 连接池大小
}

// LogConfig 日志配置
//
// 定义日志输出方式、日志级别和文件轮转策略
type LogConfig struct {
	FileEnable    bool   `yaml:"fileEnable"`                                          // 是否启用文件日志
	FileName      string `yaml:"fileName"`                                            // 日志文件路径
	FileLevel     string `yaml:"fileLevel" validate:"oneof=debug info warn error"`   // 文件日志级别
	ConsoleEnable bool   `yaml:"consoleEnable"`                                       // 是否启用控制台日志
	ConsoleLevel  string `yaml:"consoleLevel" validate:"oneof=debug info warn error"` // 控制台日志级别
	MaxSize       int    `yaml:"maxSize" validate:"min=1"`                            // 单个日志文件最大大小（MB）
	MaxBackups    int    `yaml:"maxBackups" validate:"min=1"`                         // 保留的旧日志文件最大数量
	MaxAges       int    `yaml:"maxAges" validate:"min=1"`                            // 保留旧日志文件的最大天数
	Compress      bool   `yaml:"compress"`                                            // 是否压缩旧日志文件
	JSONEnable    bool   `yaml:"jsonEnable"`                                          // 是否使用JSON格式日志
}

// GeoIPConfig GeoIP配置
//
// 定义GeoIP数据库的下载和配置
type GeoIPConfig struct {
	Name     string `yaml:"name" validate:"required"`             // GeoIP数据库文件名
	Path     string `yaml:"path" validate:"required"`             // GeoIP数据库存储路径
	URL      string `yaml:"url" validate:"required,url"`          // GeoIP数据库下载地址
	Duration int    `yaml:"duration" validate:"min=1"`             // GeoIP数据库更新间隔（小时）
	Scope    int    `yaml:"scope" validate:"oneof=0 1"`            // GeoIP查询范围：0=国家，1=城市
}

// GetAddr 获取服务器地址
//
// 返回格式化的服务器监听地址（host:port）
//
// 返回:
//   string: 如果IP为空则返回":port"，否则返回"ip:port"
//
// 示例:
//   addr := cfg.Server.GetAddr()  // ":9091" 或 "0.0.0.0:9091"
func (s *ServerConfig) GetAddr() string {
	if s.IP == "" {
		return fmt.Sprintf(":%d", s.Port)
	}
	return fmt.Sprintf("%s:%d", s.IP, s.Port)
}

// IsDebug 是否为调试模式
//
// 判断当前运行模式是否为调试模式
//
// 返回:
//   bool: debug或test模式返回true，release模式返回false
//
// 示例:
//   if cfg.App.IsDebug() {
//       log.Println("Debug mode enabled")
//   }
func (a *AppConfig) IsDebug() bool {
	return a.RunMode == "debug" || a.RunMode == "test"
}

// GetRunMode 获取运行模式（标准化）
//
// 将运行模式标准化为debug或release
//
// 返回:
//   string: "debug" 或 "release"
//
// 示例:
//   mode := cfg.App.GetRunMode()  // "debug" 或 "release"
func (a *AppConfig) GetRunMode() string {
	switch strings.ToLower(a.RunMode) {
	case "debug", "test":
		return "debug"
	default:
		return "release"
	}
}

// GetLogFilePath 获取日志文件完整路径
//
// 获取日志文件的绝对路径
//
// 返回:
//   string: 日志文件的绝对路径
//   error: 如果文件名为空或获取绝对路径失败则返回错误
//
// 示例:
//   path, err := cfg.Log.GetLogFilePath()
//   if err != nil {
//       log.Fatal(err)
//   }
func (l *LogConfig) GetLogFilePath() (string, error) {
	if !l.FileEnable {
		return "", nil
	}

	if l.FileName == "" {
		return "", fmt.Errorf("log file name is empty")
	}

	// 如果是相对路径，转换为绝对路径
	if !filepath.IsAbs(l.FileName) {
		abs, err := filepath.Abs(l.FileName)
		if err != nil {
			return "", fmt.Errorf("failed to get absolute path: %w", err)
		}
		return abs, nil
	}

	return l.FileName, nil
}

// EnsureLogDir 确保日志目录存在
//
// 如果日志文件目录不存在则创建目录
//
// 返回:
//   error: 目录创建失败则返回错误
//
// 示例:
//   if err := cfg.Log.EnsureLogDir(); err != nil {
//       log.Fatal(err)
//   }
func (l *LogConfig) EnsureLogDir() error {
	if !l.FileEnable {
		return nil
	}

	filePath, err := l.GetLogFilePath()
	if err != nil {
		return err
	}

	if filePath == "" {
		return nil
	}

	logDir := filepath.Dir(filePath)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return fmt.Errorf("failed to create log directory: %w", err)
	}

	return nil
}

// GetGeoDBFilePath 获取GeoIP数据库文件完整路径
//
// 拼接GeoIP数据库目录和文件名
//
// 返回:
//   string: GeoIP数据库文件的完整路径
//
// 示例:
//   path := cfg.GeoIP.GetGeoDBFilePath()  // "./data/geoip.mmdb"
func (g *GeoIPConfig) GetGeoDBFilePath() string {
	path := g.Path

	// 如果路径为空，直接返回文件名
	if path == "" {
		return g.Name
	}

	// 确保路径以 / 结尾
	if !strings.HasSuffix(path, "/") {
		path += "/"
	}

	return path + g.Name
}

// IsCountryMode 是否为国家模式
//
// 判断GeoIP查询模式是否为国家级别
//
// 返回:
//   bool: Scope=0时返回true
//
// 示例:
//   if cfg.GeoIP.IsCountryMode() {
//       // 只查询国家信息
//   }
func (g *GeoIPConfig) IsCountryMode() bool {
	return g.Scope == 0
}

// IsCityMode 是否为城市模式
//
// 判断GeoIP查询模式是否为城市级别
//
// 返回:
//   bool: Scope=1时返回true
//
// 示例:
//   if cfg.GeoIP.IsCityMode() {
//       // 查询城市详细信息
//   }
func (g *GeoIPConfig) IsCityMode() bool {
	return g.Scope == 1
}

// GlobalConfig 全局配置实例（单例）
//
// 该变量保存应用程序的全局配置，在配置加载后初始化
var GlobalConfig *Config
