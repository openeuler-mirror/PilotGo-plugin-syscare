package controller

import (
	"encoding/json"
	"io"
	"time"

	"gitee.com/openeuler/PilotGo-plugin-syscare/server/dao"
	"gitee.com/openeuler/PilotGo-plugin-syscare/server/service"
	"gitee.com/openeuler/PilotGo/sdk/logger"
	"gitee.com/openeuler/PilotGo/sdk/response"
	"github.com/gin-gonic/gin"
)

func ScriptResult(c *gin.Context) {
	j, err := io.ReadAll(c.Request.Body)
	if err != nil {
		logger.Error("没获取到：%s", err.Error())
		response.Fail(c, nil, err.Error())
		return
	}
	logger.Info(string(j))
	result := &struct {
		TaskId       string `json:"taskId"`
		IP           string `json:"ip"`
		KernelSrcRpm string `json:"kernelSrcRpm"`
		KernelRpm    string `json:"kernelRpm"`
		RetCode      int    `json:"retcode"`
		Stdout       string `json:"stdout"`
		BuildLog     string `json:"buildLog"`
	}{}
	err = json.Unmarshal(j, &result)
	if err != nil {
		logger.Error("解析出错:%v", result)
		response.Fail(c, nil, err.Error())
		return
	}

	var buildStatus string
	if result.BuildLog != "" {
		buildStatus = service.Failure
	} else {
		buildStatus = service.Complete
	}
	if err := service.UpdateWarmInfo(result.TaskId, &dao.WarmList{
		BuildStatus:     buildStatus,
		EndTime:         time.Now().Format("2006-01-02 15:04:05"),
		HotPatchName:    result.KernelRpm,
		PatchKernelName: result.KernelSrcRpm,
		ExitCode:        result.RetCode,
		Stdout:          result.Stdout,
		BuildLog:        result.BuildLog,
	}); err != nil {
		response.Fail(c, nil, err.Error())
		return
	}

	if agent := service.MyTask.GetTaskAgent(result.IP); agent != nil { // 监控当前agent正在执行任务的数目
		agent.DecrementCurrentTasksCount()
	}

	// 将远端的软件包下载到本机服务器
	if err := service.DownloadRpmsFromAgent(result.IP, result.TaskId, result.KernelSrcRpm); err != nil {
		response.Fail(c, nil, err.Error())
		return
	}
	if err := service.DownloadRpmsFromAgent(result.IP, result.TaskId, result.KernelRpm); err != nil {
		response.Fail(c, nil, err.Error())
		return
	}
}
