package main

import (
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/oschwald/geoip2-golang"
)

//go:embed static/*
var staticFiles embed.FS

// Config represents server configuration
// kept for compatibility with tests and docs.
type Config struct {
	Port         string   `json:"port"`
	DatabasePath string   `json:"database_path"`
	UpdateCron   string   `json:"update_cron"`
	AllowedDBs   []string `json:"allowed_dbs"`
}

type DatabaseSourceStatus struct {
	Path      string    `json:"path"`
	Type      string    `json:"type"`
	LastLoad  time.Time `json:"last_load"`
	LastMod   time.Time `json:"last_mod"`
	Available bool      `json:"available"`
}

type databaseSource struct {
	path      string
	dbType    string
	reader    *geoip2.Reader
	lastLoad  time.Time
	lastMod   time.Time
	available bool
}

type DatabasePool struct {
	mu      sync.RWMutex
	sources []*databaseSource
}

type GeoIPResponse struct {
	IP           string            `json:"ip"`
	SourceDB     string            `json:"source_db"`
	Country      CountryInfo       `json:"country"`
	City         CityInfo          `json:"city,omitempty"`
	Location     LocationInfo      `json:"location"`
	Continent    ContinentInfo     `json:"continent"`
	Subdivisions []SubdivisionInfo `json:"subdivisions,omitempty"`
	Postal       PostalInfo        `json:"postal,omitempty"`
	Traits       TraitsInfo        `json:"traits"`
}

type CountryInfo struct {
	IsoCode string            `json:"iso_code"`
	Names   map[string]string `json:"names"`
}

type CityInfo struct {
	GeoNameID uint              `json:"geoname_id,omitempty"`
	Names     map[string]string `json:"names"`
}

type LocationInfo struct {
	Latitude       float64 `json:"latitude"`
	Longitude      float64 `json:"longitude"`
	TimeZone       string  `json:"time_zone,omitempty"`
	AccuracyRadius uint16  `json:"accuracy_radius,omitempty"`
}

type ContinentInfo struct {
	Code  string            `json:"code"`
	Names map[string]string `json:"names"`
}

type SubdivisionInfo struct {
	IsoCode string            `json:"iso_code,omitempty"`
	Names   map[string]string `json:"names"`
}

type PostalInfo struct {
	Code string `json:"code,omitempty"`
}

type TraitsInfo struct {
	IsAnonymousProxy    bool `json:"is_anonymous_proxy"`
	IsSatelliteProvider bool `json:"is_satellite_provider"`
}

type StatusResponse struct {
	Status       string                 `json:"status"`
	DatabasePath string                 `json:"database_path,omitempty"`
	DatabaseType string                 `json:"database_type,omitempty"`
	Databases    []DatabaseSourceStatus `json:"databases"`
	Uptime       string                 `json:"uptime"`
	HotReload    bool                   `json:"hot_reload"`
	FallbackMode bool                   `json:"fallback_mode"`
}

type BatchLookupRequest struct {
	IPs []string `json:"ips"`
}

type BatchLookupResponse struct {
	Results map[string]*GeoIPResponse `json:"results"`
	Errors  map[string]string         `json:"errors"`
	Summary struct {
		Total   int `json:"total"`
		Success int `json:"success"`
		Failed  int `json:"failed"`
	} `json:"summary"`
}

type cacheEntry struct {
	data      *GeoIPResponse
	expiresAt time.Time
}

type ResponseCache struct {
	mu      sync.RWMutex
	entries map[string]cacheEntry
	ttl     time.Duration
}

type RateLimiter struct {
	mu           sync.Mutex
	requests     map[string]*rateRecord
	maxPerMinute int
}

type rateRecord struct {
	windowStart time.Time
	count       int
}

type UsageStats struct {
	TotalLookups  atomic.Uint64
	BatchRequests atomic.Uint64
	CacheHits     atomic.Uint64
	CacheMisses   atomic.Uint64
	LookupErrors  atomic.Uint64
	RateLimited   atomic.Uint64
	AuthRejected  atomic.Uint64
	HotReloads    atomic.Uint64
	PerDBLookups  sync.Map // map[string]*atomic.Uint64
}

