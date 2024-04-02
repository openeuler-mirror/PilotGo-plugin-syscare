package controller

import (
	"gitee.com/openeuler/PilotGo-plugin-syscare/server/service"
	"gitee.com/openeuler/PilotGo/sdk/response"
	"github.com/gin-gonic/gin"
)

func QueryWarmLists(c *gin.Context) {
	query := &response.PaginationQ{}
	err := c.ShouldBindQuery(query)
	if err != nil {
		response.Fail(c, gin.H{"status": false}, err.Error())
		return
	}

	lists, total, err := service.QueryWarmLists(query)
	if err != nil {
		response.Fail(c, nil, err.Error())
		return
	}
	response.DataPagination(c, lists, total, query)

}
