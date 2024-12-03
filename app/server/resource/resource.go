/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo-plugin-syscare licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: zhanghan2021 <zhanghan@kylinos.cn>
 * Date: Wed Apr 3 17:22:38 2024 +0800
 */
//go:build !production
// +build !production

package resource

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func StaticRouter(router *gin.Engine) {
	router.Static("/plugin/syscare/static", "../../web/dist/static")
	router.StaticFile("/plugin/syscare", "../../web/dist/index.html")

	// 解决页面刷新404的问题
	router.NoRoute(func(c *gin.Context) {
		if !strings.HasPrefix(c.Request.RequestURI, "/plugin/syscare/*path") {
			c.File("../../web/dist/index.html")
			return
		}
		c.Status(http.StatusNotFound)
	})
}
