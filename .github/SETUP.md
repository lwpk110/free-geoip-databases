# GitHub Actions 设置指南

本项目使用 GitHub Actions 自动下载和发布 MaxMind GeoLite2 数据库。

## 🔧 配置步骤

### 1. 获取 MaxMind License Key

1. 访问 [MaxMind 注册页面](https://www.maxmind.com/en/geolite2/signup) 注册账号
2. 验证邮箱后登录
3. 进入 [License Keys 管理页面](https://www.maxmind.com/en/accounts/current/license-key)
4. 点击 "Generate new license key"
5. 输入描述（例如：GitHub Actions）
6. 选择 "No" for GeoIP Update Program
7. 复制生成的 License Key（只显示一次，请妥善保存！）

### 2. 在 GitHub 仓库中添加 Secret

1. 进入你的 GitHub 仓库
2. 点击 **Settings** (设置)
3. 在左侧菜单中选择 **Secrets and variables** > **Actions**
4. 点击 **New repository secret**
5. 添加以下 Secret：
   - **Name**: `MAXMIND_LICENSE_KEY`
   - **Value**: 你从 MaxMind 获取的 License Key
6. 点击 **Add secret**

### 3. 启用 GitHub Actions

1. 进入仓库的 **Actions** 标签页
2. 如果看到提示，点击 "I understand my workflows, go ahead and enable them"
3. 你应该能看到两个工作流：
   - **Update GeoLite2 Database** - 自动下载和发布数据库
   - **Test GeoIP Database** - 测试数据库完整性

### 4. 启用 Release 权限

确保 GitHub Actions 有权限创建 Release：

1. 进入 **Settings** > **Actions** > **General**
2. 滚动到 **Workflow permissions** 部分
3. 选择 **Read and write permissions**
4. 勾选 **Allow GitHub Actions to create and approve pull requests**
5. 点击 **Save**

## 🚀 工作流说明

### Update GeoLite2 Database

**触发条件**：
- 定时任务：每周二和周五 UTC 时间 10:00（北京时间 18:00）
- 手动触发：在 Actions 页面点击 "Run workflow"

**功能**：
1. 下载 GeoLite2-City、GeoLite2-Country 和 GeoLite2-ASN 数据库
2. 创建新的 Release，标签格式为日期（例如：20250105）
3. 将三个数据库文件作为 Release 资产上传

### Test GeoIP Database

**触发条件**：
- 新 Release 发布后自动运行
- 手动触发：在 Actions 页面点击 "Run workflow"

**功能**：
1. 从最新 Release 下载数据库
2. 设置 Go 环境
3. 运行测试脚本验证数据库完整性

## 📅 更新频率

根据 MaxMind 的更新策略：
- GeoLite2 数据库通常每周二更新
- 本项目设置为每周二和周五检查更新
- 你可以随时手动触发工作流

## 🔍 手动触发工作流

1. 进入仓库的 **Actions** 标签页
2. 在左侧选择要运行的工作流
3. 点击右侧的 **Run workflow** 按钮
4. 选择分支（通常是 main）
5. 点击绿色的 **Run workflow** 按钮

## 📥 使用已发布的数据库

### 方法一：通过 Release 页面下载

1. 访问仓库的 [Releases 页面](../../releases)
2. 选择最新的 Release
3. 在 Assets 部分下载所需的 `.mmdb` 文件

### 方法二：通过命令行下载

```bash
# 下载 GeoLite2-City 数据库
curl -L -o GeoLite2-City.mmdb \
  https://github.com/YOUR_USERNAME/YOUR_REPO/releases/latest/download/GeoLite2-City.mmdb

# 下载 GeoLite2-Country 数据库
curl -L -o GeoLite2-Country.mmdb \
  https://github.com/YOUR_USERNAME/YOUR_REPO/releases/latest/download/GeoLite2-Country.mmdb

# 下载 GeoLite2-ASN 数据库
curl -L -o GeoLite2-ASN.mmdb \
  https://github.com/YOUR_USERNAME/YOUR_REPO/releases/latest/download/GeoLite2-ASN.mmdb
```

记得替换 `YOUR_USERNAME` 和 `YOUR_REPO` 为你的实际仓库信息。

### 方法三：使用 GitHub API

```bash
# 获取最新 Release 的下载链接
curl -s https://api.github.com/repos/YOUR_USERNAME/YOUR_REPO/releases/latest | \
  jq -r '.assets[] | select(.name == "GeoLite2-City.mmdb") | .browser_download_url'
```

## 🛠️ 故障排查

### 工作流失败

1. **License Key 错误**
   - 检查 Secret 中的 `MAXMIND_LICENSE_KEY` 是否正确
   - 确认 License Key 没有过期
   - 访问 MaxMind 账户验证 License Key 状态

2. **权限错误**
   - 确认已启用 "Read and write permissions"
   - 检查 GITHUB_TOKEN 权限

3. **下载失败**
   - MaxMind 服务器可能暂时不可用，稍后重试
   - 检查网络连接

### 查看日志

1. 进入 **Actions** 标签页
2. 点击失败的工作流运行
3. 点击失败的步骤查看详细日志

## 📝 注意事项

1. **遵守许可协议**：使用 GeoLite2 数据库需遵守 [MaxMind EULA](https://www.maxmind.com/en/geolite2/eula)
2. **私有仓库**：如果仓库是私有的，确保有足够的 Actions 配额
3. **存储限制**：定期清理旧的 Release 以节省存储空间
4. **数据准确性**：GeoLite2 是免费版本，精确度可能不如商业版本

## 🔗 相关链接

- [MaxMind 官方网站](https://www.maxmind.com)
- [GeoLite2 数据库文档](https://dev.maxmind.com/geoip/geolite2-free-geolocation-data)
- [GitHub Actions 文档](https://docs.github.com/en/actions)
- [softprops/action-gh-release](https://github.com/softprops/action-gh-release)

## 💡 自定义配置

### 修改更新频率

编辑 `.github/workflows/update-geolite2.yml`：

```yaml
schedule:
  # 自定义 cron 表达式
  - cron: '0 10 * * 2,5'  # 每周二、五 UTC 10:00
```

### 只下载特定数据库

如果只需要某个数据库，可以删除工作流中对应的下载步骤。

### 添加通知

可以在工作流中添加通知步骤，例如发送邮件或 Slack 消息。
