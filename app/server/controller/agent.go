package controller

import (
	"gitee.com/openeuler/PilotGo-plugin-syscare/server/agentmanager"
	"gitee.com/openeuler/PilotGo-plugin-syscare/server/service"
	"gitee.com/openeuler/PilotGo/sdk/response"
	"github.com/gin-gonic/gin"
)

func GetAgentsHandler(c *gin.Context) {
	IsPage := c.Query("paged")
	switch IsPage {
	case "false":
		agents, err := service.QueryAgents()
		if err != nil {
			response.Fail(c, nil, err.Error())
			return
		}
		response.Success(c, agents, "查询到主机列表")
	default:
		query := &response.PaginationQ{}
		err := c.ShouldBindQuery(query)
		if err != nil {
			response.Fail(c, gin.H{"status": false}, err.Error())
			return
		}

		search := c.Query("search")
		agents, total, err := service.SearchAgents(search, query)
		if err != nil {
			response.Fail(c, nil, err.Error())
			return
		}
		response.DataPagination(c, agents, total, query)
	}
}

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
