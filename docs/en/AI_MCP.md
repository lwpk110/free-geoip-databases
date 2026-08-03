# AI MCP Integration

This repository now ships an MCP server binary for AI agents:

- Binary: `geoip-mcp` (`cmd/geoip-mcp`)
- Transport: stdio JSON-RPC compatible with MCP tool patterns
- Default API target: `http://localhost:8080` (override with `GEOIP_API_URL`)

## Exposed MCP tools

- `geoip_lookup` - single IP lookup
- `geoip_batch_lookup` - batch lookup
- `geoip_status` - service status
- `geoip_metrics` - observability metrics
- `geoip_hot_reload` - trigger DB hot reload

## Quick start

```bash
# Terminal 1
make run

# Terminal 2
make mcp
```

Then connect your AI agent MCP client to this stdio server.
