package agentmanager

import (
	"errors"

	"gitee.com/openeuler/PilotGo-plugin-syscare/server/dao"
)

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
