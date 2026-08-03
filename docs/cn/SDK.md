# SDK 快速接入

## OpenAPI

- 规范文件：`docs/openapi/geoip-api.yaml`

可使用 OpenAPI 生成器生成 Go/Python/Node SDK。

## 直接 HTTP 调用示例

### Go
```go
resp, _ := http.Get("http://localhost:8080/api/v1/lookup/8.8.8.8")
```

### Python
```python
requests.get("http://localhost:8080/api/v1/lookup/8.8.8.8")
```

### Node.js
```js
await fetch("http://localhost:8080/api/v1/lookup/8.8.8.8");
```
