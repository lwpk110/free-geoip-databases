# 下载脚本

本目录包含用于下载 GeoIP 数据库的脚本。

## 脚本列表

### download_geolite2.sh

下载 MaxMind GeoLite2 数据库。

**使用方法**:
```bash
./scripts/download_geolite2.sh <YOUR_LICENSE_KEY>
```

**需要**:
- MaxMind License Key（从 https://www.maxmind.com/en/geolite2/signup 获取）

**下载内容**:
- GeoLite2-City.mmdb

### download_dbip.sh

下载 DB-IP 免费数据库。

**使用方法**:
```bash
# 下载所有数据库
./scripts/download_dbip.sh all

# 下载单个数据库
./scripts/download_dbip.sh city
./scripts/download_dbip.sh country
./scripts/download_dbip.sh asn

# 下载指定月份的数据库
./scripts/download_dbip.sh city 2024-10
```

**特点**:
- 无需注册或 License Key
- CC BY 4.0 授权，可商用（需署名）

## 注意事项

- 下载的数据库文件会保存到项目根目录
- 数据库文件已在 `.gitignore` 中排除，不会被提交到 Git
