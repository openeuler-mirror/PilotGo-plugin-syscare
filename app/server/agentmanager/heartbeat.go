package agentmanager

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"gitee.com/openeuler/PilotGo-plugin-syscare/server/config"
	"gitee.com/openeuler/PilotGo/sdk/logger"
	"gitee.com/openeuler/PilotGo/sdk/utils/httputils"
)

const (
	heartbeatInterval = 30 * time.Second
)

var AgentHeartbeatMap = &sync.Map{}

type AgentConnect struct {
	IP                string    `json:"ip"`
	Connect           bool      `json:"connect"`
	LastHeartbeatTime time.Time `json:"lastheartbeat"`
}

func checkAgentsHeartbeat() {
	for {
		broadcastHeartbeat() // 广播发送心跳
		checkAndRebind()
		time.Sleep(heartbeatInterval)
	}
}

func checkAndRebind() {
	agents, err := globalAgentManager.getAgents()
	if err != nil {
		logger.Error("get agents failed: %v", err.Error())
	}

	for _, a := range agents {
		agent_status, ok := AgentHeartbeatMap.Load(a.UUID)
		if !ok {
			logger.Error("Error get %v last heartbeat: %v", a.IP, err)
			continue
		}

		if !agent_status.(*AgentConnect).Connect || time.Since(agent_status.(*AgentConnect).LastHeartbeatTime) > heartbeatInterval+1*time.Second {
			err := bindServerAndUUID(a.IP, a.UUID)
			if err != nil {
				logger.Error("rebind agent and syscare server failed: %v", err.Error())
				value := &AgentConnect{
					IP:                a.IP,
					Connect:           false,
					LastHeartbeatTime: agent_status.(*AgentConnect).LastHeartbeatTime,
				}
				AgentHeartbeatMap.Store(a.UUID, value)
			} else {
				agentHeartbeatUpdateSuccess(a.IP, a.UUID)
			}
		}
	}
}

// 获取到agent的基本信息
func broadcastHeartbeat() {
	agents, err := globalAgentManager.getAgents()
	if err != nil {
		logger.Error("get agents failed: %v", err.Error())
	}
	for _, a := range agents {
		url := "http://" + a.IP + ":" + config.Config().AgentServer.Port + "/plugin_agent_manage/heartbeat"
		resp, err := httputils.Get(url, nil)
		if err != nil {
			logger.Error("get agent info error: %s", err.Error())
			continue
		}
		res := &struct {
			Code int    `json:"code"`
			Data string `json:"data"`
		}{}
		if err := json.Unmarshal(resp.Body, res); err != nil {
			logger.Error("parse struct failed: %v", err.Error())
			continue
		}

		if res.Code != http.StatusOK {
			agent_status, ok := AgentHeartbeatMap.Load(a.UUID)
			if !ok {
				logger.Error("Error get %v last heartbeat: %v", a.IP, err)
				continue
			}
			value := &AgentConnect{
				IP:                a.IP,
				Connect:           false,
				LastHeartbeatTime: agent_status.(*AgentConnect).LastHeartbeatTime,
			}
			AgentHeartbeatMap.Store(a.UUID, value)
		} else {
			agentHeartbeatUpdateSuccess(a.IP, a.UUID)
		}
	}

}

func agentHeartbeatUpdateSuccess(ip, uuid string) {
	value := &AgentConnect{
		IP:                ip,
		Connect:           true,
		LastHeartbeatTime: time.Now(),
	}
	AgentHeartbeatMap.Store(uuid, value)
}
func agentHeartbeatUpdateFail(ip, uuid string) {
	value := &AgentConnect{
		IP:                ip,
		Connect:           false,
		LastHeartbeatTime: time.Now(),
	}
	AgentHeartbeatMap.Store(uuid, value)
}
