package api

import (
	"net/http"

	"down-nexus-api/internal/core"
	"down-nexus-api/pkg/clients"

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
	// 获取查询参数
	clientID := c.Query("client_id")
	category := c.Query("category")
	state := c.Query("state")
	tag := c.Query("tag")

	// 调用核心服务获取种子
	torrents, err := h.service.GetFilteredTorrents(clientID, category, state, tag)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to get torrents: " + err.Error(),
		})
		return
	}

	// 构建响应数据
	response := gin.H{
		"success": true,
		"data":    torrents,
		"count":   len(torrents),
	}

	// 返回 JSON 响应
	c.JSON(http.StatusOK, response)
}

// GetTorrentDetails 获取种子详细信息的处理器
func (h *TorrentHandler) GetTorrentDetails(c *gin.Context) {
	hash := c.Param("hash")
	clientID := c.Query("client_id")

	if clientID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "client_id is required",
		})
		return
	}

	// 调用核心服务获取种子详细信息
	torrent, err := h.service.GetTorrentDetails(clientID, hash)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to get torrent details: " + err.Error(),
		})
		return
	}

	// 构建响应数据
	response := gin.H{
		"success": true,
		"data":    torrent,
	}

	// 返回 JSON 响应
	c.JSON(http.StatusOK, response)
}

// AddTorrentRequest 添加种子的请求结构
type AddTorrentRequest struct {
	ClientID   string   `json:"client_id" binding:"required"`
	URL        string   `json:"url" binding:"required"`
	SavePath   string   `json:"save_path"`
	Category   string   `json:"category"`
	Tags       []string `json:"tags"`
	Priority   int      `json:"priority"`
	Sequential bool     `json:"sequential"`
	FirstLast  bool     `json:"first_last_piece"`
	Paused     bool     `json:"paused"`
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

	// 创建种子选项
	options := clients.TorrentOptions{
		URL:        req.URL,
		SavePath:   req.SavePath,
		Category:   req.Category,
		Tags:       req.Tags,
		Priority:   req.Priority,
		Sequential: req.Sequential,
		FirstLast:  req.FirstLast,
		Paused:     req.Paused,
	}

	// 调用核心服务添加种子
	err := h.service.AddTorrentWithOptions(req.ClientID, options)
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
	// 直接从数据库获取客户端配置
	clientConfigs, err := h.service.GetClientConfigs()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to get client configs: " + err.Error(),
		})
		return
	}

	// 转换为 API 响应格式
	clientList := make([]map[string]interface{}, 0, len(clientConfigs))
	for _, config := range clientConfigs {
		client := map[string]interface{}{
			"id":       config.ClientID,
			"name":     config.ClientID, // 可以后续优化为更友好的名称
			"type":     config.Type,
			"host":     config.Host,
			"username": config.Username,
			"enabled":  config.Enabled,
			"status":   "configured", // 表示已配置
		}
		clientList = append(clientList, client)
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    clientList,
		"count":   len(clientList),
	})
}

// TorrentControlRequest 控制单个种子的请求结构
type TorrentControlRequest struct {
	ClientID string `json:"client_id" binding:"required"`
	Hash     string `json:"hash" binding:"required"`
}

// TorrentBatchControlRequest 控制多个种子的请求结构
type TorrentBatchControlRequest struct {
	ClientID string   `json:"client_id" binding:"required"`
	Hashes   []string `json:"hashes" binding:"required"`
}

// DeleteTorrentRequest 删除种子的请求结构
type DeleteTorrentRequest struct {
	ClientID    string `json:"client_id" binding:"required"`
	Hash        string `json:"hash" binding:"required"`
	DeleteFiles bool   `json:"delete_files"`
}

// DeleteBatchTorrentRequest 批量删除种子的请求结构
type DeleteBatchTorrentRequest struct {
	ClientID    string   `json:"client_id" binding:"required"`
	Hashes      []string `json:"hashes" binding:"required"`
	DeleteFiles bool     `json:"delete_files"`
}

