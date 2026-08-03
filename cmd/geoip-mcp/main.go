package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

type rpcRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      any             `json:"id,omitempty"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

type rpcResponse struct {
	JSONRPC string    `json:"jsonrpc"`
	ID      any       `json:"id,omitempty"`
	Result  any       `json:"result,omitempty"`
	Error   *rpcError `json:"error,omitempty"`
}

type rpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type toolsCallParams struct {
	Name      string          `json:"name"`
	Arguments json.RawMessage `json:"arguments"`
}

type lookupArgs struct {
	IP     string `json:"ip"`
	APIKey string `json:"api_key,omitempty"`
}

type batchArgs struct {
	IPs    []string `json:"ips"`
	APIKey string   `json:"api_key,omitempty"`
}

func main() {
	apiURL := strings.TrimRight(getEnv("GEOIP_API_URL", "http://localhost:8080"), "/")
	log.SetOutput(io.Discard)

	decoder := json.NewDecoder(os.Stdin)
	encoder := json.NewEncoder(os.Stdout)

	for {
		var req rpcRequest
		if err := decoder.Decode(&req); err != nil {
			if err == io.EOF {
				return
			}
			_ = encoder.Encode(rpcResponse{JSONRPC: "2.0", Error: &rpcError{Code: -32700, Message: err.Error()}})
			continue
		}

		res := handle(apiURL, req)
		if err := encoder.Encode(res); err != nil {
			return
		}
	}
}

func handle(apiURL string, req rpcRequest) rpcResponse {
	resp := rpcResponse{JSONRPC: "2.0", ID: req.ID}
	switch req.Method {
	case "initialize":
		resp.Result = map[string]any{
			"protocolVersion": "2024-11-05",
			"serverInfo":      map[string]string{"name": "geoip-mcp", "version": "1.0.0"},
			"capabilities":    map[string]any{"tools": map[string]any{}},
		}
	case "tools/list":
		resp.Result = map[string]any{"tools": []map[string]any{
			{"name": "geoip_lookup", "description": "Lookup single IP", "inputSchema": map[string]any{"type": "object", "properties": map[string]any{"ip": map[string]string{"type": "string"}, "api_key": map[string]string{"type": "string"}}, "required": []string{"ip"}}},
			{"name": "geoip_batch_lookup", "description": "Batch lookup", "inputSchema": map[string]any{"type": "object", "properties": map[string]any{"ips": map[string]any{"type": "array", "items": map[string]string{"type": "string"}}, "api_key": map[string]string{"type": "string"}}, "required": []string{"ips"}}},
			{"name": "geoip_status", "description": "Get service status", "inputSchema": map[string]any{"type": "object", "properties": map[string]any{}}},
			{"name": "geoip_metrics", "description": "Get usage metrics", "inputSchema": map[string]any{"type": "object", "properties": map[string]any{}}},
			{"name": "geoip_hot_reload", "description": "Trigger database hot reload", "inputSchema": map[string]any{"type": "object", "properties": map[string]any{"api_key": map[string]string{"type": "string"}}}},
		}}
	case "tools/call":
		var p toolsCallParams
		if err := json.Unmarshal(req.Params, &p); err != nil {
			resp.Error = &rpcError{Code: -32602, Message: err.Error()}
			break
		}
		out, err := callTool(apiURL, p)
		if err != nil {
			resp.Error = &rpcError{Code: -32000, Message: err.Error()}
			break
		}
		resp.Result = map[string]any{"content": []map[string]string{{"type": "text", "text": out}}}
	default:
		resp.Error = &rpcError{Code: -32601, Message: fmt.Sprintf("method not found: %s", req.Method)}
	}
	return resp
}

func callTool(apiURL string, p toolsCallParams) (string, error) {
	switch p.Name {
	case "geoip_lookup":
		var a lookupArgs
		if err := json.Unmarshal(p.Arguments, &a); err != nil {
			return "", err
		}
		return requestJSON(http.MethodGet, fmt.Sprintf("%s/api/v1/lookup/%s", apiURL, a.IP), nil, a.APIKey)
	case "geoip_batch_lookup":
		var a batchArgs
		if err := json.Unmarshal(p.Arguments, &a); err != nil {
			return "", err
		}
		body, _ := json.Marshal(map[string]any{"ips": a.IPs})
		return requestJSON(http.MethodPost, apiURL+"/api/v1/lookup/batch", body, a.APIKey)
	case "geoip_status":
		return requestJSON(http.MethodGet, apiURL+"/api/v1/status", nil, "")
	case "geoip_metrics":
		return requestJSON(http.MethodGet, apiURL+"/api/v1/metrics", nil, "")
	case "geoip_hot_reload":
		var a lookupArgs
		_ = json.Unmarshal(p.Arguments, &a)
		return requestJSON(http.MethodPost, apiURL+"/api/v1/admin/reload", []byte(`{}`), a.APIKey)
	default:
		return "", fmt.Errorf("unknown tool: %s", p.Name)
	}
}

func requestJSON(method, url string, body []byte, apiKey string) (string, error) {
	var reader io.Reader
	if body != nil {
		reader = bytes.NewReader(body)
	}
	req, err := http.NewRequest(method, url, reader)
	if err != nil {
		return "", err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if apiKey != "" {
		req.Header.Set("X-API-Key", apiKey)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func getEnv(name, fallback string) string {
	if v := strings.TrimSpace(os.Getenv(name)); v != "" {
		return v
	}
	return fallback
}