type MetricsResponse struct {
	TotalLookups  uint64            `json:"total_lookups"`
	BatchRequests uint64            `json:"batch_requests"`
	CacheHits     uint64            `json:"cache_hits"`
	CacheMisses   uint64            `json:"cache_misses"`
	CacheHitRate  string            `json:"cache_hit_rate"`
	LookupErrors  uint64            `json:"lookup_errors"`
	RateLimited   uint64            `json:"rate_limited"`
	AuthRejected  uint64            `json:"auth_rejected"`
	HotReloads    uint64            `json:"hot_reloads"`
	PerDBLookups  map[string]uint64 `json:"per_db_lookups"`
	Uptime        string            `json:"uptime"`
}

var (
	dbPool         *DatabasePool
	responseCache  *ResponseCache
	rateLimiter    *RateLimiter
	usageStats     *UsageStats
	serverStart    time.Time
	apiKeys        map[string]struct{}
	apiAuthEnabled bool
	hotReloadEvery time.Duration
)

func init() {
	serverStart = time.Now()
}

func parseDatabasePaths() []string {
	pathsVar := strings.TrimSpace(os.Getenv("GEOIP_DB_PATHS"))
	if pathsVar == "" {
		if one := strings.TrimSpace(os.Getenv("GEOIP_DB_PATH")); one != "" {
			return []string{one}
		}
		return []string{"GeoLite2-City.mmdb"}
	}

	seen := map[string]struct{}{}
	var paths []string
	for _, part := range strings.Split(pathsVar, ",") {
		p := strings.TrimSpace(part)
		if p == "" {
			continue
		}
		if _, ok := seen[p]; ok {
			continue
		}
		seen[p] = struct{}{}
		paths = append(paths, p)
	}
	if len(paths) == 0 {
		paths = []string{"GeoLite2-City.mmdb"}
	}
	return paths
}

func detectDBType(path string) string {
	name := strings.ToLower(filepath.Base(path))
	switch {
	case strings.Contains(name, "dbip"):
		return "DB-IP"
	case strings.Contains(name, "geolite"):
		return "GeoLite2"
	default:
		return "Unknown"
	}
}

func NewDatabasePool(paths []string) *DatabasePool {
	sources := make([]*databaseSource, 0, len(paths))
	for _, p := range paths {
		sources = append(sources, &databaseSource{path: p, dbType: detectDBType(p)})
	}
	return &DatabasePool{sources: sources}
}

func (p *DatabasePool) LoadAll() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	var available int
	for _, src := range p.sources {
		if err := p.loadOneLocked(src); err != nil {
			log.Printf("failed to load database %s: %v", src.path, err)
			continue
		}
		available++
	}

	if available == 0 {
		return fmt.Errorf("no database could be loaded")
	}
	return nil
}

func (p *DatabasePool) loadOneLocked(src *databaseSource) error {
	fi, err := os.Stat(src.path)
	if err != nil {
		src.available = false
		return err
	}

	reader, err := geoip2.Open(src.path)
	if err != nil {
		src.available = false
		return err
	}

	if src.reader != nil {
		src.reader.Close()
	}

	src.reader = reader
	src.lastLoad = time.Now()
	src.lastMod = fi.ModTime()
	src.available = true
	log.Printf("database loaded: %s (type=%s)", src.path, src.dbType)
	return nil
}

func (p *DatabasePool) ReloadChanged(force bool) (int, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	reloaded := 0
	for _, src := range p.sources {
		fi, err := os.Stat(src.path)
		if err != nil {
			log.Printf("hot reload stat failed %s: %v", src.path, err)
			src.available = false
			continue
		}

		if !force && src.available && !fi.ModTime().After(src.lastMod) {
			continue
		}

		if err := p.loadOneLocked(src); err != nil {
			log.Printf("hot reload failed %s: %v", src.path, err)
			continue
		}
		reloaded++
	}

	if reloaded == 0 && force {
		return 0, fmt.Errorf("no database reloaded")
	}
	return reloaded, nil
}

func (p *DatabasePool) Lookup(ipStr string) (*GeoIPResponse, error) {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return nil, fmt.Errorf("invalid IP address: %s", ipStr)
	}

	p.mu.RLock()
	defer p.mu.RUnlock()

	var lastErr error
	for _, src := range p.sources {
		if !src.available || src.reader == nil {
			continue
		}
		record, err := src.reader.City(ip)
		if err != nil {
			lastErr = err
			continue
		}
		resp := mapRecordToResponse(ipStr, src.dbType, record)
		incrementPerDB(src.dbType)
		return resp, nil
	}

	if lastErr != nil {
		return nil, fmt.Errorf("lookup failed on all databases: %w", lastErr)
	}
	return nil, errors.New("no available database")
}

