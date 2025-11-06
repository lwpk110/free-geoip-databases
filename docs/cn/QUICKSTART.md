# 快速开始指南

[ 中文 | [English](../en/QUICKSTART.md) ]

本指南将帮助你快速设置和运行 GeoIP 查询工具，并配置自动更新功能。

## 📋 前置要求

- Git
- Go 1.21 或更高版本
- GitHub 账号（用于自动更新功能）
- MaxMind 账号（用于获取 License Key）

## 🚀 5 分钟快速开始

### 步骤 1: 克隆仓库

```bash
git clone https://github.com/YOUR_USERNAME/geoip.git
cd geoip
```

### 步骤 2: 获取数据库

#### 选项 A：从 Release 下载（推荐）

```bash
# 如果仓库已经有 Release，直接下载
curl -L -o GeoLite2-City.mmdb \
  https://github.com/YOUR_USERNAME/geoip/releases/latest/download/GeoLite2-City.mmdb
```

#### 选项 B：手动下载

1. 访问 [MaxMind 注册](https://www.maxmind.com/en/geolite2/signup) 注册账号
2. 生成 License Key
3. 运行下载脚本：

```bash
chmod +x download_geolite2.sh
./download_geolite2.sh YOUR_LICENSE_KEY
```

### 步骤 3: 安装依赖并运行

```bash
# 安装 Go 依赖
go mod tidy

# 运行程序
go run main.go
```

你应该能看到类似这样的输出：

```text
=== 开始查询多个 IP 地址的地理位置 ===

查询 IP: 8.8.8.8
  国家: United States (美国)
  省份/州:
  城市:
  时区: America/Chicago
  坐标: 37.7510, -97.8220
  邮编:
--------------------
...
```

## 🤖 配置自动更新（可选）

如果你想让数据库自动更新，请按照以下步骤操作：

### 步骤 1: Fork 仓库

点击 GitHub 页面右上角的 "Fork" 按钮

### 步骤 2: 获取 MaxMind License Key

1. 访问 [https://www.maxmind.com/en/geolite2/signup](https://www.maxmind.com/en/geolite2/signup)
2. 注册并验证邮箱
3. 登录后访问 [https://www.maxmind.com/en/accounts/current/license-key](https://www.maxmind.com/en/accounts/current/license-key)
4. 点击 "Generate new license key"
5. 复制生成的 License Key

### 步骤 3: 添加 GitHub Secret

1. 进入你 Fork 的仓库
2. 点击 **Settings** → **Secrets and variables** → **Actions**
3. 点击 **New repository secret**
4. 添加：
   - Name: `MAXMIND_LICENSE_KEY`
   - Value: 你的 License Key
5. 点击 **Add secret**

### 步骤 4: 启用工作流权限

1. 在 **Settings** → **Actions** → **General**
2. 滚动到 **Workflow permissions**
3. 选择 **Read and write permissions**
4. 勾选 **Allow GitHub Actions to create and approve pull requests**
5. 点击 **Save**

### 步骤 5: 手动触发第一次更新

1. 进入 **Actions** 标签页
2. 选择 "Update GeoLite2 Database" 工作流
3. 点击 **Run workflow**
4. 选择 `main` 分支
5. 点击绿色的 **Run workflow** 按钮

等待几分钟，你应该能在 **Releases** 页面看到新发布的数据库！

## 📅 自动更新说明

配置完成后：

- 工作流会在每周二和周五 UTC 10:00（北京时间 18:00）自动运行
- 每次更新会创建新的 Release
- Release 包含三个数据库文件：
  - GeoLite2-City.mmdb
  - GeoLite2-Country.mmdb
  - GeoLite2-ASN.mmdb

## 💡 使用技巧

### 在你的项目中使用

```go
package main

import (
    "fmt"
    "log"
    "net"

    "github.com/oschwald/geoip2-golang"
)

func main() {
    db, err := geoip2.Open("GeoLite2-City.mmdb")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    ip := net.ParseIP("1.1.1.1")
    record, err := db.City(ip)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("国家: %s\n", record.Country.Names["zh-CN"])
    fmt.Printf("城市: %s\n", record.City.Names["zh-CN"])
}
```

### 从最新 Release 自动下载

在你的脚本中添加：

```bash
#!/bin/bash

REPO="YOUR_USERNAME/geoip"
DB_FILE="GeoLite2-City.mmdb"

echo "下载最新的 GeoLite2 数据库..."
curl -L -o "$DB_FILE" \
  "https://github.com/${REPO}/releases/latest/download/${DB_FILE}"

echo "✓ 下载完成: $DB_FILE"
```

### 定期更新本地数据库

创建一个 cron 任务：

```bash
# 编辑 crontab
crontab -e

# 添加这一行（每周三凌晨 2 点更新）
0 2 * * 3 cd /path/to/your/project && curl -L -o GeoLite2-City.mmdb https://github.com/YOUR_USERNAME/geoip/releases/latest/download/GeoLite2-City.mmdb
```

## 🔍 故障排查

### "数据库文件未找到"

确保 `GeoLite2-City.mmdb` 文件在程序运行目录中：

```bash
ls -lh GeoLite2-City.mmdb
```

### "License Key 无效"

1. 检查 License Key 是否正确复制（没有多余空格）
2. 确认 License Key 没有过期
3. 访问 MaxMind 账户验证状态

### "工作流运行失败"

1. 检查 GitHub Actions 日志
2. 确认 Secret 已正确添加
3. 验证工作流权限已启用

## 📚 更多资源

- [完整 README](README.md)
- [下载指南](DOWNLOAD_GUIDE.md)
- [GitHub Actions 设置详解](.github/SETUP.md)
- [MaxMind 官方文档](https://dev.maxmind.com/geoip/geolite2-free-geolocation-data)

## 🆘 需要帮助？

- 查看 [Issues](../../issues) 页面
- 提交新的 [Issue](../../issues/new)
- 阅读 [MaxMind 文档](https://dev.maxmind.com/geoip/)

---

祝你使用愉快！🎉
