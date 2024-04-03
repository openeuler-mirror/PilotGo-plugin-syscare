package service

import (
	"errors"
	"io"
	"net/http"
	"os"

	"gitee.com/openeuler/PilotGo-plugin-syscare/server/config"
)

func DownloadRpmsFromAgent(ip, taskId, rpm string) error {
	url := "http://" + ip + ":" + config.Config().AgentServer.Port + "/plugin_agent_manage/download/" + rpm + "?path=" + taskId

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New(resp.Status)
	}

	storagekDir := config.Config().Storage.Path + taskId
	out, err := os.Create(storagekDir + "/" + rpm) // 保存文件的路径
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}
	return nil
}
