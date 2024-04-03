package router

import (
	"net/http"
	"strings"

	"gitee.com/openeuler/PilotGo-plugin-syscare/server/config"
	"gitee.com/openeuler/PilotGo-plugin-syscare/server/controller"
	"gitee.com/openeuler/PilotGo/sdk/logger"
	"github.com/gin-gonic/gin"
)

func HttpServerInit(conf *config.HttpServer) error {
	go func() {
		r := setupRouter()
		logger.Info("start http service on: http://%s", conf.Addr)
		if err := r.Run(conf.Addr); err != nil {
			logger.Error("start http server failed:%v", err)
		}

	}()

	return nil
}
func setupRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(logger.RequestLogger())
	router.Use(gin.Recovery())

	registerAPIs(router)
	staticRouter(router)

	return router
}
func registerAPIs(router *gin.Engine) {
	logger.Debug("router register")
	agent := router.Group("/plugin/syscare")
	{
		agent.GET("agentList", controller.GetAgentsHandler)
		agent.POST("addAgent", controller.AddAgentHandler)
		agent.DELETE("/:ip", controller.DeleteAgentHandler)
		agent.GET("buildEnv", controller.AgentBuildEnv)
	}

	warm := router.Group("/plugin/syscare")
	{
		warm.POST("addWarm", controller.AddWarmList)
		warm.GET("lists", controller.QueryWarmLists)
		warm.GET("/delete", controller.DeleteWarmList)
	}

	task := router.Group("/plugin/syscare")
	{
		task.POST("scriptResult", controller.ScriptResult)
	}

	fileservice := router.Group("/plugin/syscare")
	{
		fileservice.POST("/upload", controller.Upload)
		fileservice.GET("/download/:filename", controller.Download)
	}
}

func staticRouter(router *gin.Engine) {
	router.Static("/plugin/syscare/static", "../web/dist/static")
	router.StaticFile("/plugin/syscare", "../web/dist/index.html")

	// 解决页面刷新404的问题
	router.NoRoute(func(c *gin.Context) {
		logger.Debug("process noroute: %s", c.Request.URL)
		if !strings.HasPrefix(c.Request.RequestURI, "/plugin/syscare/*path") {
			c.File("../web/dist/index.html")
			return
		}
		c.Status(http.StatusNotFound)
	})
}