// SetTorrentCategoryRequest 设置种子分类的请求结构
type SetTorrentCategoryRequest struct {
	ClientID string `json:"client_id" binding:"required"`
	Hash     string `json:"hash" binding:"required"`
	Category string `json:"category"`
}

// SetTorrentTagsRequest 设置种子标签的请求结构
type SetTorrentTagsRequest struct {
	ClientID string   `json:"client_id" binding:"required"`
	Hash     string   `json:"hash" binding:"required"`
	Tags     []string `json:"tags"`
}

// SetTorrentLimitsRequest 设置种子限速的请求结构
type SetTorrentLimitsRequest struct {
	ClientID      string `json:"client_id" binding:"required"`
	Hash          string `json:"hash" binding:"required"`
	DownloadLimit int64  `json:"download_limit"` // -1表示无限制
	UploadLimit   int64  `json:"upload_limit"`   // -1表示无限制
}

// SetTorrentPriorityRequest 设置种子优先级的请求结构
type SetTorrentPriorityRequest struct {
	ClientID string `json:"client_id" binding:"required"`
	Hash     string `json:"hash" binding:"required"`
	Priority int    `json:"priority"` // 0-7: 0为最高优先级，7为最低优先级
}

// ClientConfigCreateRequest 创建客户端配置的请求结构
type ClientConfigCreateRequest struct {
	ClientID            string                 `json:"client_id" binding:"required"`
	Type                string                 `json:"type" binding:"required"`
	Host                string                 `json:"host" binding:"required"`
	Username            string                 `json:"username" binding:"required"`
	Password            string                 `json:"password" binding:"required"`
	DisplayName         string                 `json:"display_name"`
	DefaultSavePath     string                 `json:"default_save_path"`
	DefaultCategory     string                 `json:"default_category"`
	DefaultTags         []string               `json:"default_tags"`
	GlobalDownloadLimit int64                  `json:"global_download_limit"`
	GlobalUploadLimit   int64                  `json:"global_upload_limit"`
	AutoConnectOnStart  bool                   `json:"auto_connect_on_start"`
	AdvancedSettings    map[string]interface{} `json:"advanced_settings"`
	Enabled             bool                   `json:"enabled"`
}

// ClientConfigUpdateRequest 更新客户端配置的请求结构
type ClientConfigUpdateRequest struct {
	DisplayName         *string                `json:"display_name"`
	Host                *string                `json:"host"`
	Username            *string                `json:"username"`
	Password            *string                `json:"password"`
	DefaultSavePath     *string                `json:"default_save_path"`
	DefaultCategory     *string                `json:"default_category"`
	DefaultTags         []string               `json:"default_tags"`
	GlobalDownloadLimit *int64                 `json:"global_download_limit"`
	GlobalUploadLimit   *int64                 `json:"global_upload_limit"`
	AutoConnectOnStart  *bool                  `json:"auto_connect_on_start"`
	AdvancedSettings    map[string]interface{} `json:"advanced_settings"`
	Enabled             *bool                  `json:"enabled"`
}

// GlobalSpeedLimitRequest 设置全局限速的请求结构
type GlobalSpeedLimitRequest struct {
	ClientID      string `json:"client_id" binding:"required"`
	DownloadLimit int64  `json:"download_limit"` // -1表示无限制
	UploadLimit   int64  `json:"upload_limit"`   // -1表示无限制
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

// PauseTorrents 批量暂停种子的处理器
func (h *TorrentHandler) PauseTorrents(c *gin.Context) {
	// 解析请求体
	var req TorrentBatchControlRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request format: " + err.Error(),
		})
		return
	}

	// 调用核心服务批量暂停种子
	err := h.service.PauseTorrents(req.ClientID, req.Hashes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to pause torrents: " + err.Error(),
		})
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Torrents paused successfully",
		"count":   len(req.Hashes),
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

