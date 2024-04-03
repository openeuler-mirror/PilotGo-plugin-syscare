package router

import (
	"net/http"
	"time"

	"gitee.com/openeuler/PilotGo-plugin-syscare/agent/client"
	"gitee.com/openeuler/PilotGo-plugin-syscare/agent/config"
	"gitee.com/openeuler/PilotGo-plugin-syscare/agent/controller"
	"gitee.com/openeuler/PilotGo-plugin-syscare/utils"
	"gitee.com/openeuler/PilotGo/sdk/logger"
	"github.com/gin-gonic/gin"
)

func HttpServerInit() error {
	_, err := client.DefaultClient()
	if err != nil {
		return err
	}
	ip, err := utils.GetLocalIP()
	if err != nil {
		return err
	}
	httpAddr := ip + ":" + config.Config().Server.Port

	go func() {
		r := setupRouter()
		logger.Info("start http service on: http://%s", httpAddr)
		if err := r.Run(httpAddr); err != nil {
			logger.Error("start http server failed:%v", err)
		}

	}()

	return nil
}
func setupRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(requestLogger())
	router.Use()
	router.Use(gin.Recovery())

	registerAPIs(router)

	return router
}
func registerAPIs(router *gin.Engine) {
	logger.Debug("router register")

	// 提供syscare agent基本信息
	mg := router.Group("/plugin_agent_manage/")
	{
		mg.GET("/heartbeat", controller.Heartbeat)
		mg.GET("/info", func(c *gin.Context) { c.JSON(http.StatusOK, client.GetClient()) })
		mg.PUT("/bind", controller.BindHandler)     // 绑定agent server和uuid
		mg.PUT("/unbind", controller.UnBindHandler) // 解绑agent server和uuid
	}
	storage := router.Group("/plugin_agent_manage/")
	{
		storage.GET("/buildEnv", controller.BuildEnv)
	}
	run := router.Group("/plugin_agent_manage/")
	{
		run.POST("/run", controller.RunCommandHandler)
	}
	fileservice := router.Group("/plugin_agent_manage/")
	{
		fileservice.GET("/download/:filename", controller.Download)
	}
}

func requestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 开始时间
		startTime := time.Now()

		// 处理请求
		c.Next()

		// 结束时间
		endTime := time.Now()

		// 执行时间
		latencyTime := endTime.Sub(startTime)

		// 请求方式
		reqMethod := c.Request.Method

		// 请求路由
		reqUri := c.Request.RequestURI

		// 状态码
		statusCode := c.Writer.Status()

		// 请求IP
		clientIP := c.ClientIP()

		// 日志格式
		if reqUri != "/plugin_agent_manage/heartbeat" {
			logger.Debug("status_code:%d latency_time:%s client_ip:%s req_method:%s req_uri:%s",
				statusCode,
				latencyTime,
				clientIP,
				reqMethod,
				reqUri,
			)
		}
	}
}
