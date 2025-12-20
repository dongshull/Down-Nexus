package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ==================== 种子校验/重新汇报 ====================

// RecheckTorrent 校验种子
func (h *TorrentHandler) RecheckTorrent(c *gin.Context) {
	var req struct {
		ClientID string `json:"client_id" binding:"required"`
		Hash     string `json:"hash" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	client := h.service.GetClient(req.ClientID)
	if client == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "客户端不存在"})
		return
	}

	if err := client.RecheckTorrent(req.Hash); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "校验已开始"})
}

// RecheckTorrents 批量校验种子
func (h *TorrentHandler) RecheckTorrents(c *gin.Context) {
	var req struct {
		ClientID string   `json:"client_id" binding:"required"`
		Hashes   []string `json:"hashes" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	client := h.service.GetClient(req.ClientID)
	if client == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "客户端不存在"})
		return
	}

	if err := client.RecheckTorrents(req.Hashes); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "批量校验已开始"})
}

// ReannounceTorrent 重新汇报种子
func (h *TorrentHandler) ReannounceTorrent(c *gin.Context) {
	var req struct {
		ClientID string `json:"client_id" binding:"required"`
		Hash     string `json:"hash" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	client := h.service.GetClient(req.ClientID)
	if client == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "客户端不存在"})
		return
	}

	if err := client.ReannounceTorrent(req.Hash); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "重新汇报已开始"})
}

// ReannounceTorrents 批量重新汇报种子
func (h *TorrentHandler) ReannounceTorrents(c *gin.Context) {
	var req struct {
		ClientID string   `json:"client_id" binding:"required"`
		Hashes   []string `json:"hashes" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	client := h.service.GetClient(req.ClientID)
	if client == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "客户端不存在"})
		return
	}

	if err := client.ReannounceTorrents(req.Hashes); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "批量重新汇报已开始"})
}

// ==================== 种子位置/名称 ====================

// SetTorrentLocation 设置种子保存位置
func (h *TorrentHandler) SetTorrentLocation(c *gin.Context) {
	var req struct {
		ClientID string `json:"client_id" binding:"required"`
		Hash     string `json:"hash" binding:"required"`
		Location string `json:"location" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	client := h.service.GetClient(req.ClientID)
	if client == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "客户端不存在"})
		return
	}

	if err := client.SetTorrentLocation(req.Hash, req.Location); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "保存位置已设置"})
}

// SetTorrentName 设置种子名称
func (h *TorrentHandler) SetTorrentName(c *gin.Context) {
	var req struct {
		ClientID string `json:"client_id" binding:"required"`
		Hash     string `json:"hash" binding:"required"`
		Name     string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	client := h.service.GetClient(req.ClientID)
	if client == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "客户端不存在"})
		return
	}

	if err := client.SetTorrentName(req.Hash, req.Name); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "名称已设置"})
}

// ==================== 种子高级设置 ====================

// SetForceStart 设置强制开始
func (h *TorrentHandler) SetForceStart(c *gin.Context) {
	var req struct {
		ClientID string `json:"client_id" binding:"required"`
		Hash     string `json:"hash" binding:"required"`
		Enabled  bool   `json:"enabled"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	client := h.service.GetClient(req.ClientID)
	if client == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "客户端不存在"})
		return
	}

	if err := client.SetForceStart(req.Hash, req.Enabled); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "强制开始已设置"})
}

// SetAutoManagement 设置自动管理
func (h *TorrentHandler) SetAutoManagement(c *gin.Context) {
	var req struct {
		ClientID string `json:"client_id" binding:"required"`
		Hash     string `json:"hash" binding:"required"`
		Enabled  bool   `json:"enabled"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	client := h.service.GetClient(req.ClientID)
	if client == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "客户端不存在"})
		return
	}

	if err := client.SetAutoManagement(req.Hash, req.Enabled); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "自动管理已设置"})
}

// SetSequentialDownload 设置顺序下载
func (h *TorrentHandler) SetSequentialDownload(c *gin.Context) {
	var req struct {
		ClientID string `json:"client_id" binding:"required"`
		Hash     string `json:"hash" binding:"required"`
		Enabled  bool   `json:"enabled"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	client := h.service.GetClient(req.ClientID)
	if client == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "客户端不存在"})
		return
	}

	if err := client.SetSequentialDownload(req.Hash, req.Enabled); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "顺序下载已设置"})
}

