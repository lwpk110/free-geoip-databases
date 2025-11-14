# Free GeoIP Databases - Auto-Updated 🌍

**[ [中文](README.md) | English ]**

**📚 Documentation: [English](docs/en/) | [中文](docs/cn/)**

## 📦 Database Updates

[![GeoLite2](https://github.com/lwpk110/free-geoip-databases/actions/workflows/update-geolite2.yml/badge.svg)](https://github.com/lwpk110/free-geoip-databases/actions/workflows/update-geolite2.yml)
[![DB-IP](https://github.com/lwpk110/free-geoip-databases/actions/workflows/update-dbip.yml/badge.svg)](https://github.com/lwpk110/free-geoip-databases/actions/workflows/update-dbip.yml)

## 🧪 Quality Checks

[![GeoLite2 Tests](https://github.com/lwpk110/free-geoip-databases/actions/workflows/test-geolite.yml/badge.svg)](https://github.com/lwpk110/free-geoip-databases/actions/workflows/test-geolite.yml)
[![DB-IP Tests](https://github.com/lwpk110/free-geoip-databases/actions/workflows/test-dbip.yml/badge.svg)](https://github.com/lwpk110/free-geoip-databases/actions/workflows/test-dbip.yml)

---

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![GitHub Stars](https://img.shields.io/github/stars/lwpk110/free-geoip-databases?style=flat-square)](https://github.com/lwpk110/free-geoip-databases/stargazers)
[![GitHub Forks](https://img.shields.io/github/forks/lwpk110/free-geoip-databases?style=flat-square)](https://github.com/lwpk110/free-geoip-databases/network/members)

> **Free, auto-updated GeoIP database collection** - Automated download and release service for MaxMind GeoLite2, DB-IP and other databases

No MaxMind account registration needed, no License Key required - download the latest databases directly from [Releases](../../releases)!

## 🎯 Why Choose This Project?

- ✅ **Completely Free** - No registration, no License Key required (DB-IP version)
- 🤖 **Auto-Updated** - GeoLite2 updates weekly (Tue, Fri), DB-IP updates monthly (1st, 15th)
- 📦 **Ready to Use** - Download `.mmdb` files and start using immediately
- 🗂️ **Multiple Databases** - City, Country, ASN and more
- 🌐 **Multiple Sources** - Both MaxMind GeoLite2 and DB-IP databases
- ✅ **Quality Assured** - Automated testing ensures database integrity
- 🚀 **Self-Hosted API** - Ready-to-use GeoIP query API service with Web UI

## 🆕 New Feature: Self-Hosted GeoIP API Service

This project now provides a Golang-based GeoIP query API service with:

- 🌐 **RESTful API** - Standard HTTP API endpoints
- 🖥️ **Web Interface** - User-friendly query interface
- 🔄 **Multi-Database Support** - Compatible with GeoLite2 and DB-IP
- 🚀 **Easy Deployment** - Docker support, one-click start
- 📊 **Real-time Monitoring** - Database status and system monitoring

### Quick Start API Service

```bash
# 1. Clone repository
git clone https://github.com/lwpk110/free-geoip-databases.git
cd free-geoip-databases

# 2. Download database (using DB-IP, no registration required)
./scripts/download_dbip.sh all

# 3. Start service
make run

# Or use Docker
make docker-run
```

Visit http://localhost:8080 to use the Web interface for IP queries!

For detailed documentation, see [GeoIP API Documentation](cmd/geoip-api/README.md).

## 📥 Quick Download

### MaxMind GeoLite2 Databases

Visit the [**Releases page**](../../releases/latest) to download the latest database files:

| Database | Description | File Size | License | Download |
|----------|-------------|-----------|---------|----------|
| **GeoLite2-City** | City-level geolocation data | ~70 MB | GeoLite2 EULA | [Releases page](../../releases?q=geolite2) |
| **GeoLite2-Country** | Country-level geolocation data | ~6 MB | GeoLite2 EULA | [Releases page](../../releases?q=geolite2) |
| **GeoLite2-ASN** | ASN network operator data | ~8 MB | GeoLite2 EULA | [Releases page](../../releases?q=geolite2) |

> **Note**: GeoLite2 database filenames include date information (e.g., `GeoLite2-City-20251105.mmdb`). Please visit the Releases page to select and download the latest version.

### DB-IP Databases

Check [**DB-IP Releases**](../../releases?q=dbip&expanded=true) to download DB-IP databases:

| Database | Description | File Size | License | Download |
|----------|-------------|-----------|---------|----------|
| **DB-IP City Lite** | City-level geolocation data | ~130 MB | CC BY 4.0 | [Releases page](../../releases?q=dbip) |
| **DB-IP Country Lite** | Country-level geolocation data | ~7 MB | CC BY 4.0 | [Releases page](../../releases?q=dbip) |
| **DB-IP ASN Lite** | ASN network operator data | ~9 MB | CC BY 4.0 | [Releases page](../../releases?q=dbip) |

> **Note**: DB-IP database filenames include year-month information (e.g., `dbip-city-lite-2025-11.mmdb`). Please visit the Releases page to select and download the latest version.

### Command Line Download

Since filenames include date information, we recommend getting the latest databases via:

#### Method 1: Use project download scripts

```bash
# MaxMind GeoLite2 (requires License Key)
./scripts/download_geolite2.sh <YOUR_LICENSE_KEY>

# DB-IP (no License Key required)
./scripts/download_dbip.sh all
```

#### Method 2: Download from GitHub Releases

```bash
# Visit Releases page to select the latest version
# GeoLite2: https://github.com/lwpk110/free-geoip-databases/releases?q=geolite2
# DB-IP: https://github.com/lwpk110/free-geoip-databases/releases?q=dbip

# Or use GitHub CLI
gh release list --repo lwpk110/free-geoip-databases
gh release download <tag-name> --repo lwpk110/free-geoip-databases
```

## 🚀 Usage Examples

### Go

```go
package main

import (
    "fmt"
    "log"
    "net"

    "github.com/oschwald/geoip2-golang"
)

func main() {
    // Open database
    db, err := geoip2.Open("GeoLite2-City.mmdb")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    // Query IP
    ip := net.ParseIP("8.8.8.8")
    record, err := db.City(ip)
    if err != nil {
        log.Fatal(err)
    }

    // Print results
    fmt.Printf("Country: %s\n", record.Country.Names["en"])
    fmt.Printf("City: %s\n", record.City.Names["en"])
    fmt.Printf("Latitude, Longitude: %.4f, %.4f\n", record.Location.Latitude, record.Location.Longitude)
}
```

### Python

```python
import geoip2.database

# Open database
with geoip2.database.Reader('GeoLite2-City.mmdb') as reader:
    # Query IP
    response = reader.city('8.8.8.8')

    # Print results
    print(f"Country: {response.country.names['en']}")
    print(f"City: {response.city.names.get('en', 'N/A')}")
    print(f"Latitude, Longitude: {response.location.latitude}, {response.location.longitude}")
```

### Node.js

```javascript
const maxmind = require('maxmind');

// Open database
maxmind.open('GeoLite2-City.mmdb').then(lookup => {
    // Query IP
    const result = lookup.get('8.8.8.8');

    // Print results
    console.log(`Country: ${result.country.names['en']}`);
    console.log(`City: ${result.city.names['en']}`);
    console.log(`Latitude, Longitude: ${result.location.latitude}, ${result.location.longitude}`);
});
```

## 🔄 Update Frequency

### MaxMind GeoLite2

- **Auto-update**: Every Tuesday and Friday (UTC 10:00 / Beijing Time 18:00)
- **Data Source**: Official MaxMind GeoLite2 databases
- **Update Strategy**: Follows MaxMind's official update cycle

### DB-IP

- **Auto-update**: 1st and 15th of each month (UTC 10:00 / Beijing Time 18:00)
- **Data Source**: Official DB-IP free databases
- **Update Strategy**: Downloads the latest version of the current month
- **License**: Creative Commons Attribution 4.0 (CC BY 4.0)

## 📊 Database Comparison

### MaxMind GeoLite2 vs DB-IP

| Feature | GeoLite2 | DB-IP Lite |
|---------|----------|-----------|
| **License** | MaxMind EULA required | CC BY 4.0 (more permissive) |
| **Update Frequency** | Twice weekly | Once monthly |
| **Data Accuracy** | ⭐⭐⭐⭐ | ⭐⭐⭐⭐ |
| **File Format** | .mmdb | .mmdb (compatible) |
| **Usage Restrictions** | Non-commercial use preferred | Commercial use with attribution |
| **Registration** | Required (handled by this project) | Not required |

### Supported Database Types

| Database Type | Information Included | File Size |
|--------------|---------------------|-----------|
| **City** | Country, province/state, city, postal code, lat/long, timezone | ~70 MB |
| **Country** | Country, continent | ~6 MB |
| **ASN** | Autonomous System Number, network operator | ~8 MB |

## 🛠️ Running Local Test Tools

This project provides simple Go language query tools to test database files.

### Project Structure

```text
free-geoip-databases/
├── .github/workflows/    # GitHub Actions automation workflows
│   ├── update-geolite2.yml
│   ├── update-dbip.yml
│   └── test-database.yml
├── scripts/              # Download scripts
│   ├── download_geolite2.sh
│   └── download_dbip.sh
├── examples/             # Example code
│   ├── query/            # IP query examples
│   │   └── main.go
│   └── test/             # Test programs
│       └── test_cities.go
├── docs/                 # Detailed documentation
│   ├── QUICKSTART.md
│   └── TESTING.md
├── README.md
├── LICENSE
├── go.mod
└── go.sum
```

### Quick Start

```bash
# 1. Clone repository
git clone https://github.com/lwpk110/free-geoip-databases.git
cd free-geoip-databases

# 2. Download database files
./scripts/download_geolite2.sh <YOUR_LICENSE_KEY>
# Or use DB-IP (no License Key required)
./scripts/download_dbip.sh all

# 3. Install dependencies
go mod tidy

# 4. Run query example
cd examples/query
go run main.go

# 5. Run tests
cd ../test
go run test_cities.go
```

For more details, see [examples/README_EN.md](examples/README_EN.md) and [docs/en/TESTING.md](docs/en/TESTING.md).

## ⚙️ Fork This Project for Auto-Updates

### Configure MaxMind GeoLite2 Auto-Updates

If you want to fork this project and implement your own GeoLite2 auto-updates:

1. **Fork this repository**

2. **Add MaxMind License Key**
   - Visit [MaxMind registration](https://www.maxmind.com/en/geolite2/signup)
   - Generate a License Key
   - In your repository: Settings → Secrets and variables → Actions
   - Add Secret: `MAXMIND_LICENSE_KEY`

3. **Add GitHub Token**
   - In your repository: Settings → Secrets and variables → Actions
   - Add Secret: `GEOIP_ACCESS_TOKEN` (use a Personal Access Token with repo permissions)

4. **Enable GitHub Actions**
   - Workflows will run automatically
   - Or trigger manually: Actions → Update GeoLite2 Database → Run workflow

### Configure DB-IP Auto-Updates

DB-IP databases require no configuration, just:

1. **Fork this repository**

2. **Add GitHub Token**
   - In your repository: Settings → Secrets and variables → Actions
   - Add Secret: `GEOIP_ACCESS_TOKEN`

3. **Enable GitHub Actions**
   - DB-IP workflows will run automatically
   - Or trigger manually: Actions → Update DB-IP Database → Run workflow

> **Tip**: DB-IP doesn't require account registration or License Key - completely free and open!

## 📋 License and Attribution

- **Project Code**: MIT License
- **GeoLite2 Databases**: Provided by MaxMind, must comply with [MaxMind EULA](https://www.maxmind.com/en/geolite2/eula)
  - This product includes GeoLite2 data created by MaxMind, available from [https://www.maxmind.com](https://www.maxmind.com)
- **DB-IP Databases**: Provided by DB-IP, licensed under [Creative Commons Attribution 4.0](https://creativecommons.org/licenses/by/4.0/). Attribution required: Contains data from [https://db-ip.com](https://db-ip.com)

### Important Notes

- Both GeoLite2 and DB-IP Lite are free versions with lower accuracy than commercial versions
- GeoLite2: For learning, testing, and non-commercial use only. Purchase MaxMind commercial license for commercial use
- DB-IP: Commercial use allowed with attribution required
- For higher accuracy or commercial support, consider purchasing commercial database versions

## 🌟 Star History

If this project helps you, please give it a ⭐ Star! This motivates me to continue maintaining and updating.

## 🤝 Contributing

Issues and Pull Requests are welcome!

## 📞 Support

- 🐛 [Report Bugs](../../issues)
- 💡 [Feature Requests](../../issues)
- ⭐ [Star This Project](../../stargazers)
