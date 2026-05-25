# AI MCP 集成

仓库新增 MCP 服务端能力，可直接被 AI Agent 调用：

- 二进制：`geoip-mcp`（`cmd/geoip-mcp`）
- 通信：stdio JSON-RPC（MCP 工具调用模式）
- 默认 API 地址：`http://localhost:8080`（可通过 `GEOIP_API_URL` 覆盖）

## 提供的 MCP 工具

- `geoip_lookup`：单个 IP 查询
- `geoip_batch_lookup`：批量查询
- `geoip_status`：服务状态
- `geoip_metrics`：观测指标
- `geoip_hot_reload`：触发数据库热更新

## 快速开始

```bash
# 终端1
make run

# 终端2
make mcp
```

然后在你的 AI Agent MCP 客户端中接入该 stdio 服务。
