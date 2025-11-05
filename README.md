# Free GeoIP Databases - Auto-Updated 🌍

[![Update GeoLite2 Database](https://github.com/lwpk110/free-geoip-databases/actions/workflows/update-geolite2.yml/badge.svg)](https://github.com/lwpk110/free-geoip-databases/actions/workflows/update-geolite2.yml)
[![Update DB-IP Database](https://github.com/lwpk110/free-geoip-databases/actions/workflows/update-dbip.yml/badge.svg)](https://github.com/lwpk110/free-geoip-databases/actions/workflows/update-dbip.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

> **免费、自动更新的 GeoIP 数据库集合** - 提供 MaxMind GeoLite2、DB-IP 等数据库的自动化下载和发布服务

无需注册 MaxMind 账号，无需申请 License Key，直接从 [Releases](../../releases) 下载最新数据库！

## 🎯 为什么选择这个项目？

- ✅ **完全免费** - 无需注册，无需 License Key（DB-IP 版本）
- 🤖 **自动更新** - GeoLite2 每周更新（周二、周五），DB-IP 每月更新（1号、15号）
- 📦 **开箱即用** - 直接下载 `.mmdb` 文件即可使用
-  **多种数据库** - 提供 City、Country、ASN 等多种数据库
- 🌐 **多种来源** - 同时提供 MaxMind GeoLite2 和 DB-IP 数据库
- ✅ **质量保证** - 自动化测试确保数据库完整性

## 📥 快速下载

### MaxMind GeoLite2 数据库

访问 [**Releases 页面**](../../releases/latest) 下载最新数据库文件：

| 数据库 | 说明 | 文件大小 | 授权 | 下载链接 |
|--------|------|---------|------|----------|
| **GeoLite2-City** | 城市级别地理位置数据 | ~70 MB | GeoLite2 EULA | [下载](../../releases/latest/download/GeoLite2-City.mmdb) |
| **GeoLite2-Country** | 国家级别地理位置数据 | ~6 MB | GeoLite2 EULA | [下载](../../releases/latest/download/GeoLite2-Country.mmdb) |
| **GeoLite2-ASN** | ASN 网络运营商数据 | ~8 MB | GeoLite2 EULA | [下载](../../releases/latest/download/GeoLite2-ASN.mmdb) |

### DB-IP 数据库

查看 [**DB-IP Releases**](../../releases?q=dbip&expanded=true) 下载 DB-IP 数据库：

| 数据库 | 说明 | 文件大小 | 授权 | 下载链接 |
|--------|------|---------|------|----------|
| **DB-IP City Lite** | 城市级别地理位置数据 | ~130 MB | CC BY 4.0 | [查看 Releases](../../releases?q=dbip) |
| **DB-IP Country Lite** | 国家级别地理位置数据 | ~7 MB | CC BY 4.0 | [查看 Releases](../../releases?q=dbip) |
| **DB-IP ASN Lite** | ASN 网络运营商数据 | ~9 MB | CC BY 4.0 | [查看 Releases](../../releases?q=dbip) |

> **注意**: DB-IP 数据库文件名包含年月信息，如 `dbip-city-lite-2024-11.mmdb`

### 命令行下载

```bash
# 下载 City 数据库
curl -L -o GeoLite2-City.mmdb \
  https://github.com/lwpk110/free-geoip-databases/releases/latest/download/GeoLite2-City.mmdb

# 下载 Country 数据库
curl -L -o GeoLite2-Country.mmdb \
  https://github.com/lwpk110/free-geoip-databases/releases/latest/download/GeoLite2-Country.mmdb

# 下载 ASN 数据库
curl -L -o GeoLite2-ASN.mmdb \
  https://github.com/lwpk110/free-geoip-databases/releases/latest/download/GeoLite2-ASN.mmdb
```

或使用项目提供的下载脚本：

```bash
# MaxMind GeoLite2
./scripts/download_geolite2.sh <YOUR_LICENSE_KEY>

# DB-IP (无需 License Key)
./scripts/download_dbip.sh all
```

## 🚀 使用示例

### Go 语言

```go
package main

import (
    "fmt"
    "log"
    "net"

    "github.com/oschwald/geoip2-golang"
)

func main() {
    // 打开数据库
    db, err := geoip2.Open("GeoLite2-City.mmdb")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    // 查询 IP
    ip := net.ParseIP("8.8.8.8")
    record, err := db.City(ip)
    if err != nil {
        log.Fatal(err)
    }

    // 输出结果
    fmt.Printf("国家: %s\n", record.Country.Names["zh-CN"])
    fmt.Printf("城市: %s\n", record.City.Names["zh-CN"])
    fmt.Printf("经纬度: %.4f, %.4f\n", record.Location.Latitude, record.Location.Longitude)
}
```

### Python

```python
import geoip2.database

# 打开数据库
with geoip2.database.Reader('GeoLite2-City.mmdb') as reader:
    # 查询 IP
    response = reader.city('8.8.8.8')

    # 输出结果
    print(f"国家: {response.country.names['zh-CN']}")
    print(f"城市: {response.city.names.get('zh-CN', 'N/A')}")
    print(f"经纬度: {response.location.latitude}, {response.location.longitude}")
```

### Node.js

```javascript
const maxmind = require('maxmind');

// 打开数据库
maxmind.open('GeoLite2-City.mmdb').then(lookup => {
    // 查询 IP
    const result = lookup.get('8.8.8.8');

    // 输出结果
    console.log(`国家: ${result.country.names['zh-CN']}`);
    console.log(`城市: ${result.city.names['zh-CN']}`);
    console.log(`经纬度: ${result.location.latitude}, ${result.location.longitude}`);
});
```

## 🔄 更新频率

### MaxMind GeoLite2
- **自动更新**: 每周二和周五（UTC 10:00 / 北京时间 18:00）
- **数据源**: MaxMind 官方 GeoLite2 数据库
- **更新策略**: 跟随 MaxMind 官方更新周期

### DB-IP
- **自动更新**: 每月1号和15号（UTC 10:00 / 北京时间 18:00）
- **数据源**: DB-IP 官方免费数据库
- **更新策略**: 下载当月最新版本
- **授权**: Creative Commons Attribution 4.0 (CC BY 4.0)

## 📊 数据库对比

### MaxMind GeoLite2 vs DB-IP

| 特性 | GeoLite2 | DB-IP Lite |
|------|----------|-----------|
| **授权** | 需遵守 MaxMind EULA | CC BY 4.0（更自由） |
| **更新频率** | 每周2次 | 每月1次 |
| **数据准确性** | ⭐⭐⭐⭐ | ⭐⭐⭐⭐ |
| **文件格式** | .mmdb | .mmdb（兼容） |
| **使用限制** | 非商业用途优先 | 署名即可商用 |
| **注册要求** | 需要（本项目已处理） | 无需注册 |

### 支持的数据库类型

| 数据库类型 | 包含信息 | 文件大小 |
|-----------|---------|----------|
| **City** | 国家、省份/州、城市、邮编、经纬度、时区 | ~70 MB |
| **Country** | 国家、大洲 | ~6 MB |
| **ASN** | 自治系统号、网络运营商 | ~8 MB |

## 🛠️ 本地运行测试工具

本项目提供了简单的 Go 语言查询工具，可以测试数据库文件。

### 项目结构

```
free-geoip-databases/
├── .github/workflows/    # GitHub Actions 自动化工作流
│   ├── update-geolite2.yml
│   ├── update-dbip.yml
│   └── test-database.yml
├── scripts/              # 下载脚本
│   ├── download_geolite2.sh
│   └── download_dbip.sh
├── examples/             # 示例代码
│   ├── query/            # IP 查询示例
│   │   └── main.go
│   └── test/             # 测试程序
│       └── test_cities.go
├── docs/                 # 详细文档
│   ├── QUICKSTART.md
│   └── TESTING.md
├── README.md
├── LICENSE
├── go.mod
└── go.sum
```

### 快速开始

```bash
# 1. 克隆仓库
git clone https://github.com/lwpk110/free-geoip-databases.git
cd free-geoip-databases

# 2. 下载数据库文件
./scripts/download_geolite2.sh <YOUR_LICENSE_KEY>
# 或使用 DB-IP (无需 License Key)
./scripts/download_dbip.sh all

# 3. 安装依赖
go mod tidy

# 4. 运行查询示例
cd examples/query
go run main.go

# 5. 运行测试
cd ../test
go run test_cities.go
```

更多详情请查看 [examples/README.md](examples/README.md) 和 [docs/TESTING.md](docs/TESTING.md)。

## ⚙️ Fork 本项目实现自动更新

### 配置 MaxMind GeoLite2 自动更新

如果你想 Fork 本项目并实现自己的 GeoLite2 自动更新：

1. **Fork 本仓库**

2. **添加 MaxMind License Key**
   - 访问 [MaxMind 注册](https://www.maxmind.com/en/geolite2/signup)
   - 生成 License Key
   - 在你的仓库中：Settings → Secrets and variables → Actions
   - 添加 Secret: `MAXMIND_LICENSE_KEY`

3. **添加 GitHub Token**
   - 在你的仓库中：Settings → Secrets and variables → Actions
   - 添加 Secret: `GEOIP_ACCESS_TOKEN`（使用有 repo 权限的 Personal Access Token）

4. **启用 GitHub Actions**
   - 工作流将自动运行
   - 或手动触发：Actions → Update GeoLite2 Database → Run workflow

### 配置 DB-IP 自动更新

DB-IP 数据库无需任何配置，只需：

1. **Fork 本仓库**

2. **添加 GitHub Token**
   - 在你的仓库中：Settings → Secrets and variables → Actions
   - 添加 Secret: `GEOIP_ACCESS_TOKEN`

3. **启用 GitHub Actions**
   - DB-IP 工作流会自动运行
   - 或手动触发：Actions → Update DB-IP Database → Run workflow

> **提示**: DB-IP 不需要注册账号或 License Key，完全免费且开放！

## 📋 许可与声明

- **项目代码**: MIT License
- **GeoLite2 数据库**: 由 MaxMind 提供，需遵守 [MaxMind EULA](https://www.maxmind.com/en/geolite2/eula)
  - This product includes GeoLite2 data created by MaxMind, available from [https://www.maxmind.com](https://www.maxmind.com)
- **DB-IP 数据库**: 由 DB-IP 提供，使用 [Creative Commons Attribution 4.0](https://creativecommons.org/licenses/by/4.0/) 授权
  - 使用时需要署名：Contains data from https://db-ip.com

### 重要提示

- GeoLite2 和 DB-IP Lite 都是免费版本，精确度低于商业版
- GeoLite2: 仅供学习、测试和非商业用途，商业用途请购买 MaxMind 商业授权
- DB-IP: 可商用，但需要署名（Attribution required）
- 如需更高精确度或商业支持，请考虑购买商业版数据库

## 🌟 Star History

如果这个项目对你有帮助，请给一个 ⭐ Star！这将激励我持续维护和更新。

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

## � 支持

- 🐛 [报告 Bug](../../issues)
- 💡 [功能建议](../../issues)
- ⭐ [Star 本项目](../../stargazers)
