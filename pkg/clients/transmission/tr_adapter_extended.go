package transmission

import (
	"context"
	"fmt"

	"down-nexus-api/pkg/clients"

	tr "github.com/hekmon/transmissionrpc/v2"
)

// ==================== 种子添加 ====================

// AddTorrentWithOptions 添加带选项的种子
func (tc *TransmissionClient) AddTorrentWithOptions(options clients.TorrentOptions) error {
	addArgs := tr.TorrentAddPayload{
		Filename: &options.URL,
	}

	if options.SavePath != "" {
		addArgs.DownloadDir = &options.SavePath
	}
	if options.Paused {
		addArgs.Paused = &options.Paused
	}
	if options.DownloadLimit > 0 {
		limitKBps := options.DownloadLimit / 1024
		addArgs.BandwidthPriority = &limitKBps
	}

	torrent, err := tc.client.TorrentAdd(context.Background(), addArgs)
	if err != nil {
		return err
	}

	// 如果设置了标签
	if len(options.Tags) > 0 && torrent.ID != nil {
		err = tc.client.TorrentSet(context.Background(), tr.TorrentSetPayload{
			IDs:    []int64{*torrent.ID},
			Labels: options.Tags,
		})
		if err != nil {
			return fmt.Errorf("设置标签失败: %w", err)
		}
	}

	return nil
}

// AddTorrentFromFile 从文件添加种子
func (tc *TransmissionClient) AddTorrentFromFile(filePath string, options map[string]string) error {
	_, err := tc.client.TorrentAddFile(context.Background(), filePath)
	return err
}

// ==================== 批量操作 ====================

// PauseTorrents 批量暂停种子
func (tc *TransmissionClient) PauseTorrents(hashes []string) error {
	return tc.client.TorrentStopHashes(context.Background(), hashes)
}

// ResumeTorrents 批量恢复种子
func (tc *TransmissionClient) ResumeTorrents(hashes []string) error {
	return tc.client.TorrentStartHashes(context.Background(), hashes)
}

// DeleteTorrents 批量删除种子
func (tc *TransmissionClient) DeleteTorrents(hashes []string, deleteFiles bool) error {
	ids, err := tc.hashesToIDs(hashes)
	if err != nil {
		return err
	}
	return tc.client.TorrentRemove(context.Background(), tr.TorrentRemovePayload{
		IDs:             ids,
		DeleteLocalData: deleteFiles,
	})
}

// ==================== 种子操作 ====================

// RecheckTorrent 校验种子
func (tc *TransmissionClient) RecheckTorrent(hash string) error {
	return tc.client.TorrentVerifyHashes(context.Background(), []string{hash})
}

// RecheckTorrents 批量校验种子
func (tc *TransmissionClient) RecheckTorrents(hashes []string) error {
	return tc.client.TorrentVerifyHashes(context.Background(), hashes)
}

// ReannounceTorrent 重新汇报种子
func (tc *TransmissionClient) ReannounceTorrent(hash string) error {
	return tc.client.TorrentReannounceHashes(context.Background(), []string{hash})
}

// ReannounceTorrents 批量重新汇报种子
func (tc *TransmissionClient) ReannounceTorrents(hashes []string) error {
	return tc.client.TorrentReannounceHashes(context.Background(), hashes)
}

// SetTorrentLocation 设置种子保存位置
func (tc *TransmissionClient) SetTorrentLocation(hash string, location string) error {
	return tc.client.TorrentSetLocationHash(context.Background(), hash, location, true)
}

// SetTorrentName 设置种子名称（通过重命名路径）
func (tc *TransmissionClient) SetTorrentName(hash string, name string) error {
	torrents, err := tc.client.TorrentGetAllForHashes(context.Background(), []string{hash})
	if err != nil {
		return err
	}
	if len(torrents) == 0 {
		return fmt.Errorf("种子未找到: %s", hash)
	}
	if torrents[0].Name != nil {
		return tc.client.TorrentRenamePathHash(context.Background(), hash, *torrents[0].Name, name)
	}
	return fmt.Errorf("无法获取种子名称")
}

// SetForceStart 设置强制开始（Transmission 使用 StartNow）
func (tc *TransmissionClient) SetForceStart(hash string, enabled bool) error {
	if enabled {
		return tc.client.TorrentStartNowHashes(context.Background(), []string{hash})
	}
	return tc.client.TorrentStartHashes(context.Background(), []string{hash})
}

