package logger

import (
	"game_slots_vsn/pkg/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
)

//Level logger Level
type Level string

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
	Logger *zap.SugaredLogger
)

func InitLogger(c *config.LogConfig) {
	cores := make([]zapcore.Core, 0)

	if c.ConsoleEnable {
		lvl := zap.NewAtomicLevel()
		lvl.SetLevel(getZapLevel(Level(c.ConsoleLevel)))

		writeSync := zapcore.Lock(os.Stdout)
		core := zapcore.NewCore(getEncoder(c.JsonEnable), writeSync, lvl)
		cores = append(cores, core)
	}

	if c.FileEnable {
		lvl := zap.NewAtomicLevel()
		lvl.SetLevel(getZapLevel(Level(c.FileLevel)))
		writerSync := zapcore.AddSync(&lumberjack.Logger{
			Filename:   c.FileName,
			MaxSize:    c.MaxSize,
			Compress:   c.Compress,
			MaxBackups: c.MaxBackups,
			MaxAge:     c.MaxAges,
		})
		core := zapcore.NewCore(getEncoder(c.JsonEnable), writerSync, lvl)
		cores = append(cores, core)
	}
	multiCore := zapcore.NewTee(cores...)
	Logger = zap.New(multiCore,
		zap.AddStacktrace(zapcore.ErrorLevel),
		zap.AddCaller(),
		zap.AddCallerSkip(1),
	).Sugar()

	gawd, _ := os.Getwd()
	Logger.Info("pwd", gawd)
	defer Logger.Sync()
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
