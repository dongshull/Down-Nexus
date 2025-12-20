package qbittorrent

import (
	"fmt"
	"strconv"
	"strings"

	"down-nexus-api/pkg/clients"
)

// ==================== 种子添加 ====================

// AddTorrentWithOptions 添加带选项的种子
func (qc *QbitClient) AddTorrentWithOptions(options clients.TorrentOptions) error {
	addOptions := map[string]string{}

	if options.SavePath != "" {
		addOptions["savepath"] = options.SavePath
	}
	if options.Category != "" {
		addOptions["category"] = options.Category
	}
	if len(options.Tags) > 0 {
		addOptions["tags"] = strings.Join(options.Tags, ",")
	}
	if options.Sequential {
		addOptions["sequentialDownload"] = "true"
	}
	if options.FirstLast {
		addOptions["firstLastPiecePrio"] = "true"
	}
	if options.SkipChecking {
		addOptions["skip_checking"] = "true"
	}
	if options.Paused {
		addOptions["paused"] = "true"
	}
	if options.AutoManagement {
		addOptions["autoTMM"] = "true"
	}
	if options.DownloadLimit > 0 {
		addOptions["dlLimit"] = strconv.FormatInt(options.DownloadLimit, 10)
	}
	if options.UploadLimit > 0 {
		addOptions["upLimit"] = strconv.FormatInt(options.UploadLimit, 10)
	}
	if options.RatioLimit > 0 {
		addOptions["ratioLimit"] = strconv.FormatFloat(options.RatioLimit, 'f', 2, 64)
	}
	if options.SeedingTimeLimit > 0 {
		addOptions["seedingTimeLimit"] = strconv.Itoa(options.SeedingTimeLimit)
	}
	if options.Rename != "" {
		addOptions["rename"] = options.Rename
	}
	if options.ContentLayout != "" {
		addOptions["contentLayout"] = options.ContentLayout
	}

	// 如果有种子文件内容，使用 AddTorrentFromMemory
	if len(options.TorrentFile) > 0 {
		return qc.client.AddTorrentFromMemory(options.TorrentFile, addOptions)
	}

	return qc.client.AddTorrentFromUrl(options.URL, addOptions)
}

// AddTorrentFromFile 从文件添加种子
func (qc *QbitClient) AddTorrentFromFile(filePath string, options map[string]string) error {
	return qc.client.AddTorrentFromFile(filePath, options)
}

// ==================== 批量操作 ====================

// PauseTorrents 批量暂停种子
func (qc *QbitClient) PauseTorrents(hashes []string) error {
	return qc.client.Pause(hashes)
}

// ResumeTorrents 批量恢复种子
func (qc *QbitClient) ResumeTorrents(hashes []string) error {
	return qc.client.Resume(hashes)
}

// DeleteTorrents 批量删除种子
func (qc *QbitClient) DeleteTorrents(hashes []string, deleteFiles bool) error {
	return qc.client.DeleteTorrents(hashes, deleteFiles)
}

// ==================== 种子操作 ====================

// RecheckTorrent 校验种子
func (qc *QbitClient) RecheckTorrent(hash string) error {
	return qc.client.Recheck([]string{hash})
}

// RecheckTorrents 批量校验种子
func (qc *QbitClient) RecheckTorrents(hashes []string) error {
	return qc.client.Recheck(hashes)
}

// ReannounceTorrent 重新汇报种子
func (qc *QbitClient) ReannounceTorrent(hash string) error {
	return qc.client.ReAnnounceTorrents([]string{hash})
}

// ReannounceTorrents 批量重新汇报种子
func (qc *QbitClient) ReannounceTorrents(hashes []string) error {
	return qc.client.ReAnnounceTorrents(hashes)
}

// SetTorrentLocation 设置种子保存位置
func (qc *QbitClient) SetTorrentLocation(hash string, location string) error {
	return qc.client.SetLocation([]string{hash}, location)
}

// SetTorrentName 设置种子名称
func (qc *QbitClient) SetTorrentName(hash string, name string) error {
	return qc.client.SetTorrentName(hash, name)
}

