package controller

import (
	"gitee.com/openeuler/PilotGo-plugin-syscare/server/service"
	"gitee.com/openeuler/PilotGo/sdk/response"
	"github.com/gin-gonic/gin"
)

func AddWarmList(c *gin.Context) {
	if err := c.Request.ParseMultipartForm(10 << 20); err != nil { // 解析表单设置最大内存为 10MB
		response.Fail(c, nil, err.Error())
		return
	}

	ip := c.Request.FormValue("ip")
	buildKernel := c.Request.FormValue("buildKernel")
	buildDebugInfo := c.Request.FormValue("buildDebugInfo")
	patchType := c.Request.FormValue("patchType")
	buildVersion := c.Request.FormValue("version")
	patchDescription := c.Request.FormValue("patchDescription")
	patchVersion := c.Request.FormValue("patchVersion")
	patchRelease := c.Request.FormValue("patchRelease")

	files, ok := c.Request.MultipartForm.File["upload"]
	if !ok {
		response.Fail(c, nil, "未获取到patch文件")
		return
	}

	if err := service.CreateWarmList(c, ip, buildVersion, patchDescription, patchVersion, patchRelease, buildKernel, buildDebugInfo, patchType, files); err != nil {
		response.Fail(c, nil, err.Error())
		return
	}
	response.Success(c, nil, "成功创建热补丁任务")
}

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

func DeleteWarmList(c *gin.Context) {
	id := c.Query("id")
	err := service.DeleteWarmList(id)
	if err != nil {
		response.Fail(c, nil, err.Error())
		return
	}
	response.Success(c, nil, "删除成功")
}