// SetAutoManagement 设置自动管理（Transmission 不支持）
func (tc *TransmissionClient) SetAutoManagement(hash string, enabled bool) error {
	return nil // Transmission 没有自动管理功能
}

// SetSequentialDownload 设置顺序下载（Transmission 不支持）
func (tc *TransmissionClient) SetSequentialDownload(hash string, enabled bool) error {
	return nil // Transmission 不支持顺序下载
}

// SetFirstLastPiecePriority 设置首尾片优先（Transmission 不支持）
func (tc *TransmissionClient) SetFirstLastPiecePriority(hash string, enabled bool) error {
	return nil // Transmission 不支持首尾片优先
}

// SetSuperSeeding 设置超级做种（Transmission 不支持）
func (tc *TransmissionClient) SetSuperSeeding(hash string, enabled bool) error {
	return nil // Transmission 不支持超级做种
}

// SetShareLimit 设置分享限制
func (tc *TransmissionClient) SetShareLimit(hash string, ratioLimit float64, seedingTimeLimit int64) error {
	torrentID, err := tc.getTorrentIDByHash(hash)
	if err != nil {
		return err
	}
	return tc.client.TorrentSet(context.Background(), tr.TorrentSetPayload{
		IDs:            []int64{torrentID},
		SeedRatioLimit: &ratioLimit,
	})
}

// ==================== 种子属性 ====================

// SetTorrentCategory 设置种子分类
func (tc *TransmissionClient) SetTorrentCategory(hash, category string) error {
	torrentID, err := tc.getTorrentIDByHash(hash)
	if err != nil {
		return err
	}
	return tc.client.TorrentSet(context.Background(), tr.TorrentSetPayload{
		IDs:    []int64{torrentID},
		Labels: []string{category},
	})
}

// SetTorrentTags 设置种子标签
func (tc *TransmissionClient) SetTorrentTags(hash string, tags []string) error {
	torrentID, err := tc.getTorrentIDByHash(hash)
	if err != nil {
		return err
	}
	return tc.client.TorrentSet(context.Background(), tr.TorrentSetPayload{
		IDs:    []int64{torrentID},
		Labels: tags,
	})
}

// AddTorrentTags 添加种子标签
func (tc *TransmissionClient) AddTorrentTags(hash string, tags []string) error {
	torrents, err := tc.client.TorrentGetAllForHashes(context.Background(), []string{hash})
	if err != nil {
		return err
	}
	if len(torrents) == 0 {
		return fmt.Errorf("种子未找到: %s", hash)
	}

	// 合并现有标签和新标签
	existingTags := torrents[0].Labels
	allTags := append(existingTags, tags...)

	torrentID, _ := tc.getTorrentIDByHash(hash)
	return tc.client.TorrentSet(context.Background(), tr.TorrentSetPayload{
		IDs:    []int64{torrentID},
		Labels: allTags,
	})
}

// RemoveTorrentTags 移除种子标签
func (tc *TransmissionClient) RemoveTorrentTags(hash string, tags []string) error {
	torrents, err := tc.client.TorrentGetAllForHashes(context.Background(), []string{hash})
	if err != nil {
		return err
	}
	if len(torrents) == 0 {
		return fmt.Errorf("种子未找到: %s", hash)
	}

	// 过滤掉要移除的标签
	tagSet := make(map[string]bool)
	for _, t := range tags {
		tagSet[t] = true
	}
	var newTags []string
	for _, t := range torrents[0].Labels {
		if !tagSet[t] {
			newTags = append(newTags, t)
		}
	}

	torrentID, _ := tc.getTorrentIDByHash(hash)
	return tc.client.TorrentSet(context.Background(), tr.TorrentSetPayload{
		IDs:    []int64{torrentID},
		Labels: newTags,
	})
}

// SetTorrentTrackers 设置种子Tracker
func (tc *TransmissionClient) SetTorrentTrackers(hash string, trackers []string) error {
	torrentID, err := tc.getTorrentIDByHash(hash)
	if err != nil {
		return err
	}
	return tc.client.TorrentSet(context.Background(), tr.TorrentSetPayload{
		IDs:        []int64{torrentID},
		TrackerAdd: trackers,
	})
}

