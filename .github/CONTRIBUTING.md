# 贡献指南

感谢你对本项目感兴趣！本文档将帮助你了解如何为项目做出贡献。

## 项目结构

```
free-geoip-databases/
├── .github/
│   ├── agents/          # GitHub Copilot 智能助手配置
│   └── workflows/       # GitHub Actions 自动化
├── scripts/              # 下载脚本
├── examples/             # 示例代码
│   ├── query/           # 查询示例
│   └── test/            # 测试程序
├── docs/                # 文档
└── [数据库文件]          # *.mmdb (gitignored)
```

## 如何贡献

### 报告问题

如果你发现了 bug 或有功能建议：

1. 检查 [Issues](../../issues) 页面，避免重复
2. 创建新 Issue，提供详细信息：
   - Bug: 重现步骤、预期行为、实际行为
   - 功能建议: 用例说明、预期效果

### 提交代码

1. **Fork 本仓库**

2. **创建分支**
   ```bash
   git checkout -b feature/your-feature-name
   # 或
   git checkout -b fix/your-bug-fix
   ```

3. **进行更改**
   - 遵循现有代码风格
   - 添加必要的注释
   - 更新相关文档

4. **测试更改**
   ```bash
   # 测试下载脚本
   ./scripts/download_dbip.sh city

   # 测试示例代码
   cd examples/test
   go run test_cities.go
   ```

5. **提交更改**
   ```bash
   git add .
   git commit -m "type: 简短描述

   详细说明（可选）"
   ```

   提交类型：
   - `feat`: 新功能
   - `fix`: Bug 修复
   - `docs`: 文档更新
   - `style`: 代码格式调整
   - `refactor`: 代码重构
   - `test`: 测试相关
   - `chore`: 构建/工具更改

6. **推送并创建 Pull Request**
   ```bash
   git push origin feature/your-feature-name
   ```

## 开发指南

### 添加新的数据库源

如果要添加新的 GeoIP 数据库源：

1. 在 `scripts/` 创建下载脚本
2. 在 `.github/workflows/` 创建自动更新工作流
3. 更新 `README.md` 添加说明
4. 更新 `CHANGELOG.md` 记录变更

### 添加新的示例

在 `examples/` 目录创建新的子目录：

```bash
mkdir -p examples/your-example
cd examples/your-example
# 创建示例代码
```

记得在 `examples/README.md` 中添加说明。

### 使用智能助手

本项目配置了 GitHub Copilot 智能助手，可以帮助你更好地参与项目贡献：

- **资深产品经理助手** (`.github/agents/senior-product-manager.md`)
  - 帮助进行产品功能规划和迭代
  - 提供用户体验优化建议
  - 进行竞品分析和功能借鉴
  - 提出创新性解决方案

如果你有产品相关的想法或需求，可以咨询智能助手获取建议。

### 更新文档

- 用户文档: `docs/`
- API 文档: 在代码中添加注释
- README: 保持简洁，详细内容放在 `docs/`

## 代码规范

### Go 代码

```go
// 使用 gofmt 格式化
go fmt ./...

// 使用有意义的变量名
var ipAddress string // Good
var ip string         // OK
var i string          // Bad

// 添加注释说明函数用途
// getPublicIP 获取当前主机的公网 IP 地址
func getPublicIP() (string, error) {
    // ...
}
```

### Shell 脚本

```bash
#!/bin/bash

# 使用有意义的变量名
LICENSE_KEY="$1"  # Good
KEY="$1"          # Less clear

# 添加错误处理
if [ -z "$LICENSE_KEY" ]; then
    echo "Error: License key required"
    exit 1
fi

# 使用引号包裹变量
curl -o "$OUTPUT_FILE" "$URL"
```

### Markdown 文档

- 使用清晰的标题层级
- 代码块指定语言
- 添加适当的链接
- 保持行宽合理（建议 80-120 字符）

## 测试

### 本地测试

运行测试确保更改不会破坏现有功能：

```bash
# 测试数据库查询
cd examples/test
go run test_cities.go

# 测试下载脚本
./scripts/download_dbip.sh city
```

### GitHub Actions

推送到仓库后，检查 Actions 标签确保所有工作流通过。

## 发布流程

（仅限维护者）

1. 更新 `CHANGELOG.md`
2. 创建 Git tag: `git tag -a vX.Y.Z -m "Release vX.Y.Z"`
3. 推送 tag: `git push origin vX.Y.Z`
4. GitHub Actions 会自动创建 Release

## 问题？

如有任何问题，欢迎：

- 创建 [Issue](../../issues)
- 在 [Discussions](../../discussions) 提问
- 查看现有的 [Pull Requests](../../pulls)

感谢你的贡献！🎉
