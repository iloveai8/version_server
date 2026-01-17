// Package log 提供结构化日志功能
//
// 该包基于uber-go/zap实现高性能结构化日志：
// - JSON格式输出（生产环境）
// - 支持多级别日志（Debug/Info/Warn/Error/Fatal）
// - 支持请求ID等自定义字段
// - 线程安全
// - 支持全局Logger和子Logger
//
// 使用示例：
//   log.InitLogger("info")
//   log.Info("服务器启动", log.String("addr", ":9091"))
//   log.Error("请求失败", log.WithError(err), log.String("path", "/api/v1/version"))
package log

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	globalLogger *Logger
)

// Logger 日志记录器
//
// 封装zap.SugaredLogger，提供更友好的API
type Logger struct {
	logger *zap.SugaredLogger
}

// Field 日志字段
//
// 用于结构化日志中的键值对
type Field struct {
	Key   string
	Value interface{}
}

// InitLogger 初始化全局Logger
//
// 根据指定的日志级别创建全局Logger实例
//
// 参数:
//   level: 日志级别，可选值：debug/info/warn/error
//
// 返回:
//   error: 日志级别无效时返回错误
//
// 示例:
//   err := log.InitLogger("info")
//   if err != nil {
//       os.Exit(1)
//   }
func InitLogger(level string) error {
	// 检查空字符串
	if level == "" {
		return fmt.Errorf("日志级别不能为空")
	}

	// 解析日志级别
	var zapLevel zapcore.Level
	if err := zapLevel.UnmarshalText([]byte(level)); err != nil {
		return fmt.Errorf("无效的日志级别: %w", err)
	}

	// 编码器配置
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,  // 小写级别名
		EncodeTime:     zapcore.ISO8601TimeEncoder,     // ISO8601时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,     // 短调用者格式
	}

	// 创建Core
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(os.Stdout),
		zapLevel,
	)

	// 创建Logger
	zapLogger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	globalLogger = &Logger{logger: zapLogger.Sugar()}

	return nil
}

// GetLogger 获取全局Logger
//
// 返回全局Logger实例，如果未初始化则返回一个默认Logger
//
// 返回:
//   *Logger: 全局Logger实例
//
// 示例:
//   logger := log.GetLogger()
//   logger.Info("使用全局Logger")
func GetLogger() *Logger {
	if globalLogger == nil {
		// 如果未初始化，创建一个默认的info级别Logger
		_ = InitLogger("info")
	}
	return globalLogger
}

// Debug 输出调试级别日志
//
// 用于开发调试时的详细信息，生产环境默认不输出
//
// 参数:
//   msg: 日志消息
//   fields: 日志字段键值对
//
// 示例:
//   log.Debug("用户详情", log.String("userId", "12345"), log.Int("age", 25))
func (l *Logger) Debug(msg string, fields ...Field) {
	l.logger.Debugw(msg, toZapFields(fields...)...)
}

// Info 输出信息级别日志
//
// 用于记录常规的业务流程信息
//
// 参数:
//   msg: 日志消息
//   fields: 日志字段键值对
//
// 示例:
//   log.Info("服务器启动", log.String("addr", ":9091"))
func (l *Logger) Info(msg string, fields ...Field) {
	l.logger.Infow(msg, toZapFields(fields...)...)
}

// Warn 输出警告级别日志
//
// 用于记录可能的问题，但不会影响系统运行
//
// 参数:
//   msg: 日志消息
//   fields: 日志字段键值对
//
// 示例:
//   log.Warn("配置使用默认值", log.String("key", "timeout"))
func (l *Logger) Warn(msg string, fields ...Field) {
	l.logger.Warnw(msg, toZapFields(fields...)...)
}

// Error 输出错误级别日志
//
// 用于记录错误信息，需要关注但不影响系统继续运行
//
// 参数:
//   msg: 日志消息
//   fields: 日志字段键值对
//
// 示例:
//   log.Error("Redis查询失败", log.WithError(err), log.String("key", "user:123"))
func (l *Logger) Error(msg string, fields ...Field) {
	l.logger.Errorw(msg, toZapFields(fields...)...)
}

// Fatal 输出致命错误日志并退出程序
//
// 用于记录无法恢复的错误，记录后调用os.Exit(1)退出程序
//
// 参数:
//   msg: 日志消息
//   fields: 日志字段键值对
//
// 示例:
//   log.Fatal("配置文件加载失败", log.WithError(err))
func (l *Logger) Fatal(msg string, fields ...Field) {
	l.logger.Fatalw(msg, toZapFields(fields...)...)
}

