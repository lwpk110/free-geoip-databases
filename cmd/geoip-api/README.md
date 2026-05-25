# GeoIP API Server

基于 Golang 实现的自托管 GeoIP 查询服务，提供 RESTful API 和 Web 界面。

## 功能特性

- ✅ **RESTful API** - 提供标准的 HTTP API 接口
- ✅ **Web 界面** - 用户友好的查询界面
- ✅ **多数据库支持** - 兼容 GeoLite2 和 DB-IP 数据库
- ✅ **实时查询** - 快速的 IP 地理位置查询
- ✅ **完整信息** - 返回国家、城市、经纬度、时区等详细信息
- ✅ **轻量部署** - 单一二进制文件，易于部署
- ✅ **热更新** - 数据库文件更新后自动热加载，无需重启
- ✅ **回退策略** - 多数据库优先级查询，失败自动 fallback
- ✅ **批量查询** - `POST /api/v1/lookup/batch`
- ✅ **可观测性** - `GET /api/v1/metrics` 输出运行指标
- ✅ **安全控制** - API Key 鉴权 + 限流

## 快速开始

### 1. 准备数据库文件

首先需要下载 GeoIP 数据库文件（GeoLite2 或 DB-IP）：

```bash
# 下载 DB-IP 数据库（推荐，无需注册）
cd /home/runner/work/free-geoip-databases/free-geoip-databases
./scripts/download_dbip.sh all

# 或下载 GeoLite2 数据库（需要 License Key）
./scripts/download_geolite2.sh YOUR_LICENSE_KEY
```

### 2. 编译并运行

```bash
# 编译
cd cmd/geoip-api
go build -o geoip-api

# 运行（默认使用当前目录下的 GeoLite2-City.mmdb）
./geoip-api

# 或指定数据库路径
GEOIP_DB_PATH=/path/to/GeoLite2-City.mmdb ./geoip-api

# 指定端口
PORT=8888 ./geoip-api
```

### 3. 访问服务

- **Web 界面**: http://localhost:8080
- **API 文档**: http://localhost:8080/api/v1/

## API 使用

### 查询 IP 地址

**请求**:
```bash
GET /api/v1/lookup/:ip
GET /api/v1/lookup?ip=8.8.8.8
```

**示例**:
```bash
curl http://localhost:8080/api/v1/lookup/8.8.8.8
```

**响应**:
```json
{
  "ip": "8.8.8.8",
  "country": {
    "iso_code": "US",
    "names": {
      "en": "United States",
      "zh-CN": "美国"
    }
  },
  "city": {
    "geoname_id": 0,
    "names": {}
  },
  "location": {
    "latitude": 37.751,
    "longitude": -97.822,
    "time_zone": "America/Chicago",
    "accuracy_radius": 1000
  },
  "continent": {
    "code": "NA",
    "names": {
      "en": "North America",
      "zh-CN": "北美洲"
    }
  },
  "traits": {
    "is_anonymous_proxy": false,
    "is_satellite_provider": false
  }
}
```

### 查询系统状态

**请求**:
```bash
GET /api/v1/status
```

**示例**:
```bash
curl http://localhost:8080/api/v1/status
```

**响应**:
```json
{
  "status": "healthy",
  "database_path": "GeoLite2-City.mmdb",
  "database_type": "GeoLite2",
  "last_loaded": "2025-11-13T04:00:00Z",
  "uptime": "1h30m45s"
}
```

### 健康检查

**请求**:
```bash
GET /api/v1/health
```

**响应**:
```json
{
  "status": "ok"
}
```

### 批量查询

```bash
curl -X POST http://localhost:8080/api/v1/lookup/batch \
  -H 'Content-Type: application/json' \
  -d '{"ips":["8.8.8.8","1.1.1.1"]}'
```

### 观测指标

```bash
curl http://localhost:8080/api/v1/metrics
```

### 手动触发热更新

```bash
curl -X POST http://localhost:8080/api/v1/admin/reload
```

## 环境变量

| 变量 | 说明 | 默认值 |
|------|------|--------|
| `GEOIP_DB_PATH` | 单个 GeoIP 数据库文件路径 | `GeoLite2-City.mmdb` |
| `GEOIP_DB_PATHS` | 多数据库路径（逗号分隔，按优先级） | 空 |
| `PORT` | 服务监听端口 | `8080` |
| `GIN_MODE` | Gin 运行模式 (debug/release) | `release` |
| `HOT_RELOAD_INTERVAL` | 热更新轮询周期 | `60s` |
| `CACHE_TTL` | 查询缓存 TTL | `300s` |
| `RATE_LIMIT_PER_MIN` | 每分钟请求限额（按 key/IP） | `240` |
| `API_KEYS` | API Key 列表（逗号分隔） | 空（不启用） |

## Docker 部署

### 使用 Dockerfile

```dockerfile
FROM golang:1.23-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN cd cmd/geoip-api && go build -o /geoip-api

FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/
COPY --from=builder /geoip-api .

# 复制数据库文件（需要提前下载）
COPY GeoLite2-City.mmdb .

EXPOSE 8080
CMD ["./geoip-api"]
```

### 使用 Docker Compose

```yaml
version: '3.8'

services:
  geoip-api:
    build: .
    ports:
      - "8080:8080"
    environment:
      - GEOIP_DB_PATH=/data/GeoLite2-City.mmdb
      - PORT=8080
    volumes:
      - ./GeoLite2-City.mmdb:/data/GeoLite2-City.mmdb:ro
    restart: unless-stopped
```

运行:
```bash
docker-compose up -d
```

## 性能优化

### 1. 使用缓存

可以在 API 前面添加 Redis 或其他缓存层，缓存查询结果。

### 2. 使用 CDN

对于静态资源，可以使用 CDN 加速。

### 3. 负载均衡

对于高并发场景，可以部署多个实例，使用 Nginx 或其他负载均衡器。

## 数据库更新

### 手动更新

```bash
# 下载最新数据库
./scripts/download_dbip.sh all

# 重启服务（未来版本将支持热重载）
killall geoip-api
./geoip-api
```

### 自动更新（计划中）

未来版本将支持：
- 定时自动下载最新数据库
- 数据库热重载，无需重启服务
- 多数据库版本管理

## 开发

### 运行开发服务器

```bash
cd cmd/geoip-api
GIN_MODE=debug go run main.go
```

### 测试 API

```bash
# 测试查询
curl http://localhost:8080/api/v1/lookup/8.8.8.8

# 测试状态
curl http://localhost:8080/api/v1/status

# 测试健康检查
curl http://localhost:8080/api/v1/health
```

## 许可证

MIT License

## 致谢

- 数据来源: [MaxMind GeoLite2](https://www.maxmind.com) / [DB-IP](https://db-ip.com)
- GeoIP2 库: [oschwald/geoip2-golang](https://github.com/oschwald/geoip2-golang)
- Web 框架: [Gin](https://github.com/gin-gonic/gin)
