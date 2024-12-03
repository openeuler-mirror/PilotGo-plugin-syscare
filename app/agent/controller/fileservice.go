/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo-plugin-syscare licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: zhanghan2021 <zhanghan@kylinos.cn>
 * Date: Tue Apr 2 17:45:45 2024 +0800
 */
package controller

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"gitee.com/openeuler/PilotGo-plugin-syscare/agent/config"
	"gitee.com/openeuler/PilotGo/sdk/response"
	"github.com/gin-gonic/gin"
)

func Download(c *gin.Context) {
	filename := c.Param("filename")

	// 获取下载文件的路径，以taskId作为存储文件夹名字
	taskId := c.Query("path")
	downloadPath := config.Config().Storage.Work + taskId
	filePath := filepath.Join(downloadPath, filename) // 构建完整的文件路径

	_, err := os.Stat(filePath) // 检查文件是否存在
	if err != nil {
		response.Fail(c, gin.H{"error": "文件不存在"}, "文件下载失败")
		return
	}

	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename)) // 设置下载文件的响应头
	c.Header("Content-Type", "application/octet-stream")                              // 设置文件下载的响应类型

	file, err := os.Open(filePath) // 打开文件并将其内容写入响应体
	if err != nil {
		response.Fail(c, gin.H{"error": err.Error()}, "文件下载失败")
		return
	}
	defer file.Close()

	_, err = io.Copy(c.Writer, file)
	if err != nil {
		response.Fail(c, gin.H{"error": err.Error()}, "文件下载失败")
		return
	}
}
