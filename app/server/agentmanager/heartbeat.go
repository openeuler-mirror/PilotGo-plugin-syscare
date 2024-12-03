/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo-plugin-syscare licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: zhanghan2021 <zhanghan@kylinos.cn>
 * Date: Mon Mar 18 17:27:35 2024 +0800
 */
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
		time.Sleep(heartbeatInterval)
	}
}

// 获取到agent的基本信息
func broadcastHeartbeat() {
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

		url := "http://" + a.IP + ":" + config.Config().AgentServer.Port + "/plugin_agent_manage/heartbeat"
		resp, err := httputils.Get(url, nil)
		if err != nil {
			logger.Error("get agent info error: %s", err.Error())
			agentHeartbeatUpdate(a.IP, a.UUID, agent_status)
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

		if res.Code == http.StatusOK && res.Data != "" {
			agentHeartbeatUpdateSuccess(a.IP, a.UUID)
		} else if res.Code != http.StatusOK || res.Data == "" {
			err := bindServerAndUUID(a.IP, a.UUID)
			if err != nil {
				logger.Error("rebind agent and syscare server failed: %v", err.Error())
				agentHeartbeatUpdate(a.IP, a.UUID, agent_status)
			} else {
				agentHeartbeatUpdateSuccess(a.IP, a.UUID)
			}
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
func agentHeartbeatUpdate(ip, uuid string, agent_status any) {
	value := &AgentConnect{
		IP:                ip,
		Connect:           false,
		LastHeartbeatTime: agent_status.(*AgentConnect).LastHeartbeatTime,
	}
	AgentHeartbeatMap.Store(uuid, value)
}
