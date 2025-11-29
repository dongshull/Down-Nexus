package api

import (
	"down-nexus-api/internal/core"
	"github.com/gin-gonic/gin"
)

// SetupRouter 设置路由器并返回 Gin 引擎
func SetupRouter(service *core.TorrentService) *gin.Engine {
	// 创建 Gin 路由器
	router := gin.Default()

	// 创建处理器
	handler := NewTorrentHandler(service)

	// 添加 CORS 中间件
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// API 路由组
	v1 := router.Group("/api/v1")
	{
		// 种子相关路由
		torrents := v1.Group("/torrents")
		{
			torrents.GET("", handler.GetTorrents)           // 获取所有种子
			torrents.POST("", handler.AddTorrent)            // 添加种子
			torrents.POST("/pause", handler.PauseTorrent)    // 暂停种子
			torrents.POST("/resume", handler.ResumeTorrent)   // 恢复种子
			torrents.DELETE("", handler.DeleteTorrent)       // 删除种子
		}

		// 客户端相关路由
		clients := v1.Group("/clients")
		{
			clients.GET("", handler.GetClients)              // 获取所有客户端
		}
	}

	// 健康检查路由
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "down-nexus-api",
		})
	})

	// 根路径
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to Down-Nexus API",
			"version": "1.0.0",
			"endpoints": gin.H{
				"health":         "/health",
				"torrents":       "/api/v1/torrents",
				"add_torrent":    "/api/v1/torrents (POST)",
				"pause_torrent":  "/api/v1/torrents/pause (POST)",
				"resume_torrent": "/api/v1/torrents/resume (POST)",
				"delete_torrent": "/api/v1/torrents (DELETE)",
				"clients":        "/api/v1/clients",
			},
		})
	})

	return router
}