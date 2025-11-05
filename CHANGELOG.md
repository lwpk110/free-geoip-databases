# Changelog

All notable changes to this project will be documented in this file.

## [Unreleased]

### Changed - 2025-11-05

**项目结构重组** - 改进文件组织，使项目更清晰易用

- 📁 **创建 `scripts/` 目录**: 移动所有下载脚本
  - `download_geolite2.sh` → `scripts/download_geolite2.sh`
  - `download_dbip.sh` → `scripts/download_dbip.sh`

- 📁 **创建 `examples/` 目录**: 组织示例代码
  - `main.go` → `examples/query/main.go` (查询示例)
  - `test_cities.go` → `examples/test/test_cities.go` (测试程序)

- 📁 **创建 `docs/` 目录**: 集中存放文档
  - `QUICKSTART.md` → `docs/QUICKSTART.md`
  - `TESTING.md` → `docs/TESTING.md`

- 📝 **更新所有路径引用**:
  - 示例代码中的数据库路径更新为相对路径 `../../GeoLite2-City.mmdb`
  - GitHub Actions 工作流更新以适配新的文件结构
  - README.md 添加了项目结构说明

- 📚 **新增 README 文件**:
  - `scripts/README.md` - 下载脚本使用说明
  - `examples/README.md` - 示例代码说明
  - `docs/README.md` - 文档索引

### Added

- ✅ DB-IP 自动更新工作流 (`.github/workflows/update-dbip.yml`)
- ✅ DB-IP 下载脚本 (`scripts/download_dbip.sh`)
- ✅ 测试程序 (`examples/test/test_cities.go`)
- ✅ 测试文档 (`docs/TESTING.md`)

### Fixed

- 🐛 修复 GitHub Actions 测试工作流中 `test_cities.go: no such file or directory` 错误

## [Initial Release] - 2024

### Added

- 🚀 MaxMind GeoLite2 自动更新工作流
- 📥 GeoLite2 下载脚本
- 📝 基础文档和使用示例
- ⚙️ Go 语言查询示例
