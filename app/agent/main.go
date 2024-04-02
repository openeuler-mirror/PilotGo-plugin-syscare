package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"gitee.com/openeuler/PilotGo-plugin-syscare/agent/config"
	"gitee.com/openeuler/PilotGo-plugin-syscare/agent/router"
	"gitee.com/openeuler/PilotGo-plugin-syscare/agent/service"
	"gitee.com/openeuler/PilotGo/sdk/logger"
)

func main() {
	if err := config.Init(); err != nil {
		fmt.Println("failed to load configure, exit..", err)
		os.Exit(-1)
	}
	if err := logger.Init(config.Config().Logopts); err != nil {
		fmt.Printf("logger init failed, please check the config file: %s", err)
		os.Exit(-1)
	}
	if err := service.MakeDir(config.Config().Storage.Path); err != nil {
		fmt.Printf("storage path init failed, error: %s", err)
		os.Exit(-1)
	}
	logger.Info("Thanks to choose PilotGo-plugin-syscare-agent!")

	if err := router.HttpServerInit(); err != nil {
		logger.Error("http server init failed, error: %s", err)
		os.Exit(-1)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	for {
		s := <-c
		switch s {
		case syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
			logger.Info("signal interrupted: %s", s.String())
			goto EXIT
		default:
			logger.Info("unknown signal: %s", s.String())
		}
	}

EXIT:
	logger.Info("exit system, bye~")
}
