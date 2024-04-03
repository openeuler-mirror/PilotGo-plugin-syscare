package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"gitee.com/openeuler/PilotGo-plugin-syscare/server/agentmanager"
	"gitee.com/openeuler/PilotGo-plugin-syscare/server/config"
	"gitee.com/openeuler/PilotGo-plugin-syscare/server/dao"
	"gitee.com/openeuler/PilotGo-plugin-syscare/server/db"
	"gitee.com/openeuler/PilotGo-plugin-syscare/server/router"
	"gitee.com/openeuler/PilotGo-plugin-syscare/server/service"
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

	if err := db.MysqldbInit(config.Config().Mysql); err != nil {
		logger.Error("mysql db init failed, please check again: %s", err)
		os.Exit(-1)
	}
	logger.Info("Thanks to choose PilotGo-plugin-syscare-server!")

	if err := router.HttpServerInit(config.Config().HttpServer); err != nil {
		logger.Error("http server init failed, error: %s", err)
		os.Exit(-1)
	}

	// 启动所有功能模块服务
	if err := startServices(); err != nil {
		logger.Error("start services error: %s", err)
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

func startServices() error {
	db.MySQL().AutoMigrate(dao.WarmList{})

	if err := agentmanager.Init(); err != nil {
		return err
	}

	service.CreateTaskQueue()
	return nil
}
