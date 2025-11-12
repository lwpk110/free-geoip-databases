# Example Code

This directory contains example programs using GeoIP databases.

## Directory Structure

```text
examples/
├── query/        # IP geolocation query example
│   └── main.go
└── test/         # Database test program
    └── test_cities.go
```

## Usage Instructions

### 1. Query Example (query)

Query detailed geolocation information for IP addresses.

```bash
# Enter example directory
cd examples/query

# Run example (default uses GeoLite2-City.mmdb file in root directory)
go run main.go

# Or specify another database file (supports both GeoLite2 and DB-IP)
GEOIP_DB_PATH="../../dbip-city-lite-2025-11.mmdb" go run main.go
```

**Features**:

- Automatically get current host's public IP
- Display complete geolocation information (continent, country, province, city, coordinates, etc.)
- Output JSON format data
- **Supports both GeoLite2 and DB-IP databases** (specify via `GEOIP_DB_PATH` environment variable)

### 2. Test Program (test)

Test if the database is working properly.

```bash
# Enter test directory
cd examples/test

# Run tests (default uses GeoLite2-City.mmdb file in root directory)
go run test_cities.go

# Or specify another database file (supports both GeoLite2 and DB-IP)
GEOIP_DB_PATH="../../dbip-city-lite-2025-11.mmdb" go run test_cities.go
```

**Features**:

- Test queries for multiple well-known DNS IPs
- Verify database integrity
- Display test result statistics

## Prerequisites

1. **Download Database File**

Use download script to get database:

```bash
# Run from project root directory
./scripts/download_geolite2.sh <YOUR_LICENSE_KEY>
# or
./scripts/download_dbip.sh all
```

Or download from GitHub Releases:

```bash
curl -L -o GeoLite2-City.mmdb \
  https://github.com/lwpk110/free-geoip-databases/releases/latest/download/GeoLite2-City.mmdb
```

1. **Install Dependencies**

```bash
go mod download
```

## Database File Location

Example programs expect database file to be in the **project root directory**:

- `../../GeoLite2-City.mmdb` (viewed from examples/query or examples/test)

## Develop Your Own Program

Reference these examples to create your own GeoIP query program:

```go
package main

import (
    "github.com/oschwald/geoip2-golang"
    "net"
)

func main() {
    db, _ := geoip2.Open("GeoLite2-City.mmdb")
    defer db.Close()

    ip := net.ParseIP("8.8.8.8")
    record, _ := db.City(ip)

    println(record.Country.Names["en"])
}
```

## More Examples

For more language examples (Python, Node.js, PHP, etc.), please refer to:

- [Main README](../README_EN.md)
- [Quick Start Guide](../docs/en/QUICKSTART.md)
