package main

import (
	"fmt"
	"os"

	"gitee.com/openeuler/PilotGo-plugin-syscare/server/config"
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
	logger.Info("Thanks to choose PilotGo-plugin-syscare-server!")
}
