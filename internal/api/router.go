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
			torrents.GET("", handler.GetTorrents)                  // 获取所有种子
			torrents.GET("/:hash", handler.GetTorrentDetails)      // 获取种子详细信息
			torrents.POST("", handler.AddTorrent)                  // 添加种子
			torrents.POST("/pause", handler.PauseTorrent)          // 暂停单个种子
			torrents.POST("/pause/batch", handler.PauseTorrents)   // 批量暂停种子
			torrents.POST("/resume", handler.ResumeTorrent)        // 恢复单个种子
			torrents.POST("/resume/batch", handler.ResumeTorrents) // 批量恢复种子
			torrents.DELETE("", handler.DeleteTorrent)             // 删除单个种子
			torrents.DELETE("/batch", handler.DeleteTorrents)      // 批量删除种子
			torrents.POST("/category", handler.SetTorrentCategory) // 设置种子分类
			torrents.POST("/tags", handler.SetTorrentTags)         // 设置种子标签
			torrents.POST("/limits", handler.SetTorrentLimits)     // 设置种子限速
			torrents.POST("/priority", handler.SetTorrentPriority) // 设置种子优先级

			// 扩展功能
			torrents.POST("/recheck", handler.RecheckTorrent)                        // 校验种子
			torrents.POST("/recheck/batch", handler.RecheckTorrents)                 // 批量校验种子
			torrents.POST("/reannounce", handler.ReannounceTorrent)                  // 重新汇报种子
			torrents.POST("/reannounce/batch", handler.ReannounceTorrents)           // 批量重新汇报种子
			torrents.POST("/location", handler.SetTorrentLocation)                   // 设置种子保存位置
			torrents.POST("/name", handler.SetTorrentName)                           // 设置种子名称
			torrents.POST("/force-start", handler.SetForceStart)                     // 设置强制开始
			torrents.POST("/auto-management", handler.SetAutoManagement)             // 设置自动管理
			torrents.POST("/sequential", handler.SetSequentialDownload)              // 设置顺序下载
			torrents.POST("/first-last-priority", handler.SetFirstLastPiecePriority) // 设置首尾片优先
			torrents.POST("/super-seeding", handler.SetSuperSeeding)                 // 设置超级做种
			torrents.POST("/share-limit", handler.SetShareLimit)                     // 设置分享限制

			// 种子文件操作
			torrents.GET("/:hash/files", handler.GetTorrentFiles)     // 获取种子文件列表
			torrents.POST("/files/priority", handler.SetFilePriority) // 设置文件优先级
			torrents.POST("/files/rename", handler.RenameFile)        // 重命名文件
			torrents.POST("/folders/rename", handler.RenameFolder)    // 重命名文件夹

			// 种子详细信息
			torrents.GET("/:hash/trackers", handler.GetTorrentTrackers)      // 获取种子Tracker列表
			torrents.POST("/trackers/add", handler.AddTorrentTrackers)       // 添加种子Tracker
			torrents.POST("/trackers/remove", handler.RemoveTorrentTrackers) // 移除种子Tracker
			torrents.GET("/:hash/peers", handler.GetTorrentPeers)            // 获取种子Peer列表
			torrents.GET("/:hash/properties", handler.GetTorrentProperties)  // 获取种子详细属性
			torrents.POST("/piece-states", handler.GetTorrentPieceStates)    // 获取种子分片状态
		}

		// 客户端相关路由
		clients := v1.Group("/clients")
		{
			clients.GET("", handler.GetClients)                     // 获取所有客户端
			clients.GET("/:id", handler.GetClientConfig)            // 获取特定客户端配置
			clients.POST("", handler.CreateClientConfig)            // 创建客户端配置
			clients.PUT("/:id", handler.UpdateClientConfig)         // 更新客户端配置
			clients.DELETE("/:id", handler.DeleteClientConfig)      // 删除客户端配置
			clients.POST("/:id/test", handler.TestClientConnection) // 测试客户端连接
		}

		// 分类相关路由
		v1.GET("/categories", handler.GetCategories)       // 获取分类列表
		v1.POST("/categories", handler.CreateCategory)     // 创建分类
		v1.DELETE("/categories", handler.DeleteCategories) // 删除分类

		// 标签相关路由
		v1.GET("/tags", handler.GetTags)       // 获取标签列表
		v1.POST("/tags", handler.CreateTags)   // 创建标签
		v1.DELETE("/tags", handler.DeleteTags) // 删除标签

		// 全局限速相关路由
		v1.GET("/global-limits", handler.GetGlobalLimits)                           // 获取全局限速
		v1.POST("/global-limits", handler.SetGlobalSpeedLimits)                     // 设置全局限速
		v1.GET("/alternative-limits", handler.GetAlternativeSpeedLimits)            // 获取备用速度限制状态
		v1.POST("/alternative-limits/toggle", handler.ToggleAlternativeSpeedLimits) // 切换备用速度限制

		// 传输信息
		v1.GET("/transfer-info", handler.GetTransferInfo) // 获取传输信息
		v1.GET("/free-space", handler.GetFreeSpace)       // 获取磁盘剩余空间

		// 服务器操作
		v1.GET("/server-info", handler.GetServerInfo)            // 获取服务器信息
		v1.GET("/default-save-path", handler.GetDefaultSavePath) // 获取默认保存路径
		v1.POST("/shutdown", handler.ShutdownClient)             // 关闭客户端

		// 日志
		v1.GET("/logs", handler.GetLogs) // 获取日志
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
				"health": "/health",
				"torrents": gin.H{
					"list":         "/api/v1/torrents (GET)",
					"details":      "/api/v1/torrents/:hash (GET)",
					"add":          "/api/v1/torrents (POST)",
					"pause":        "/api/v1/torrents/pause (POST)",
					"pause_batch":  "/api/v1/torrents/pause/batch (POST)",
					"resume":       "/api/v1/torrents/resume (POST)",
					"resume_batch": "/api/v1/torrents/resume/batch (POST)",
					"delete":       "/api/v1/torrents (DELETE)",
					"delete_batch": "/api/v1/torrents/delete/batch (DELETE)",
					"category":     "/api/v1/torrents/category (POST)",
					"tags":         "/api/v1/torrents/tags (POST)",
					"limits":       "/api/v1/torrents/limits (POST)",
					"priority":     "/api/v1/torrents/priority (POST)",
				},
				"clients": gin.H{
					"list":    "/api/v1/clients (GET)",
					"details": "/api/v1/clients/:id (GET)",
					"create":  "/api/v1/clients (POST)",
					"update":  "/api/v1/clients/:id (PUT)",
					"delete":  "/api/v1/clients/:id (DELETE)",
					"test":    "/api/v1/clients/:id/test (POST)",
				},
				"categories": gin.H{
					"list":   "/api/v1/categories (GET)",
					"create": "/api/v1/categories (POST)",
					"delete": "/api/v1/categories (DELETE)",
				},
				"tags": gin.H{
					"list":   "/api/v1/tags (GET)",
					"create": "/api/v1/tags (POST)",
					"delete": "/api/v1/tags (DELETE)",
				},
				"global_limits": "/api/v1/global-limits (POST)",
			},
		})
	})

	return router
}