// AddTorrentTrackers 添加种子Tracker
func (tc *TransmissionClient) AddTorrentTrackers(hash string, trackers []string) error {
	return tc.SetTorrentTrackers(hash, trackers)
}

// RemoveTorrentTrackers 移除种子Tracker
func (tc *TransmissionClient) RemoveTorrentTrackers(hash string, trackerURLs []string) error {
	torrents, err := tc.client.TorrentGetAllForHashes(context.Background(), []string{hash})
	if err != nil {
		return err
	}
	if len(torrents) == 0 {
		return fmt.Errorf("种子未找到: %s", hash)
	}

	// 找出要删除的 tracker IDs
	var trackerIDs []int64
	urlSet := make(map[string]bool)
	for _, url := range trackerURLs {
		urlSet[url] = true
	}

	if torrents[0].Trackers != nil {
		for _, t := range torrents[0].Trackers {
			if t != nil && urlSet[t.Announce] {
				trackerIDs = append(trackerIDs, int64(t.ID))
			}
		}
	}

	if len(trackerIDs) == 0 {
		return nil
	}

	torrentID, _ := tc.getTorrentIDByHash(hash)
	return tc.client.TorrentSet(context.Background(), tr.TorrentSetPayload{
		IDs:           []int64{torrentID},
		TrackerRemove: trackerIDs,
	})
}

// SetTorrentDownloadLimit 设置种子下载限速
func (tc *TransmissionClient) SetTorrentDownloadLimit(hash string, limit int64) error {
	torrentID, err := tc.getTorrentIDByHash(hash)
	if err != nil {
		return err
	}

	limitKBps := limit / 1024
	enabled := limit > 0

	return tc.client.TorrentSet(context.Background(), tr.TorrentSetPayload{
		IDs:             []int64{torrentID},
		DownloadLimited: &enabled,
		DownloadLimit:   &limitKBps,
	})
}

// SetTorrentUploadLimit 设置种子上传限速
func (tc *TransmissionClient) SetTorrentUploadLimit(hash string, limit int64) error {
	torrentID, err := tc.getTorrentIDByHash(hash)
	if err != nil {
		return err
	}

	limitKBps := limit / 1024
	enabled := limit > 0

	return tc.client.TorrentSet(context.Background(), tr.TorrentSetPayload{
		IDs:           []int64{torrentID},
		UploadLimited: &enabled,
		UploadLimit:   &limitKBps,
	})
}

// GetTorrentDownloadLimit 获取种子下载限速
func (tc *TransmissionClient) GetTorrentDownloadLimit(hash string) (int64, error) {
	torrents, err := tc.client.TorrentGetAllForHashes(context.Background(), []string{hash})
	if err != nil {
		return -1, err
	}
	if len(torrents) == 0 {
		return -1, fmt.Errorf("种子未找到: %s", hash)
	}
	if torrents[0].DownloadLimited != nil && *torrents[0].DownloadLimited && torrents[0].DownloadLimit != nil {
		return *torrents[0].DownloadLimit * 1024, nil
	}
	return -1, nil
}

// GetTorrentUploadLimit 获取种子上传限速
func (tc *TransmissionClient) GetTorrentUploadLimit(hash string) (int64, error) {
	torrents, err := tc.client.TorrentGetAllForHashes(context.Background(), []string{hash})
	if err != nil {
		return -1, err
	}
	if len(torrents) == 0 {
		return -1, fmt.Errorf("种子未找到: %s", hash)
	}
	if torrents[0].UploadLimited != nil && *torrents[0].UploadLimited && torrents[0].UploadLimit != nil {
		return *torrents[0].UploadLimit * 1024, nil
	}
	return -1, nil
}

// SetTorrentPriority 设置种子优先级
func (tc *TransmissionClient) SetTorrentPriority(hash string, priority int) error {
	torrentID, err := tc.getTorrentIDByHash(hash)
	if err != nil {
		return err
	}

	// Transmission 优先级: -1=低, 0=正常, 1=高
	var priorityInt int64
	switch {
	case priority <= 3:
		priorityInt = 1
	case priority >= 7:
		priorityInt = -1
	default:
		priorityInt = 0
	}

	return tc.client.TorrentSet(context.Background(), tr.TorrentSetPayload{
		IDs:               []int64{torrentID},
		BandwidthPriority: &priorityInt,
	})
}

