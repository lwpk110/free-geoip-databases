# 测试说明

## 本地测试

### 测试 GeoLite2 数据库

```bash
# 1. 确保已下载数据库文件
./download_geolite2.sh <YOUR_LICENSE_KEY>

# 或者从 GitHub Releases 下载
curl -L -o GeoLite2-City.mmdb \
  https://github.com/lwpk110/free-geoip-databases/releases/latest/download/GeoLite2-City.mmdb

# 2. 运行测试
go run test_cities.go
```

### 测试 DB-IP 数据库

```bash
# 1. 下载 DB-IP 数据库
./download_dbip.sh city

# 2. 重命名文件为标准格式（如果需要）
mv dbip-city-lite-*.mmdb GeoLite2-City.mmdb

# 3. 运行测试
go run test_cities.go
```

## 自动化测试

### GitHub Actions 测试

项目包含自动化测试工作流 (`.github/workflows/test-database.yml`)：

- **触发条件**:
  - 每次创建新的 Release 后自动运行
  - 可以手动触发

- **测试内容**:
  - 下载最新的 GeoLite2-City.mmdb
  - 运行 IP 地理位置查询测试
  - 验证数据库完整性

- **手动触发**:
  1. 进入仓库的 Actions 页面
  2. 选择 "Test GeoIP Database" 工作流
  3. 点击 "Run workflow"

## 测试覆盖的 IP 地址

`test_cities.go` 测试以下知名 DNS 服务器的 IP：

- Google DNS (8.8.8.8)
- Cloudflare DNS (1.1.1.1)
- 114 DNS (114.114.114.114)
- 阿里云 DNS (223.5.5.5)
- 中国联通 DNS (202.106.0.20)
- OpenDNS (208.67.222.222)
- CleanBrowsing DNS (185.228.168.9)
- Yandex DNS (77.88.8.8)

## 添加自定义测试

编辑 `test_cities.go` 中的 `testIPs` 数组：

```go
var testIPs = []struct {
	ip          string
	description string
}{
	{"your.ip.address.here", "Your Description"},
	// ... 更多测试 IP
}
```

## 预期输出

成功的测试输出示例：

```
========================================
GeoIP 数据库测试
========================================

✓ 8.8.8.8 (Google DNS (美国))
  国家: 美国 (US)
  城市: N/A
  经纬度: 37.7510, -97.8220

...

========================================
测试完成: 成功 8/8
✓ 所有测试通过！
========================================
```

## 故障排除

### 数据库文件未找到

```
❌ 无法打开数据库: open GeoLite2-City.mmdb: no such file or directory
```

**解决方案**: 确保 `GeoLite2-City.mmdb` 文件在当前目录中。

### 某些 IP 查询失败

某些 IP 地址可能在免费数据库中没有详细信息，这是正常的。免费版数据库的覆盖范围有限。

### GitHub Actions 失败

检查以下几点：
1. 是否有有效的 Release（包含 `GeoLite2-City.mmdb` 文件）
2. `GEOIP_ACCESS_TOKEN` 是否正确配置
3. 网络连接是否正常
