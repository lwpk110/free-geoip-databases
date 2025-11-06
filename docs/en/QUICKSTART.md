# Quick Start Guide

[ [中文](../cn/QUICKSTART.md) | English ]

This guide will help you quickly set up and run GeoIP query tools and configure auto-update functionality.

## 📋 Prerequisites

- Git
- Go 1.21 or higher
- GitHub account (for auto-update feature)
- MaxMind account (for obtaining License Key)

## 🚀 5-Minute Quick Start

### Step 1: Clone the Repository

```bash
git clone https://github.com/YOUR_USERNAME/geoip.git
cd geoip
```

### Step 2: Get Database

#### Option A: Download from Release (Recommended)

```bash
# If the repository already has a Release, download directly
curl -L -o GeoLite2-City.mmdb \
  https://github.com/YOUR_USERNAME/geoip/releases/latest/download/GeoLite2-City.mmdb
```

#### Option B: Manual Download

1. Visit [MaxMind signup](https://www.maxmind.com/en/geolite2/signup) to register
2. Generate License Key
3. Run download script:

```bash
chmod +x download_geolite2.sh
./download_geolite2.sh YOUR_LICENSE_KEY
```

### Step 3: Install Dependencies and Run

```bash
# Install Go dependencies
go mod tidy

# Run program
go run main.go
```

You should see output similar to:

```text
=== Starting IP Geolocation Query ===

Query IP: 8.8.8.8
  Country: United States
  Province/State:
  City:
  Timezone: America/Chicago
  Coordinates: 37.7510, -97.8220
  Postal Code:
--------------------
...
```

## 🤖 Configure Auto-Update (Optional)

If you want the database to auto-update, follow these steps:

### Step 1: Fork the Repository

Click the "Fork" button in the upper right corner of the GitHub page

### Step 2: Get MaxMind License Key

1. Visit [https://www.maxmind.com/en/geolite2/signup](https://www.maxmind.com/en/geolite2/signup)
2. Register and verify your email
3. After login, visit [https://www.maxmind.com/en/accounts/current/license-key](https://www.maxmind.com/en/accounts/current/license-key)
4. Click "Generate new license key"
5. Copy the generated License Key

### Step 3: Add GitHub Secret

1. Go to your forked repository
2. Click **Settings** → **Secrets and variables** → **Actions**
3. Click **New repository secret**
4. Add:
   - Name: `MAXMIND_LICENSE_KEY`
   - Value: Your License Key
5. Click **Add secret**

### Step 4: Enable Workflow Permissions

1. In **Settings** → **Actions** → **General**
2. Scroll to **Workflow permissions**
3. Select **Read and write permissions**
4. Check **Allow GitHub Actions to create and approve pull requests**
5. Click **Save**

### Step 5: Manually Trigger First Update

1. Go to **Actions** tab
2. Select "Update GeoLite2 Database" workflow
3. Click **Run workflow**
4. Select `main` branch
5. Click the green **Run workflow** button

Wait a few minutes, and you should see the newly released database in the **Releases** page!

## 📅 Auto-Update Schedule

After configuration:

- Workflow runs automatically every Tuesday and Friday at UTC 10:00 (Beijing Time 18:00)
- Each update creates a new Release
- Release contains three database files:
  - GeoLite2-City.mmdb
  - GeoLite2-Country.mmdb
  - GeoLite2-ASN.mmdb

## 💡 Usage Tips

### Use in Your Project

```go
package main

import (
    "fmt"
    "log"
    "net"

    "github.com/oschwald/geoip2-golang"
)

func main() {
    db, err := geoip2.Open("GeoLite2-City.mmdb")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    ip := net.ParseIP("1.1.1.1")
    record, err := db.City(ip)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Country: %s\n", record.Country.Names["en"])
    fmt.Printf("City: %s\n", record.City.Names["en"])
}
```

### Auto-Download from Latest Release

Add to your script:

```bash
#!/bin/bash

REPO="YOUR_USERNAME/geoip"
DB_FILE="GeoLite2-City.mmdb"

echo "Downloading latest GeoLite2 database..."
curl -L -o "$DB_FILE" \
  "https://github.com/${REPO}/releases/latest/download/${DB_FILE}"

echo "✓ Download complete: $DB_FILE"
```

### Regularly Update Local Database

Create a cron task:

```bash
# Edit crontab
crontab -e

# Add this line (update every Wednesday at 2 AM)
0 2 * * 3 cd /path/to/your/project && curl -L -o GeoLite2-City.mmdb https://github.com/YOUR_USERNAME/geoip/releases/latest/download/GeoLite2-City.mmdb
```

## 🔍 Troubleshooting

### "Database file not found"

Ensure the `GeoLite2-City.mmdb` file is in the program's working directory:

```bash
ls -lh GeoLite2-City.mmdb
```

### "Invalid License Key"

1. Check if License Key is copied correctly (no extra spaces)
2. Confirm License Key hasn't expired
3. Visit MaxMind account to verify status

### "Workflow run failed"

1. Check GitHub Actions logs
2. Confirm Secret is added correctly
3. Verify workflow permissions are enabled

## 📚 More Resources

- [Complete README](README.md)
- [Download Guide](DOWNLOAD_GUIDE.md)
- [GitHub Actions Setup Details](.github/SETUP.md)
- [MaxMind Official Documentation](https://dev.maxmind.com/geoip/geolite2-free-geolocation-data)

## 🆘 Need Help?

- Check [Issues](../../issues) page
- Submit a new [Issue](../../issues/new)
- Read [MaxMind Documentation](https://dev.maxmind.com/geoip/)

---

Enjoy! 🎉
