package agentmanager

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"gitee.com/openeuler/PilotGo-plugin-syscare/server/config"
	"gitee.com/openeuler/PilotGo-plugin-syscare/server/dao"
	"gitee.com/openeuler/PilotGo-plugin-syscare/server/db"
	"gitee.com/openeuler/PilotGo/sdk/logger"
	"gitee.com/openeuler/PilotGo/sdk/utils/httputils"
	"github.com/google/uuid"
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
	MaxTaskNum        int    `json:"maxTaskNum"`
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

// 从DB中恢复agent信息
func (am *AgentManager) recovery() error {
	agents, err := dao.QueryAgents()
	if err != nil {
		logger.Error("failed to recovery agent info from db")
		return nil
	}

	for _, a := range agents {
		logger.Debug("recovery agent: %s %s", a.UUID, a.IP)
		err := am.updateAgent(a.UUID, a.IP)
		if err != nil {
			logger.Error("failed to update agent info: %s", a.IP)
			am.Lock()
			am.Agents = append(am.Agents, &Agent{
				UUID:            a.UUID,
				IP:              a.IP,
				MaxTaskNum:      a.MaxTaskNum,
				Platform:        a.Platform,
				PlatformVersion: a.PlatformVersion,
				OsVersion:       a.OsVersion,
				KernelVersion:   a.KernelVersion,
				KernelArch:      a.KernelArch,
				Uptime:          a.Uptime,
				CpuModelName:    a.CpuModelName,
				CpuNum:          a.CpuNum,
			})
			am.Unlock()
			agentHeartbeatUpdateFail(a.IP, a.UUID)
		}
	}

	logger.Debug("finish recovery")
	return nil
}

// 获取所有agent
func (am *AgentManager) getAgents() ([]*Agent, error) {
	result := []*Agent{}
	am.Lock()
	for _, v := range am.Agents {
		a := v.clone()
		result = append(result, a)
	}
	am.Unlock()

	return result, nil
}

// 根据ip查询最新的agent信息，更新到指定的uuid记录当中
func (am *AgentManager) updateAgent(uuid string, ip string) error {
	info, err := getAgentInfo(ip)
	if err != nil {
		logger.Error("failed to get agent info:%s", err.Error())
		return err
	}
	info.UUID = uuid
	err = bindServerAndUUID(ip, uuid)
	if err != nil {
		logger.Error("update bind agent error:%s", err.Error())
		return err
	}

	a := &Agent{
		UUID:            uuid,
		IP:              ip,
		MaxTaskNum:      info.MaxTaskNum,
		Platform:        info.Platform,
		PlatformVersion: info.PlatformVersion,
		OsVersion:       info.OsVersion,
		KernelVersion:   info.KernelVersion,
		KernelArch:      info.KernelArch,
		Uptime:          info.Uptime,
		CpuModelName:    info.CpuModelName,
		CpuNum:          info.CpuNum,
	}
	if err := dao.UpdateAgentInfo(toAgentDao(a)); err != nil {
		return err
	}
	agentHeartbeatUpdateSuccess(ip, a.UUID)

	am.Lock()
	am.Agents = append(am.Agents, a)
	am.Unlock()

	return nil
}
func (am *AgentManager) deleteAgent(ip string) error {
	err := unbind(ip)
	if err != nil {
		return err
	}

	if err := dao.DeleteAgent(ip); err != nil {
		logger.Error("failed to delete agent info:%s", err.Error())
		return err
	}

	am.Lock()
	index := 0
	uuid := ""
	for i, v := range am.Agents {
		if v.IP == ip {
			index = i
			uuid = v.UUID
			break
		}
	}

	if index == 0 {
		am.Agents = am.Agents[1:]
	} else if index == len(am.Agents)-1 {
		am.Agents = am.Agents[:index]
	} else {
		am.Agents = append(am.Agents[:index], am.Agents[index+1:]...)
	}

	AgentHeartbeatMap.Delete(uuid)
	am.Unlock()
	return nil
}
func (am *AgentManager) addAgent(ip string) (*Agent, error) {
	a, err := getAgentInfo(ip)
	if err != nil {
		logger.Error("first get agent info error: %s", err.Error())
		return nil, err
	}
	if a.UUID == "" {
		a.UUID = uuid.New().String()
	} else {
		return a, errors.New("the added host has been bound")
	}

	err = bindServerAndUUID(ip, a.UUID)
	if err != nil {
		logger.Error("first bind agent error:%s", err.Error())
		return nil, err
	}

	if err := dao.SaveAgent(toAgentDao(a)); err != nil {
		return nil, err
	}

	agentHeartbeatUpdateSuccess(ip, a.UUID)

	am.Lock()
	am.Agents = append(am.Agents, a)
	am.Unlock()

	return a, nil
}

// 获取到agent的基本信息
func getAgentInfo(ip string) (*Agent, error) {
	url := "http://" + ip + ":" + config.Config().AgentServer.Port + "/plugin_agent_manage/info"
	logger.Debug("agent url is: %s", url)

	client := &http.Client{
		Timeout: 3 * time.Second,
		Transport: &http.Transport{
			DisableKeepAlives: true,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}

	AgentInfo := &Agent{}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error("error reading response body: %v", err.Error())
		return nil, err
	}
	err = json.Unmarshal(body, AgentInfo)
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

func toAgentDao(a *Agent) *dao.Agents {
	return &dao.Agents{
		UUID:            a.UUID,
		IP:              a.IP,
		MaxTaskNum:      a.MaxTaskNum,
		Platform:        a.Platform,
		PlatformVersion: a.PlatformVersion,
		OsVersion:       a.OsVersion,
		KernelVersion:   a.KernelVersion,
		KernelArch:      a.KernelArch,
		Uptime:          a.Uptime,
		CpuModelName:    a.CpuModelName,
		CpuNum:          a.CpuNum,
	}
}

func (a *Agent) clone() *Agent {
	result := &Agent{
		UUID:            a.UUID,
		IP:              a.IP,
		MaxTaskNum:      a.MaxTaskNum,
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
