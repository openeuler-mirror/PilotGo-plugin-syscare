package service

import (
	"mime/multipart"
	"strings"
	"time"

	"gitee.com/openeuler/PilotGo-plugin-syscare/server/config"
	"gitee.com/openeuler/PilotGo-plugin-syscare/server/dao"
	"gitee.com/openeuler/PilotGo/sdk/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	// BuildStatus
	pending  = "准备中"
	building = "构建中"
	Failure  = "失败"
	Complete = "完成"
)

func CreateWarmList(c *gin.Context, ip, buildVersion, patchDescription, patchVersion, patchRelease, buildKernelSrc, buildDebugInfo, patchType string, Patchs []*multipart.FileHeader) error {
	taskId := uuid.New().String()
	workDir := config.Config().Storage.Path + taskId
	if err := MakeDir(workDir); err != nil {
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

func QueryWarmLists(query *response.PaginationQ) (interface{}, int, error) {
	lists, total, err := dao.QueryWarmLists(query)
	if err != nil {
		return nil, 0, err
	}
	return lists, int(total), nil

}
