package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/oschwald/geoip2-golang"
)

//go:embed static/*
var staticFiles embed.FS

// Config represents the server configuration
type Config struct {
	Port         string   `json:"port"`
	DatabasePath string   `json:"database_path"`
	UpdateCron   string   `json:"update_cron"`
	AllowedDBs   []string `json:"allowed_dbs"`
}

// DatabaseManager manages GeoIP database
type DatabaseManager struct {
	mu       sync.RWMutex
	db       *geoip2.Reader
	dbPath   string
	dbType   string
	lastLoad time.Time
}

// GeoIPResponse represents the API response
type GeoIPResponse struct {
	IP        string            `json:"ip"`
	Country   CountryInfo       `json:"country"`
	City      CityInfo          `json:"city,omitempty"`
	Location  LocationInfo      `json:"location"`
	Continent ContinentInfo     `json:"continent"`
	Subdivisions []SubdivisionInfo `json:"subdivisions,omitempty"`
	Postal    PostalInfo        `json:"postal,omitempty"`
	Traits    TraitsInfo        `json:"traits"`
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

// StatusResponse represents database status
type StatusResponse struct {
	Status       string    `json:"status"`
	DatabasePath string    `json:"database_path"`
	DatabaseType string    `json:"database_type"`
	LastLoaded   time.Time `json:"last_loaded"`
	Uptime       string    `json:"uptime"`
}

var (
	dbManager  *DatabaseManager
	serverStart time.Time
)

func init() {
	serverStart = time.Now()
}

// NewDatabaseManager creates a new database manager
func NewDatabaseManager(dbPath string) (*DatabaseManager, error) {
	dm := &DatabaseManager{
		dbPath: dbPath,
	}
	
	if err := dm.LoadDatabase(); err != nil {
		return nil, err
	}
	
	return dm, nil
}

// LoadDatabase loads or reloads the GeoIP database
func (dm *DatabaseManager) LoadDatabase() error {
	dm.mu.Lock()
	defer dm.mu.Unlock()

	// Close existing database if any
	if dm.db != nil {
		dm.db.Close()
	}

	// Determine database type from filename
	dbFilename := filepath.Base(dm.dbPath)
	if strings.Contains(strings.ToLower(dbFilename), "dbip") {
		dm.dbType = "DB-IP"
	} else if strings.Contains(strings.ToLower(dbFilename), "geolite") {
		dm.dbType = "GeoLite2"
	} else {
		dm.dbType = "Unknown"
	}

	// Open database
	db, err := geoip2.Open(dm.dbPath)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	dm.db = db
	dm.lastLoad = time.Now()
	
	log.Printf("Database loaded successfully: %s (Type: %s)", dm.dbPath, dm.dbType)
	return nil
}

// Lookup performs a GeoIP lookup
func (dm *DatabaseManager) Lookup(ipStr string) (*GeoIPResponse, error) {
	dm.mu.RLock()
	defer dm.mu.RUnlock()

	if dm.db == nil {
		return nil, fmt.Errorf("database not loaded")
	}

	ip := net.ParseIP(ipStr)
	if ip == nil {
		return nil, fmt.Errorf("invalid IP address: %s", ipStr)
	}

	record, err := dm.db.City(ip)
	if err != nil {
		return nil, fmt.Errorf("lookup failed: %w", err)
	}

	response := &GeoIPResponse{
		IP: ipStr,
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
		Postal: PostalInfo{
			Code: record.Postal.Code,
		},
		Traits: TraitsInfo{
			IsAnonymousProxy:    record.Traits.IsAnonymousProxy,
			IsSatelliteProvider: record.Traits.IsSatelliteProvider,
		},
	}

	// Add subdivisions if present
	for _, sub := range record.Subdivisions {
		response.Subdivisions = append(response.Subdivisions, SubdivisionInfo{
			IsoCode: sub.IsoCode,
			Names:   sub.Names,
		})
	}

	return response, nil
}

// GetStatus returns the database status
func (dm *DatabaseManager) GetStatus() StatusResponse {
	dm.mu.RLock()
	defer dm.mu.RUnlock()

	status := "healthy"
	if dm.db == nil {
		status = "not loaded"
	}

	uptime := time.Since(serverStart).Round(time.Second)

	return StatusResponse{
		Status:       status,
		DatabasePath: dm.dbPath,
		DatabaseType: dm.dbType,
		LastLoaded:   dm.lastLoad,
		Uptime:       uptime.String(),
	}
}

// API Handlers
func lookupHandler(c *gin.Context) {
	ipParam := c.Param("ip")
	if ipParam == "" {
		ipParam = c.Query("ip")
	}

	if ipParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "IP parameter required"})
		return
	}

	result, err := dbManager.Lookup(ipParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func statusHandler(c *gin.Context) {
	status := dbManager.GetStatus()
	c.JSON(http.StatusOK, status)
}

func healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func main() {
	// Get database path from environment or use default
	dbPath := os.Getenv("GEOIP_DB_PATH")
	if dbPath == "" {
		dbPath = "GeoLite2-City.mmdb"
	}

	// Check if database file exists
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		log.Printf("Warning: Database file not found at %s", dbPath)
		log.Printf("Please download a GeoIP database first")
		log.Printf("You can use: ./scripts/download_geolite2.sh or ./scripts/download_dbip.sh")
	}

	// Initialize database manager
	var err error
	dbManager, err = NewDatabaseManager(dbPath)
	if err != nil {
		log.Fatalf("Failed to initialize database manager: %v", err)
	}

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Set gin mode
	if os.Getenv("GIN_MODE") == "" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Setup router
	router := gin.Default()

	// API routes
	api := router.Group("/api/v1")
	{
		api.GET("/lookup/:ip", lookupHandler)
		api.GET("/lookup", lookupHandler)
		api.GET("/status", statusHandler)
		api.GET("/health", healthHandler)
	}

	// Serve static files
	staticFS, err := fs.Sub(staticFiles, "static")
	if err != nil {
		log.Fatalf("Failed to load static files: %v", err)
	}
	router.StaticFS("/static", http.FS(staticFS))

	// Serve index page
	router.GET("/", func(c *gin.Context) {
		data, err := staticFiles.ReadFile("static/index.html")
		if err != nil {
			c.String(http.StatusInternalServerError, "Failed to load index page")
			return
		}
		c.Data(http.StatusOK, "text/html; charset=utf-8", data)
	})

	log.Printf("Starting GeoIP API Server on port %s", port)
	log.Printf("Database: %s", dbPath)
	log.Printf("Web UI: http://localhost:%s", port)
	log.Printf("API: http://localhost:%s/api/v1/lookup/<ip>", port)
	
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
