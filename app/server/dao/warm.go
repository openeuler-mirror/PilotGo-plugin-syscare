package dao

import (
	"gitee.com/openeuler/PilotGo-plugin-syscare/server/db"
	"gitee.com/openeuler/PilotGo/sdk/response"
)

type WarmList struct {
	ID           int    `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	CreateTime   string `json:"creatTime"`
	BuildMachine string `json:"ip"`
	BuildStatus  string `json:"status"`
	EndTime      string `json:"endTime"`
	// 构建值
	BuildVersion     string `json:"buildVersion"`
	BuildKernel      string `json:"buildKernel"`
	BuildDebugInfo   string `json:"buildDebugInfo"`
	PatchDescription string `json:"patchDescription"`
	PatchVersion     string `json:"patchVersion"`
	PatchRelease     string `json:"patchRelease"`
	PatchType        string `json:"patchType"`
	Patchs           string `json:"patchs"`
	// 热补丁包
	TaskId          string `json:"taskId"` // 1：热补丁任务编号 2：热补丁存储文件目录
	HotPatchName    string `json:"hotPatch"`
	PatchKernelName string `json:"patchKernel"`
	// 返回值
	ExitCode int    `json:"exitCode"`
	Stdout   string `json:"stdout"`
	BuildLog string `json:"buildLog"`
}

func QueryWarmLists(query *response.PaginationQ) ([]*WarmList, int64, error) {
	var lists []*WarmList
	if err := db.MySQL().Limit(query.PageSize).Offset((query.Page - 1) * query.PageSize).Find(&lists).Error; err != nil {
		return nil, 0, err
	}

	var total int64
	if err := db.MySQL().Model(&lists).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	return lists, total, nil
}
