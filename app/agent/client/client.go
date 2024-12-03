/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo-plugin-syscare licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: zhanghan2021 <zhanghan@kylinos.cn>
 * Date: Mon Mar 18 17:10:03 2024 +0800
 */
package client

import (
	"sync"

	"gitee.com/openeuler/PilotGo-plugin-syscare/agent/config"
	uos "gitee.com/openeuler/PilotGo-plugin-syscare/utils/os"
)

type Client struct {
	AgentServer string //
	UUID        string
	lock        sync.Mutex

	IP              string
	MaxTaskNum      int    //最大任务并行数
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
		MaxTaskNum:      config.Config().Task.MaxTaskNum,
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