// ==================== 种子文件操作 ====================

// GetTorrentFiles 获取种子文件列表
func (tc *TransmissionClient) GetTorrentFiles(hash string) ([]clients.TorrentFile, error) {
	torrents, err := tc.client.TorrentGetAllForHashes(context.Background(), []string{hash})
	if err != nil {
		return nil, err
	}
	if len(torrents) == 0 {
		return nil, fmt.Errorf("种子未找到: %s", hash)
	}

	var result []clients.TorrentFile
	if torrents[0].Files != nil {
		for i, file := range torrents[0].Files {
			if file != nil {
				wanted := true
				if torrents[0].Wanted != nil && i < len(torrents[0].Wanted) {
					wanted = torrents[0].Wanted[i]
				}

				priority := 1 // 默认普通
				if torrents[0].Priorities != nil && i < len(torrents[0].Priorities) {
					priority = int(torrents[0].Priorities[i])
				}

				var progress float64 = 0
				if file.Length > 0 {
					progress = float64(file.BytesCompleted) / float64(file.Length)
				}

				result = append(result, clients.TorrentFile{
					Index:      i,
					Name:       file.Name,
					Size:       file.Length,
					Progress:   progress,
					Priority:   priority,
					Downloaded: file.BytesCompleted,
					Wanted:     wanted,
				})
			}
		}
	}
	return result, nil
}

// SetFilePriority 设置文件优先级
func (tc *TransmissionClient) SetFilePriority(hash string, fileIDs []int, priority int) error {
	torrentID, err := tc.getTorrentIDByHash(hash)
	if err != nil {
		return err
	}

	ids := make([]int64, len(fileIDs))
	for i, id := range fileIDs {
		ids[i] = int64(id)
	}

	payload := tr.TorrentSetPayload{
		IDs: []int64{torrentID},
	}

	switch priority {
	case 0: // 不下载
		payload.FilesUnwanted = ids
	case 7, 6: // 高优先级
		payload.PriorityHigh = ids
		payload.FilesWanted = ids
	case -1: // 低优先级
		payload.PriorityLow = ids
		payload.FilesWanted = ids
	default: // 普通优先级
		payload.PriorityNormal = ids
		payload.FilesWanted = ids
	}

	return tc.client.TorrentSet(context.Background(), payload)
}

// RenameFile 重命名文件
func (tc *TransmissionClient) RenameFile(hash string, oldPath string, newPath string) error {
	return tc.client.TorrentRenamePathHash(context.Background(), hash, oldPath, newPath)
}

// RenameFolder 重命名文件夹
func (tc *TransmissionClient) RenameFolder(hash string, oldPath string, newPath string) error {
	return tc.client.TorrentRenamePathHash(context.Background(), hash, oldPath, newPath)
}

// ==================== 种子详细信息 ====================

// GetTorrentTrackers 获取种子Tracker列表
func (tc *TransmissionClient) GetTorrentTrackers(hash string) ([]clients.TrackerInfo, error) {
	torrents, err := tc.client.TorrentGetAllForHashes(context.Background(), []string{hash})
	if err != nil {
		return nil, err
	}
	if len(torrents) == 0 {
		return nil, fmt.Errorf("种子未找到: %s", hash)
	}

	var result []clients.TrackerInfo
	if torrents[0].TrackerStats != nil {
		for _, t := range torrents[0].TrackerStats {
			if t != nil {
				result = append(result, clients.TrackerInfo{
					URL:           t.Announce,
					Status:        0,
					StatusMessage: t.LastAnnounceResult,
					Tier:          int(t.Tier),
					Peers:         int(t.LastAnnouncePeerCount),
					Seeds:         int(t.SeederCount),
					Leechers:      int(t.LeecherCount),
					Downloaded:    int(t.DownloadCount),
					Message:       t.LastAnnounceResult,
				})
			}
		}
	}
	return result, nil
}

// GetTorrentPeers 获取种子Peer列表
func (tc *TransmissionClient) GetTorrentPeers(hash string) ([]clients.PeerInfo, error) {
	torrents, err := tc.client.TorrentGetAllForHashes(context.Background(), []string{hash})
	if err != nil {
		return nil, err
	}
	if len(torrents) == 0 {
		return nil, fmt.Errorf("种子未找到: %s", hash)
	}

	var result []clients.PeerInfo
	if torrents[0].Peers != nil {
		for _, p := range torrents[0].Peers {
			if p != nil {
				result = append(result, clients.PeerInfo{
					IP:            p.Address,
					Port:          int(p.Port),
					Client:        p.ClientName,
					Progress:      p.Progress,
					DownloadSpeed: int64(p.RateToClient),
					UploadSpeed:   int64(p.RateToPeer),
					Flags:         p.FlagStr,
				})
			}
		}
	}
	return result, nil
}