// SetFirstLastPiecePriority 设置首尾片优先
func (h *TorrentHandler) SetFirstLastPiecePriority(c *gin.Context) {
	var req struct {
		ClientID string `json:"client_id" binding:"required"`
		Hash     string `json:"hash" binding:"required"`
		Enabled  bool   `json:"enabled"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	client := h.service.GetClient(req.ClientID)
	if client == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "客户端不存在"})
		return
	}

	if err := client.SetFirstLastPiecePriority(req.Hash, req.Enabled); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "首尾片优先已设置"})
}

// SetSuperSeeding 设置超级做种
func (h *TorrentHandler) SetSuperSeeding(c *gin.Context) {
	var req struct {
		ClientID string `json:"client_id" binding:"required"`
		Hash     string `json:"hash" binding:"required"`
		Enabled  bool   `json:"enabled"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	client := h.service.GetClient(req.ClientID)
	if client == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "客户端不存在"})
		return
	}

	if err := client.SetSuperSeeding(req.Hash, req.Enabled); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "超级做种已设置"})
}

// SetShareLimit 设置分享限制
func (h *TorrentHandler) SetShareLimit(c *gin.Context) {
	var req struct {
		ClientID         string  `json:"client_id" binding:"required"`
		Hash             string  `json:"hash" binding:"required"`
		RatioLimit       float64 `json:"ratio_limit"`
		SeedingTimeLimit int64   `json:"seeding_time_limit"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	client := h.service.GetClient(req.ClientID)
	if client == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "客户端不存在"})
		return
	}

	if err := client.SetShareLimit(req.Hash, req.RatioLimit, req.SeedingTimeLimit); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "分享限制已设置"})
}

// ==================== 种子文件操作 ====================

// GetTorrentFiles 获取种子文件列表
func (h *TorrentHandler) GetTorrentFiles(c *gin.Context) {
	clientID := c.Query("client_id")
	hash := c.Param("hash")

	if clientID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "client_id 参数必填"})
		return
	}

	client := h.service.GetClient(clientID)
	if client == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "客户端不存在"})
		return
	}

	files, err := client.GetTorrentFiles(hash)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, files)
}

// SetFilePriority 设置文件优先级
func (h *TorrentHandler) SetFilePriority(c *gin.Context) {
	var req struct {
		ClientID string `json:"client_id" binding:"required"`
		Hash     string `json:"hash" binding:"required"`
		FileIDs  []int  `json:"file_ids" binding:"required"`
		Priority int    `json:"priority"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	client := h.service.GetClient(req.ClientID)
	if client == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "客户端不存在"})
		return
	}

	if err := client.SetFilePriority(req.Hash, req.FileIDs, req.Priority); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "文件优先级已设置"})
}

// RenameFile 重命名文件
func (h *TorrentHandler) RenameFile(c *gin.Context) {
	var req struct {
		ClientID string `json:"client_id" binding:"required"`
		Hash     string `json:"hash" binding:"required"`
		OldPath  string `json:"old_path" binding:"required"`
		NewPath  string `json:"new_path" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	client := h.service.GetClient(req.ClientID)
	if client == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "客户端不存在"})
		return
	}

	if err := client.RenameFile(req.Hash, req.OldPath, req.NewPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "文件已重命名"})
}

// RenameFolder 重命名文件夹
func (h *TorrentHandler) RenameFolder(c *gin.Context) {
	var req struct {
		ClientID string `json:"client_id" binding:"required"`
		Hash     string `json:"hash" binding:"required"`
		OldPath  string `json:"old_path" binding:"required"`
		NewPath  string `json:"new_path" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	client := h.service.GetClient(req.ClientID)
	if client == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "客户端不存在"})
		return
	}

	if err := client.RenameFolder(req.Hash, req.OldPath, req.NewPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "文件夹已重命名"})
}

