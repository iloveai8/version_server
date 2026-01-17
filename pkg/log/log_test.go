package log

import (
	"bytes"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// TestInitLogger 测试日志初始化
func TestInitLogger(t *testing.T) {
	tests := []struct {
		name    string
		level   string
		wantErr bool
	}{
		{"有效级别debug", "debug", false},
		{"有效级别info", "info", false},
		{"有效级别warn", "warn", false},
		{"有效级别error", "error", false},
		{"无效级别", "invalid", true},
		{"空级别", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := InitLogger(tt.level)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, globalLogger)
			}
		})
	}
}

// TestGetLogger 测试获取全局Logger
func TestGetLogger(t *testing.T) {
	// 重置全局Logger
	globalLogger = nil

	logger := GetLogger()
	assert.NotNil(t, logger)

	// 再次调用应该返回同一个实例
	logger2 := GetLogger()
	assert.Same(t, logger, logger2)
}

// TestLoggerLevels 测试日志级别输出
func TestLoggerLevels(t *testing.T) {
	// 创建一个自定义的写入缓冲区
	var buf bytes.Buffer

	// 创建测试Logger
	testLogger := createTestLogger(&buf)

	tests := []struct {
		name   string
		level  string
		logMsg string
		logFn  func(msg string, fields ...Field)
	}{
		{"Debug级别", "debug", "这是debug日志", testLogger.Debug},
		{"Info级别", "info", "这是info日志", testLogger.Info},
		{"Warn级别", "warn", "这是warn日志", testLogger.Warn},
		{"Error级别", "error", "这是error日志", testLogger.Error},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf.Reset()

			// 写入日志
			tt.logFn(tt.logMsg)

			// 验证日志输出
			output := buf.String()
			assert.Contains(t, output, tt.logMsg)
			assert.Contains(t, output, tt.level)
		})
	}
}

// TestLoggerWithFields 测试带字段的日志
func TestLoggerWithFields(t *testing.T) {
	var buf bytes.Buffer
	testLogger := createTestLogger(&buf)

	tests := []struct {
		name    string
		msg     string
		fields  []Field
		wantKey []string
	}{
		{
			name:    "字符串字段",
			msg:     "用户登录",
			fields:  []Field{String("username", "admin"), String("ip", "127.0.0.1")},
			wantKey: []string{"username", "ip"},
		},
		{
			name:    "整数字段",
			msg:     "处理完成",
			fields:  []Field{Int("count", 100), Int("total", 200)},
			wantKey: []string{"count", "total"},
		},
		{
			name:    "布尔字段",
			msg:     "系统状态",
			fields:  []Field{Bool("success", true), Bool("enabled", false)},
			wantKey: []string{"success", "enabled"},
		},
		{
			name:    "混合字段",
			msg:     "请求处理",
			fields:  []Field{
				String("method", "GET"),
				String("path", "/api/v1/version"),
				Int("status", 200),
				Int64("duration", 1234567890),
				Float64("time", 1.234),
				Bool("cache", true),
			},
			wantKey: []string{"method", "path", "status", "duration", "time", "cache"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf.Reset()

			testLogger.Info(tt.msg, tt.fields...)

			// 解析JSON日志
			output := buf.String()
			var logEntry map[string]interface{}
			err := json.Unmarshal([]byte(output), &logEntry)
			assert.NoError(t, err)

			// 验证消息
			assert.Equal(t, tt.msg, logEntry["msg"])

			// 验证字段存在
			for _, key := range tt.wantKey {
				_, exists := logEntry[key]
				assert.True(t, exists, "字段 %s 应该存在", key)
			}
		})
	}
}

// TestWithRequestID 测试请求ID功能
func TestWithRequestID(t *testing.T) {
	var buf bytes.Buffer
	testLogger := createTestLogger(&buf)

	// 创建带请求ID的Logger
	requestID := "req-12345-67890"
	loggerWithID := testLogger.WithRequestID(requestID)

	// 使用带请求ID的Logger记录日志
	loggerWithID.Info("处理请求")

	// 验证日志输出
	output := buf.String()
	var logEntry map[string]interface{}
	err := json.Unmarshal([]byte(output), &logEntry)
	assert.NoError(t, err)

	// 验证请求ID存在
	assert.Equal(t, requestID, logEntry["requestId"])
	assert.Equal(t, "处理请求", logEntry["msg"])
}

// TestWith测试With功能
func TestWith(t *testing.T) {
	var buf bytes.Buffer
	testLogger := createTestLogger(&buf)

	// 创建带自定义字段的Logger
	serviceName := "user-api"
	version := "1.0.0"
	loggerWithFields := testLogger.With(
		String("service", serviceName),
		String("version", version),
	)

	// 记录多条日志，都应该包含service和version字段
	loggerWithFields.Info("服务启动")
	loggerWithFields.Info("处理请求")

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	assert.Len(t, lines, 2)

	// 验证第一条日志
	var log1 map[string]interface{}
	assert.NoError(t, json.Unmarshal([]byte(lines[0]), &log1))
	assert.Equal(t, serviceName, log1["service"])
	assert.Equal(t, version, log1["version"])
	assert.Equal(t, "服务启动", log1["msg"])

	// 验证第二条日志
	var log2 map[string]interface{}
	assert.NoError(t, json.Unmarshal([]byte(lines[1]), &log2))
	assert.Equal(t, serviceName, log2["service"])
	assert.Equal(t, version, log2["version"])
	assert.Equal(t, "处理请求", log2["msg"])
}

