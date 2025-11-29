package models

// UnifiedTorrent 统一种子模型
type UnifiedTorrent struct {
	ClientID      string  `json:"client_id"`
	Name          string  `json:"name"`
	Hash          string  `json:"hash"`
	Size          int64   `json:"size"`
	State         string  `json:"state"`
	Progress      float64 `json:"progress"`
	DownloadSpeed int64   `json:"download_speed"`
	UploadSpeed   int64   `json:"upload_speed"`
	Downloaded    int64   `json:"downloaded"`
	Uploaded      int64   `json:"uploaded"`
	ETA           int64   `json:"eta"`
}