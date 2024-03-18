package agentmanager

import (
	"sync"

	"gitee.com/openeuler/PilotGo-plugin-syscare/server/dao"
	"gitee.com/openeuler/PilotGo-plugin-syscare/server/db"
)

func Init() error {
	if err := db.MySQL().AutoMigrate(&dao.Agents{}); err != nil {
		return err
	}

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
