package config

import (
	"flag"

	"gitee.com/openeuler/PilotGo/sdk/logger"
	"gitee.com/openeuler/PilotGo/sdk/utils/config"
)

type HttpServer struct {
	Addr string `yaml:"addr"`
}
type AgentServer struct {
	Port string `yaml:"port"`
}
type MysqlDBInfo struct {
	HostName string `yaml:"host"`
	Port     int    `yaml:"port"`
	UserName string `yaml:"user"`
	Password string `yaml:"password"`
	DataBase string `yaml:"database"`
}
type Storage struct {
	Path string `yaml:"path"`
}
type ServerConfig struct {
	HttpServer  *HttpServer     `yaml:"http_server"`
	AgentServer *AgentServer    `yaml:"agent_server"`
	Logopts     *logger.LogOpts `yaml:"log"`
	Mysql       *MysqlDBInfo    `yaml:"mysql"`
	Storage     *Storage        `yaml:"storage"`
}

var config_file string
var global_config ServerConfig

func Init() error {
	flag.StringVar(&config_file, "conf", "./config_server.yaml", "pilotgo-plugin-syscare configuration file")
	flag.Parse()
	return config.Load(config_file, &global_config)
}

func Config() *ServerConfig {
	return &global_config
}
