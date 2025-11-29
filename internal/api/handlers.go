package api

import (
	"net/http"

	"down-nexus-api/internal/core"
	"github.com/gin-gonic/gin"
)

type TorrentHandler struct {
	service *core.TorrentService
}

func NewTorrentHandler(s *core.TorrentService) *TorrentHandler {
	return &TorrentHandler{
		service: s,
	}
}

// GetTorrents 获取所有种子的处理器
func (h *TorrentHandler) GetTorrents(c *gin.Context) {
	// 调用核心服务获取所有种子
	torrents := h.service.GetAllTorrents()

	// 构建响应数据
	response := gin.H{
		"success": true,
		"data":    torrents,
		"count":   len(torrents),
	}

	// 返回 JSON 响应
	c.JSON(http.StatusOK, response)
}

// AddTorrentRequest 添加种子的请求结构
type AddTorrentRequest struct {
	MagnetURL string `json:"magnetURL" binding:"required"`
	ClientID  string `json:"clientID" binding:"required"`
}

// AddTorrent 添加种子的处理器
func (h *TorrentHandler) AddTorrent(c *gin.Context) {
	// 解析请求体
	var req AddTorrentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request format: " + err.Error(),
		})
		return
	}

	// 调用核心服务添加种子
	err := h.service.AddTorrent(req.MagnetURL, req.ClientID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to add torrent: " + err.Error(),
		})
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Torrent added successfully",
	})
}

// GetClients 获取所有客户端信息的处理器
func (h *TorrentHandler) GetClients(c *gin.Context) {
	// 获取所有种子来推断客户端信息
	torrents := h.service.GetAllTorrents()
	
	// 使用 map 去重收集客户端信息
	clients := make(map[string]map[string]interface{})
	for _, torrent := range torrents {
		if _, exists := clients[torrent.ClientID]; !exists {
			clients[torrent.ClientID] = map[string]interface{}{
				"id":     torrent.ClientID,
				"name":   torrent.ClientID, // 可以后续优化为更友好的名称
				"status": "online",
			}
		}
	}

	// 转换为切片
	clientList := make([]map[string]interface{}, 0, len(clients))
	for _, client := range clients {
		clientList = append(clientList, client)
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    clientList,
		"count":   len(clientList),
	})
}

// TorrentControlRequest 控制种子的请求结构
type TorrentControlRequest struct {
	ClientID string `json:"clientID" binding:"required"`
	Hash     string `json:"hash" binding:"required"`
}

// DeleteTorrentRequest 删除种子的请求结构
type DeleteTorrentRequest struct {
	ClientID   string `json:"clientID" binding:"required"`
	Hash       string `json:"hash" binding:"required"`
	DeleteFiles bool  `json:"deleteFiles"`
}

// PauseTorrent 暂停种子的处理器
func (h *TorrentHandler) PauseTorrent(c *gin.Context) {
	// 解析请求体
	var req TorrentControlRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request format: " + err.Error(),
		})
		return
	}

	// 调用核心服务暂停种子
	err := h.service.PauseTorrent(req.ClientID, req.Hash)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to pause torrent: " + err.Error(),
		})
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Torrent paused successfully",
	})
}

// ResumeTorrent 恢复种子的处理器
func (h *TorrentHandler) ResumeTorrent(c *gin.Context) {
	// 解析请求体
	var req TorrentControlRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request format: " + err.Error(),
		})
		return
	}

	// 调用核心服务恢复种子
	err := h.service.ResumeTorrent(req.ClientID, req.Hash)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to resume torrent: " + err.Error(),
		})
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Torrent resumed successfully",
	})
}

// DeleteTorrent 删除种子的处理器
func (h *TorrentHandler) DeleteTorrent(c *gin.Context) {
	// 解析请求体
	var req DeleteTorrentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request format: " + err.Error(),
		})
		return
	}

	// 调用核心服务删除种子
	err := h.service.DeleteTorrent(req.ClientID, req.Hash, req.DeleteFiles)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to delete torrent: " + err.Error(),
		})
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Torrent deleted successfully",
	})
}