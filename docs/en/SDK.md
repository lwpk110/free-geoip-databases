# SDK Quick Integration

## OpenAPI

- Spec file: `docs/openapi/geoip-api.yaml`

Use this file with your OpenAPI generator to build Go/Python/Node clients.

## Direct HTTP examples

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