// ==================== 种子详细信息 ====================

// GetTorrentTrackers 获取种子Tracker列表
func (h *TorrentHandler) GetTorrentTrackers(c *gin.Context) {
	clientID := c.Query("client_id")
	hash := c.Param("hash")

	if clientID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "client_id 参数必填"})
		return
	}

	client := h.service.GetClient(clientID)
	if client == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "客户端不存在"})
		return
	}

	trackers, err := client.GetTorrentTrackers(hash)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, trackers)
}

// AddTorrentTrackers 添加种子Tracker
func (h *TorrentHandler) AddTorrentTrackers(c *gin.Context) {
	var req struct {
		ClientID string   `json:"client_id" binding:"required"`
		Hash     string   `json:"hash" binding:"required"`
		Trackers []string `json:"trackers" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	client := h.service.GetClient(req.ClientID)
	if client == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "客户端不存在"})
		return
	}

	if err := client.AddTorrentTrackers(req.Hash, req.Trackers); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tracker 已添加"})
}

// RemoveTorrentTrackers 移除种子Tracker
func (h *TorrentHandler) RemoveTorrentTrackers(c *gin.Context) {
	var req struct {
		ClientID string   `json:"client_id" binding:"required"`
		Hash     string   `json:"hash" binding:"required"`
		Trackers []string `json:"trackers" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	client := h.service.GetClient(req.ClientID)
	if client == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "客户端不存在"})
		return
	}

	if err := client.RemoveTorrentTrackers(req.Hash, req.Trackers); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tracker 已移除"})
}

// GetTorrentPeers 获取种子Peer列表
func (h *TorrentHandler) GetTorrentPeers(c *gin.Context) {
	clientID := c.Query("client_id")
	hash := c.Param("hash")

	if clientID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "client_id 参数必填"})
		return
	}

	client := h.service.GetClient(clientID)
	if client == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "客户端不存在"})
		return
	}

	peers, err := client.GetTorrentPeers(hash)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, peers)
}

// GetTorrentProperties 获取种子详细属性
func (h *TorrentHandler) GetTorrentProperties(c *gin.Context) {
	clientID := c.Query("client_id")
	hash := c.Param("hash")

	if clientID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "client_id 参数必填"})
		return
	}

	client := h.service.GetClient(clientID)
	if client == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "客户端不存在"})
		return
	}

	props, err := client.GetTorrentProperties(hash)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, props)
}

// ==================== 全局限速 ====================

// GetGlobalLimits 获取全局限速
func (h *TorrentHandler) GetGlobalLimits(c *gin.Context) {
	clientID := c.Query("client_id")

	if clientID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "client_id 参数必填"})
		return
	}

	client := h.service.GetClient(clientID)
	if client == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "客户端不存在"})
		return
	}

	downloadLimit, _ := client.GetGlobalDownloadLimit()
	uploadLimit, _ := client.GetGlobalUploadLimit()
	altEnabled, _ := client.GetAlternativeSpeedLimitsEnabled()

	c.JSON(http.StatusOK, gin.H{
		"download_limit":                   downloadLimit,
		"upload_limit":                     uploadLimit,
		"alternative_speed_limits_enabled": altEnabled,
	})
}

