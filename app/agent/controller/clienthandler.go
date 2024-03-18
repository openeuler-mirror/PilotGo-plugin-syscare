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