func (p *DatabasePool) Statuses() []DatabaseSourceStatus {
	p.mu.RLock()
	defer p.mu.RUnlock()

	statuses := make([]DatabaseSourceStatus, 0, len(p.sources))
	for _, src := range p.sources {
		statuses = append(statuses, DatabaseSourceStatus{
			Path:      src.path,
			Type:      src.dbType,
			LastLoad:  src.lastLoad,
			LastMod:   src.lastMod,
			Available: src.available,
		})
	}
	return statuses
}

func mapRecordToResponse(ipStr, source string, record *geoip2.City) *GeoIPResponse {
	response := &GeoIPResponse{
		IP:       ipStr,
		SourceDB: source,
		Country: CountryInfo{
			IsoCode: record.Country.IsoCode,
			Names:   record.Country.Names,
		},
		City: CityInfo{
			GeoNameID: record.City.GeoNameID,
			Names:     record.City.Names,
		},
		Location: LocationInfo{
			Latitude:       record.Location.Latitude,
			Longitude:      record.Location.Longitude,
			TimeZone:       record.Location.TimeZone,
			AccuracyRadius: record.Location.AccuracyRadius,
		},
		Continent: ContinentInfo{
			Code:  record.Continent.Code,
			Names: record.Continent.Names,
		},
		Postal: PostalInfo{Code: record.Postal.Code},
		Traits: TraitsInfo{
			IsAnonymousProxy:    record.Traits.IsAnonymousProxy,
			IsSatelliteProvider: record.Traits.IsSatelliteProvider,
		},
	}

	for _, sub := range record.Subdivisions {
		response.Subdivisions = append(response.Subdivisions, SubdivisionInfo{IsoCode: sub.IsoCode, Names: sub.Names})
	}
	return response
}

func NewResponseCache(ttl time.Duration) *ResponseCache {
	return &ResponseCache{entries: make(map[string]cacheEntry), ttl: ttl}
}

func (c *ResponseCache) Get(ip string) (*GeoIPResponse, bool) {
	if c == nil || c.ttl <= 0 {
		return nil, false
	}
	c.mu.RLock()
	entry, ok := c.entries[ip]
	c.mu.RUnlock()
	if !ok || time.Now().After(entry.expiresAt) {
		if ok {
			c.mu.Lock()
			delete(c.entries, ip)
			c.mu.Unlock()
		}
		return nil, false
	}
	return entry.data, true
}

func (c *ResponseCache) Set(ip string, data *GeoIPResponse) {
	if c == nil || c.ttl <= 0 {
		return
	}
	c.mu.Lock()
	c.entries[ip] = cacheEntry{data: data, expiresAt: time.Now().Add(c.ttl)}
	c.mu.Unlock()
}

func (c *ResponseCache) Size() int {
	if c == nil {
		return 0
	}
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.entries)
}

func NewRateLimiter(maxPerMinute int) *RateLimiter {
	if maxPerMinute <= 0 {
		maxPerMinute = 120
	}
	return &RateLimiter{requests: make(map[string]*rateRecord), maxPerMinute: maxPerMinute}
}

func (rl *RateLimiter) Allow(key string) bool {
	now := time.Now()
	rl.mu.Lock()
	defer rl.mu.Unlock()

	rec, ok := rl.requests[key]
	if !ok {
		rl.requests[key] = &rateRecord{windowStart: now, count: 1}
		return true
	}
	if now.Sub(rec.windowStart) >= time.Minute {
		rec.windowStart = now
		rec.count = 1
		return true
	}
	if rec.count >= rl.maxPerMinute {
		return false
	}
	rec.count++
	return true
}

func parseAPIKeys() map[string]struct{} {
	result := map[string]struct{}{}
	for _, k := range strings.Split(os.Getenv("API_KEYS"), ",") {
		k = strings.TrimSpace(k)
		if k != "" {
			result[k] = struct{}{}
		}
	}
	return result
}

