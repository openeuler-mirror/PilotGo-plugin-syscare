/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo-plugin-syscare licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: zhanghan2021 <zhanghan@kylinos.cn>
 * Date: Fri Mar 22 16:34:55 2024 +0800
 */
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
