/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo-plugin-syscare licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: zhanghan2021 <zhanghan@kylinos.cn>
 * Date: Mon Mar 18 17:15:21 2024 +0800
 */
package controller

import (
	"strings"

	"gitee.com/openeuler/PilotGo-plugin-syscare/agent/client"
	"gitee.com/openeuler/PilotGo/sdk/response"
	"github.com/gin-gonic/gin"
)

func BindHandler(c *gin.Context) {
	port := c.Query("port")
	uuid := c.Query("uuid")

	server := strings.Split(c.Request.RemoteAddr, ":")[0] + ":" + port
	if client.GlobalClient.AgentServer == "" {
		client.GlobalClient.AgentServer = server
		client.GlobalClient.UUID = uuid
	} else if client.GlobalClient.AgentServer != "" && client.GlobalClient.AgentServer != server {
		response.Fail(c, client.GetClient(), "已有syscare-server与此插件绑定")
		return
	}
	response.Success(c, nil, "bind server success")
}

func UnBindHandler(c *gin.Context) {
	client.GlobalClient.AgentServer = ""
	client.GlobalClient.UUID = ""
	response.Success(c, nil, "unbind server success")
}

func Heartbeat(c *gin.Context) {
	ok := client.GlobalClient.IsBind()
	if ok {
		response.Success(c, "bind", "")
	} else {
		response.Fail(c, nil, "")
	}
}
