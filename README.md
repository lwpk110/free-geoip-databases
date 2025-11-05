# Free GeoIP Database - Auto-Updated 🌍

[![Update GeoLite2 Database](https://github.com/YOUR_USERNAME/YOUR_REPO/actions/workflows/update-geolite2.yml/badge.svg)](https://github.com/YOUR_USERNAME/YOUR_REPO/actions/workflows/update-geolite2.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

> **免费、自动更新的 GeoIP 数据库** - 提供 MaxMind GeoLite2 和 DB-IP 数据库的自动化下载和发布服务

无需注册 MaxMind 账号，无需申请 License Key，直接从 [Releases](../../releases) 下载最新数据库！

## 🎯 为什么选择这个项目？

- ✅ **完全免费** - 无需注册，无需 License Key
- 🤖 **自动更新** - 每周自动跟随 MaxMind 官方更新（周二、周五）
- � **开箱即用** - 直接下载 `.mmdb` 文件即可使用
- �🚀 **多种数据库** - 提供 City、Country、ASN 等多种数据库
- ✅ **质量保证** - 自动化测试确保数据库完整性

## � 快速下载

访问 [**Releases 页面**](../../releases/latest) 下载最新数据库文件：

| 数据库 | 说明 | 下载链接 |
|--------|------|----------|
| **GeoLite2-City** | 城市级别地理位置数据 | [下载](../../releases/latest/download/GeoLite2-City.mmdb) |
| **GeoLite2-Country** | 国家级别地理位置数据 | [下载](../../releases/latest/download/GeoLite2-Country.mmdb) |
| **GeoLite2-ASN** | ASN 网络运营商数据 | [下载](../../releases/latest/download/GeoLite2-ASN.mmdb) |

### 命令行下载

```bash
# 下载 City 数据库
curl -L -o GeoLite2-City.mmdb \
  https://github.com/YOUR_USERNAME/YOUR_REPO/releases/latest/download/GeoLite2-City.mmdb

# 下载 Country 数据库
curl -L -o GeoLite2-Country.mmdb \
  https://github.com/YOUR_USERNAME/YOUR_REPO/releases/latest/download/GeoLite2-Country.mmdb

# 下载 ASN 数据库
curl -L -o GeoLite2-ASN.mmdb \
  https://github.com/YOUR_USERNAME/YOUR_REPO/releases/latest/download/GeoLite2-ASN.mmdb
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

- **自动更新**: 每周二和周五（UTC 10:00 / 北京时间 18:00）
- **数据源**: MaxMind 官方 GeoLite2 数据库
- **更新策略**: 跟随 MaxMind 官方更新周期

## 📊 支持的数据库

| 数据库类型 | 包含信息 | 文件大小 |
|-----------|---------|----------|
| **City** | 国家、省份/州、城市、邮编、经纬度、时区 | ~70 MB |
| **Country** | 国家、大洲 | ~6 MB |
| **ASN** | 自治系统号、网络运营商 | ~8 MB |

## 🛠️ 本地运行测试工具

本项目提供了一个简单的 Go 语言查询工具，可以测试数据库文件：

```bash
# 1. 克隆仓库
git clone https://github.com/YOUR_USERNAME/YOUR_REPO.git
cd YOUR_REPO

# 2. 下载数据库文件（见上方下载说明）

# 3. 安装依赖
go mod tidy

# 4. 运行测试工具
go run main.go
```

## ⚙️ Fork 本项目实现自动更新

如果你想 Fork 本项目并实现自己的自动更新：

1. **Fork 本仓库**

2. **添加 MaxMind License Key**
   - 访问 [MaxMind 注册](https://www.maxmind.com/en/geolite2/signup)
   - 生成 License Key
   - 在你的仓库中：Settings → Secrets and variables → Actions
   - 添加 Secret: `MAXMIND_LICENSE_KEY`

3. **启用 GitHub Actions**
   - 工作流将自动运行
   - 或手动触发：Actions → Update GeoLite2 Database → Run workflow

## 📋 许可与声明

- **项目代码**: MIT License
- **GeoLite2 数据库**: 由 MaxMind 提供，需遵守 [MaxMind EULA](https://www.maxmind.com/en/geolite2/eula)
- **使用声明**: This product includes GeoLite2 data created by MaxMind, available from [https://www.maxmind.com](https://www.maxmind.com)

### 重要提示

- GeoLite2 是免费版本，精确度低于商业版 GeoIP2
- 仅供学习、测试和非商业用途
- 商业用途请购买 MaxMind 商业授权

## 🌟 Star History

如果这个项目对你有帮助，请给一个 ⭐ Star！这将激励我持续维护和更新。

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

## � 支持

- 🐛 [报告 Bug](../../issues)
- 💡 [功能建议](../../issues)
- ⭐ [Star 本项目](../../stargazers)
