package core

import (
	"sync"

	"down-nexus-api/internal/models"
	"down-nexus-api/pkg/clients"

	"gorm.io/gorm"
)

type TorrentService struct {
	clients []clients.DownloaderClient
	Db      *gorm.DB
}

func NewTorrentService(clients []clients.DownloaderClient, db *gorm.DB) *TorrentService {
	return &TorrentService{
		clients: clients,
		Db:      db,
	}
}

func (ts *TorrentService) GetAllTorrents() []models.UnifiedTorrent {
	var allTorrents []models.UnifiedTorrent
	var mutex sync.Mutex
	var wg sync.WaitGroup

	// 并发调用每个下载器的 GetTorrents 方法
	for _, client := range ts.clients {
		wg.Add(1)
		go func(c clients.DownloaderClient) {
			defer wg.Done()

			torrents, err := c.GetTorrents()
			if err != nil {
				// 在实际应用中，可能需要记录错误日志
				// 这里暂时忽略错误，继续处理其他下载器
				return
			}

			// 使用互斥锁安全地合并结果
			mutex.Lock()
			allTorrents = append(allTorrents, torrents...)
			mutex.Unlock()
		}(client)
	}

	// 等待所有 goroutine 完成
	wg.Wait()

	return allTorrents
}

func (ts *TorrentService) AddTorrent(magnetURL string, clientID string) error {
	// 遍历所有客户端，找到匹配的 clientID
	for _, client := range ts.clients {
		// 由于我们需要检查 clientID，我们需要先获取种子来确定客户端ID
		// 但更好的方法是在客户端接口中添加 GetClientID() 方法
		// 目前我们通过获取一个种子来检查 ClientID
		torrents, err := client.GetTorrents()
		if err != nil {
			continue
		}

		// 如果这个客户端有种子且 ClientID 匹配
		if len(torrents) > 0 && torrents[0].ClientID == clientID {
			return client.AddTorrent(magnetURL)
		}

		// 如果客户端没有种子，我们需要创建一个临时种子来检查 ClientID
		// 这种情况比较复杂，可能需要修改接口设计
	}

	// 如果没有找到匹配的客户端，返回错误
	return &ClientNotFoundError{ClientID: clientID}
}

func (ts *TorrentService) PauseTorrent(clientID string, hash string) error {
	// 遍历所有客户端，找到匹配的 clientID
	for _, client := range ts.clients {
		if client.GetClientID() == clientID {
			return client.PauseTorrent(hash)
		}
	}

	// 如果没有找到匹配的客户端，返回错误
	return &ClientNotFoundError{ClientID: clientID}
}

func (ts *TorrentService) ResumeTorrent(clientID string, hash string) error {
	// 遍历所有客户端，找到匹配的 clientID
	for _, client := range ts.clients {
		if client.GetClientID() == clientID {
			return client.ResumeTorrent(hash)
		}
	}

	// 如果没有找到匹配的客户端，返回错误
	return &ClientNotFoundError{ClientID: clientID}
}

func (ts *TorrentService) DeleteTorrent(clientID string, hash string, deleteFiles bool) error {
	// 遍历所有客户端，找到匹配的 clientID
	for _, client := range ts.clients {
		if client.GetClientID() == clientID {
			return client.DeleteTorrent(hash, deleteFiles)
		}
	}

	// 如果没有找到匹配的客户端，返回错误
	return &ClientNotFoundError{ClientID: clientID}
}

// 自定义错误类型
type ClientNotFoundError struct {
	ClientID string
}

func (e *ClientNotFoundError) Error() string {
	return "client not found: " + e.ClientID
}

// GetClientConfigs 获取所有客户端配置
func (ts *TorrentService) GetClientConfigs() ([]models.ClientConfig, error) {
	var configs []models.ClientConfig
	err := ts.Db.Find(&configs).Error
	return configs, err
}

// GetClient 根据客户端ID获取客户端实例
func (ts *TorrentService) GetClient(clientID string) clients.DownloaderClient {
	for _, client := range ts.clients {
		if client.GetClientID() == clientID {
			return client
		}
	}
	return nil
}
