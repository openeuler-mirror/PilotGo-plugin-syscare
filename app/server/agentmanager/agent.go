package agentmanager

import (
	"errors"

	"gitee.com/openeuler/PilotGo-plugin-syscare/server/dao"
	"gitee.com/openeuler/PilotGo/sdk/logger"
)

// 获取所有的agent
func GetAgentsWithStatus() ([]*Agent, error) {
	agents, err := globalAgentManager.getAgents()
	if err != nil {
		logger.Error("failed to read plugin info from db:%s", err.Error())
		return nil, err
	}
	result := []*Agent{}
	for _, a := range agents {
		agent_status, ok := AgentHeartbeatMap.Load(a.UUID)
		if !ok {
			logger.Error("Error get %v last heartbeat: %v", a.IP, err)
			continue
		}
		a.Connect = agent_status.(*AgentConnect).Connect
		a.LastHeartbeatTime = agent_status.(*AgentConnect).LastHeartbeatTime.Format("2006-01-02 15:04:05")
		result = append(result, a)
	}
	return result, nil
}

func AddAgent(ip string) (*Agent, error) {
	if ok, err := dao.IsExistIP(ip); err != nil {
		return &Agent{}, err
	} else if err == nil && ok {
		return &Agent{}, errors.New("该主机已添加")
	}

	a, err := globalAgentManager.addAgent(ip)
	if err != nil {
		return &Agent{}, err
	}
	return a, err
}
