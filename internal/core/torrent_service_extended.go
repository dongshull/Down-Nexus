package core

import (
	"fmt"
	"sync"

	"down-nexus-api/internal/models"
	"down-nexus-api/pkg/clients"
)

// GetFilteredTorrents 获取过滤后的种子
func (ts *TorrentService) GetFilteredTorrents(clientID, category, state, tag string) ([]models.UnifiedTorrent, error) {
	var torrents []models.UnifiedTorrent
	var mutex sync.Mutex
	var wg sync.WaitGroup

	// 并发调用每个下载器的 GetTorrents 方法
	for _, client := range ts.clients {
		wg.Add(1)
		go func(c clients.DownloaderClient) {
			defer wg.Done()

			// 如果指定了客户端ID，但不是当前客户端，跳过
			if clientID != "" && c.GetClientID() != clientID {
				return
			}

			t, err := c.GetTorrents()
			if err != nil {
				// 在实际应用中，可能需要记录错误日志
				// 这里暂时忽略错误，继续处理其他下载器
				return
			}

			// 过滤种子
			var filteredTorrents []models.UnifiedTorrent
			for _, torrent := range t {
				if category != "" && torrent.Category != category {
					continue
				}
				if state != "" && torrent.State != state {
					continue
				}
				if tag != "" && !containsTag(torrent.Tags, tag) {
					continue
				}
				filteredTorrents = append(filteredTorrents, torrent)
			}

			// 使用互斥锁安全地合并结果
			mutex.Lock()
			torrents = append(torrents, filteredTorrents...)
			mutex.Unlock()
		}(client)
	}

	// 等待所有 goroutine 完成
	wg.Wait()

	return torrents, nil
}

// GetTorrentDetails 获取种子详细信息
func (ts *TorrentService) GetTorrentDetails(clientID, hash string) (*models.UnifiedTorrent, error) {
	// 查找指定的客户端
	for _, client := range ts.clients {
		if client.GetClientID() == clientID {
			torrents, err := client.GetTorrents()
			if err != nil {
				return nil, err
			}

			// 查找指定哈希的种子
			for _, torrent := range torrents {
				if torrent.Hash == hash {
					return &torrent, nil
				}
			}

			return nil, fmt.Errorf("torrent not found: %s", hash)
		}
	}

	return nil, &ClientNotFoundError{ClientID: clientID}
}

// AddTorrentWithOptions 添加带选项的种子
func (ts *TorrentService) AddTorrentWithOptions(clientID string, options clients.TorrentOptions) error {
	// 查找指定的客户端
	for _, client := range ts.clients {
		if client.GetClientID() == clientID {
			return client.AddTorrentWithOptions(options)
		}
	}

	return &ClientNotFoundError{ClientID: clientID}
}

// PauseTorrents 批量暂停种子
func (ts *TorrentService) PauseTorrents(clientID string, hashes []string) error {
	// 查找指定的客户端
	for _, client := range ts.clients {
		if client.GetClientID() == clientID {
			// 批量操作
			for _, hash := range hashes {
				err := client.PauseTorrent(hash)
				if err != nil {
					return err
				}
			}
			return nil
		}
	}

	return &ClientNotFoundError{ClientID: clientID}
}

// ResumeTorrents 批量恢复种子
func (ts *TorrentService) ResumeTorrents(clientID string, hashes []string) error {
	// 查找指定的客户端
	for _, client := range ts.clients {
		if client.GetClientID() == clientID {
			// 批量操作
			for _, hash := range hashes {
				err := client.ResumeTorrent(hash)
				if err != nil {
					return err
				}
			}
			return nil
		}
	}

	return &ClientNotFoundError{ClientID: clientID}
}

// DeleteTorrents 批量删除种子
func (ts *TorrentService) DeleteTorrents(clientID string, hashes []string, deleteFiles bool) error {
	// 查找指定的客户端
	for _, client := range ts.clients {
		if client.GetClientID() == clientID {
			// 批量操作
			for _, hash := range hashes {
				err := client.DeleteTorrent(hash, deleteFiles)
				if err != nil {
					return err
				}
			}
			return nil
		}
	}

	return &ClientNotFoundError{ClientID: clientID}
}

// SetTorrentCategory 设置种子分类
func (ts *TorrentService) SetTorrentCategory(clientID, hash, category string) error {
	// 查找指定的客户端
	for _, client := range ts.clients {
		if client.GetClientID() == clientID {
			return client.SetTorrentCategory(hash, category)
		}
	}

	return &ClientNotFoundError{ClientID: clientID}
}

// SetTorrentTags 设置种子标签
func (ts *TorrentService) SetTorrentTags(clientID, hash string, tags []string) error {
	// 查找指定的客户端
	for _, client := range ts.clients {
		if client.GetClientID() == clientID {
			return client.SetTorrentTags(hash, tags)
		}
	}

	return &ClientNotFoundError{ClientID: clientID}
}

// SetTorrentDownloadLimit 设置种子下载限速
func (ts *TorrentService) SetTorrentDownloadLimit(clientID, hash string, limit int64) error {
	// 查找指定的客户端
	for _, client := range ts.clients {
		if client.GetClientID() == clientID {
			return client.SetTorrentDownloadLimit(hash, limit)
		}
	}

	return &ClientNotFoundError{ClientID: clientID}
}

// SetTorrentUploadLimit 设置种子上传限速
func (ts *TorrentService) SetTorrentUploadLimit(clientID, hash string, limit int64) error {
	// 查找指定的客户端
	for _, client := range ts.clients {
		if client.GetClientID() == clientID {
			return client.SetTorrentUploadLimit(hash, limit)
		}
	}

	return &ClientNotFoundError{ClientID: clientID}
}

