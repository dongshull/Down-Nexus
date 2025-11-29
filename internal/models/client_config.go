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
}