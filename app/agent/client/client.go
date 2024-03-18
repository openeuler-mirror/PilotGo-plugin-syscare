package client

import (
	"sync"

	uos "gitee.com/openeuler/PilotGo-plugin-syscare/utils/os"
)

type Client struct {
	AgentServer string //
	UUID        string
	lock        sync.Mutex

	IP              string
	Platform        string //系统平台
	PlatformVersion string //系统版本
	OsVersion       string //可读性良好的OS具体版本
	KernelVersion   string //内核版本
	KernelArch      string //内核支持架构
	Uptime          string //系统最新启动时间
	CpuModelName    string //cpu架构名称
	CpuNum          int    //cpu核数
}

var GlobalClient *Client

func DefaultClient() (*Client, error) {
	cpuInfo, err := uos.GetCPUInfo()
	if err != nil {
		return GlobalClient, err
	}
	hostInfo, err := uos.GetHostInfo()
	if err != nil {
		return GlobalClient, err
	}
	GlobalClient = &Client{
		IP:              hostInfo.IP,
		Platform:        hostInfo.Platform,
		PlatformVersion: hostInfo.PlatformVersion,
		OsVersion:       hostInfo.OsVersion,
		KernelVersion:   hostInfo.KernelVersion,
		KernelArch:      hostInfo.KernelArch,
		Uptime:          hostInfo.Uptime,
		CpuModelName:    cpuInfo.ModelName,
		CpuNum:          cpuInfo.CpuNum,
	}

	return GlobalClient, nil
}

func GetClient() *Client {
	return GlobalClient
}

func (client *Client) Server() string {
	return client.AgentServer
}

// is client bind syscare server?
func (c *Client) IsBind() bool {
	c.lock.Lock()
	defer c.lock.Unlock()

	return c.AgentServer != ""
}
