package internal

import (
	`game_slots_vsn/internal/api`
	`os`
	`os/signal`
	`syscall`
	`time`
)

func Run() {
	//dao.Init()
	go api.ServerStart()

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
