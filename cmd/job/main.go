package main

import (
	"github.com/g-Halo/go-server/internal/job"
	"github.com/g-Halo/go-server/internal/job/conf"
	"github.com/g-Halo/go-server/pkg/logger"
	"os"
	"os/signal"
	"syscall"
)

func init() {
	// 初始化日志文件
	logger.InitLogger("./job.log", "debug")
}

func main() {
	if err := conf.Init(); err != nil {
		logger.Debug(err)
		return
	}

	j := job.New(conf.Conf)
	go j.ConsumerHandle()

	// Let's allow our queues to drain properly during shutdown.
	// We'll create a channel to listen for SIGINT (Ctrl+C) to signal
	// to our application to gracefully shutdown.
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			_ = j.Close()
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}