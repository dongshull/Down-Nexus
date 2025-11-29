package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"down-nexus-api/internal/api"
	"down-nexus-api/internal/core"
	"down-nexus-api/internal/models"
	"down-nexus-api/pkg/clients"
	"down-nexus-api/pkg/clients/qbittorrent"
	"down-nexus-api/pkg/clients/transmission"
	"down-nexus-api/pkg/database"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

func main() {
	// 加载 .env 文件
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  未找到 .env 文件，使用默认环境变量")
	}
	
	// 设置 Gin 为发布模式，隐藏调试信息
	gin.SetMode(gin.ReleaseMode)
	
	// 精美的启动横幅
	printBanner()

	// 初始化数据库
	fmt.Println("🗄️  正在初始化数据库...")
	
	// 从环境变量构建 PostgreSQL 连接字符串
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		getEnv("DB_HOST", "localhost"),
		getEnv("DB_USER", "downnexus"),
		getEnv("DB_PASSWORD", "downnexus"),
		getEnv("DB_NAME", "downnexus"),
		getEnv("DB_PORT", "5432"),
		getEnv("DB_SSLMODE", "disable"),
		getEnv("DB_TIMEZONE", "Asia/Shanghai"),
	)
	
	db, err := database.InitDB(dsn)
	if err != nil {
		log.Fatalf("❌ 数据库初始化失败: %v", err)
	}
	fmt.Println("   ✨ PostgreSQL 数据库连接成功")

	// 检查数据库配置
	if err := checkDatabaseConfig(db); err != nil {
		log.Fatalf("❌ 数据库配置检查失败: %v", err)
	}

	// 从数据库加载客户端配置
	fmt.Println("🔧 正在从数据库加载客户端配置...")
	adapters, err := loadClientsFromDB(db)
	if err != nil {
		log.Fatalf("❌ 客户端加载失败: %v", err)
	}

	// 创建核心服务
	torrentService := core.NewTorrentService(adapters)
	fmt.Println("🎯 核心服务初始化完成")

	// 设置路由器
	router := api.SetupRouter(torrentService)
	fmt.Println("🌐 API 路由配置完成")

	// 启动服务器
	portNum := getEnv("SERVER_PORT", "8081")
	port := ":" + portNum
	
	printServerInfo(portNum)
	printAPIInfo()

	// 启动 HTTP 服务器
	fmt.Println("🚀 服务器正在启动...")
	if err := router.Run(port); err != nil {
		log.Fatal("💥 服务器启动失败:", err)
	}
}

// checkDatabaseConfig 检查数据库配置
func checkDatabaseConfig(db *gorm.DB) error {
	var count int64
	db.Model(&models.ClientConfig{}).Count(&count)
	
	if count == 0 {
		fmt.Println("   ⚠️  数据库为空，请手动配置客户端")
		fmt.Println("   💡 参考 config.example.json 文件进行配置")
		fmt.Println("   📝 您可以通过以下方式添加配置：")
		fmt.Println("      1. 直接修改 SQLite 数据库")
		fmt.Println("      2. 使用配置文件导入功能（开发中）")
		fmt.Println("      3. 使用 API 接口（开发中）")
		return fmt.Errorf("database is empty, please configure clients manually")
	}
	
	fmt.Printf("   ✨ 发现 %d 个客户端配置\n", count)
	return nil
}

