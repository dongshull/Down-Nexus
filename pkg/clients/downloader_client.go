package clients

import "down-nexus-api/internal/models"

type DownloaderClient interface {
	// ==================== 基础功能 ====================
	GetTorrents() ([]models.UnifiedTorrent, error)
	GetTorrentDetails(hash string) (*models.UnifiedTorrent, error)
	AddTorrent(magnetURL string) error
	AddTorrentWithOptions(options TorrentOptions) error
	AddTorrentFromFile(filePath string, options map[string]string) error
	PauseTorrent(hash string) error
	PauseTorrents(hashes []string) error
	ResumeTorrent(hash string) error
	ResumeTorrents(hashes []string) error
	DeleteTorrent(hash string, deleteFiles bool) error
	DeleteTorrents(hashes []string, deleteFiles bool) error
	GetClientID() string

	// ==================== 种子操作 ====================
	RecheckTorrent(hash string) error
	RecheckTorrents(hashes []string) error
	ReannounceTorrent(hash string) error
	ReannounceTorrents(hashes []string) error
	SetTorrentLocation(hash string, location string) error
	SetTorrentName(hash string, name string) error
	SetForceStart(hash string, enabled bool) error
	SetAutoManagement(hash string, enabled bool) error
	SetSequentialDownload(hash string, enabled bool) error
	SetFirstLastPiecePriority(hash string, enabled bool) error
	SetSuperSeeding(hash string, enabled bool) error
	SetShareLimit(hash string, ratioLimit float64, seedingTimeLimit int64) error

	// ==================== 种子属性 ====================
	SetTorrentCategory(hash string, category string) error
	SetTorrentTags(hash string, tags []string) error
	AddTorrentTags(hash string, tags []string) error
	RemoveTorrentTags(hash string, tags []string) error
	SetTorrentTrackers(hash string, trackers []string) error
	AddTorrentTrackers(hash string, trackers []string) error
	RemoveTorrentTrackers(hash string, trackerURLs []string) error
	SetTorrentDownloadLimit(hash string, limit int64) error
	SetTorrentUploadLimit(hash string, limit int64) error
	GetTorrentDownloadLimit(hash string) (int64, error)
	GetTorrentUploadLimit(hash string) (int64, error)
	SetTorrentPriority(hash string, priority int) error

	// ==================== 种子文件操作 ====================
	GetTorrentFiles(hash string) ([]TorrentFile, error)
	SetFilePriority(hash string, fileIDs []int, priority int) error
	RenameFile(hash string, oldPath string, newPath string) error
	RenameFolder(hash string, oldPath string, newPath string) error

	// ==================== 种子详细信息 ====================
	GetTorrentTrackers(hash string) ([]TrackerInfo, error)
	GetTorrentPeers(hash string) ([]PeerInfo, error)
	GetTorrentProperties(hash string) (*TorrentProperties, error)
	GetTorrentPieceStates(hash string) ([]int, error)

	// ==================== 分类管理 ====================
	GetCategories() (map[string]Category, error)
	CreateCategory(name string, savePath string) error
	EditCategory(name string, savePath string) error
	RemoveCategories(names []string) error

	// ==================== 标签管理 ====================
	GetTags() ([]string, error)
	CreateTags(tags []string) error
	DeleteTags(tags []string) error

	// ==================== 全局限速 ====================
	GetGlobalDownloadLimit() (int64, error)
	GetGlobalUploadLimit() (int64, error)
	SetGlobalDownloadLimit(limit int64) error
	SetGlobalUploadLimit(limit int64) error
	GetAlternativeSpeedLimitsEnabled() (bool, error)
	ToggleAlternativeSpeedLimits() error

	// ==================== 传输信息 ====================
	GetTransferInfo() (*TransferInfo, error)
	GetFreeSpace(path string) (int64, error)

	// ==================== 服务器操作 ====================
	GetServerInfo() (map[string]interface{}, error)
	GetDefaultSavePath() (string, error)
	Shutdown() error

	// ==================== 日志 ====================
	GetLogs(normal, warning, critical bool, lastKnownID int64) ([]LogEntry, error)

	// ==================== RSS (可选) ====================
	// GetRSSFeeds() (map[string]interface{}, error)
	// AddRSSFeed(url string, path string) error
	// RemoveRSSFeed(path string) error
}

// ==================== 数据结构 ====================

// TorrentOptions 种子添加选项
type TorrentOptions struct {
	URL                  string   `json:"url"`                     // 种子URL或磁力链接
	TorrentFile          []byte   `json:"torrent_file,omitempty"`  // 种子文件内容
	SavePath             string   `json:"save_path"`               // 保存路径
	Category             string   `json:"category"`                // 分类
	Tags                 []string `json:"tags"`                    // 标签
	Priority             int      `json:"priority"`                // 优先级
	Sequential           bool     `json:"sequential"`              // 顺序下载
	FirstLast            bool     `json:"first_last_piece"`        // 优先下载首尾分片
	SkipChecking         bool     `json:"skip_checking"`           // 跳过校验
	Paused               bool     `json:"paused"`                  // 添加后暂停
	DownloadLimit        int64    `json:"download_limit"`          // 下载限速 (B/s)
	UploadLimit          int64    `json:"upload_limit"`            // 上传限速 (B/s)
	RatioLimit           float64  `json:"ratio_limit"`             // 分享率限制
	SeedingTimeLimit     int      `json:"seeding_time_limit"`      // 做种时间限制(秒)
	AutoManagement       bool     `json:"auto_management"`         // 自动管理
	ContentLayout        string   `json:"content_layout"`          // 内容布局: Original, Subfolder, NoSubfolder
	RootFolder           bool     `json:"root_folder"`             // 创建根文件夹
	Rename               string   `json:"rename"`                  // 重命名种子
	UploadLimitPerPeer   int64    `json:"upload_limit_per_peer"`   // 每个连接的上传限速
	DownloadLimitPerPeer int64    `json:"download_limit_per_peer"` // 每个连接的下载限速
}

