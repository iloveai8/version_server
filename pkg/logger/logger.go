package logger

import (
	"game_slots_vsn/pkg/config"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"time"
)

//Level logger Level
type Level string

type Configuration struct {
	//  是否开启控制台日志输出
	ConsoleStdoutEnable bool
	//  控制台日志是否是 JSON 格式
	ConsoleStdoutIsJSONFormat bool
	// Level 控制台日志等级
	ConsoleStdoutLevel Level

	// 是否开启控制台日志输出
	FileStdoutEnable bool
	// 控制台日志是否是 JSON 格式
	FileStdoutIsJSONFormat bool
	// Level 控制台日志等级
	FileStdoutLevel Level
	// 写入文件位置
	FileStdoutFileLocation string
	FileStdoutLogMaxSize   int
	FileStdoutCompress     bool
	FileStdoutLogMaxAge    int
}

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
	Logger  *zap.Logger
	//Logger *zap.SugaredLogger
)

//
//// InitLogger 初始化Logger
//func InitLogger(c *config.LogConfig) {
//	encoder := doGetEncoder()
//	writeSync := doGetLogWriter(c.FileName, c.MaxSize, c.MaxBackups, c.MaxAges)
//
//	var l = new(zapcore.Level)
//	err := l.UnmarshalText([]byte(c.Level))
//	if err != nil {
//		return
//	}
//	core := zapcore.NewCore(encoder, writeSync, l)
//	Logger = zap.New(core, zap.AddCaller())
//}
//
//func doGetEncoder() zapcore.Encoder {
//	encoderConfig := zap.NewProductionEncoderConfig()
//	encoderConfig.TimeKey = "time"
//	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
//	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
//	return zapcore.NewJSONEncoder(encoderConfig)
//}
//
//func doGetLogWriter(fileName string, maxSize, maxBackup, maxAge int) zapcore.WriteSyncer {
//	lumberJackLogger := &lumberjack.Logger{
//		Filename:   fileName,
//		MaxSize:    maxSize,
//		MaxBackups: maxBackup,
//		MaxAge:     maxAge,
//	}
//	return zapcore.AddSync(lumberJackLogger)
//}

func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		c.Next()
		costTime := time.Since(start)
		Logger.Info(path,
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
			zap.Duration("costTime", costTime),
		)
	}
}
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
	)
}

func getEncoder(JsonEnable bool) zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.LowercaseLevelEncoder
	if JsonEnable {
		return zapcore.NewJSONEncoder(encoderConfig)
	} else {
		return zapcore.NewConsoleEncoder(encoderConfig)
	}
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