// ToggleAlternativeSpeedLimits 切换备用速度限制
func (h *TorrentHandler) ToggleAlternativeSpeedLimits(c *gin.Context) {
	var req struct {
		ClientID string `json:"client_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	client := h.service.GetClient(req.ClientID)
	if client == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "客户端不存在"})
		return
	}

	if err := client.ToggleAlternativeSpeedLimits(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "备用速度限制已切换"})
}

// ==================== 传输信息 ====================

// GetTransferInfo 获取传输信息
func (h *TorrentHandler) GetTransferInfo(c *gin.Context) {
	clientID := c.Query("client_id")

	if clientID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "client_id 参数必填"})
		return
	}

	client := h.service.GetClient(clientID)
	if client == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "客户端不存在"})
		return
	}

	info, err := client.GetTransferInfo()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, info)
}

// GetFreeSpace 获取磁盘剩余空间
func (h *TorrentHandler) GetFreeSpace(c *gin.Context) {
	clientID := c.Query("client_id")
	path := c.Query("path")

	if clientID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "client_id 参数必填"})
		return
	}

	client := h.service.GetClient(clientID)
	if client == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "客户端不存在"})
		return
	}

	space, err := client.GetFreeSpace(path)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"free_space": space,
		"path":       path,
	})
}

// ==================== 服务器信息 ====================

// GetServerInfo 获取服务器信息
func (h *TorrentHandler) GetServerInfo(c *gin.Context) {
	clientID := c.Query("client_id")

	if clientID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "client_id 参数必填"})
		return
	}

	client := h.service.GetClient(clientID)
	if client == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "客户端不存在"})
		return
	}

	info, err := client.GetServerInfo()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, info)
}

// GetDefaultSavePath 获取默认保存路径
func (h *TorrentHandler) GetDefaultSavePath(c *gin.Context) {
	clientID := c.Query("client_id")

	if clientID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "client_id 参数必填"})
		return
	}

	client := h.service.GetClient(clientID)
	if client == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "客户端不存在"})
		return
	}

	path, err := client.GetDefaultSavePath()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"default_save_path": path,
	})
}

// ShutdownClient 关闭客户端
func (h *TorrentHandler) ShutdownClient(c *gin.Context) {
	var req struct {
		ClientID string `json:"client_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	client := h.service.GetClient(req.ClientID)
	if client == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "客户端不存在"})
		return
	}

	if err := client.Shutdown(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "客户端已关闭"})
}

// ==================== 日志 ====================

// GetLogs 获取日志
func (h *TorrentHandler) GetLogs(c *gin.Context) {
	clientID := c.Query("client_id")
	lastKnownIDStr := c.DefaultQuery("last_known_id", "0")

	if clientID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "client_id 参数必填"})
		return
	}

	lastKnownID, _ := strconv.ParseInt(lastKnownIDStr, 10, 64)

	client := h.service.GetClient(clientID)
	if client == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "客户端不存在"})
		return
	}

	logs, err := client.GetLogs(true, true, true, lastKnownID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, logs)
}

// ==================== 种子分片信息 ====================

// GetTorrentPieceStates 获取种子分片状态
func (h *TorrentHandler) GetTorrentPieceStates(c *gin.Context) {
	var req struct {
		ClientID string `json:"client_id" binding:"required"`
		Hash     string `json:"hash" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	client := h.service.GetClient(req.ClientID)
	if client == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "客户端不存在"})
		return
	}

	states, err := client.GetTorrentPieceStates(req.Hash)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    states,
		"count":   len(states),
	})
}

// ==================== 备用速度限制 ====================

// GetAlternativeSpeedLimits 获取备用速度限制状态
func (h *TorrentHandler) GetAlternativeSpeedLimits(c *gin.Context) {
	clientID := c.Query("client_id")

	if clientID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "client_id 参数必填",
		})
		return
	}

	enabled, err := h.service.GetAlternativeSpeedLimitsEnabled(clientID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to get alternative speed limits: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"enabled": enabled,
		},
	})
}

// ToggleAlternativeSpeedLimits 切换备用速度限制
func (h *TorrentHandler) ToggleAlternativeSpeedLimits(c *gin.Context) {
	var req struct {
		ClientID string `json:"client_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request format: " + err.Error(),
		})
		return
	}

	err := h.service.ToggleAlternativeSpeedLimits(req.ClientID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to toggle alternative speed limits: " + err.Error(),
		})
		return
	}

	// 获取新的状态
	enabled, _ := h.service.GetAlternativeSpeedLimitsEnabled(req.ClientID)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Alternative speed limits toggled successfully",
		"data": gin.H{
			"enabled": enabled,
		},
	})
}