// GetTorrentProperties 获取种子详细属性
func (tc *TransmissionClient) GetTorrentProperties(hash string) (*clients.TorrentProperties, error) {
	torrents, err := tc.client.TorrentGetAllForHashes(context.Background(), []string{hash})
	if err != nil {
		return nil, err
	}
	if len(torrents) == 0 {
		return nil, fmt.Errorf("种子未找到: %s", hash)
	}

	t := torrents[0]
	props := &clients.TorrentProperties{
		Hash: hash,
	}

	if t.Name != nil {
		props.Name = *t.Name
	}
	if t.DownloadDir != nil {
		props.SavePath = *t.DownloadDir
	}
	if t.DateCreated != nil {
		props.CreationDate = t.DateCreated.Unix()
	}
	if t.PieceSize != nil {
		props.PieceSize = int64(*t.PieceSize)
	}
	if t.Comment != nil {
		props.Comment = *t.Comment
	}
	if t.CorruptEver != nil {
		props.TotalWasted = *t.CorruptEver
	}
	if t.UploadedEver != nil {
		props.TotalUploaded = *t.UploadedEver
	}
	if t.DownloadedEver != nil {
		props.TotalDownloaded = *t.DownloadedEver
	}
	if t.UploadRatio != nil {
		props.ShareRatio = *t.UploadRatio
	}
	if t.AddedDate != nil {
		props.AdditionDate = t.AddedDate.Unix()
	}
	if t.DoneDate != nil {
		props.CompletionDate = t.DoneDate.Unix()
	}
	if t.Creator != nil {
		props.CreatedBy = *t.Creator
	}
	if t.PeersConnected != nil {
		props.PeersCount = int(*t.PeersConnected)
	}

	return props, nil
}

// GetTorrentPieceStates 获取种子分片状态
func (tc *TransmissionClient) GetTorrentPieceStates(hash string) ([]int, error) {
	// Transmission 的 Pieces 是一个 base64 编码的位图
	return nil, fmt.Errorf("Transmission 不支持获取分片状态")
}

// ==================== 分类管理 ====================

// GetCategories 获取分类列表（使用标签模拟）
func (tc *TransmissionClient) GetCategories() (map[string]clients.Category, error) {
	torrents, err := tc.client.TorrentGetAll(context.Background())
	if err != nil {
		return nil, err
	}

	categories := make(map[string]clients.Category)
	for _, t := range torrents {
		if t.Labels != nil {
			for _, label := range t.Labels {
				categories[label] = clients.Category{
					Name:     label,
					SavePath: "",
				}
			}
		}
	}
	return categories, nil
}

// CreateCategory 创建分类
func (tc *TransmissionClient) CreateCategory(name, savePath string) error {
	return nil // Transmission 标签自动创建
}

// EditCategory 编辑分类
func (tc *TransmissionClient) EditCategory(name, savePath string) error {
	return nil // Transmission 不支持
}

// RemoveCategories 删除分类
func (tc *TransmissionClient) RemoveCategories(names []string) error {
	// 从所有种子中移除指定标签
	torrents, err := tc.client.TorrentGetAll(context.Background())
	if err != nil {
		return err
	}

	nameSet := make(map[string]bool)
	for _, n := range names {
		nameSet[n] = true
	}

	for _, t := range torrents {
		if t.Labels != nil && t.ID != nil {
			var newLabels []string
			needUpdate := false
			for _, label := range t.Labels {
				if nameSet[label] {
					needUpdate = true
				} else {
					newLabels = append(newLabels, label)
				}
			}
			if needUpdate {
				tc.client.TorrentSet(context.Background(), tr.TorrentSetPayload{
					IDs:    []int64{*t.ID},
					Labels: newLabels,
				})
			}
		}
	}
	return nil
}

// ==================== 标签管理 ====================

