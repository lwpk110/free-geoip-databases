# GeoIP API Server

A self-hosted GeoIP query service implemented in Golang, providing RESTful API and Web interface.

## Features

- ✅ **RESTful API** - Standard HTTP API endpoints
- ✅ **Web Interface** - User-friendly query interface
- ✅ **Multi-Database Support** - Compatible with GeoLite2 and DB-IP databases
- ✅ **Real-time Query** - Fast IP geolocation lookup
- ✅ **Complete Information** - Returns country, city, coordinates, timezone, and more
- ✅ **Lightweight Deployment** - Single binary, easy to deploy

## Quick Start

### 1. Prepare Database File

First, download a GeoIP database file (GeoLite2 or DB-IP):

```bash
# Download DB-IP database (recommended, no registration required)
cd /home/runner/work/free-geoip-databases/free-geoip-databases
./scripts/download_dbip.sh all

# Or download GeoLite2 database (requires License Key)
./scripts/download_geolite2.sh YOUR_LICENSE_KEY
```

### 2. Build and Run

```bash
# Build
cd cmd/geoip-api
go build -o geoip-api

# Run (uses GeoLite2-City.mmdb in current directory by default)
./geoip-api

# Or specify database path
GEOIP_DB_PATH=/path/to/GeoLite2-City.mmdb ./geoip-api

# Specify port
PORT=8888 ./geoip-api
```

### 3. Access Service

- **Web Interface**: http://localhost:8080
- **API Documentation**: http://localhost:8080/api/v1/

## API Usage

### Query IP Address

**Request**:
```bash
GET /api/v1/lookup/:ip
GET /api/v1/lookup?ip=8.8.8.8
```

**Example**:
```bash
curl http://localhost:8080/api/v1/lookup/8.8.8.8
```

**Response**:
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

### Query System Status

**Request**:
```bash
GET /api/v1/status
```

**Example**:
```bash
curl http://localhost:8080/api/v1/status
```

**Response**:
```json
{
  "status": "healthy",
  "database_path": "GeoLite2-City.mmdb",
  "database_type": "GeoLite2",
  "last_loaded": "2025-11-13T04:00:00Z",
  "uptime": "1h30m45s"
}
```

### Health Check

**Request**:
```bash
GET /api/v1/health
```

**Response**:
```json
{
  "status": "ok"
}
```

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `GEOIP_DB_PATH` | GeoIP database file path | `GeoLite2-City.mmdb` |
| `PORT` | Service listening port | `8080` |
| `GIN_MODE` | Gin running mode (debug/release) | `release` |

## Docker Deployment

### Using Dockerfile

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

# Copy database file (download first)
COPY GeoLite2-City.mmdb .

EXPOSE 8080
CMD ["./geoip-api"]
```

### Using Docker Compose

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

Run:
```bash
docker-compose up -d
```

## Performance Optimization

### 1. Use Cache

You can add Redis or other cache layer in front of the API to cache query results.

### 2. Use CDN

For static resources, you can use CDN for acceleration.

### 3. Load Balancing

For high concurrency scenarios, you can deploy multiple instances and use Nginx or other load balancers.

## Database Updates

### Manual Update

```bash
# Download latest database
./scripts/download_dbip.sh all

# Restart service (future versions will support hot reload)
killall geoip-api
./geoip-api
```

### Automatic Update (Planned)

Future versions will support:
- Scheduled automatic download of latest database
- Database hot reload without service restart
- Multiple database version management

## Development

### Run Development Server

```bash
cd cmd/geoip-api
GIN_MODE=debug go run main.go
```

### Test API

```bash
# Test query
curl http://localhost:8080/api/v1/lookup/8.8.8.8

# Test status
curl http://localhost:8080/api/v1/status

# Test health check
curl http://localhost:8080/api/v1/health
```

## License

MIT License

## Acknowledgments

- Data source: [MaxMind GeoLite2](https://www.maxmind.com) / [DB-IP](https://db-ip.com)
- GeoIP2 library: [oschwald/geoip2-golang](https://github.com/oschwald/geoip2-golang)
- Web framework: [Gin](https://github.com/gin-gonic/gin)
