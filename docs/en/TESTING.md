# Testing Guide

[ [中文](../cn/TESTING.md) | English ]

## Local Testing

### Test GeoLite2 Database

```bash
# 1. Ensure database file is downloaded
./download_geolite2.sh <YOUR_LICENSE_KEY>

# Or download from GitHub Releases
curl -L -o GeoLite2-City.mmdb \
  https://github.com/lwpk110/free-geoip-databases/releases/latest/download/GeoLite2-City.mmdb

# 2. Run tests
go run test_cities.go
```

### Test DB-IP Database

```bash
# 1. Download DB-IP database
./download_dbip.sh city

# 2. Rename file to standard format (if needed)
mv dbip-city-lite-*.mmdb GeoLite2-City.mmdb

# 3. Run tests
go run test_cities.go
```

## Automated Testing

### GitHub Actions Testing

The project includes automated testing workflow (`.github/workflows/test-database.yml`):

- **Trigger Conditions**:
  - Automatically runs after each new Release is created
  - Can be triggered manually

- **Test Content**:
  - Download latest GeoLite2-City.mmdb
  - Run IP geolocation query tests
  - Verify database integrity

- **Manual Trigger**:
  1. Go to repository's Actions page
  2. Select "Test GeoIP Database" workflow
  3. Click "Run workflow"

## Tested IP Addresses

`test_cities.go` tests the following well-known DNS server IPs:

- Google DNS (8.8.8.8)
- Cloudflare DNS (1.1.1.1)
- 114 DNS (114.114.114.114)
- Alibaba DNS (223.5.5.5)
- China Unicom DNS (202.106.0.20)
- OpenDNS (208.67.222.222)
- CleanBrowsing DNS (185.228.168.9)
- Yandex DNS (77.88.8.8)

## Add Custom Tests

Edit the `testIPs` array in `test_cities.go`:

```go
var testIPs = []struct {
    ip          string
    description string
}{
    {"your.ip.address.here", "Your Description"},
    // ... more test IPs
}
```

## Expected Output

Example of successful test output:

```text
========================================
GeoIP Database Test
========================================

✓ 8.8.8.8 (Google DNS (United States))
  Country: United States (US)
  City: N/A
  Lat/Long: 37.7510, -97.8220

...

========================================
Test Complete: Success 8/8
✓ All tests passed!
========================================
```

## Troubleshooting

### Database file not found

```text
❌ Cannot open database: open GeoLite2-City.mmdb: no such file or directory
```

**Solution**: Ensure `GeoLite2-City.mmdb` file is in the current directory.

### Some IP queries fail

Some IP addresses may not have detailed information in the free database, which is normal. The free version database has limited coverage.

### GitHub Actions failure

Check the following:

1. Is there a valid Release (containing `GeoLite2-City.mmdb` file)
2. Is `GEOIP_ACCESS_TOKEN` configured correctly
3. Is network connection normal
