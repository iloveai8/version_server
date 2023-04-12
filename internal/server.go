package internal

import (
	"fmt"
	"game_slots_vsn/internal/routers"
	"game_slots_vsn/pkg/logger"
	"game_slots_vsn/pkg/server"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Run() {
	go func() {
		addr := fmt.Sprintf("%s:%d", server.Server.S.IP, server.Server.S.Port)
		engine := routers.Register()
		if err := engine.Run(addr); err != nil {
			logger.ErrorF("Listen addr:%s err:%v\n", err)
		}
	}()
	shutdown()
}

func Stop() {
}

func shutdown() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			time.Sleep(time.Second * 1)
			return
		case syscall.SIGHUP:
			time.Sleep(time.Second * 1)
			return
		}
	}
}
