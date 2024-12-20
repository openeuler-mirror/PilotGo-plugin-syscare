/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo-plugin-syscare licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: zhanghan2021 <zhanghan@kylinos.cn>
 * Date: Mon Mar 11 14:49:16 2024 +0800
 */
package router

import (
	"gitee.com/openeuler/PilotGo-plugin-syscare/server/config"
	"gitee.com/openeuler/PilotGo-plugin-syscare/server/controller"
	"gitee.com/openeuler/PilotGo-plugin-syscare/server/resource"
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
	// 绑定前端静态资源handler
	resource.StaticRouter(router)

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
		fileservice.GET("/download/:filename", controller.Download)
	}
}
