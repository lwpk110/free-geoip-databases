package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/oschwald/geoip2-golang"
)

type batchRequest struct {
	IPs []string `json:"ips"`
}

func main() {
	apiURL := flag.String("api", "", "GeoIP API base URL, e.g. http://localhost:8080")
	dbPath := flag.String("db", "", "Local mmdb path for offline lookup")
	ip := flag.String("ip", "", "Single IP lookup")
	batchFile := flag.String("batch-file", "", "File containing IPs (one per line)")
	apiKey := flag.String("api-key", "", "API key for protected endpoints")
	flag.Parse()

	if *ip == "" && *batchFile == "" {
		log.Fatal("provide -ip or -batch-file")
	}

	if *dbPath != "" {
		runLocal(*dbPath, *ip, *batchFile)
		return
	}

	if *apiURL == "" {
		*apiURL = "http://localhost:8080"
	}
	runRemote(strings.TrimRight(*apiURL, "/"), *ip, *batchFile, *apiKey)
}

func runLocal(dbPath, ip, batchFile string) {
	db, err := geoip2.Open(dbPath)
	if err != nil {
		log.Fatalf("open db failed: %v", err)
	}
	defer db.Close()

	ips, err := collectIPs(ip, batchFile)
	if err != nil {
		log.Fatal(err)
	}

	results := map[string]any{}
	errorsMap := map[string]string{}
	for _, rawIP := range ips {
		parsed := net.ParseIP(rawIP)
		if parsed == nil {
			errorsMap[rawIP] = "invalid ip"
			continue
		}
		record, err := db.City(parsed)
		if err != nil {
			errorsMap[rawIP] = err.Error()
			continue
		}
		results[rawIP] = map[string]any{
			"country":   record.Country.Names,
			"city":      record.City.Names,
			"location":  record.Location,
			"continent": record.Continent.Names,
		}
	}

	out := map[string]any{"mode": "local", "db": dbPath, "results": results, "errors": errorsMap}
	prettyPrint(out)
}

func runRemote(apiURL, ip, batchFile, apiKey string) {
	ips, err := collectIPs(ip, batchFile)
	if err != nil {
		log.Fatal(err)
	}

	client := &http.Client{}
	if len(ips) == 1 {
		url := fmt.Sprintf("%s/api/v1/lookup/%s", apiURL, ips[0])
		req, _ := http.NewRequest(http.MethodGet, url, nil)
		if apiKey != "" {
			req.Header.Set("X-API-Key", apiKey)
		}
		resp, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		var out any
		_ = json.Unmarshal(body, &out)
		prettyPrint(out)
		return
	}

	payload, _ := json.Marshal(batchRequest{IPs: ips})
	req, _ := http.NewRequest(http.MethodPost, apiURL+"/api/v1/lookup/batch", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	if apiKey != "" {
		req.Header.Set("X-API-Key", apiKey)
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var out any
	_ = json.Unmarshal(body, &out)
	prettyPrint(out)
}

func collectIPs(singleIP, batchFile string) ([]string, error) {
	if batchFile == "" {
		return []string{singleIP}, nil
	}
	b, err := os.ReadFile(batchFile)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(b), "\n")
	ips := make([]string, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			ips = append(ips, line)
		}
	}
	if singleIP != "" {
		ips = append([]string{singleIP}, ips...)
	}
	if len(ips) == 0 {
		return nil, fmt.Errorf("no ips found")
	}
	return ips, nil
}

func prettyPrint(v any) {
	b, _ := json.MarshalIndent(v, "", "  ")
	fmt.Println(string(b))
}
