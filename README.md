# Down-Nexus API （开发中）

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
- PostgreSQL 12+
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

首次运行时，系统会创建空的数据库。您需要手动配置客户端：

1. **直接修改数据库** (推荐)
   ```sql
   INSERT INTO client_configs (client_id, type, host, username, password, enabled) 
   VALUES ('qb-home', 'qbittorrent', 'http://localhost:8080', 'your_username', 'your_password', true);
   ```

2. **参考配置文件**: `config.example.json`

3. **使用 API 接口** (开发中)

⚠️ **注意**: 首次运行时如果数据库为空，服务器会退出并提示配置客户端。

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
├── data/               # 数据文件目录 (已弃用)
└── config.example.json # 配置文件示例
```

## 开发说明

### 添加新的客户端支持

1. 在 `pkg/clients/` 下创建新的适配器
2. 实现 `DownloaderClient` 接口
3. 在 `loadClientsFromDB` 函数中添加新的 case

### 数据库配置

#### PostgreSQL 设置

1. **创建数据库和用户**:
```sql
CREATE DATABASE downnexus;
CREATE USER downnexus WITH PASSWORD 'downnexus';
GRANT ALL PRIVILEGES ON DATABASE downnexus TO downnexus;
```

2. **客户端配置表**:
```sql
CREATE TABLE client_configs (
    id SERIAL PRIMARY KEY,
    client_id TEXT UNIQUE NOT NULL,
    type TEXT NOT NULL,
    host TEXT NOT NULL,
    username TEXT NOT NULL,
    password TEXT NOT NULL,
    enabled BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
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