/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo-plugin-syscare licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: zhanghan2021 <zhanghan@kylinos.cn>
 * Date: Wed Apr 3 09:24:34 2024 +0800
 */
package service

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"gitee.com/openeuler/PilotGo-plugin-syscare/server/config"
	"gitee.com/openeuler/PilotGo/sdk/logger"
	"gitee.com/openeuler/PilotGo/sdk/utils/httputils"
)

func DownloadRpmsFromAgent(ip, taskId, rpm string) error {
	url := "http://" + ip + ":" + config.Config().AgentServer.Port + "/plugin_agent_manage/download/" + rpm + "?path=" + taskId
	logger.Info("下载地址：%v", url)
	resp, err := httputils.Get(url, nil)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New("未获取到文件")
	}

	storagekDir := config.Config().Storage.Path + taskId
	out, err := os.Create(storagekDir + "/" + rpm) // 保存文件的路径
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, ioutil.NopCloser(bytes.NewReader(resp.Body)))
	if err != nil {
		return err
	}
	return nil
}
