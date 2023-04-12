package logger

import (
	"fmt"
	"game_slots_vsn/pkg/consts"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
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

type Level string

type setting struct {
	ConsoleEnable bool   `yaml:"consoleEnable"`
	ConsoleLevel  string `yaml:"consoleLevel"`

	FileEnable     bool   `yaml:"fileEnable"`
	FileName       string `yaml:"fileName"`
	FileLevel      string `yaml:"fileLevel"`
	FileMaxSize    int    `yaml:"maxSize"`
	FileMaxBackups int    `yaml:"maxBackups"`
	FileMaxAges    int    `yaml:"maxAges"`
	FileCompress   bool   `yaml:"compress"`
	FileJsonEnable bool   `yaml:"jsonEnable"`
}

type logger struct {
	sugar *zap.SugaredLogger
	S     *setting
}

var Logger = &logger{}

func SetUp() {
	ls := &setting{}
	err := viper.UnmarshalKey(consts.ConfigLogger, ls)
	if err != nil {
		panic(err)
	}
	Logger.setup(ls)
}

func (l *logger) setup(ls *setting) {
	fmt.Printf("logger setting:%v\n", *ls)
	cores := make([]zapcore.Core, 0)
	if ls.ConsoleEnable {
		lvl := zap.NewAtomicLevel()
		lvl.SetLevel(getZapLevel(Level(ls.ConsoleLevel)))

		writeSync := zapcore.Lock(os.Stdout)
		core := zapcore.NewCore(getEncoder(ls.FileJsonEnable), writeSync, lvl)
		cores = append(cores, core)
	}

	if ls.FileEnable {
		lvl := zap.NewAtomicLevel()
		lvl.SetLevel(getZapLevel(Level(ls.FileLevel)))
		writerSync := zapcore.AddSync(&lumberjack.Logger{
			Filename:   ls.FileName,
			MaxSize:    ls.FileMaxSize,
			Compress:   ls.FileCompress,
			MaxBackups: ls.FileMaxBackups,
			MaxAge:     ls.FileMaxAges,
		})
		core := zapcore.NewCore(getEncoder(ls.FileJsonEnable), writerSync, lvl)
		cores = append(cores, core)
	}
	multiCore := zapcore.NewTee(cores...)
	sugar := zap.New(multiCore,
		zap.AddStacktrace(zapcore.ErrorLevel),
		zap.AddCaller(),
		zap.AddCallerSkip(1),
	).Sugar()
	defer sugar.Sync()
	Logger.S = ls
	Logger.sugar = sugar
}

func GinZap() gin.HandlerFunc {
	return ginzap.Ginzap(Logger.sugar.Desugar(), time.RFC3339, true)
}

func GinZapWithSkipPaths() gin.HandlerFunc {
	return ginzap.GinzapWithConfig(Logger.sugar.Desugar(),
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
	return ginzap.RecoveryWithZap(Logger.sugar.Desugar(), true)
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
	Logger.sugar.Debugf(template, args...)
}

func Debug(args ...interface{}) {
	Logger.sugar.Debug(args...)
}

func InfoF(template string, args ...interface{}) {
	Logger.sugar.Infof(template, args...)
}
func Info(args ...interface{}) {
	Logger.sugar.Info(args...)
}

func WarnF(template string, args ...interface{}) {
	Logger.sugar.Warnf(template, args...)
}

func Warn(args ...interface{}) {
	Logger.sugar.Warn(args...)
}

func ErrorF(template string, args ...interface{}) {
	Logger.sugar.Errorf(template, args...)
}

func Error(args ...interface{}) {
	Logger.sugar.Error(args...)
}

func FatalF(template string, args ...interface{}) {
	Logger.sugar.Fatalf(template, args...)
}

func Fatal(args ...interface{}) {
	Logger.sugar.Fatal(args...)
}
