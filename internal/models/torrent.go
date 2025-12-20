package models

import "time"

// UnifiedTorrent 统一种子模型
type UnifiedTorrent struct {
	ClientID      string    `json:"client_id"`
	Name          string    `json:"name"`
	Hash          string    `json:"hash"`
	Size          int64     `json:"size"`
	State         string    `json:"state"`
	Progress      float64   `json:"progress"`
	DownloadSpeed int64     `json:"download_speed"`
	UploadSpeed   int64     `json:"upload_speed"`
	Downloaded    int64     `json:"downloaded"`
	Uploaded      int64     `json:"uploaded"`
	ETA           int64     `json:"eta"`
	AddedOn       time.Time `json:"added_on"`
	CompletedOn   time.Time `json:"completed_on"`
	SavePath      string    `json:"save_path"`
	Category      string    `json:"category"`
	Tags          []string  `json:"tags"`
	Tracker       string    `json:"tracker"`
	DownloadLimit int64     `json:"download_limit"` // 下载限速 (字节/秒)，-1表示无限制
	UploadLimit   int64     `json:"upload_limit"`   // 上传限速 (字节/秒)，-1表示无限制
	Ratio         float64   `json:"ratio"`          // 分享率
	Priority      int       `json:"priority"`       // 优先级
}