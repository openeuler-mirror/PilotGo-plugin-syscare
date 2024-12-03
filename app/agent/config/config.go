/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo-plugin-syscare licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: zhanghan2021 <zhanghan@kylinos.cn>
 * Date: Tue Mar 5 15:43:27 2024 +0800
 */
package config

import (
	"flag"

	"gitee.com/openeuler/PilotGo/sdk/logger"
	"gitee.com/openeuler/PilotGo/sdk/utils/config"
)

type Server struct {
	Port string `yaml:"port"`
}

type Task struct {
	MaxTaskNum int `yaml:"max_task_num"`
}
type Storage struct {
	Path string `yaml:"path"`
	Work string `yaml:"work"`
}
type AgentConfig struct {
	Server  *Server         `yaml:"server"`
	Task    *Task           `yaml:"task"`
	Logopts *logger.LogOpts `yaml:"log"`
	Storage Storage         `yaml:"storage"`
}

var config_file string
var global_config AgentConfig

func Init() error {
	flag.StringVar(&config_file, "conf", "./config_agent.yaml", "pilotgo-plugin-syscare configuration file")
	flag.Parse()
	return config.Load(config_file, &global_config)
}

func Config() *AgentConfig {
	return &global_config
}
