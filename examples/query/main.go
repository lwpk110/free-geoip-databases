package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/oschwald/geoip2-golang"
)

// 获取当前主机的外网 IP
func getPublicIP() (string, error) {
	resp, err := http.Get("https://api.ipify.org?format=text")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	ip, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(ip), nil
}

// 打印所有 GeoIP2 记录字段
func printAllFields(ipStr string, record *geoip2.City) {
	fmt.Printf("\n%s\n", strings.Repeat("=", 100))
	fmt.Printf("IP 地址: %s - 完整字段信息\n", ipStr)
	fmt.Printf("%s\n\n", strings.Repeat("=", 100))

	// 大陆信息
	fmt.Println("【大陆信息】")
	fmt.Printf("  代码: %s\n", record.Continent.Code)
	fmt.Printf("  GeoName ID: %d\n", record.Continent.GeoNameID)
	fmt.Printf("  名称 (英文): %s\n", record.Continent.Names["en"])
	fmt.Printf("  名称 (中文): %s\n", record.Continent.Names["zh-CN"])
	fmt.Println()

	// 国家信息
	fmt.Println("【国家信息】")
	fmt.Printf("  ISO 代码: %s\n", record.Country.IsoCode)
	fmt.Printf("  GeoName ID: %d\n", record.Country.GeoNameID)
	fmt.Printf("  名称 (英文): %s\n", record.Country.Names["en"])
	fmt.Printf("  名称 (中文): %s\n", record.Country.Names["zh-CN"])
	fmt.Printf("  是否在欧盟: %v\n", record.Country.IsInEuropeanUnion)
	fmt.Println()

	// 注册国家信息
	fmt.Println("【注册国家信息】")
	fmt.Printf("  ISO 代码: %s\n", record.RegisteredCountry.IsoCode)
	fmt.Printf("  GeoName ID: %d\n", record.RegisteredCountry.GeoNameID)
	fmt.Printf("  名称 (英文): %s\n", record.RegisteredCountry.Names["en"])
	fmt.Printf("  名称 (中文): %s\n", record.RegisteredCountry.Names["zh-CN"])
	fmt.Printf("  是否在欧盟: %v\n", record.RegisteredCountry.IsInEuropeanUnion)
	fmt.Println()

	// 代表国家信息
	if record.RepresentedCountry.IsoCode != "" {
		fmt.Println("【代表国家信息】")
		fmt.Printf("  ISO 代码: %s\n", record.RepresentedCountry.IsoCode)
		fmt.Printf("  GeoName ID: %d\n", record.RepresentedCountry.GeoNameID)
		fmt.Printf("  名称 (英文): %s\n", record.RepresentedCountry.Names["en"])
		fmt.Printf("  类型: %s\n", record.RepresentedCountry.Type)
		fmt.Println()
	}

	// 行政区划信息 (省/州)
	if len(record.Subdivisions) > 0 {
		fmt.Println("【行政区划信息】")
		for i, subdivision := range record.Subdivisions {
			fmt.Printf("  级别 %d:\n", i+1)
			fmt.Printf("    ISO 代码: %s\n", subdivision.IsoCode)
			fmt.Printf("    GeoName ID: %d\n", subdivision.GeoNameID)
			fmt.Printf("    名称 (英文): %s\n", subdivision.Names["en"])
			fmt.Printf("    名称 (中文): %s\n", subdivision.Names["zh-CN"])
		}
		fmt.Println()
	}

	// 城市信息
	fmt.Println("【城市信息】")
	fmt.Printf("  GeoName ID: %d\n", record.City.GeoNameID)
	fmt.Printf("  名称 (英文): %s\n", record.City.Names["en"])
	fmt.Printf("  名称 (中文): %s\n", record.City.Names["zh-CN"])
	fmt.Println()

	// 邮政编码
	fmt.Println("【邮政信息】")
	fmt.Printf("  邮政编码: %s\n", record.Postal.Code)
	fmt.Println()

	// 地理位置信息
	fmt.Println("【地理位置信息】")
	fmt.Printf("  纬度: %.6f\n", record.Location.Latitude)
	fmt.Printf("  经度: %.6f\n", record.Location.Longitude)
	fmt.Printf("  精度半径 (km): %d\n", record.Location.MetroCode)
	fmt.Printf("  时区: %s\n", record.Location.TimeZone)
	fmt.Printf("  精度半径: %d\n", record.Location.AccuracyRadius)
	fmt.Println()

	// 特征信息
	fmt.Println("【特征信息】")
	fmt.Printf("  匿名代理: %v\n", record.Traits.IsAnonymousProxy)
	fmt.Printf("  卫星供应商: %v\n", record.Traits.IsSatelliteProvider)
	fmt.Println()

	// JSON 格式输出所有数据
	fmt.Println("【完整 JSON 数据】")
	jsonData, err := json.MarshalIndent(record, "", "  ")
	if err != nil {
		fmt.Printf("JSON 序列化错误: %v\n", err)
	} else {
		fmt.Println(string(jsonData))
	}

	fmt.Printf("\n%s\n\n", strings.Repeat("=", 100))
}

func main() {
	// 支持通过环境变量指定数据库文件路径
	// 默认使用 GeoLite2-City.mmdb (相对于项目根目录)
	dbPath := os.Getenv("GEOIP_DB_PATH")
	if dbPath == "" {
		dbPath = "../../GeoLite2-City.mmdb"
	}

	// 打开数据库
	db, err := geoip2.Open(dbPath)
	if err != nil {
		log.Fatalf("无法打开数据库 %s: %v", dbPath, err)
	}
	defer db.Close()

	fmt.Printf("使用数据库: %s\n\n", dbPath)

	// 获取当前主机的外网 IP
	fmt.Println("正在获取当前主机外网 IP...")
	publicIP, err := getPublicIP()
	if err != nil {
		log.Fatalf("获取外网 IP 失败: %v", err)
	}

	fmt.Printf("当前主机外网 IP: %s\n", publicIP)

	// 解析 IP
	ip := net.ParseIP(publicIP)
	if ip == nil {
		log.Fatalf("无效的 IP 地址: %s", publicIP)
	}

	// 查询 IP 信息
	record, err := db.City(ip)
	if err != nil {
		log.Fatalf("查询 IP %s 时出错: %v", publicIP, err)
	}

	// 打印所有字段信息
	printAllFields(publicIP, record)
}