// TestWithError 测试错误字段
func TestWithError(t *testing.T) {
	var buf bytes.Buffer
	testLogger := createTestLogger(&buf)

	t.Run("非空错误", func(t *testing.T) {
		buf.Reset()
		err := errors.New("数据库连接失败")
		testLogger.Error("操作失败", WithError(err))

		output := buf.String()
		var logEntry map[string]interface{}
		jsonErr := json.Unmarshal([]byte(output), &logEntry)
		assert.NoError(t, jsonErr)

		// 验证错误字段存在
		assert.Contains(t, logEntry["error"], "数据库连接失败")
		assert.Equal(t, "操作失败", logEntry["msg"])
	})

	t.Run("空错误", func(t *testing.T) {
		buf.Reset()
		testLogger.Error("测试", WithError(nil))

		output := buf.String()
		var logEntry map[string]interface{}
		jsonErr := json.Unmarshal([]byte(output), &logEntry)
		assert.NoError(t, jsonErr)

		// 验证错误字段为nil
		assert.Nil(t, logEntry["error"])
	})
}

// TestGlobalConvenienceFunctions 测试全局便捷函数
func TestGlobalConvenienceFunctions(t *testing.T) {
	// 确保全局Logger已初始化
	_ = InitLogger("info")

	t.Run("Debug全局函数", func(t *testing.T) {
		Debug("调试信息", String("key", "value"))
	})

	t.Run("Info全局函数", func(t *testing.T) {
		Info("信息", String("key", "value"))
	})

	t.Run("Warn全局函数", func(t *testing.T) {
		Warn("警告", String("key", "value"))
	})

	t.Run("Error全局函数", func(t *testing.T) {
		Error("错误", String("key", "value"))
	})

	t.Run("WithRequestID全局函数", func(t *testing.T) {
		logger := WithRequestID("global-req-123")
		assert.NotNil(t, logger)
		logger.Info("测试")
	})

	t.Run("With全局函数", func(t *testing.T) {
		logger := With(String("component", "test"))
		assert.NotNil(t, logger)
		logger.Info("测试")
	})
}

// TestFieldConstructors 测试字段构造函数
func TestFieldConstructors(t *testing.T) {
	tests := []struct {
		name   string
		field  Field
		wantKey string
		wantValue interface{}
	}{
		{"String字段", String("name", "test"), "name", "test"},
		{"Int字段", Int("count", 42), "count", 42},
		{"Int64字段", Int64("big", 12345678901234), "big", int64(12345678901234)},
		{"Float64字段", Float64("rate", 3.14), "rate", 3.14},
		{"Bool字段true", Bool("flag", true), "flag", true},
		{"Bool字段false", Bool("flag", false), "flag", false},
		{"Any字段", Any("data", map[string]int{"a": 1}), "data", map[string]int{"a": 1}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.wantKey, tt.field.Key)
			assert.Equal(t, tt.wantValue, tt.field.Value)
		})
	}
}

// TestToZapFields 测试toZapFields函数
func TestToZapFields(t *testing.T) {
	fields := []Field{
		String("key1", "value1"),
		Int("key2", 123),
		Bool("key3", true),
	}

	result := toZapFields(fields...)

	// 验证返回的是交替的键值对
	assert.Equal(t, 6, len(result)) // 3个字段 = 6个元素
	assert.Equal(t, "key1", result[0])
	assert.Equal(t, "value1", result[1])
	assert.Equal(t, "key2", result[2])
	assert.Equal(t, 123, result[3])
	assert.Equal(t, "key3", result[4])
	assert.Equal(t, true, result[5])
}

// createTestLogger 创建测试用的Logger
func createTestLogger(buf *bytes.Buffer) *Logger {
	// 创建一个写入到buffer的Core
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(buf),
		zapcore.DebugLevel,
	)

	zapLogger := zap.New(core)
	return &Logger{logger: zapLogger.Sugar()}
}

// TestLoggerChaining 测试Logger链式调用
func TestLoggerChaining(t *testing.T) {
	var buf bytes.Buffer
	baseLogger := createTestLogger(&buf)

	// 链式创建带多个字段的Logger
	logger := baseLogger.
		With(String("service", "api")).
		With(String("version", "1.0")).
		WithRequestID("chain-test-123")

	logger.Info("链式调用测试")

	// 验证所有字段都存在
	output := buf.String()
	var logEntry map[string]interface{}
	err := json.Unmarshal([]byte(output), &logEntry)
	assert.NoError(t, err)

	assert.Equal(t, "api", logEntry["service"])
	assert.Equal(t, "1.0", logEntry["version"])
	assert.Equal(t, "chain-test-123", logEntry["requestId"])
	assert.Equal(t, "链式调用测试", logEntry["msg"])
}

// TestLogFormat 测试日志输出格式
func TestLogFormat(t *testing.T) {
	var buf bytes.Buffer
	testLogger := createTestLogger(&buf)

	testLogger.Info("格式测试",
		String("user", "admin"),
		Int("age", 30),
		Bool("active", true),
	)

	output := buf.String()
	var logEntry map[string]interface{}
	err := json.Unmarshal([]byte(output), &logEntry)
	assert.NoError(t, err)

	// 验证必须的字段存在
	requiredFields := []string{"time", "level", "msg"}
	for _, field := range requiredFields {
		_, exists := logEntry[field]
		assert.True(t, exists, "字段 %s 必须存在", field)
	}

	// 验证字段值
	assert.Equal(t, "格式测试", logEntry["msg"])
	assert.Equal(t, "info", logEntry["level"])
	assert.Equal(t, "admin", logEntry["user"])
	assert.Equal(t, float64(30), logEntry["age"]) // JSON数字解析为float64
	assert.Equal(t, true, logEntry["active"])
}