// SetForceStart 设置强制开始
func (qc *QbitClient) SetForceStart(hash string, enabled bool) error {
	return qc.client.SetForceStart([]string{hash}, enabled)
}

// SetAutoManagement 设置自动管理
func (qc *QbitClient) SetAutoManagement(hash string, enabled bool) error {
	return qc.client.SetAutoManagement([]string{hash}, enabled)
}

// SetSequentialDownload 设置顺序下载
func (qc *QbitClient) SetSequentialDownload(hash string, enabled bool) error {
	// ToggleTorrentSequentialDownload 是切换操作，需要先检查当前状态
	return qc.client.ToggleTorrentSequentialDownload([]string{hash})
}

// SetFirstLastPiecePriority 设置首尾片优先
func (qc *QbitClient) SetFirstLastPiecePriority(hash string, enabled bool) error {
	return qc.client.ToggleFirstLastPiecePrio([]string{hash})
}

// SetSuperSeeding 设置超级做种
func (qc *QbitClient) SetSuperSeeding(hash string, enabled bool) error {
	return qc.client.SetTorrentSuperSeeding([]string{hash}, enabled)
}

// SetShareLimit 设置分享限制
func (qc *QbitClient) SetShareLimit(hash string, ratioLimit float64, seedingTimeLimit int64) error {
	return qc.client.SetTorrentShareLimit([]string{hash}, ratioLimit, seedingTimeLimit, -1)
}

// ==================== 种子属性 ====================

// SetTorrentCategory 设置种子分类
func (qc *QbitClient) SetTorrentCategory(hash, category string) error {
	return qc.client.SetCategory([]string{hash}, category)
}

// SetTorrentTags 设置种子标签（替换所有标签）
func (qc *QbitClient) SetTorrentTags(hash string, tags []string) error {
	// 先移除所有标签，再添加新标签
	return qc.client.AddTags([]string{hash}, strings.Join(tags, ","))
}

// AddTorrentTags 添加种子标签
func (qc *QbitClient) AddTorrentTags(hash string, tags []string) error {
	return qc.client.AddTags([]string{hash}, strings.Join(tags, ","))
}

// RemoveTorrentTags 移除种子标签
func (qc *QbitClient) RemoveTorrentTags(hash string, tags []string) error {
	return qc.client.RemoveTags([]string{hash}, strings.Join(tags, ","))
}

// SetTorrentTrackers 设置种子Tracker（添加）
func (qc *QbitClient) SetTorrentTrackers(hash string, trackers []string) error {
	return qc.client.AddTrackers(hash, strings.Join(trackers, "\n"))
}

// AddTorrentTrackers 添加种子Tracker
func (qc *QbitClient) AddTorrentTrackers(hash string, trackers []string) error {
	return qc.client.AddTrackers(hash, strings.Join(trackers, "\n"))
}

// RemoveTorrentTrackers 移除种子Tracker
func (qc *QbitClient) RemoveTorrentTrackers(hash string, trackerURLs []string) error {
	// qBittorrent 没有直接的删除tracker方法，需要通过编辑实现
	// 这里暂时返回不支持
	return fmt.Errorf("qBittorrent 不支持直接删除单个 tracker，请使用 Web UI")
}

// SetTorrentDownloadLimit 设置种子下载限速
func (qc *QbitClient) SetTorrentDownloadLimit(hash string, limit int64) error {
	return qc.client.SetTorrentDownloadLimit([]string{hash}, limit)
}

// SetTorrentUploadLimit 设置种子上传限速
func (qc *QbitClient) SetTorrentUploadLimit(hash string, limit int64) error {
	return qc.client.SetTorrentUploadLimit([]string{hash}, limit)
}

// GetTorrentDownloadLimit 获取种子下载限速
func (qc *QbitClient) GetTorrentDownloadLimit(hash string) (int64, error) {
	limits, err := qc.client.GetTorrentDownloadLimit([]string{hash})
	if err != nil {
		return -1, err
	}
	if limit, ok := limits[hash]; ok {
		return limit, nil
	}
	return -1, fmt.Errorf("种子未找到: %s", hash)
}