// GetTags 获取标签列表
func (tc *TransmissionClient) GetTags() ([]string, error) {
	torrents, err := tc.client.TorrentGetAll(context.Background())
	if err != nil {
		return nil, err
	}

	tagSet := make(map[string]bool)
	for _, t := range torrents {
		if t.Labels != nil {
			for _, label := range t.Labels {
				tagSet[label] = true
			}
		}
	}

	var tags []string
	for tag := range tagSet {
		tags = append(tags, tag)
	}
	return tags, nil
}

// CreateTags 创建标签
func (tc *TransmissionClient) CreateTags(tags []string) error {
	return nil // Transmission 标签自动创建
}

// DeleteTags 删除标签
func (tc *TransmissionClient) DeleteTags(tags []string) error {
	return tc.RemoveCategories(tags)
}

// ==================== 全局限速 ====================

// GetGlobalDownloadLimit 获取全局下载限速
func (tc *TransmissionClient) GetGlobalDownloadLimit() (int64, error) {
	session, err := tc.client.SessionArgumentsGet(context.Background(), []string{
		"speed-limit-down", "speed-limit-down-enabled",
	})
	if err != nil {
		return -1, err
	}
	if session.SpeedLimitDownEnabled != nil && *session.SpeedLimitDownEnabled {
		if session.SpeedLimitDown != nil {
			return *session.SpeedLimitDown * 1024, nil
		}
	}
	return -1, nil
}

// GetGlobalUploadLimit 获取全局上传限速
func (tc *TransmissionClient) GetGlobalUploadLimit() (int64, error) {
	session, err := tc.client.SessionArgumentsGet(context.Background(), []string{
		"speed-limit-up", "speed-limit-up-enabled",
	})
	if err != nil {
		return -1, err
	}
	if session.SpeedLimitUpEnabled != nil && *session.SpeedLimitUpEnabled {
		if session.SpeedLimitUp != nil {
			return *session.SpeedLimitUp * 1024, nil
		}
	}
	return -1, nil
}

// SetGlobalDownloadLimit 设置全局下载限速
func (tc *TransmissionClient) SetGlobalDownloadLimit(limit int64) error {
	enabled := limit > 0
	limitKBps := limit / 1024
	return tc.client.SessionArgumentsSet(context.Background(), tr.SessionArguments{
		SpeedLimitDownEnabled: &enabled,
		SpeedLimitDown:        &limitKBps,
	})
}

// SetGlobalUploadLimit 设置全局上传限速
func (tc *TransmissionClient) SetGlobalUploadLimit(limit int64) error {
	enabled := limit > 0
	limitKBps := limit / 1024
	return tc.client.SessionArgumentsSet(context.Background(), tr.SessionArguments{
		SpeedLimitUpEnabled: &enabled,
		SpeedLimitUp:        &limitKBps,
	})
}

// GetAlternativeSpeedLimitsEnabled 获取备用速度限制状态
func (tc *TransmissionClient) GetAlternativeSpeedLimitsEnabled() (bool, error) {
	session, err := tc.client.SessionArgumentsGet(context.Background(), []string{"alt-speed-enabled"})
	if err != nil {
		return false, err
	}
	if session.AltSpeedEnabled != nil {
		return *session.AltSpeedEnabled, nil
	}
	return false, nil
}

// ToggleAlternativeSpeedLimits 切换备用速度限制
func (tc *TransmissionClient) ToggleAlternativeSpeedLimits() error {
	current, err := tc.GetAlternativeSpeedLimitsEnabled()
	if err != nil {
		return err
	}
	newValue := !current
	return tc.client.SessionArgumentsSet(context.Background(), tr.SessionArguments{
		AltSpeedEnabled: &newValue,
	})
}

// ==================== 传输信息 ====================

// GetTransferInfo 获取传输信息
func (tc *TransmissionClient) GetTransferInfo() (*clients.TransferInfo, error) {
	stats, err := tc.client.SessionStats(context.Background())
	if err != nil {
		return nil, err
	}

	return &clients.TransferInfo{
		DownloadSpeed:       stats.DownloadSpeed,
		UploadSpeed:         stats.UploadSpeed,
		TotalPeersConnected: int(stats.PausedTorrentCount + stats.ActiveTorrentCount),
	}, nil
}

