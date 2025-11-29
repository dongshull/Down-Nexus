package qbittorrent

import (
	"fmt"
	"down-nexus-api/internal/models"
	qb "github.com/autobrr/go-qbittorrent"
)

type QbitClient struct {
	client   *qb.Client
	clientID string
}

func NewQbitClient(host, username, password, clientID string) (*QbitClient, error) {
	// 不添加 API 路径，让 go-qbittorrent 库自己处理
	cfg := qb.Config{
		Host:     host,
		Username: username,
		Password: password,
	}
	qbClient := qb.NewClient(cfg)
	
	// Login to qBittorrent
	err := qbClient.Login()
	if err != nil {
		return nil, fmt.Errorf("qBittorrent 登录失败: %w", err)
	}
	
	return &QbitClient{
		client:   qbClient,
		clientID: clientID,
	}, nil
}

func (qc *QbitClient) GetTorrents() ([]models.UnifiedTorrent, error) {
	torrents, err := qc.client.GetTorrents(qb.TorrentFilterOptions{})
	if err != nil {
		return nil, err
	}
	
	var unifiedTorrents []models.UnifiedTorrent
	for _, torrent := range torrents {
		unifiedTorrent := models.UnifiedTorrent{
			ClientID:      qc.clientID,
			Name:          torrent.Name,
			Hash:          torrent.Hash,
			Size:          torrent.Size,
			State:         string(torrent.State),
			Progress:      torrent.Progress,
			DownloadSpeed: torrent.DlSpeed,
			UploadSpeed:   torrent.UpSpeed,
			Downloaded:    torrent.Downloaded,
			Uploaded:      torrent.Uploaded,
			ETA:           torrent.ETA,
		}
		unifiedTorrents = append(unifiedTorrents, unifiedTorrent)
	}
	
	return unifiedTorrents, nil
}

func (qc *QbitClient) AddTorrent(magnetURL string) error {
	// Use qBittorrent's AddTorrentFromUrl method
	options := map[string]string{}
	
	return qc.client.AddTorrentFromUrl(magnetURL, options)
}

func (qc *QbitClient) PauseTorrent(hash string) error {
	return qc.client.Pause([]string{hash})
}

func (qc *QbitClient) ResumeTorrent(hash string) error {
	return qc.client.Resume([]string{hash})
}

func (qc *QbitClient) DeleteTorrent(hash string, deleteFiles bool) error {
	return qc.client.DeleteTorrents([]string{hash}, deleteFiles)
}

func (qc *QbitClient) GetClientID() string {
	return qc.clientID
}