func getClientKey(c *gin.Context) string {
	apiKey := strings.TrimSpace(c.GetHeader("X-API-Key"))
	if apiKey == "" {
		apiKey = strings.TrimSpace(c.Query("api_key"))
	}
	if apiKey != "" {
		return "apikey:" + apiKey
	}
	ip := c.ClientIP()
	if ip == "" {
		ip = "unknown"
	}
	return "ip:" + ip
}

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !apiAuthEnabled {
			c.Next()
			return
		}

		if strings.HasSuffix(c.Request.URL.Path, "/health") || strings.HasSuffix(c.Request.URL.Path, "/status") {
			c.Next()
			return
		}

		key := strings.TrimSpace(c.GetHeader("X-API-Key"))
		if key == "" {
			key = strings.TrimSpace(c.Query("api_key"))
		}
		if _, ok := apiKeys[key]; !ok {
			usageStats.AuthRejected.Add(1)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing or invalid API key"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func rateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if rateLimiter == nil {
			c.Next()
			return
		}
		key := getClientKey(c)
		if !rateLimiter.Allow(key) {
			usageStats.RateLimited.Add(1)
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "rate limit exceeded"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func cachedLookup(ip string) (*GeoIPResponse, error) {
	usageStats.TotalLookups.Add(1)
	if resp, ok := responseCache.Get(ip); ok {
		usageStats.CacheHits.Add(1)
		return resp, nil
	}
	usageStats.CacheMisses.Add(1)
	resp, err := dbPool.Lookup(ip)
	if err != nil {
		usageStats.LookupErrors.Add(1)
		return nil, err
	}
	responseCache.Set(ip, resp)
	return resp, nil
}

func lookupHandler(c *gin.Context) {
	ipParam := strings.TrimSpace(c.Param("ip"))
	if ipParam == "" {
		ipParam = strings.TrimSpace(c.Query("ip"))
	}
	if ipParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "IP parameter required"})
		return
	}

	result, err := cachedLookup(ipParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func batchLookupHandler(c *gin.Context) {
	var req BatchLookupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	if len(req.IPs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ips must not be empty"})
		return
	}
	if len(req.IPs) > 1000 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ips max length is 1000"})
		return
	}

	usageStats.BatchRequests.Add(1)
	resp := BatchLookupResponse{
		Results: make(map[string]*GeoIPResponse),
		Errors:  make(map[string]string),
	}
	for _, ip := range req.IPs {
		ip = strings.TrimSpace(ip)
		if ip == "" {
			continue
		}
		result, err := cachedLookup(ip)
		if err != nil {
			resp.Errors[ip] = err.Error()
			continue
		}
		resp.Results[ip] = result
	}
	resp.Summary.Total = len(req.IPs)
	resp.Summary.Success = len(resp.Results)
	resp.Summary.Failed = len(resp.Errors)
	c.JSON(http.StatusOK, resp)
}

func statusHandler(c *gin.Context) {
	statuses := dbPool.Statuses()
	healthy := false
	primaryPath := ""
	primaryType := ""
	for _, s := range statuses {
		if s.Available {
			healthy = true
			if primaryPath == "" {
				primaryPath = s.Path
				primaryType = s.Type
			}
			break
		}
	}
	state := "degraded"
	if healthy {
		state = "healthy"
	}
	c.JSON(http.StatusOK, StatusResponse{
		Status:       state,
		DatabasePath: primaryPath,
		DatabaseType: primaryType,
		Databases:    statuses,
		Uptime:       time.Since(serverStart).Round(time.Second).String(),
		HotReload:    hotReloadEvery > 0,
		FallbackMode: true,
	})
}

func healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func metricsHandler(c *gin.Context) {
	lookups := usageStats.TotalLookups.Load()
	hits := usageStats.CacheHits.Load()
	hitRate := "0%"
	if lookups > 0 {
		hitRate = fmt.Sprintf("%.2f%%", float64(hits)*100/float64(lookups))
	}

	perDB := map[string]uint64{}
	usageStats.PerDBLookups.Range(func(key, value any) bool {
		k := key.(string)
		v := value.(*atomic.Uint64)
		perDB[k] = v.Load()
		return true
	})
	if len(perDB) == 0 {
		for _, st := range dbPool.Statuses() {
			if _, ok := perDB[st.Type]; !ok {
				perDB[st.Type] = 0
			}
		}
	}

	c.JSON(http.StatusOK, MetricsResponse{
		TotalLookups:  lookups,
		BatchRequests: usageStats.BatchRequests.Load(),
		CacheHits:     hits,
		CacheMisses:   usageStats.CacheMisses.Load(),
		CacheHitRate:  hitRate,
		LookupErrors:  usageStats.LookupErrors.Load(),
		RateLimited:   usageStats.RateLimited.Load(),
		AuthRejected:  usageStats.AuthRejected.Load(),
		HotReloads:    usageStats.HotReloads.Load(),
		PerDBLookups:  perDB,
		Uptime:        time.Since(serverStart).Round(time.Second).String(),
	})
}

