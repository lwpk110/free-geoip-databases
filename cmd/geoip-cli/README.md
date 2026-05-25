# geoip-cli

CLI for GeoIP lookup in two modes:

- Local mode: read `.mmdb` directly with `-db`
- Remote mode: call API server with `-api`

## Examples

```bash
# single IP via API
geoip-cli -api http://localhost:8080 -ip 8.8.8.8

# batch via API
geoip-cli -api http://localhost:8080 -batch-file ips.txt

# local db lookup
geoip-cli -db ./GeoLite2-City.mmdb -ip 1.1.1.1
```