// Category 分类信息
type Category struct {
	Name     string `json:"name"`
	SavePath string `json:"save_path"`
}

// TorrentFile 种子文件信息
type TorrentFile struct {
	Index      int     `json:"index"`      // 文件索引
	Name       string  `json:"name"`       // 文件名
	Size       int64   `json:"size"`       // 文件大小
	Progress   float64 `json:"progress"`   // 下载进度
	Priority   int     `json:"priority"`   // 优先级: 0=不下载, 1=普通, 6=高, 7=最高
	Downloaded int64   `json:"downloaded"` // 已下载大小
	Wanted     bool    `json:"wanted"`     // 是否需要下载
}

// TrackerInfo Tracker信息
type TrackerInfo struct {
	URL           string `json:"url"`            // Tracker URL
	Status        int    `json:"status"`         // 状态码
	StatusMessage string `json:"status_message"` // 状态消息
	Tier          int    `json:"tier"`           // 层级
	Peers         int    `json:"peers"`          // Peers数量
	Seeds         int    `json:"seeds"`          // Seeds数量
	Leechers      int    `json:"leechers"`       // Leechers数量
	Downloaded    int    `json:"downloaded"`     // 已完成数
	Message       string `json:"message"`        // 消息
}

// PeerInfo Peer信息
type PeerInfo struct {
	IP            string  `json:"ip"`             // IP地址
	Port          int     `json:"port"`           // 端口
	Client        string  `json:"client"`         // 客户端名称
	Progress      float64 `json:"progress"`       // 下载进度
	DownloadSpeed int64   `json:"download_speed"` // 下载速度
	UploadSpeed   int64   `json:"upload_speed"`   // 上传速度
	Downloaded    int64   `json:"downloaded"`     // 已下载
	Uploaded      int64   `json:"uploaded"`       // 已上传
	Connection    string  `json:"connection"`     // 连接类型
	Flags         string  `json:"flags"`          // 标志
	Country       string  `json:"country"`        // 国家代码
	CountryName   string  `json:"country_name"`   // 国家名称
}

// TorrentProperties 种子详细属性
type TorrentProperties struct {
	Hash              string  `json:"hash"`
	Name              string  `json:"name"`
	SavePath          string  `json:"save_path"`
	CreationDate      int64   `json:"creation_date"`
	PieceSize         int64   `json:"piece_size"`
	Comment           string  `json:"comment"`
	TotalWasted       int64   `json:"total_wasted"`
	TotalUploaded     int64   `json:"total_uploaded"`
	TotalDownloaded   int64   `json:"total_downloaded"`
	UploadedSession   int64   `json:"uploaded_session"`
	DownloadedSession int64   `json:"downloaded_session"`
	UploadLimit       int64   `json:"upload_limit"`
	DownloadLimit     int64   `json:"download_limit"`
	TimeElapsed       int64   `json:"time_elapsed"`
	SeedingTime       int64   `json:"seeding_time"`
	ConnectionsCount  int     `json:"connections_count"`
	ConnectionsLimit  int     `json:"connections_limit"`
	ShareRatio        float64 `json:"share_ratio"`
	AdditionDate      int64   `json:"addition_date"`
	CompletionDate    int64   `json:"completion_date"`
	CreatedBy         string  `json:"created_by"`
	LastActivity      int64   `json:"last_activity"`
	PeersCount        int     `json:"peers_count"`
	SeedsCount        int     `json:"seeds_count"`
	TotalPeers        int     `json:"total_peers"`
	TotalSeeds        int     `json:"total_seeds"`
}

// TransferInfo 传输信息
type TransferInfo struct {
	DownloadSpeed          int64  `json:"download_speed"`           // 当前下载速度
	UploadSpeed            int64  `json:"upload_speed"`             // 当前上传速度
	DownloadedBytes        int64  `json:"downloaded_bytes"`         // 总下载量
	UploadedBytes          int64  `json:"uploaded_bytes"`           // 总上传量
	DownloadSpeedLimit     int64  `json:"download_speed_limit"`     // 下载限速
	UploadSpeedLimit       int64  `json:"upload_speed_limit"`       // 上传限速
	DHT                    int    `json:"dht_nodes"`                // DHT节点数
	ConnectionStatus       string `json:"connection_status"`        // 连接状态
	TotalPeersConnected    int    `json:"total_peers_connected"`    // 连接的peers数
	TotalBuffersSize       int64  `json:"total_buffers_size"`       // 缓冲区大小
	TotalWasted            int64  `json:"total_wasted"`             // 浪费的数据
	AlternativeSpeedLimits bool   `json:"alternative_speed_limits"` // 备用速度限制
	FreeSpaceOnDisk        int64  `json:"free_space_on_disk"`       // 磁盘剩余空间
}

// LogEntry 日志条目
type LogEntry struct {
	ID        int64  `json:"id"`        // 日志ID
	Message   string `json:"message"`   // 日志消息
	Timestamp int64  `json:"timestamp"` // 时间戳
	Type      int    `json:"type"`      // 类型: 1=普通, 2=信息, 4=警告, 8=严重
}