// loadClientsFromDB 从数据库加载客户端配置并创建适配器
func loadClientsFromDB(db *gorm.DB) ([]clients.DownloaderClient, error) {
	var configs []models.ClientConfig
	
	// 查询所有启用的配置
	if err := db.Where("enabled = ?", true).Find(&configs).Error; err != nil {
		return nil, fmt.Errorf("failed to query client configs: %w", err)
	}

	var adapters []clients.DownloaderClient
	
	// 遍历配置创建客户端适配器
	for _, config := range configs {
		var client clients.DownloaderClient
		var err error
		
		switch config.Type {
		case "qbittorrent":
			client, err = qbittorrent.NewQbitClient(config.Host, config.Username, config.Password, config.ClientID)
		case "transmission":
			client, err = transmission.NewTransmissionClient(config.Host, config.Username, config.Password, config.ClientID)
		default:
			log.Printf("⚠️  未知的客户端类型: %s (ID: %s)", config.Type, config.ClientID)
			continue
		}
		
		if err != nil {
			log.Printf("❌ 创建客户端失败 [%s]: %v", config.ClientID, err)
			continue
		}
		
		adapters = append(adapters, client)
		fmt.Printf("   ✨ %s (%s) 已连接\n", config.Type, config.ClientID)
	}
	
	if len(adapters) == 0 {
		return nil, fmt.Errorf("no valid client adapters were created")
	}
	
	return adapters, nil
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// printBanner 打印精美的启动横幅
func printBanner() {
	fmt.Println()
	fmt.Println("╔══════════════════════════════════════════════════════════════╗")
	fmt.Println("║                                                              ║")
	fmt.Println("║                    🌟 Down-Nexus API 🌟                      ║")
	fmt.Println("║                                                              ║")
	fmt.Println("║                   多客户端种子管理系统                        ║")
	fmt.Println("║                                                              ║")
	fmt.Println("║                     Version: v1.0.0                          ║")
	fmt.Println("║                                                              ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════╝")
	fmt.Println()
}

// printServerInfo 打印服务器信息
func printServerInfo(portNum string) {
	fmt.Println()
	fmt.Println("┌──────────────────────────────────────────────────────────────┐")
	fmt.Println("│                    🌐 服务器访问地址                           │")
	fmt.Println("├──────────────────────────────────────────────────────────────┤")
	fmt.Printf("│  📍 本机地址:  %-45s │\n", fmt.Sprintf("http://localhost:%s/", portNum))
	
	// 获取内网IP地址
	addrs, err := net.InterfaceAddrs()
	if err == nil {
		for _, addr := range addrs {
			if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
				if ipNet.IP.To4() != nil {
					ipStr := ipNet.IP.String()
					// 过滤掉 198.18.0.1 这个IP地址
					if ipStr != "198.18.0.1" {
						fmt.Printf("│  🌍 内网地址:  %-45s │\n", fmt.Sprintf("http://%s:%s/", ipStr, portNum))
					}
				}
			}
		}
	}
	fmt.Println("└──────────────────────────────────────────────────────────────┘")
}

// printAPIInfo 打印API接口信息
func printAPIInfo() {
	fmt.Println()
	fmt.Println("┌──────────────────────────────────────────────────────────────┐")
	fmt.Println("│                      📚 API 接口列表                           │")
	fmt.Println("├──────────────────────────────────────────────────────────────┤")
	fmt.Println("│  🏠 基础接口:                                                 │")
	fmt.Println("│     GET  /                    - 欢迎页面                       │")
	fmt.Println("│     GET  /health              - 健康检查                       │")
	fmt.Println("│                                                              │")
	fmt.Println("│  🔥 种子管理:                                                 │")
	fmt.Println("│     GET  /api/v1/torrents     - 获取所有种子                   │")
	fmt.Println("│     POST /api/v1/torrents     - 添加种子                       │")
	fmt.Println("│     POST /api/v1/torrents/pause   - 暂停种子                   │")
	fmt.Println("│     POST /api/v1/torrents/resume  - 恢复种子                   │")
	fmt.Println("│     DELETE /api/v1/torrents   - 删除种子                       │")
	fmt.Println("│                                                              │")
	fmt.Println("│  🔧 客户端管理:                                                │")
	fmt.Println("│     GET  /api/v1/clients      - 获取客户端列表                 │")
	fmt.Println("│                                                              │")
	fmt.Println("│  📖 完整文档: API_DOCS.md                                      │")
	fmt.Println("└──────────────────────────────────────────────────────────────┘")
	fmt.Println()
}