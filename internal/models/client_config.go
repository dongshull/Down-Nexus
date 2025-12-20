package models

import (
	"gorm.io/gorm"
)

// ClientConfig 客户端配置模型
// 用于存储各种下载客户端的连接信息
type ClientConfig struct {
	gorm.Model
	// ClientID 客户端唯一标识符，用于区分不同的客户端实例
	ClientID string `gorm:"uniqueIndex;not null" json:"client_id"`
	// Type 客户端类型：qbittorrent 或 transmission
	Type string `gorm:"not null" json:"type"`
	// Host 客户端服务器地址，包含端口号
	Host string `gorm:"not null" json:"host"`
	// Username 登录用户名
	Username string `gorm:"not null" json:"username"`
	// Password 登录密码
	Password string `gorm:"not null" json:"password"`
	// Enabled 是否启用该客户端配置
	Enabled bool `gorm:"default:true" json:"enabled"`
	// DisplayName 客户端显示名称，用于UI展示
	DisplayName string `json:"display_name"`
	// DefaultSavePath 默认保存路径
	DefaultSavePath string `json:"default_save_path"`
	// DefaultCategory 默认分类
	DefaultCategory string `json:"default_category"`
	// DefaultTags 默认标签
	DefaultTags []string `gorm:"serializer:json" json:"default_tags"`
	// GlobalDownloadLimit 全局下载限速 (字节/秒)，-1表示无限制
	GlobalDownloadLimit int64 `json:"global_download_limit"`
	// GlobalUploadLimit 全局上传限速 (字节/秒)，-1表示无限制
	GlobalUploadLimit int64 `json:"global_upload_limit"`
	// AutoConnectOnStart 启动时自动连接
	AutoConnectOnStart bool `gorm:"default:true" json:"auto_connect_on_start"`
	// AdvancedSettings 高级设置
	AdvancedSettings map[string]interface{} `gorm:"serializer:json" json:"advanced_settings"`
}