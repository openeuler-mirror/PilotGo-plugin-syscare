package service

import (
	"gitee.com/openeuler/PilotGo-plugin-syscare/agent/client"
	"gitee.com/openeuler/PilotGo-plugin-syscare/utils"
	"gitee.com/openeuler/PilotGo/sdk/logger"
	"gitee.com/openeuler/PilotGo/sdk/utils/httputils"
)

func RunCommand(taskId, ip, workDir string, args []string) error {
	var buildLog string
	var kernelSrcRpm string
	var kernelRpm string

	exitCode, stdout, _, err := utils.RunCommand("cd "+workDir+" && "+"kylin-warm build", args...)
	logger.Info("%v", stdout)
	if err != nil {
		logger.Error("执行热补丁命令失败：%v", err.Error())
		kernelRpm = ""
		kernelSrcRpm = ""
	}
	kernelSrcRpm, kernelRpm, err = ScanHotPatchRpm(workDir)
	if err != nil {
		logger.Error("搜索agent制作成功的热补丁包失败: %v", err.Error())
	}

	if exitCode != 0 { // 制作热补丁失败，读取build.log
		buildLog, err = GetBuildLog(workDir)
		if err != nil {
			logger.Error("读取build.log日志失败:%v", err.Error())
		}
	}

	data := &struct {
		TaskId       string `json:"taskId"`
		IP           string `json:"ip"`
		KernelSrcRpm string `json:"kernelSrcRpm"`
		KernelRpm    string `json:"kernelRpm"`
		RetCode      int    `json:"retcode"`
		Stdout       string `json:"stdout"`
		BuildLog     string `json:"buildLog"`
	}{
		TaskId:       taskId,
		IP:           ip,
		KernelSrcRpm: kernelSrcRpm,
		KernelRpm:    kernelRpm,
		RetCode:      exitCode,
		Stdout:       stdout,
		BuildLog:     buildLog,
	}
	url := "http://" + client.GlobalClient.AgentServer + "/plugin/syscare/scriptResult"
	_, err = httputils.Post(url, &httputils.Params{
		Body: data,
	})
	if err != nil {
		return err
	}

	return nil
}
func CommandSplice(buildKernel, buildDebugInfo, patchDescription, patchVersion, patchRelease, patchType string, patchs []string) []string {
	args := []string{}

	args = append(args, "--source", buildKernel)
	args = append(args, "--debuginfo", buildDebugInfo)
	args = append(args, "--patch-name", patchType)
	args = append(args, "--patch")
	args = append(args, patchs...)

	if patchDescription != "" {
		args = append(args, "--patch-description", patchDescription)
	}
	if patchVersion != "" {
		args = append(args, "--patch-version", patchVersion)
	}
	if patchRelease != "" {
		args = append(args, "--patch-release", patchRelease)
	}

	logger.Info("run command args: %v", args)
	return args
}
