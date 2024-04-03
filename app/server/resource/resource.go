//go:build !production
// +build !production

package resource

import (
	"net/http"
	"strings"

	"gitee.com/openeuler/PilotGo/sdk/logger"
	"github.com/gin-gonic/gin"
)

func StaticRouter(router *gin.Engine) {
	router.Static("/plugin/syscare/static", "../../web/dist/static")
	router.StaticFile("/plugin/syscare", "../../web/dist/index.html")

	// 解决页面刷新404的问题
	router.NoRoute(func(c *gin.Context) {
		logger.Debug("process noroute: %s", c.Request.URL)
		if !strings.HasPrefix(c.Request.RequestURI, "/plugin/syscare/*path") {
			c.File("../../web/dist/index.html")
			return
		}
		c.Status(http.StatusNotFound)
	})
}
