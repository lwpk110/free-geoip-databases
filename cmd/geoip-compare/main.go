package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"github.com/oschwald/geoip2-golang"
)

type summary struct {
	SampleSize        int              `json:"sample_size"`
	GeoLiteFound      int              `json:"geolite_found"`
	DBIPFound         int              `json:"dbip_found"`
	BothFound         int              `json:"both_found"`
	CountryMismatch   int              `json:"country_mismatch"`
	ContinentMismatch int              `json:"continent_mismatch"`
	Samples           []map[string]any `json:"samples"`
}

func main() {
	geolitePath := flag.String("geolite", "", "GeoLite2 mmdb path")
	dbipPath := flag.String("dbip", "", "DB-IP mmdb path")
	ipsFile := flag.String("ips", "", "IP list file, one per line")
	flag.Parse()

	if *geolitePath == "" || *dbipPath == "" || *ipsFile == "" {
		log.Fatal("usage: geoip-compare -geolite <path> -dbip <path> -ips <file>")
	}

	ips, err := readIPs(*ipsFile)
	if err != nil {
		log.Fatal(err)
	}
	gdb, err := geoip2.Open(*geolitePath)
	if err != nil {
		log.Fatalf("open geolite: %v", err)
	}
	defer gdb.Close()
	dbip, err := geoip2.Open(*dbipPath)
	if err != nil {
		log.Fatalf("open dbip: %v", err)
	}
	defer dbip.Close()

	out := summary{SampleSize: len(ips), Samples: make([]map[string]any, 0, 10)}
	for _, ipStr := range ips {
		ip := net.ParseIP(ipStr)
		if ip == nil {
			continue
		}

		g, gErr := gdb.City(ip)
		d, dErr := dbip.City(ip)

		if gErr == nil {
			out.GeoLiteFound++
		}
		if dErr == nil {
			out.DBIPFound++
		}
		if gErr == nil && dErr == nil {
			out.BothFound++
			if g.Country.IsoCode != d.Country.IsoCode {
				out.CountryMismatch++
			}
			if g.Continent.Code != d.Continent.Code {
				out.ContinentMismatch++
			}
			if len(out.Samples) < 10 {
				out.Samples = append(out.Samples, map[string]any{
					"ip":      ipStr,
					"geolite": map[string]any{"country": g.Country.IsoCode, "continent": g.Continent.Code},
					"dbip":    map[string]any{"country": d.Country.IsoCode, "continent": d.Continent.Code},
				})
			}
		}
	}

	b, _ := json.MarshalIndent(out, "", "  ")
	fmt.Println(string(b))
}

func readIPs(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var ips []string
	s := bufio.NewScanner(f)
	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		if line != "" {
			ips = append(ips, line)
		}
	}
	return ips, s.Err()
}
