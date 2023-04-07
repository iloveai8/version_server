package logger

import (
	"game_slots_vsn/pkg/setting"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"time"
)

const (
	//DebugLevel has verbose message
	DebugLevel Level = "debug"
	//InfoLevel is default log level
	InfoLevel Level = "info"
	//WarnLevel is for logging messages about possible issues
	WarnLevel Level = "warn"
	//ErrorLevel is for logging errors
	ErrorLevel Level = "error"
	//FatalLevel is for logging fatal messages. The system shutdown after logging the message.
	FatalLevel Level = "fatal"
)

var (
	sugar *zap.SugaredLogger
)

type (
	Level string
)

func Setup(c *setting.LoggerSetting) {
	cores := make([]zapcore.Core, 0)
	if c.ConsoleEnable {
		lvl := zap.NewAtomicLevel()
		lvl.SetLevel(getZapLevel(Level(c.ConsoleLevel)))

		writeSync := zapcore.Lock(os.Stdout)
		core := zapcore.NewCore(getEncoder(c.FileJsonEnable), writeSync, lvl)
		cores = append(cores, core)
	}

	if c.FileEnable {
		lvl := zap.NewAtomicLevel()
		lvl.SetLevel(getZapLevel(Level(c.FileLevel)))
		writerSync := zapcore.AddSync(&lumberjack.Logger{
			Filename:   c.FileName,
			MaxSize:    c.FileMaxSize,
			Compress:   c.FileCompress,
			MaxBackups: c.FileMaxBackups,
			MaxAge:     c.FileMaxAges,
		})
		core := zapcore.NewCore(getEncoder(c.FileJsonEnable), writerSync, lvl)
		cores = append(cores, core)
	}
	multiCore := zapcore.NewTee(cores...)
	sugar = zap.New(multiCore,
		zap.AddStacktrace(zapcore.ErrorLevel),
		zap.AddCaller(),
		zap.AddCallerSkip(1),
	).Sugar()
	defer sugar.Sync()
}

func GinZap() gin.HandlerFunc {
	return ginzap.Ginzap(sugar.Desugar(), time.RFC3339, true)
}

func GinZapWithSkipPaths() gin.HandlerFunc {
	return ginzap.GinzapWithConfig(sugar.Desugar(),
		&ginzap.Config{
			TimeFormat: time.RFC3339,
			UTC:        true,
			SkipPaths: []string{
				"/heartbeat",
				"/favicon.ico",
			},
		},
	)
}

func RecoverZap() gin.HandlerFunc {
	return ginzap.RecoveryWithZap(sugar.Desugar(), true)
}

func getEncoder(JsonEnable bool) zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.LowercaseLevelEncoder

	var encoder zapcore.Encoder
	if JsonEnable {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}
	return encoder
}

func getZapLevel(level Level) zapcore.Level {
	switch level {
	case InfoLevel:
		return zapcore.InfoLevel
	case WarnLevel:
		return zapcore.WarnLevel
	case DebugLevel:
		return zapcore.DebugLevel
	case ErrorLevel:
		return zapcore.ErrorLevel
	case FatalLevel:
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}

func DebugF(template string, args ...interface{}) {
	sugar.Debugf(template, args...)
}

func Debug(args ...interface{}) {
	sugar.Debug(args...)
}

func InfoF(template string, args ...interface{}) {
	sugar.Infof(template, args...)
}
func Info(args ...interface{}) {
	sugar.Info(args...)
}

func WarnF(template string, args ...interface{}) {
	sugar.Warnf(template, args...)
}

func Warn(args ...interface{}) {
	sugar.Warn(args...)
}

func ErrorF(template string, args ...interface{}) {
	sugar.Errorf(template, args...)
}

func Error(args ...interface{}) {
	sugar.Error(args...)
}

func FatalF(template string, args ...interface{}) {
	sugar.Fatalf(template, args...)
}

func Fatal(args ...interface{}) {
	sugar.Fatal(args...)
}