// GetTorrentUploadLimit 获取种子上传限速
func (qc *QbitClient) GetTorrentUploadLimit(hash string) (int64, error) {
	limits, err := qc.client.GetTorrentUploadLimit([]string{hash})
	if err != nil {
		return -1, err
	}
	if limit, ok := limits[hash]; ok {
		return limit, nil
	}
	return -1, fmt.Errorf("种子未找到: %s", hash)
}

// SetTorrentPriority 设置种子优先级
func (qc *QbitClient) SetTorrentPriority(hash string, priority int) error {
	switch {
	case priority <= 1:
		return qc.client.SetMaxPriority([]string{hash})
	case priority >= 9:
		return qc.client.SetMinPriority([]string{hash})
	case priority < 5:
		return qc.client.IncreasePriority([]string{hash})
	case priority > 5:
		return qc.client.DecreasePriority([]string{hash})
	}
	return nil
}

// ==================== 种子文件操作 ====================

// GetTorrentFiles 获取种子文件列表
func (qc *QbitClient) GetTorrentFiles(hash string) ([]clients.TorrentFile, error) {
	files, err := qc.client.GetFilesInformation(hash)
	if err != nil {
		return nil, err
	}

	var result []clients.TorrentFile
	for i, file := range *files {
		result = append(result, clients.TorrentFile{
			Index:    i,
			Name:     file.Name,
			Size:     file.Size,
			Progress: float64(file.Progress),
			Priority: file.Priority,
		})
	}
	return result, nil
}

// SetFilePriority 设置文件优先级
func (qc *QbitClient) SetFilePriority(hash string, fileIDs []int, priority int) error {
	ids := make([]string, len(fileIDs))
	for i, id := range fileIDs {
		ids[i] = strconv.Itoa(id)
	}
	return qc.client.SetFilePriority(hash, strings.Join(ids, "|"), priority)
}

// RenameFile 重命名文件
func (qc *QbitClient) RenameFile(hash string, oldPath string, newPath string) error {
	return qc.client.RenameFile(hash, oldPath, newPath)
}

// RenameFolder 重命名文件夹
func (qc *QbitClient) RenameFolder(hash string, oldPath string, newPath string) error {
	return qc.client.RenameFolder(hash, oldPath, newPath)
}

// ==================== 种子详细信息 ====================

// GetTorrentTrackers 获取种子Tracker列表
func (qc *QbitClient) GetTorrentTrackers(hash string) ([]clients.TrackerInfo, error) {
	trackers, err := qc.client.GetTorrentTrackers(hash)
	if err != nil {
		return nil, err
	}

	var result []clients.TrackerInfo
	for _, t := range trackers {
		result = append(result, clients.TrackerInfo{
			URL:           t.Url,
			Status:        int(t.Status),
			Tier:          0, // qBittorrent 库移除了 Tier 字段
			Peers:         t.NumPeers,
			Seeds:         t.NumSeeds,
			Leechers:      t.NumLeechers,
			Downloaded:    t.NumDownloaded,
			Message:       t.Message,
			StatusMessage: t.Message,
		})
	}
	return result, nil
}

// GetTorrentPeers 获取种子Peer列表
func (qc *QbitClient) GetTorrentPeers(hash string) ([]clients.PeerInfo, error) {
	// qBittorrent API 不提供直接获取peers的方法
	// 需要通过 SyncMainData 获取
	return nil, fmt.Errorf("请使用 Web UI 查看 Peers 列表")
}

