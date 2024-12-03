/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo-plugin-syscare licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: zhanghan2021 <zhanghan@kylinos.cn>
 * Date: Wed Apr 3 09:32:32 2024 +0800
 */
package controller

import (
	"gitee.com/openeuler/PilotGo-plugin-syscare/agent/config"
	"gitee.com/openeuler/PilotGo-plugin-syscare/agent/service"
	"gitee.com/openeuler/PilotGo-plugin-syscare/utils"
	"gitee.com/openeuler/PilotGo/sdk/logger"
	"gitee.com/openeuler/PilotGo/sdk/response"
	"github.com/gin-gonic/gin"
)

func RunCommandHandler(c *gin.Context) {
	w := &struct {
		TaskId           string   `json:"taskId"`
		Ip               string   `json:"ip"`
		BuildKernelSrc   string   `json:"buildKernel"`
		BuildDebugInfo   string   `json:"buildDebugInfo"`
		PatchDescription string   `json:"patchDescription"`
		PatchVersion     string   `json:"patchVersion"`
		PatchRelease     string   `json:"patchRelease"`
		PatchType        string   `json:"patchType"`
		Patchs           []string `json:"patchs"`
	}{}
	if err := c.BindJSON(w); err != nil {
		response.Fail(c, nil, err.Error())
		return
	}
	workDir := config.Config().Storage.Work + w.TaskId // 创建工作目录
	if err := utils.MakeDir(workDir); err != nil {
		response.Fail(c, nil, err.Error())
		return
	}
	for i, patch := range w.Patchs {
		w.Patchs[i] = workDir + "/" + patch
		if err := service.DownloadPatchsFromServer(patch, w.TaskId, workDir); err != nil { // 从服务器中下载patch到工作目录
			logger.Error("保存patch文件%v失败: %v", patch, err.Error())
			continue
		}
	}

	sourceDir := config.Config().Storage.Path + w.BuildKernelSrc
	debugInfoDir := config.Config().Storage.Path + w.BuildDebugInfo
	args := service.CommandSplice(sourceDir, debugInfoDir, w.PatchDescription, w.PatchVersion, w.PatchRelease, w.PatchType, w.Patchs)
	err := service.RunCommand(w.TaskId, w.Ip, workDir, args)
	if err != nil {
		response.Fail(c, nil, err.Error())
		return
	}
	response.Success(c, nil, "执行结束")
}
