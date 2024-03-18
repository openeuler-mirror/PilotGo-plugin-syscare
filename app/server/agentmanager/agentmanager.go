package agentmanager

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"sync"

	"gitee.com/openeuler/PilotGo-plugin-syscare/server/config"
	"gitee.com/openeuler/PilotGo-plugin-syscare/server/dao"
	"gitee.com/openeuler/PilotGo-plugin-syscare/server/db"
	"gitee.com/openeuler/PilotGo/sdk/logger"
	"gitee.com/openeuler/PilotGo/sdk/utils/httputils"
)

func Init() error {
	if err := db.MySQL().AutoMigrate(&dao.Agents{}); err != nil {
		return err
	}
	if err := globalAgentManager.recovery(); err != nil {
		return err
	}

	// 检查所有agent状态
	go checkAgentsHeartbeat()
	return nil
}

type Agent struct {
	ID                int    `json:"id"`
	UUID              string `json:"uuid"`
	IP                string `json:"ip"`
	Platform          string `json:"platform"`
	PlatformVersion   string `json:"platformVersion"`
	OsVersion         string `json:"osVersion"`
	KernelVersion     string `json:"kernelVersion"`
	KernelArch        string `json:"kernelArch"`
	Uptime            string `json:"upTime"`
	CpuModelName      string `json:"cpuModelName"`
	CpuNum            int    `json:"cpuNum"`
	Connect           bool   `json:"connect"`
	LastHeartbeatTime string `json:"lastheartbeat"`
}

type AgentManager struct {
	sync.Mutex
	Agents []*Agent
}

var globalAgentManager = &AgentManager{
	Mutex:  sync.Mutex{},
	Agents: []*Agent{},
}

// 获取到agent的基本信息
func getAgentInfo(ip string) (*Agent, error) {
	url := "http://" + ip + ":" + config.Config().AgentServer.Port + "/plugin_agent_manage/info"
	logger.Debug("agent url is: %s", url)

	resp, err := httputils.Get(url, nil)
	if err != nil {
		return nil, err
	}

	AgentInfo := &Agent{}
	err = json.Unmarshal(resp.Body, AgentInfo)
	if err != nil {
		logger.Error("unmarshal get agent info error: %s", err.Error())
		return nil, err
	}

	return AgentInfo, nil
}

// syscare server 和agent 绑定
func bindServerAndUUID(ip, uuid string) error {
	port := strings.Split(config.Config().HttpServer.Addr, ":")[1]

	url := "http://" + ip + ":" + config.Config().AgentServer.Port + "/plugin_agent_manage/bind?port=" + port + "&uuid=" + uuid
	logger.Debug("agent url is: %s", url)

	resp, err := httputils.Put(url, &httputils.Params{})
	if err != nil {
		return err
	}

	d := &struct {
		Code    int    `json:"code"`
		Message string `json:"msg"`
	}{}
	err = json.Unmarshal(resp.Body, d)
	if err != nil {
		logger.Error("unmarshal bind agent error: %s", err.Error())
		return err
	}
	if d.Code != http.StatusOK {
		return errors.New(d.Message)
	}
	return nil
}

// syscare server 和agent 解绑
func unbind(ip string) error {
	url := "http://" + ip + ":" + config.Config().AgentServer.Port + "/plugin_agent_manage/unbind"
	logger.Debug("agent url is: %s", url)

	resp, err := httputils.Put(url, &httputils.Params{})
	if err != nil {
		logger.Error("unbind agent error:%s", err.Error())
		return err
	}

	d := &struct {
		Code    int    `json:"code"`
		Message string `json:"msg"`
	}{}
	err = json.Unmarshal(resp.Body, d)
	if err != nil {
		logger.Error("unmarshal bind agent error:%s", err.Error())
		return err
	}
	if d.Code != http.StatusOK {
		return errors.New(d.Message)
	}
	return nil
}

func (a *Agent) clone() *Agent {
	result := &Agent{
		UUID:            a.UUID,
		IP:              a.IP,
		Platform:        a.Platform,
		PlatformVersion: a.PlatformVersion,
		OsVersion:       a.OsVersion,
		KernelVersion:   a.KernelVersion,
		KernelArch:      a.KernelArch,
		Uptime:          a.Uptime,
		CpuModelName:    a.CpuModelName,
		CpuNum:          a.CpuNum,
	}

	return result
}
