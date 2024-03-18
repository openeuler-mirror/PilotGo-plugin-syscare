package controller

import (
	"gitee.com/openeuler/PilotGo-plugin-syscare/server/agentmanager"
	"gitee.com/openeuler/PilotGo/sdk/response"
	"github.com/gin-gonic/gin"
)

func AddAgentHandler(c *gin.Context) {
	param := &struct {
		IP string `json:"ip"`
	}{}

	if err := c.BindJSON(param); err != nil {
		response.Fail(c, nil, err.Error())
		return
	}
	agent, err := agentmanager.AddAgent(param.IP)
	if err != nil {
		response.Fail(c, nil, err.Error())
		return
	}
	response.Success(c, agent, "添加主机成功")
}
