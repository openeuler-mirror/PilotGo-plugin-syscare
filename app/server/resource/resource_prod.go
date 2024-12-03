/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo-plugin-syscare licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: zhanghan2021 <zhanghan@kylinos.cn>
 * Date: Wed Apr 3 17:22:38 2024 +0800
 */
//go:build production
// +build production

package resource

import (
	"embed"
	"io/fs"
	"net/http"
	"strings"

	"gitee.com/openeuler/PilotGo/sdk/logger"
	"github.com/gin-gonic/gin"
)

//go:embed all:static index.html logo.svg
var StaticFiles embed.FS

func StaticRouter(router *gin.Engine) {
	sf, err := fs.Sub(StaticFiles, "static")
	if err != nil {
		logger.Error("failed to load frontend assets files: %s", err.Error())
		return
	}

	router.StaticFS("/static", http.FS(sf))
	router.GET("/", func(c *gin.Context) {
		c.FileFromFS("/", http.FS(StaticFiles))
	})
	router.GET("/logo.svg", func(c *gin.Context) {
		c.FileFromFS("/logo.svg", http.FS(StaticFiles))
	})

	// 解决页面刷新404的问题
	router.NoRoute(func(c *gin.Context) {
		if !strings.HasPrefix(c.Request.RequestURI, "/api/") && !strings.HasPrefix(c.Request.RequestURI, "/plugin/") {
			c.FileFromFS("/", http.FS(StaticFiles))
			return
		}
		c.AbortWithStatus(http.StatusNotFound)
	})
}
