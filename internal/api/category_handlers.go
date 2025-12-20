package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetCategories 获取分类列表的处理器
func (h *TorrentHandler) GetCategories(c *gin.Context) {
	clientID := c.Query("client_id")
	
	if clientID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "client_id is required",
		})
		return
	}
	
	// 调用核心服务获取分类列表
	categories, err := h.service.GetCategories(clientID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to get categories: " + err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    categories,
		"count":   len(categories),
	})
}

// CreateCategory 创建分类的处理器
func (h *TorrentHandler) CreateCategory(c *gin.Context) {
	clientID := c.Query("client_id")
	name := c.Query("name")
	savePath := c.Query("save_path")
	
	if clientID == "" || name == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "client_id and name are required",
		})
		return
	}
	
	// 调用核心服务创建分类
	err := h.service.CreateCategory(clientID, name, savePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to create category: " + err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Category created successfully",
	})
}

// DeleteCategories 删除分类的处理器
func (h *TorrentHandler) DeleteCategories(c *gin.Context) {
	clientID := c.Query("client_id")
	
	var req struct {
		Names []string `json:"names" binding:"required"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request format: " + err.Error(),
		})
		return
	}
	
	if clientID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "client_id is required",
		})
		return
	}
	
	// 调用核心服务删除分类
	err := h.service.DeleteCategories(clientID, req.Names)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to delete categories: " + err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Categories deleted successfully",
		"count":   len(req.Names),
	})
}

// GetTags 获取标签列表的处理器
func (h *TorrentHandler) GetTags(c *gin.Context) {
	clientID := c.Query("client_id")
	
	if clientID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "client_id is required",
		})
		return
	}
	
	// 调用核心服务获取标签列表
	tags, err := h.service.GetTags(clientID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to get tags: " + err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    tags,
		"count":   len(tags),
	})
}

// CreateTags 创建标签的处理器
func (h *TorrentHandler) CreateTags(c *gin.Context) {
	clientID := c.Query("client_id")
	
	var req struct {
		Tags []string `json:"tags" binding:"required"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request format: " + err.Error(),
		})
		return
	}
	
	if clientID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "client_id is required",
		})
		return
	}
	
	// 调用核心服务创建标签
	err := h.service.CreateTags(clientID, req.Tags)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to create tags: " + err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Tags created successfully",
		"count":   len(req.Tags),
	})
}

// DeleteTags 删除标签的处理器
func (h *TorrentHandler) DeleteTags(c *gin.Context) {
	clientID := c.Query("client_id")
	
	var req struct {
		Tags []string `json:"tags" binding:"required"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request format: " + err.Error(),
		})
		return
	}
	
	if clientID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "client_id is required",
		})
		return
	}
	
	// 调用核心服务删除标签
	err := h.service.DeleteTags(clientID, req.Tags)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to delete tags: " + err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Tags deleted successfully",
		"count":   len(req.Tags),
	})
}