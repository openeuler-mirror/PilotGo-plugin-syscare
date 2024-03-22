package controller

import (
	"io"
	"net/url"
	"os"
	"path/filepath"

	"gitee.com/openeuler/PilotGo-plugin-syscare/server/config"
	"gitee.com/openeuler/PilotGo/sdk/logger"
	"gitee.com/openeuler/PilotGo/sdk/response"
	"github.com/gin-gonic/gin"
)

func Upload(c *gin.Context) {
	contentType := c.Request.Header.Get("Content-Type")

	if contentType == "multipart/form-data" {
		// 直接读取request body内容
		bodyBuf, err := io.ReadAll(c.Request.Body)
		if err != nil {
			logger.Error("没获取到request body: %s", err.Error())
			response.Fail(c, gin.H{"error": err.Error()}, "获取文件request body失败")
			return
		}

		parsedURL, err := url.Parse(c.Request.RequestURI)
		if err != nil {
			logger.Error("解析 URL 错误:%v", err.Error())
			response.Fail(c, gin.H{"error": err.Error()}, "解析URI失败")
			return
		}
		filename := parsedURL.Query().Get("filename")

		uploadPath := c.DefaultQuery("path", config.Config().Storage.Path) // 获取上传文件的保存路径，可以通过path设置上传路径
		if err := os.MkdirAll(uploadPath, os.ModePerm); err != nil {       // 确保保存路径存在，如果不存在则创建
			response.Fail(c, gin.H{"error": err.Error()}, "创建保存路径失败")
			return
		}
		destination := filepath.Join(uploadPath, filename) // 上传文件的目标路径

		outFile, err := os.Create(destination) // 创建并打开目标文件
		if err != nil {
			response.Fail(c, gin.H{"error": err.Error()}, "创建目标文件失败")
			return
		}
		defer outFile.Close()

		_, err = outFile.Write(bodyBuf) // 将请求体中的二进制数据写入目标文件
		if err != nil {
			response.Fail(c, gin.H{"error": err.Error()}, "文件保存失败")
			return
		}

	} else {
		file, header, err := c.Request.FormFile("file")
		if err != nil {
			response.Fail(c, nil, "未获取到文件")
			return
		}
		defer file.Close()

		uploadPath := c.DefaultQuery("path", config.Config().Storage.Path)
		if err := makeDir(uploadPath); err != nil {
			response.Fail(c, nil, "创建存储目录失败")
			return
		}

		destination := filepath.Join(uploadPath, header.Filename)
		outFile, err := os.Create(destination)
		if err != nil {
			response.Fail(c, nil, "目标文件创建失败")
			return
		}
		defer outFile.Close()
		_, err = io.Copy(outFile, file) // 复制上传的文件到目标文件
		if err != nil {
			response.Fail(c, nil, "文件保存失败")
			return
		}
	}
	response.Success(c, nil, "文件上传成功")
}

func makeDir(storage string) error {
	if _, err := os.Stat(storage); os.IsNotExist(err) {
		err := os.MkdirAll(storage, os.ModePerm)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}
	return nil
}