func reloadHandler(c *gin.Context) {
	reloaded, err := dbPool.ReloadChanged(true)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	usageStats.HotReloads.Add(uint64(reloaded))
	c.JSON(http.StatusOK, gin.H{"reloaded": reloaded})
}

func incrementPerDB(dbType string) {
	v, _ := usageStats.PerDBLookups.LoadOrStore(dbType, &atomic.Uint64{})
	v.(*atomic.Uint64).Add(1)
}

func startHotReloadLoop() {
	if hotReloadEvery <= 0 {
		return
	}
	go func() {
		ticker := time.NewTicker(hotReloadEvery)
		defer ticker.Stop()
		for range ticker.C {
			reloaded, err := dbPool.ReloadChanged(false)
			if err != nil {
				log.Printf("hot reload loop error: %v", err)
				continue
			}
			if reloaded > 0 {
				usageStats.HotReloads.Add(uint64(reloaded))
				log.Printf("hot reloaded %d database(s)", reloaded)
			}
		}
	}()
}

func parseDurationSeconds(envName string, def int) time.Duration {
	v := strings.TrimSpace(os.Getenv(envName))
	if v == "" {
		return time.Duration(def) * time.Second
	}
	n, err := time.ParseDuration(v)
	if err == nil {
		return n
	}
	secs, err := time.ParseDuration(v + "s")
	if err == nil {
		return secs
	}
	return time.Duration(def) * time.Second
}

func parsePositiveInt(envName string, def int) int {
	v := strings.TrimSpace(os.Getenv(envName))
	if v == "" {
		return def
	}
	var n int
	if _, err := fmt.Sscanf(v, "%d", &n); err != nil || n <= 0 {
		return def
	}
	return n
}

func main() {
	paths := parseDatabasePaths()
	sort.Strings(paths)
	log.Printf("configured databases: %s", strings.Join(paths, ", "))

	dbPool = NewDatabasePool(paths)
	if err := dbPool.LoadAll(); err != nil {
		log.Fatalf("failed to initialize database pool: %v", err)
	}

	hotReloadEvery = parseDurationSeconds("HOT_RELOAD_INTERVAL", 60)
	cacheTTL := parseDurationSeconds("CACHE_TTL", 300)
	maxPerMinute := parsePositiveInt("RATE_LIMIT_PER_MIN", 240)

	responseCache = NewResponseCache(cacheTTL)
	rateLimiter = NewRateLimiter(maxPerMinute)
	usageStats = &UsageStats{}

	apiKeys = parseAPIKeys()
	apiAuthEnabled = len(apiKeys) > 0
	if apiAuthEnabled {
		log.Printf("API key auth enabled with %d keys", len(apiKeys))
	}

	startHotReloadLoop()

	port := strings.TrimSpace(os.Getenv("PORT"))
	if port == "" {
		port = "8080"
	}
	if os.Getenv("GIN_MODE") == "" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()
	api := router.Group("/api/v1")
	api.Use(authMiddleware(), rateLimitMiddleware())
	{
		api.GET("/lookup/:ip", lookupHandler)
		api.GET("/lookup", lookupHandler)
		api.POST("/lookup/batch", batchLookupHandler)
		api.GET("/status", statusHandler)
		api.GET("/health", healthHandler)
		api.GET("/metrics", metricsHandler)
		api.POST("/admin/reload", reloadHandler)
	}

	staticFS, err := fs.Sub(staticFiles, "static")
	if err != nil {
		log.Fatalf("failed to load static files: %v", err)
	}
	router.StaticFS("/static", http.FS(staticFS))
	router.GET("/", func(c *gin.Context) {
		data, err := staticFiles.ReadFile("static/index.html")
		if err != nil {
			c.String(http.StatusInternalServerError, "Failed to load index page")
			return
		}
		c.Data(http.StatusOK, "text/html; charset=utf-8", data)
	})

	log.Printf("GeoIP API Server starting on :%s", port)
	log.Printf("Web UI: http://localhost:%s", port)
	log.Printf("Lookup API: http://localhost:%s/api/v1/lookup/<ip>", port)
	log.Printf("Batch API: POST http://localhost:%s/api/v1/lookup/batch", port)
	log.Printf("MCP-friendly status API: http://localhost:%s/api/v1/metrics", port)

	if err := router.Run(":" + port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
