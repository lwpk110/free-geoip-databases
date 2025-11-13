package main

import (
	"testing"
)

func TestConfig(t *testing.T) {
	// Test Config struct exists
	config := Config{
		Port:         "8080",
		DatabasePath: "/path/to/db.mmdb",
		UpdateCron:   "0 0 * * *",
		AllowedDBs:   []string{"geolite2", "dbip"},
	}

	if config.Port != "8080" {
		t.Errorf("Expected port 8080, got %s", config.Port)
	}

	if config.DatabasePath != "/path/to/db.mmdb" {
		t.Errorf("Expected database path /path/to/db.mmdb, got %s", config.DatabasePath)
	}
}

func TestGeoIPResponse(t *testing.T) {
	// Test GeoIPResponse struct
	response := GeoIPResponse{
		IP: "8.8.8.8",
		Country: CountryInfo{
			IsoCode: "US",
			Names: map[string]string{
				"en":    "United States",
				"zh-CN": "美国",
			},
		},
		Location: LocationInfo{
			Latitude:  37.751,
			Longitude: -97.822,
		},
	}

	if response.IP != "8.8.8.8" {
		t.Errorf("Expected IP 8.8.8.8, got %s", response.IP)
	}

	if response.Country.IsoCode != "US" {
		t.Errorf("Expected country code US, got %s", response.Country.IsoCode)
	}
}

func TestStatusResponse(t *testing.T) {
	// Test StatusResponse struct
	status := StatusResponse{
		Status:       "healthy",
		DatabasePath: "/path/to/db.mmdb",
		DatabaseType: "GeoLite2",
	}

	if status.Status != "healthy" {
		t.Errorf("Expected status healthy, got %s", status.Status)
	}

	if status.DatabaseType != "GeoLite2" {
		t.Errorf("Expected database type GeoLite2, got %s", status.DatabaseType)
	}
}