// GetTorrentProperties 获取种子详细属性
func (qc *QbitClient) GetTorrentProperties(hash string) (*clients.TorrentProperties, error) {
	props, err := qc.client.GetTorrentProperties(hash)
	if err != nil {
		return nil, err
	}

	return &clients.TorrentProperties{
		Hash:              hash,
		Name:              props.Name,
		SavePath:          props.SavePath,
		CreationDate:      int64(props.CreationDate),
		PieceSize:         int64(props.PieceSize),
		Comment:           props.Comment,
		TotalWasted:       props.TotalWasted,
		TotalUploaded:     props.TotalUploaded,
		TotalDownloaded:   props.TotalDownloaded,
		UploadedSession:   props.TotalUploadedSession,
		DownloadedSession: props.TotalDownloadedSession,
		UploadLimit:       int64(props.UpLimit),
		DownloadLimit:     int64(props.DlLimit),
		TimeElapsed:       int64(props.TimeElapsed),
		SeedingTime:       int64(props.SeedingTime),
		ConnectionsCount:  props.NbConnections,
		ConnectionsLimit:  props.NbConnectionsLimit,
		ShareRatio:        props.ShareRatio,
		AdditionDate:      int64(props.AdditionDate),
		CompletionDate:    int64(props.CompletionDate),
		CreatedBy:         props.CreatedBy,
		LastActivity:      int64(props.LastSeen),
		PeersCount:        props.Peers,
		SeedsCount:        props.Seeds,
		TotalPeers:        props.PeersTotal,
		TotalSeeds:        props.SeedsTotal,
	}, nil
}

// GetTorrentPieceStates 获取种子分片状态
func (qc *QbitClient) GetTorrentPieceStates(hash string) ([]int, error) {
	states, err := qc.client.GetTorrentPieceStates(hash)
	if err != nil {
		return nil, err
	}

	result := make([]int, len(states))
	for i, s := range states {
		result[i] = int(s)
	}
	return result, nil
}

// ==================== 分类管理 ====================

// GetCategories 获取分类列表
func (qc *QbitClient) GetCategories() (map[string]clients.Category, error) {
	qbCategories, err := qc.client.GetCategories()
	if err != nil {
		return nil, err
	}

	categories := make(map[string]clients.Category)
	for _, cat := range qbCategories {
		categories[cat.Name] = clients.Category{
			Name:     cat.Name,
			SavePath: cat.SavePath,
		}
	}
	return categories, nil
}

// CreateCategory 创建分类
func (qc *QbitClient) CreateCategory(name, savePath string) error {
	return qc.client.CreateCategory(name, savePath)
}

// EditCategory 编辑分类
func (qc *QbitClient) EditCategory(name, savePath string) error {
	return qc.client.EditCategory(name, savePath)
}

// RemoveCategories 删除分类
func (qc *QbitClient) RemoveCategories(names []string) error {
	return qc.client.RemoveCategories(names)
}

// ==================== 标签管理 ====================

// GetTags 获取标签列表
func (qc *QbitClient) GetTags() ([]string, error) {
	return qc.client.GetTags()
}

// CreateTags 创建标签
func (qc *QbitClient) CreateTags(tags []string) error {
	return qc.client.CreateTags(tags)
}

// DeleteTags 删除标签
func (qc *QbitClient) DeleteTags(tags []string) error {
	return qc.client.RemoveTags([]string{}, strings.Join(tags, ","))
}

// ==================== 全局限速 ====================

// GetGlobalDownloadLimit 获取全局下载限速
func (qc *QbitClient) GetGlobalDownloadLimit() (int64, error) {
	limit, err := qc.client.GetGlobalDownloadLimit()
	if err != nil {
		return -1, err
	}
	if limit == 0 {
		return -1, nil // 无限制
	}
	return limit, nil
}

// GetGlobalUploadLimit 获取全局上传限速
func (qc *QbitClient) GetGlobalUploadLimit() (int64, error) {
	limit, err := qc.client.GetGlobalUploadLimit()
	if err != nil {
		return -1, err
	}
	if limit == 0 {
		return -1, nil // 无限制
	}
	return limit, nil
}

// SetGlobalDownloadLimit 设置全局下载限速
func (qc *QbitClient) SetGlobalDownloadLimit(limit int64) error {
	return qc.client.SetGlobalDownloadLimit(limit)
}

// SetGlobalUploadLimit 设置全局上传限速
func (qc *QbitClient) SetGlobalUploadLimit(limit int64) error {
	return qc.client.SetGlobalUploadLimit(limit)
}

