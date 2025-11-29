package transmission

import (
	"context"
	"fmt"
	"down-nexus-api/internal/models"
	tr "github.com/hekmon/transmissionrpc/v2"
)

type TransmissionClient struct {
	client   *tr.Client
	clientID string
}

func NewTransmissionClient(host, username, password, clientID string) (*TransmissionClient, error) {
	// Create Transmission client with explicit port configuration
	client, err := tr.New(host, username, password, &tr.AdvancedConfig{
		HTTPS: false,
		Port:  10002, // Explicitly set the port
	})
	if err != nil {
		return nil, fmt.Errorf("Transmission 连接失败: %w", err)
	}
	
	return &TransmissionClient{
		client:   client,
		clientID: clientID,
	}, nil
}

func (tc *TransmissionClient) GetTorrents() ([]models.UnifiedTorrent, error) {
	torrents, err := tc.client.TorrentGetAll(context.Background())
	if err != nil {
		return nil, err
	}
	
	var unifiedTorrents []models.UnifiedTorrent
	for _, torrent := range torrents {
		// Handle nil pointers safely
		var name string
		if torrent.Name != nil {
			name = *torrent.Name
		}
		
		var hash string
		if torrent.HashString != nil {
			hash = *torrent.HashString
		}
		
		var size int64
		if torrent.TotalSize != nil {
			size = int64(*torrent.TotalSize)
		}
		
		var progress float64
		if torrent.PercentDone != nil {
			progress = *torrent.PercentDone
		}
		
		var downloadSpeed int64
		if torrent.RateDownload != nil {
			downloadSpeed = *torrent.RateDownload
		}
		
		var uploadSpeed int64
		if torrent.RateUpload != nil {
			uploadSpeed = *torrent.RateUpload
		}
		
		var downloaded int64
		if torrent.DownloadedEver != nil {
			downloaded = *torrent.DownloadedEver
		}
		
		var uploaded int64
		if torrent.UploadedEver != nil {
			uploaded = *torrent.UploadedEver
		}
		
		var eta int64
		if torrent.Eta != nil {
			eta = *torrent.Eta
		}
		
		// Convert status to string
		var state string
		if torrent.Status != nil {
			state = torrent.Status.String()
		}
		
		unifiedTorrent := models.UnifiedTorrent{
			ClientID:      tc.clientID,
			Name:          name,
			Hash:          hash,
			Size:          size,
			State:         state,
			Progress:      progress,
			DownloadSpeed: downloadSpeed,
			UploadSpeed:   uploadSpeed,
			Downloaded:    downloaded,
			Uploaded:      uploaded,
			ETA:           eta,
		}
		unifiedTorrents = append(unifiedTorrents, unifiedTorrent)
	}
	
	return unifiedTorrents, nil
}

func (tc *TransmissionClient) AddTorrent(magnetURL string) error {
	// Use Transmission's TorrentAdd method
	_, err := tc.client.TorrentAdd(context.Background(), tr.TorrentAddPayload{
		Filename: &magnetURL,
	})
	
	return err
}

func (tc *TransmissionClient) PauseTorrent(hash string) error {
	// Convert hash string to int64 ID (Transmission uses numeric IDs)
	torrents, err := tc.client.TorrentGetAll(context.Background())
	if err != nil {
		return err
	}
	
	for _, torrent := range torrents {
		if torrent.HashString != nil && *torrent.HashString == hash {
			if torrent.ID != nil {
				err := tc.client.TorrentStopIDs(context.Background(), []int64{*torrent.ID})
				return err
			}
		}
	}
	
	return fmt.Errorf("torrent with hash %s not found", hash)
}

func (tc *TransmissionClient) ResumeTorrent(hash string) error {
	// Convert hash string to int64 ID
	torrents, err := tc.client.TorrentGetAll(context.Background())
	if err != nil {
		return err
	}
	
	for _, torrent := range torrents {
		if torrent.HashString != nil && *torrent.HashString == hash {
			if torrent.ID != nil {
				err := tc.client.TorrentStartIDs(context.Background(), []int64{*torrent.ID})
				return err
			}
		}
	}
	
	return fmt.Errorf("torrent with hash %s not found", hash)
}

func (tc *TransmissionClient) DeleteTorrent(hash string, deleteFiles bool) error {
	// Convert hash string to int64 ID
	torrents, err := tc.client.TorrentGetAll(context.Background())
	if err != nil {
		return err
	}
	
	for _, torrent := range torrents {
		if torrent.HashString != nil && *torrent.HashString == hash {
			if torrent.ID != nil {
				err := tc.client.TorrentRemove(context.Background(), tr.TorrentRemovePayload{
					IDs:             []int64{*torrent.ID},
					DeleteLocalData: deleteFiles,
				})
				return err
			}
		}
	}
	
	return fmt.Errorf("torrent with hash %s not found", hash)
}

func (tc *TransmissionClient) GetClientID() string {
	return tc.clientID
}