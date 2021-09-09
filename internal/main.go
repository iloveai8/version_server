package internal

import (
	`fmt`
	`game_slots_vsn/internal/api`
	`os`
	`os/signal`
	`syscall`
	`time`
)

//Run 程序入口
func Run() {
	go api.ServerStart()
	Signal()
}

//Signal 接受处理系统信号
func Signal() {
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
