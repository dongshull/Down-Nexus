package database

import (
	"fmt"

	"down-nexus-api/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// InitDB 初始化 SQLite 数据库连接
// 自动创建表结构并返回数据库实例
func InitDB(dsn string) (*gorm.DB, error) {
	// 连接 SQLite 数据库
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// 自动迁移表结构
	if err := db.AutoMigrate(&models.ClientConfig{}); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return db, nil
}