package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/oschwald/geoip2-golang"
)

// 测试用的 IP 地址列表
var testIPs = []struct {
	ip          string
	description string
}{
	{"8.8.8.8", "Google DNS (美国)"},
	{"1.1.1.1", "Cloudflare DNS (美国)"},
	{"114.114.114.114", "中国 DNS"},
	{"223.5.5.5", "阿里云 DNS (中国)"},
	{"202.106.0.20", "中国联通 DNS"},
	{"208.67.222.222", "OpenDNS (美国)"},
	{"185.228.168.9", "CleanBrowsing DNS (欧洲)"},
	{"77.88.8.8", "Yandex DNS (俄罗斯)"},
}

func main() {
	// 打开数据库 (相对于项目根目录)
	db, err := geoip2.Open("../../GeoLite2-City.mmdb")
	if err != nil {
		log.Fatalf("❌ 无法打开数据库: %v", err)
	}
	defer db.Close()

	fmt.Println("========================================")
	fmt.Println("GeoIP 数据库测试")
	fmt.Println("========================================\n")

	successCount := 0
	failCount := 0

	for _, test := range testIPs {
		ip := net.ParseIP(test.ip)
		record, err := db.City(ip)

		if err != nil {
			fmt.Printf("❌ %s (%s) - 查询失败: %v\n", test.ip, test.description, err)
			failCount++
			continue
		}

		// 获取国家和城市名称
		countryName := record.Country.Names["zh-CN"]
		if countryName == "" {
			countryName = record.Country.Names["en"]
		}

		cityName := record.City.Names["zh-CN"]
		if cityName == "" {
			cityName = record.City.Names["en"]
		}
		if cityName == "" {
			cityName = "N/A"
		}

		fmt.Printf("✓ %s (%s)\n", test.ip, test.description)
		fmt.Printf("  国家: %s (%s)\n", countryName, record.Country.IsoCode)
		fmt.Printf("  城市: %s\n", cityName)
		fmt.Printf("  经纬度: %.4f, %.4f\n", record.Location.Latitude, record.Location.Longitude)
		fmt.Println()

		successCount++
	}

	fmt.Println("========================================")
	fmt.Printf("测试完成: 成功 %d/%d", successCount, len(testIPs))

	if failCount > 0 {
		fmt.Printf(", 失败 %d\n", failCount)
		os.Exit(1)
	} else {
		fmt.Println()
		fmt.Println("✓ 所有测试通过！")
	}
	fmt.Println("========================================")
}