// ResumeTorrents 批量恢复种子的处理器
func (h *TorrentHandler) ResumeTorrents(c *gin.Context) {
	// 解析请求体
	var req TorrentBatchControlRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request format: " + err.Error(),
		})
		return
	}

	// 调用核心服务批量恢复种子
	err := h.service.ResumeTorrents(req.ClientID, req.Hashes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to resume torrents: " + err.Error(),
		})
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Torrents resumed successfully",
		"count":   len(req.Hashes),
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

// DeleteTorrents 批量删除种子的处理器
func (h *TorrentHandler) DeleteTorrents(c *gin.Context) {
	// 解析请求体
	var req DeleteBatchTorrentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request format: " + err.Error(),
		})
		return
	}

	// 调用核心服务批量删除种子
	err := h.service.DeleteTorrents(req.ClientID, req.Hashes, req.DeleteFiles)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to delete torrents: " + err.Error(),
		})
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Torrents deleted successfully",
		"count":   len(req.Hashes),
	})
}

// SetTorrentCategory 设置种子分类的处理器
func (h *TorrentHandler) SetTorrentCategory(c *gin.Context) {
	// 解析请求体
	var req SetTorrentCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request format: " + err.Error(),
		})
		return
	}

	// 调用核心服务设置种子分类
	err := h.service.SetTorrentCategory(req.ClientID, req.Hash, req.Category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to set torrent category: " + err.Error(),
		})
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Torrent category updated successfully",
	})
}

// SetTorrentTags 设置种子标签的处理器
func (h *TorrentHandler) SetTorrentTags(c *gin.Context) {
	// 解析请求体
	var req SetTorrentTagsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request format: " + err.Error(),
		})
		return
	}

	// 调用核心服务设置种子标签
	err := h.service.SetTorrentTags(req.ClientID, req.Hash, req.Tags)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to set torrent tags: " + err.Error(),
		})
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Torrent tags updated successfully",
	})
}

// SetTorrentLimits 设置种子限速的处理器
func (h *TorrentHandler) SetTorrentLimits(c *gin.Context) {
	// 解析请求体
	var req SetTorrentLimitsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request format: " + err.Error(),
		})
		return
	}

	// 调用核心服务设置种子下载限速
	err := h.service.SetTorrentDownloadLimit(req.ClientID, req.Hash, req.DownloadLimit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to set torrent download limit: " + err.Error(),
		})
		return
	}

	// 调用核心服务设置种子上传限速
	err = h.service.SetTorrentUploadLimit(req.ClientID, req.Hash, req.UploadLimit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to set torrent upload limit: " + err.Error(),
		})
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Torrent limits updated successfully",
	})
}

// SetTorrentPriority 设置种子优先级的处理器
func (h *TorrentHandler) SetTorrentPriority(c *gin.Context) {
	// 解析请求体
	var req SetTorrentPriorityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request format: " + err.Error(),
		})
		return
	}

	// 调用核心服务设置种子优先级
	err := h.service.SetTorrentPriority(req.ClientID, req.Hash, req.Priority)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to set torrent priority: " + err.Error(),
		})
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Torrent priority updated successfully",
	})
}

// SetGlobalSpeedLimits 设置全局限速的处理器
func (h *TorrentHandler) SetGlobalSpeedLimits(c *gin.Context) {
	// 解析请求体
	var req GlobalSpeedLimitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request format: " + err.Error(),
		})
		return
	}

	// 调用核心服务设置全局下载限速
	err := h.service.SetGlobalDownloadLimit(req.ClientID, req.DownloadLimit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to set global download limit: " + err.Error(),
		})
		return
	}

	// 调用核心服务设置全局上传限速
	err = h.service.SetGlobalUploadLimit(req.ClientID, req.UploadLimit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to set global upload limit: " + err.Error(),
		})
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Global speed limits updated successfully",
	})
}
