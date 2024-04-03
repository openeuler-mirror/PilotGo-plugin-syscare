package service

import (
	"encoding/json"
	"errors"
	"net/http"

	"gitee.com/openeuler/PilotGo-plugin-syscare/server/agentmanager"
	"gitee.com/openeuler/PilotGo-plugin-syscare/server/config"
	"gitee.com/openeuler/PilotGo-plugin-syscare/server/dao"
	"gitee.com/openeuler/PilotGo/sdk/response"
	"gitee.com/openeuler/PilotGo/sdk/utils/httputils"

	"gitee.com/openeuler/PilotGo/sdk/logger"
)

func QueryAgents() ([]*agentmanager.Agent, error) {
	agents, err := dao.QueryAgents()
	if err != nil {
		return nil, err
	}
	data := []*agentmanager.Agent{}
	for _, a := range agents {
		agent_status, ok := agentmanager.AgentHeartbeatMap.Load(a.UUID)
		if !ok {
			logger.Error("Error get %v last heartbeat: %v", a.IP, err)
		}

		agent := &agentmanager.Agent{
			ID:                a.ID,
			UUID:              a.UUID,
			IP:                a.IP,
			MaxTaskNum:        a.MaxTaskNum,
			Platform:          a.Platform,
			PlatformVersion:   a.PlatformVersion,
			OsVersion:         a.OsVersion,
			KernelVersion:     a.KernelVersion,
			KernelArch:        a.KernelArch,
			Uptime:            a.Uptime,
			CpuModelName:      a.CpuModelName,
			CpuNum:            a.CpuNum,
			Connect:           agent_status.(*agentmanager.AgentConnect).Connect,
			LastHeartbeatTime: agent_status.(*agentmanager.AgentConnect).LastHeartbeatTime.Format("2006-01-02 15:04:05"),
		}
		data = append(data, agent)
	}
	return data, nil

}
func SearchAgents(search string, query *response.PaginationQ) ([]*agentmanager.Agent, int, error) {
	agents, total, err := dao.SearchAgents(search, query)
	if err != nil {
		return nil, 0, err
	}
	data := []*agentmanager.Agent{}
	for _, a := range agents {
		agent_status, ok := agentmanager.AgentHeartbeatMap.Load(a.UUID)
		if !ok {
			logger.Error("Error get %v last heartbeat: %v", a.IP, err)
		}
		agent := &agentmanager.Agent{
			ID:                a.ID,
			UUID:              a.UUID,
			IP:                a.IP,
			MaxTaskNum:        a.MaxTaskNum,
			Platform:          a.Platform,
			PlatformVersion:   a.PlatformVersion,
			OsVersion:         a.OsVersion,
			KernelVersion:     a.KernelVersion,
			KernelArch:        a.KernelArch,
			Uptime:            a.Uptime,
			CpuModelName:      a.CpuModelName,
			CpuNum:            a.CpuNum,
			Connect:           agent_status.(*agentmanager.AgentConnect).Connect,
			LastHeartbeatTime: agent_status.(*agentmanager.AgentConnect).LastHeartbeatTime.Format("2006-01-02 15:04:05"),
		}
		data = append(data, agent)
	}
	return data, int(total), nil

}
func GetBuildEnv(ip string) (interface{}, error) {
	url := "http://" + ip + ":" + config.Config().AgentServer.Port + "/plugin_agent_manage/buildEnv"
	resp, err := httputils.Get(url, nil)
	if err != nil {
		return nil, err
	}

	rpms := &struct {
		Code int `json:"code"`
		Data []struct {
			Version string `json:"version"`
			Rpms    struct {
				SrcRpm    string `json:"srcRpm"`
				DebugInfo string `json:"debugInfo"`
			} `json:"rpm"`
		} `json:"data"`
		Msg string `json:"msg"`
	}{}
	err = json.Unmarshal(resp.Body, rpms)
	if err != nil {
		logger.Error("unmarshal get agent info error: %s", err.Error())
		return nil, err
	}
	if rpms.Code != http.StatusOK {
		return nil, errors.New(rpms.Msg)
	}
	return rpms.Data, nil
}