// WithRequestID 创建带请求ID的子Logger
//
// 返回一个新的Logger实例，所有日志会自动包含指定的请求ID
//
// 参数:
//   requestID: 请求ID
//
// 返回:
//   *Logger: 新的Logger实例
//
// 示例:
//   logger := log.WithRequestID("req-12345")
//   logger.Info("处理请求")  // 自动包含requestID字段
func (l *Logger) WithRequestID(requestID string) *Logger {
	return &Logger{
		logger: l.logger.With("requestId", requestID),
	}
}

// With 创建带自定义字段的子Logger
//
// 返回一个新的Logger实例，所有日志会自动包含指定的字段
//
// 参数:
//   fields: 日志字段键值对
//
// 返回:
//   *Logger: 新的Logger实例
//
// 示例:
//   logger := log.With(log.String("service", "user-api"), log.String("version", "1.0.0"))
//   logger.Info("服务启动")  // 自动包含service和version字段
func (l *Logger) With(fields ...Field) *Logger {
	return &Logger{
		logger: l.logger.With(toZapFields(fields...)...),
	}
}

// Sync 刷新日志缓冲区
//
// 在程序退出前调用，确保所有日志都被写入
//
// 返回:
//   error: 刷新失败时返回错误
//
// 示例:
//   defer log.Sync()
func (l *Logger) Sync() error {
	return l.logger.Sync()
}

// ==================== 全局便捷函数 ====================

// Debug 输出调试级别日志（全局函数）
func Debug(msg string, fields ...Field) {
	GetLogger().Debug(msg, fields...)
}

// Info 输出信息级别日志（全局函数）
func Info(msg string, fields ...Field) {
	GetLogger().Info(msg, fields...)
}

// Warn 输出警告级别日志（全局函数）
func Warn(msg string, fields ...Field) {
	GetLogger().Warn(msg, fields...)
}

// Error 输出错误级别日志（全局函数）
func Error(msg string, fields ...Field) {
	GetLogger().Error(msg, fields...)
}

// Fatal 输出致命错误日志并退出（全局函数）
func Fatal(msg string, fields ...Field) {
	GetLogger().Fatal(msg, fields...)
}

// WithRequestID 创建带请求ID的子Logger（全局函数）
func WithRequestID(requestID string) *Logger {
	return GetLogger().WithRequestID(requestID)
}

// With 创建带自定义字段的子Logger（全局函数）
func With(fields ...Field) *Logger {
	return GetLogger().With(fields...)
}

// Sync 刷新全局Logger缓冲区（全局函数）
func Sync() error {
	if globalLogger != nil {
		return globalLogger.Sync()
	}
	return nil
}

// ==================== Field构造函数 ====================

// String 创建字符串类型的日志字段
//
// 参数:
//   key: 字段键名
//   value: 字段值
//
// 返回:
//   Field: 日志字段
//
// 示例:
//   log.Info("用户登录", log.String("username", "admin"))
func String(key, value string) Field {
	return Field{Key: key, Value: value}
}

// Int 创建整数类型的日志字段
//
// 参数:
//   key: 字段键名
//   value: 字段值
//
// 返回:
//   Field: 日志字段
//
// 示例:
//   log.Info("用户年龄", log.Int("age", 25))
func Int(key string, value int) Field {
	return Field{Key: key, Value: value}
}

// Int64 创建int64类型的日志字段
func Int64(key string, value int64) Field {
	return Field{Key: key, Value: value}
}

// Float64 创建float64类型的日志字段
func Float64(key string, value float64) Field {
	return Field{Key: key, Value: value}
}

// Bool 创建布尔类型的日志字段
func Bool(key string, value bool) Field {
	return Field{Key: key, Value: value}
}

// Any 创建任意类型的日志字段
//
// 使用fmt.Sprintf格式化值
//
// 参数:
//   key: 字段键名
//   value: 字段值
//
// 返回:
//   Field: 日志字段
//
// 示例:
//   log.Info("配置", log.Any("config", cfg))
func Any(key string, value interface{}) Field {
	return Field{Key: key, Value: value}
}

// WithError 创建错误类型的日志字段
//
// 用于在日志中记录错误信息
//
// 参数:
//   err: 错误对象
//
// 返回:
//   Field: 日志字段
//
// 示例:
//   log.Error("操作失败", log.WithError(err))
func WithError(err error) Field {
	if err == nil {
		return Field{Key: "error", Value: nil}
	}
	return Field{Key: "error", Value: err.Error()}
}

// ==================== 内部辅助函数 ====================

// toZapFields 将Field切片转换为interface{}切片（用于SugaredLogger）
//
// SugaredLogger的方法接受交替的键值对（key1, value1, key2, value2, ...）
func toZapFields(fields ...Field) []interface{} {
	ifaces := make([]interface{}, 0, len(fields)*2)
	for _, f := range fields {
		ifaces = append(ifaces, f.Key, f.Value)
	}
	return ifaces
}