// GetAlternativeSpeedLimitsEnabled 获取备用速度限制状态
func (qc *QbitClient) GetAlternativeSpeedLimitsEnabled() (bool, error) {
	return qc.client.GetAlternativeSpeedLimitsMode()
}

// ToggleAlternativeSpeedLimits 切换备用速度限制
func (qc *QbitClient) ToggleAlternativeSpeedLimits() error {
	return qc.client.ToggleAlternativeSpeedLimits()
}

// ==================== 传输信息 ====================

// GetTransferInfo 获取传输信息
func (qc *QbitClient) GetTransferInfo() (*clients.TransferInfo, error) {
	info, err := qc.client.GetTransferInfo()
	if err != nil {
		return nil, err
	}

	return &clients.TransferInfo{
		DownloadSpeed:          info.DlInfoSpeed,
		UploadSpeed:            info.UpInfoSpeed,
		DownloadedBytes:        info.DlInfoData,
		UploadedBytes:          info.UpInfoData,
		DownloadSpeedLimit:     info.DlRateLimit,
		UploadSpeedLimit:       info.UpRateLimit,
		DHT:                    int(info.DHTNodes),
		ConnectionStatus:       string(info.ConnectionStatus),
		AlternativeSpeedLimits: false, // API 不提供此字段
	}, nil
}

// GetFreeSpace 获取磁盘剩余空间
func (qc *QbitClient) GetFreeSpace(path string) (int64, error) {
	return qc.client.GetFreeSpaceOnDisk()
}

// ==================== 服务器操作 ====================

// GetServerInfo 获取服务器信息
func (qc *QbitClient) GetServerInfo() (map[string]interface{}, error) {
	info := make(map[string]interface{})

	if appVersion, err := qc.client.GetAppVersion(); err == nil {
		info["app_version"] = appVersion
	}
	if apiVersion, err := qc.client.GetWebAPIVersion(); err == nil {
		info["api_version"] = apiVersion
	}
	if prefs, err := qc.client.GetAppPreferences(); err == nil {
		info["preferences"] = prefs
	}
	if buildInfo, err := qc.client.GetBuildInfo(); err == nil {
		info["build_info"] = buildInfo
	}

	return info, nil
}

// GetDefaultSavePath 获取默认保存路径
func (qc *QbitClient) GetDefaultSavePath() (string, error) {
	return qc.client.GetDefaultSavePath()
}

// Shutdown 关闭客户端
func (qc *QbitClient) Shutdown() error {
	return qc.client.Shutdown()
}

// ==================== 日志 ====================

// GetLogs 获取日志
func (qc *QbitClient) GetLogs(normal, warning, critical bool, lastKnownID int64) ([]clients.LogEntry, error) {
	logs, err := qc.client.GetLogs()
	if err != nil {
		return nil, err
	}

	var result []clients.LogEntry
	for _, log := range logs {
		if log.ID <= lastKnownID {
			continue
		}
		result = append(result, clients.LogEntry{
			ID:        log.ID,
			Message:   log.Message,
			Timestamp: log.Timestamp,
			Type:      int(log.Type),
		})
	}
	return result, nil
}

// ==================== 种子分片信息 ====================

// GetTorrentPieceStates 获取种子分片状态
func (qc *QbitClient) GetTorrentPieceStates(hash string) ([]int, error) {
	states, err := qc.client.GetPieceStates(hash)
	if err != nil {
		return nil, err
	}
	return states, nil
}

// ==================== 备用速度限制 ====================

// GetAlternativeSpeedLimitsEnabled 获取备用速度限制是否启用
func (qc *QbitClient) GetAlternativeSpeedLimitsEnabled() (bool, error) {
	prefs, err := qc.client.GetAppPreferences()
	if err != nil {
		return false, err
	}
	return prefs.AltSpeedLimitEnabled, nil
}

// ToggleAlternativeSpeedLimits 切换备用速度限制
func (qc *QbitClient) ToggleAlternativeSpeedLimits() error {
	return qc.client.ToggleSpeedLimitsMode()
}
