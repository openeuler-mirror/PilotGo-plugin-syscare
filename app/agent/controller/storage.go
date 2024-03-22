package controller

import (
	"gitee.com/openeuler/PilotGo-plugin-syscare/agent/service"
	"gitee.com/openeuler/PilotGo/sdk/response"
	"github.com/gin-gonic/gin"
)

func BuildEnv(c *gin.Context) {
	data, err := service.SearchDependentEnv()
	if err != nil {
		response.Fail(c, nil, err.Error())
		return
	}
	response.Success(c, data, "获取热补丁环境")
}
