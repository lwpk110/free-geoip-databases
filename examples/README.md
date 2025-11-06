# 示例代码

**[ 中文 | [English](README_EN.md) ]**

本目录包含使用 GeoIP 数据库的示例程序。

## 目录结构

```text
examples/
├── query/        # IP 地理位置查询示例
│   └── main.go
└── test/         # 数据库测试程序
    └── test_cities.go
```

## 使用说明

### 1. 查询示例 (query)

查询 IP 地址的详细地理位置信息。

```bash
# 进入示例目录
cd examples/query

# 运行示例（需要在根目录有 GeoLite2-City.mmdb 文件）
go run main.go
```

**功能**:

- 自动获取当前主机的公网 IP
- 显示完整的地理位置信息（大陆、国家、省份、城市、经纬度等）
- 同时输出 JSON 格式数据

### 2. 测试程序 (test)

测试数据库是否正常工作。

```bash
# 进入测试目录
cd examples/test

# 运行测试（需要在根目录有 GeoLite2-City.mmdb 文件）
go run test_cities.go
```

**功能**:

- 测试多个知名 DNS IP 的查询
- 验证数据库完整性
- 显示测试结果统计

## 前置条件

1. **下载数据库文件**

使用下载脚本获取数据库：

```bash
# 从项目根目录运行
./scripts/download_geolite2.sh <YOUR_LICENSE_KEY>
# 或
./scripts/download_dbip.sh all
```

或从 GitHub Releases 下载：

```bash
curl -L -o GeoLite2-City.mmdb \
  https://github.com/lwpk110/free-geoip-databases/releases/latest/download/GeoLite2-City.mmdb
```

1. **安装依赖**

```bash
go mod download
```

## 数据库文件位置

示例程序期望数据库文件位于**项目根目录**：

- `../../GeoLite2-City.mmdb`（从 examples/query 或 examples/test 看）

## 开发自己的程序

参考这些示例创建你自己的 GeoIP 查询程序：

```go
package main

import (
    "github.com/oschwald/geoip2-golang"
    "net"
)

func main() {
    db, _ := geoip2.Open("GeoLite2-City.mmdb")
    defer db.Close()

    ip := net.ParseIP("8.8.8.8")
    record, _ := db.City(ip)

    println(record.Country.Names["zh-CN"])
}
```
