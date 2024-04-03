package controller

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"gitee.com/openeuler/PilotGo-plugin-syscare/server/config"
	"gitee.com/openeuler/PilotGo/sdk/response"
	"github.com/gin-gonic/gin"
)

func Download(c *gin.Context) {
	filename := c.Param("filename")

	// 获取下载文件的路径，可以通过path设置
	taskId := c.Query("path")
	downloadPath := config.Config().Storage.Path + taskId
	filePath := filepath.Join(downloadPath, filename) // 构建完整的文件路径

	// 检查文件是否存在
	_, err := os.Stat(filePath)
	if err != nil {
		response.Fail(c, gin.H{"error": "文件不存在"}, "文件下载失败")
		return
	}

	// 设置下载文件的响应头
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))

	// 设置文件下载的响应类型
	c.Header("Content-Type", "application/octet-stream")

	// 打开文件并将其内容写入响应体
	file, err := os.Open(filePath)
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
