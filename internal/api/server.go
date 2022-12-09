package api

import "C"
import (
	"fmt"
	"game_slots_vsn/pkg/logger"
	"os"
	"os/signal"
	"syscall"
	"time"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type WebConf struct {
	IP   string `yaml:"ip,omitempty"`
	Port int    `yaml:"port"`
}

var (
	c = &WebConf{}
)

func Run() {
	go startWeb()
	logger.Logger.Infof("web:%v start success.", c.Port)
	sighup()
}

func startWeb() {
	gin.SetMode(gin.ReleaseMode)
	e := gin.New()
	//e.Use(ginzap.Ginzap(logger.Logger.Desugar(), time.RFC3339, true))
	e.Use(
		ginzap.GinzapWithConfig(
			logger.Logger.Desugar(),
			&ginzap.Config{
				TimeFormat: time.RFC3339,
				UTC:        true,
				SkipPaths: []string{
					"/heartbeat",
					"/favicon.ico",
				},
			},
		),
	)
	e.Use(ginzap.RecoveryWithZap(logger.Logger.Desugar(), true))

	Init(e)

	if err := viper.UnmarshalKey("web", c); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "config modify fail.", err)
		os.Exit(0)
	}
	addr := fmt.Sprintf("%s:%d", c.IP, c.Port)
	if err := e.Run(addr); err != nil {
		logger.Logger.Info("service start fail.", err)
		os.Exit(0)
	}
}

func sighup() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			_, _ = fmt.Fprintln(os.Stderr, "server quit:", s)
			time.Sleep(time.Second * 1)
			return
		case syscall.SIGHUP:
			_, _ = fmt.Fprintln(os.Stderr, "server ", s)
			time.Sleep(time.Second * 1)
			return
		}
	}
}