// GetFreeSpace 获取磁盘剩余空间
func (tc *TransmissionClient) GetFreeSpace(path string) (int64, error) {
	space, err := tc.client.FreeSpace(context.Background(), path)
	if err != nil {
		return -1, err
	}
	return int64(space), nil
}

// ==================== 服务器操作 ====================

// GetServerInfo 获取服务器信息
func (tc *TransmissionClient) GetServerInfo() (map[string]interface{}, error) {
	info := make(map[string]interface{})

	session, err := tc.client.SessionArgumentsGet(context.Background(), nil)
	if err == nil {
		if session.Version != nil {
			info["version"] = *session.Version
		}
		if session.RPCVersion != nil {
			info["rpc_version"] = *session.RPCVersion
		}
	}

	stats, err := tc.client.SessionStats(context.Background())
	if err == nil {
		info["stats"] = stats
	}

	return info, nil
}

// GetDefaultSavePath 获取默认保存路径
func (tc *TransmissionClient) GetDefaultSavePath() (string, error) {
	session, err := tc.client.SessionArgumentsGet(context.Background(), []string{"download-dir"})
	if err != nil {
		return "", err
	}
	if session.DownloadDir != nil {
		return *session.DownloadDir, nil
	}
	return "", nil
}

// Shutdown 关闭客户端
func (tc *TransmissionClient) Shutdown() error {
	return tc.client.SessionClose(context.Background())
}

// ==================== 日志 ====================

// GetLogs 获取日志（Transmission 不支持）
func (tc *TransmissionClient) GetLogs(normal, warning, critical bool, lastKnownID int64) ([]clients.LogEntry, error) {
	return nil, fmt.Errorf("Transmission 不支持日志功能")
}

// ==================== 种子分片信息 ====================

// GetTorrentPieceStates 获取种子分片状态
func (tc *TransmissionClient) GetTorrentPieceStates(hash string) ([]int, error) {
	torrents, err := tc.client.TorrentGetAllForHashes(context.Background(), []string{hash})
	if err != nil {
		return nil, err
	}
	if len(torrents) == 0 {
		return nil, fmt.Errorf("种子未找到: %s", hash)
	}

	torrent := torrents[0]
	if torrent.Pieces == nil {
		return nil, nil
	}

	// 将片段状态转换为整数数组
	var states []int
	for _, piece := range *torrent.Pieces {
		// Transmission 使用 0-3 表示片段状态
		// 0=未验证, 1=验证通过, 2=正在验证, 3=验证失败
		states = append(states, int(piece))
	}
	return states, nil
}

// ==================== 备用速度限制 ====================

// GetAlternativeSpeedLimitsEnabled 获取备用速度限制是否启用
func (tc *TransmissionClient) GetAlternativeSpeedLimitsEnabled() (bool, error) {
	session, err := tc.client.SessionArgumentsGet(context.Background(), []string{"alt-speed-enabled"})
	if err != nil {
		return false, err
	}
	if session.AltSpeedEnabled == nil {
		return false, nil
	}
	return *session.AltSpeedEnabled, nil
}

// ToggleAlternativeSpeedLimits 切换备用速度限制
func (tc *TransmissionClient) ToggleAlternativeSpeedLimits() error {
	enabled, err := tc.GetAlternativeSpeedLimitsEnabled()
	if err != nil {
		return err
	}

	newState := !enabled
	return tc.client.SessionSet(context.Background(), tr.SessionArguments{
		AltSpeedEnabled: &newState,
	})
}

// ==================== 辅助方法 ====================

// getTorrentIDByHash 通过哈希获取种子ID
func (tc *TransmissionClient) getTorrentIDByHash(hash string) (int64, error) {
	torrents, err := tc.client.TorrentGetAllForHashes(context.Background(), []string{hash})
	if err != nil {
		return 0, err
	}
	if len(torrents) == 0 || torrents[0].ID == nil {
		return 0, fmt.Errorf("种子未找到: %s", hash)
	}
	return *torrents[0].ID, nil
}

// hashesToIDs 批量转换哈希为ID
func (tc *TransmissionClient) hashesToIDs(hashes []string) ([]int64, error) {
	torrents, err := tc.client.TorrentGetAllForHashes(context.Background(), hashes)
	if err != nil {
		return nil, err
	}

	var ids []int64
	for _, t := range torrents {
		if t.ID != nil {
			ids = append(ids, *t.ID)
		}
	}
	return ids, nil
}
