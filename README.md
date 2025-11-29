# Down-Nexus API

🌟 多客户端种子管理系统

## 功能特性

- 🔄 **统一接口**: 支持 qBittorrent 和 Transmission 客户端
- 🌐 **RESTful API**: 基于 Gin 框架的高性能 API
- 🗄️ **数据库管理**: SQLite 存储客户端配置
- ⚡ **种子操作**: 获取、添加、暂停、恢复、删除种子
- 🔧 **动态配置**: 从数据库加载客户端配置

## 快速开始

### 环境要求

- Go 1.19+
- qBittorrent 或 Transmission 客户端

### 安装运行

```bash
# 克隆仓库
git clone <repository-url>
cd Down-Nexus

# 安装依赖
go mod tidy

# 运行服务器
go run cmd/server/main.go
```

### 配置客户端

首次运行时，系统会自动创建数据库并插入默认配置。您可以通过以下方式管理客户端配置：

1. **直接修改数据库** (推荐用于生产环境)
2. **使用 API 接口** (开发中)
3. **参考配置文件**: `config.example.json`

## API 接口

### 基础接口
- `GET /` - 欢迎页面
- `GET /health` - 健康检查

### 种子管理
- `GET /api/v1/torrents` - 获取所有种子
- `POST /api/v1/torrents` - 添加种子
- `POST /api/v1/torrents/pause` - 暂停种子
- `POST /api/v1/torrents/resume` - 恢复种子
- `DELETE /api/v1/torrents` - 删除种子

### 客户端管理
- `GET /api/v1/clients` - 获取客户端列表

## 项目结构

```
Down-Nexus/
├── cmd/server/          # 主服务器程序
├── internal/            # 内部包
│   ├── api/            # API 处理器和路由
│   ├── core/           # 核心业务逻辑
│   └── models/         # 数据模型
├── pkg/                # 公共包
│   ├── clients/        # 客户端适配器
│   └── database/       # 数据库连接
├── data/               # 数据库文件 (本地)
└── config.example.json # 配置文件示例
```

## 开发说明

### 添加新的客户端支持

1. 在 `pkg/clients/` 下创建新的适配器
2. 实现 `DownloaderClient` 接口
3. 在 `loadClientsFromDB` 函数中添加新的 case

### 数据库模型

客户端配置存储在 `client_configs` 表中：

```sql
CREATE TABLE client_configs (
    id INTEGER PRIMARY KEY,
    client_id TEXT UNIQUE NOT NULL,
    type TEXT NOT NULL,
    host TEXT NOT NULL,
    username TEXT NOT NULL,
    password TEXT NOT NULL,
    enabled BOOLEAN DEFAULT TRUE
);
```

## 安全注意事项

⚠️ **重要**: 
- 请勿将包含真实密码的配置文件提交到版本控制
- 生产环境请使用强密码
- 建议在防火墙后运行此服务

## 许可证

MIT License

## 贡献

欢迎提交 Issue 和 Pull Request！