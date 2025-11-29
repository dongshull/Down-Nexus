package clients

import "down-nexus-api/internal/models"

type DownloaderClient interface {
	GetTorrents() ([]models.UnifiedTorrent, error)
	AddTorrent(magnetURL string) error
	PauseTorrent(hash string) error
	ResumeTorrent(hash string) error
	DeleteTorrent(hash string, deleteFiles bool) error
	GetClientID() string
}
