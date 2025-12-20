package api

import (
	"net/http"

	"down-nexus-api/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetClientConfig 根据ID获取客户端配置的处理器
func (h *TorrentHandler) GetClientConfig(c *gin.Context) {
	clientID := c.Param("id")

	var config models.ClientConfig
	err := h.service.Db.Where("client_id = ?", clientID).First(&config).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"error":   "Client config not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to get client config: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    config,
	})
}

// CreateClientConfig 创建客户端配置的处理器
func (h *TorrentHandler) CreateClientConfig(c *gin.Context) {
	// 解析请求体
	var req ClientConfigCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request format: " + err.Error(),
		})
		return
	}

	// 检查客户端ID是否已存在
	var existingConfig models.ClientConfig
	err := h.service.Db.Where("client_id = ?", req.ClientID).First(&existingConfig).Error
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{
			"success": false,
			"error":   "Client ID already exists",
		})
		return
	}

	// 创建客户端配置
	config := models.ClientConfig{
		ClientID:            req.ClientID,
		Type:                req.Type,
		Host:                req.Host,
		Username:            req.Username,
		Password:            req.Password,
		DisplayName:         req.DisplayName,
		DefaultSavePath:     req.DefaultSavePath,
		DefaultCategory:     req.DefaultCategory,
		DefaultTags:         req.DefaultTags,
		GlobalDownloadLimit: req.GlobalDownloadLimit,
		GlobalUploadLimit:   req.GlobalUploadLimit,
		AutoConnectOnStart:  req.AutoConnectOnStart,
		AdvancedSettings:    req.AdvancedSettings,
		Enabled:             req.Enabled,
	}

	// 如果没有设置DisplayName，使用ClientID
	if config.DisplayName == "" {
		config.DisplayName = req.ClientID
	}

	// 默认启用
	if req.Enabled == false {
		config.Enabled = false
	} else {
		config.Enabled = true
	}

	err = h.service.Db.Create(&config).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to create client config: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Client config created successfully",
		"data":    config,
	})
}

// UpdateClientConfig 更新客户端配置的处理器
func (h *TorrentHandler) UpdateClientConfig(c *gin.Context) {
	clientID := c.Param("id")

	// 获取现有配置
	var config models.ClientConfig
	err := h.service.Db.Where("client_id = ?", clientID).First(&config).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"error":   "Client config not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to get client config: " + err.Error(),
		})
		return
	}

	// 解析请求体
	var req ClientConfigUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request format: " + err.Error(),
		})
		return
	}

	// 更新字段
	if req.DisplayName != nil {
		config.DisplayName = *req.DisplayName
	}
	if req.Host != nil {
		config.Host = *req.Host
	}
	if req.Username != nil {
		config.Username = *req.Username
	}
	if req.Password != nil {
		config.Password = *req.Password
	}
	if req.DefaultSavePath != nil {
		config.DefaultSavePath = *req.DefaultSavePath
	}
	if req.DefaultCategory != nil {
		config.DefaultCategory = *req.DefaultCategory
	}
	if req.DefaultTags != nil {
		config.DefaultTags = req.DefaultTags
	}
	if req.GlobalDownloadLimit != nil {
		config.GlobalDownloadLimit = *req.GlobalDownloadLimit
	}
	if req.GlobalUploadLimit != nil {
		config.GlobalUploadLimit = *req.GlobalUploadLimit
	}
	if req.AutoConnectOnStart != nil {
		config.AutoConnectOnStart = *req.AutoConnectOnStart
	}
	if req.AdvancedSettings != nil {
		config.AdvancedSettings = req.AdvancedSettings
	}
	if req.Enabled != nil {
		config.Enabled = *req.Enabled
	}

	// 保存更新
	err = h.service.Db.Save(&config).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to update client config: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Client config updated successfully",
		"data":    config,
	})
}

// DeleteClientConfig 删除客户端配置的处理器
func (h *TorrentHandler) DeleteClientConfig(c *gin.Context) {
	clientID := c.Param("id")

	// 获取现有配置
	var config models.ClientConfig
	err := h.service.Db.Where("client_id = ?", clientID).First(&config).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"error":   "Client config not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to get client config: " + err.Error(),
		})
		return
	}

	// 删除配置
	err = h.service.Db.Delete(&config).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to delete client config: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Client config deleted successfully",
	})
}

// TestClientConnection 测试客户端连接的处理器
func (h *TorrentHandler) TestClientConnection(c *gin.Context) {
	clientID := c.Param("id")

	// 调用核心服务测试客户端连接
	err := h.service.TestClientConnection(clientID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "Connection failed: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Connection successful",
	})
}
