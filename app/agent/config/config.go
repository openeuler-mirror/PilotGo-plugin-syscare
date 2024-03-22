package config

import (
	"flag"

	"gitee.com/openeuler/PilotGo/sdk/logger"
	"gitee.com/openeuler/PilotGo/sdk/utils/config"
)

type Server struct {
	Port string `yaml:"port"`
}
type Storage struct {
	Path string `yaml:"path"`
	Work string `yaml:"work"`
}
type AgentConfig struct {
	Server  *Server         `yaml:"server"`
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