// SetTorrentPriority 设置种子优先级
func (ts *TorrentService) SetTorrentPriority(clientID, hash string, priority int) error {
	// 查找指定的客户端
	for _, client := range ts.clients {
		if client.GetClientID() == clientID {
			return client.SetTorrentPriority(hash, priority)
		}
	}

	return &ClientNotFoundError{ClientID: clientID}
}

// SetGlobalDownloadLimit 设置全局下载限速
func (ts *TorrentService) SetGlobalDownloadLimit(clientID string, limit int64) error {
	// 查找指定的客户端
	for _, client := range ts.clients {
		if client.GetClientID() == clientID {
			return client.SetGlobalDownloadLimit(limit)
		}
	}

	return &ClientNotFoundError{ClientID: clientID}
}

// SetGlobalUploadLimit 设置全局上传限速
func (ts *TorrentService) SetGlobalUploadLimit(clientID string, limit int64) error {
	// 查找指定的客户端
	for _, client := range ts.clients {
		if client.GetClientID() == clientID {
			return client.SetGlobalUploadLimit(limit)
		}
	}

	return &ClientNotFoundError{ClientID: clientID}
}

// GetCategories 获取分类列表
func (ts *TorrentService) GetCategories(clientID string) (map[string]clients.Category, error) {
	// 查找指定的客户端
	for _, client := range ts.clients {
		if client.GetClientID() == clientID {
			return client.GetCategories()
		}
	}

	return nil, &ClientNotFoundError{ClientID: clientID}
}

// CreateCategory 创建分类
func (ts *TorrentService) CreateCategory(clientID, name, savePath string) error {
	// 查找指定的客户端
	for _, client := range ts.clients {
		if client.GetClientID() == clientID {
			return client.CreateCategory(name, savePath)
		}
	}

	return &ClientNotFoundError{ClientID: clientID}
}

// DeleteCategories 删除分类
func (ts *TorrentService) DeleteCategories(clientID string, names []string) error {
	// 查找指定的客户端
	for _, client := range ts.clients {
		if client.GetClientID() == clientID {
			return client.RemoveCategories(names)
		}
	}

	return &ClientNotFoundError{ClientID: clientID}
}

// GetTags 获取标签列表
func (ts *TorrentService) GetTags(clientID string) ([]string, error) {
	// 查找指定的客户端
	for _, client := range ts.clients {
		if client.GetClientID() == clientID {
			return client.GetTags()
		}
	}

	return nil, &ClientNotFoundError{ClientID: clientID}
}

// CreateTags 创建标签
func (ts *TorrentService) CreateTags(clientID string, tags []string) error {
	// 查找指定的客户端
	for _, client := range ts.clients {
		if client.GetClientID() == clientID {
			return client.CreateTags(tags)
		}
	}

	return &ClientNotFoundError{ClientID: clientID}
}

// DeleteTags 删除标签
func (ts *TorrentService) DeleteTags(clientID string, tags []string) error {
	// 查找指定的客户端
	for _, client := range ts.clients {
		if client.GetClientID() == clientID {
			return client.DeleteTags(tags)
		}
	}

	return &ClientNotFoundError{ClientID: clientID}
}

// TestClientConnection 测试客户端连接
func (ts *TorrentService) TestClientConnection(clientID string) error {
	// 查找指定的客户端
	for _, client := range ts.clients {
		if client.GetClientID() == clientID {
			// 尝试获取种子列表来测试连接
			_, err := client.GetTorrents()
			return err
		}
	}

	return &ClientNotFoundError{ClientID: clientID}
}

// containsTag 检查标签数组是否包含指定标签
func containsTag(tags []string, tag string) bool {
	for _, t := range tags {
		if t == tag {
			return true
		}
	}
	return false
}

// ==================== 日志相关 ====================

// GetLogs 获取客户端日志
func (ts *TorrentService) GetLogs(clientID string, normal, warning, critical bool, lastKnownID int64) ([]clients.LogEntry, error) {
	for _, client := range ts.clients {
		if client.GetClientID() == clientID {
			return client.GetLogs(normal, warning, critical, lastKnownID)
		}
	}
	return nil, &ClientNotFoundError{ClientID: clientID}
}

// ==================== 种子分片状态 ====================

// GetTorrentPieceStates 获取种子分片状态
func (ts *TorrentService) GetTorrentPieceStates(clientID, hash string) ([]int, error) {
	for _, client := range ts.clients {
		if client.GetClientID() == clientID {
			return client.GetTorrentPieceStates(hash)
		}
	}
	return nil, &ClientNotFoundError{ClientID: clientID}
}

// ==================== 备用速度限制 ====================

// GetAlternativeSpeedLimitsEnabled 获取备用速度限制状态
func (ts *TorrentService) GetAlternativeSpeedLimitsEnabled(clientID string) (bool, error) {
	for _, client := range ts.clients {
		if client.GetClientID() == clientID {
			return client.GetAlternativeSpeedLimitsEnabled()
		}
	}
	return false, &ClientNotFoundError{ClientID: clientID}
}

// ToggleAlternativeSpeedLimits 切换备用速度限制
func (ts *TorrentService) ToggleAlternativeSpeedLimits(clientID string) error {
	for _, client := range ts.clients {
		if client.GetClientID() == clientID {
			return client.ToggleAlternativeSpeedLimits()
		}
	}
	return &ClientNotFoundError{ClientID: clientID}
}
