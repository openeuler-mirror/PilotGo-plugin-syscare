/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo-plugin-syscare licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: zhanghan2021 <zhanghan@kylinos.cn>
 * Date: Tue Apr 2 17:42:13 2024 +0800
 */
package service

import (
	"errors"
	"mime/multipart"
	"strings"
	"time"

	"gitee.com/openeuler/PilotGo-plugin-syscare/server/config"
	"gitee.com/openeuler/PilotGo-plugin-syscare/server/dao"
	"gitee.com/openeuler/PilotGo-plugin-syscare/utils"
	"gitee.com/openeuler/PilotGo/sdk/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	// BuildStatus
	pending  = "等待中"
	building = "构建中"
	Failure  = "失败"
	Complete = "完成"
)

func CreateWarmList(c *gin.Context, ip, buildVersion, patchDescription, patchVersion, patchRelease, buildKernelSrc, buildDebugInfo, patchType string, Patchs []*multipart.FileHeader) error {
	taskId := uuid.New().String()
	workDir := config.Config().Storage.Path + taskId
	if err := utils.MakeDir(workDir); err != nil {
		return err
	}

	var patchs []string
	for _, file := range Patchs {
		c.SaveUploadedFile(file, workDir+"/"+file.Filename)
		patchs = append(patchs, file.Filename)
	}
	if err := dao.SaveWarmList(&dao.WarmList{
		CreateTime:       time.Now().Format("2006-01-02 15:04:05"),
		BuildMachine:     ip,
		BuildStatus:      pending,
		BuildVersion:     buildVersion,
		PatchDescription: patchDescription,
		PatchVersion:     patchVersion,
		PatchRelease:     patchRelease,
		BuildKernel:      buildKernelSrc,
		BuildDebugInfo:   buildDebugInfo,
		PatchType:        patchType,
		Patchs:           strings.Join(patchs, ","),
		TaskId:           taskId,
	}); err != nil {
		return err
	}

	MyTask.Enqueue(&Task{
		TaskId:           taskId,
		IP:               ip,
		BuildKernel:      buildKernelSrc,
		BuildDebugInfo:   buildDebugInfo,
		PatchDescription: patchDescription,
		PatchVersion:     patchVersion,
		PatchRelease:     patchRelease,
		PatchType:        patchType,
		Patchs:           patchs,
	})
	return nil
}

func QueryWarmLists(search string, query *response.PaginationQ) (interface{}, int, error) {
	lists, total, err := dao.QueryWarmLists(search, query)
	if err != nil {
		return nil, 0, err
	}
	return lists, int(total), nil

}

func DeleteWarmList(id string) error {
	dir, err := dao.QueryStorageDir(id)
	if err != nil {
		return err
	}

	// 删除本地服务器存储的热补丁包
	path := config.Config().Storage.Path + dir
	exitCode, _, stderr, err := utils.RunCommand("rm -rf "+path, "")
	if exitCode != 0 || stderr != "" || err != nil {
		return errors.New("删除本地服务器热补丁包失败")
	}

	err = dao.DeleteWarmlist(id)
	if err != nil {
		return err
	}
	return nil
}

func UpdateWarmInfo(taskId string, warm *dao.WarmList) error {
	err := dao.UpdateWarmInfo(taskId, warm)
	if err != nil {
		return err
	}
	return nil
}

func UpdateTaskStatusToBuilding(taskId string) error {
	w := &dao.WarmList{
		BuildStatus: building,
	}
	err := dao.UpdateTaskStatusToBuilding(taskId, w)
	if err != nil {
		return err
	}
	return nil
}
