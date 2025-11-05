#!/bin/bash

# DB-IP 数据库下载脚本
# DB-IP 提供免费的 IP 地理位置数据库，使用 CC BY 4.0 授权
# 网址: https://db-ip.com/db/download/ip-to-city-lite

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${GREEN}=== DB-IP 数据库下载工具 ===${NC}\n"

# 获取当前年月
YEAR_MONTH=$(date +%Y-%m)
echo -e "${BLUE}当前年月: ${YEAR_MONTH}${NC}\n"

# 定义数据库类型
declare -A DATABASES=(
    ["city"]="dbip-city-lite"
    ["country"]="dbip-country-lite"
    ["asn"]="dbip-asn-lite"
)

# 显示可用的数据库类型
echo -e "${YELLOW}可用的数据库类型:${NC}"
echo "  1. city    - 城市级别 IP 地理位置数据库"
echo "  2. country - 国家级别 IP 地理位置数据库"
echo "  3. asn     - ASN (自治系统号) 数据库"
echo "  4. all     - 下载所有数据库"
echo ""

# 检查是否提供了数据库类型参数
if [ -z "$1" ]; then
    echo -e "${YELLOW}用法: $0 <database_type> [year-month]${NC}"
    echo -e "${YELLOW}示例: $0 city${NC}"
    echo -e "${YELLOW}示例: $0 all${NC}"
    echo -e "${YELLOW}示例: $0 city 2024-10${NC}"
    exit 1
fi

DB_TYPE="$1"
CUSTOM_YEAR_MONTH="${2:-$YEAR_MONTH}"

echo -e "${GREEN}下载年月: ${CUSTOM_YEAR_MONTH}${NC}\n"

# 下载函数
download_database() {
    local db_name=$1
    local file_name="${db_name}-${CUSTOM_YEAR_MONTH}"
    local download_url="https://download.db-ip.com/free/${file_name}.mmdb.gz"
    local temp_file="${file_name}.mmdb.gz"
    local output_file="${file_name}.mmdb"

    echo -e "${GREEN}正在下载 ${db_name} 数据库...${NC}"
    echo "下载地址: ${download_url}"
    echo ""

    # 下载数据库
    if curl -L -o "$temp_file" "$download_url"; then
        echo -e "\n${GREEN}✓ 下载成功${NC}"
    else
        echo -e "\n${RED}✗ 下载失败，请检查网络连接或年月参数是否正确${NC}"
        return 1
    fi

    # 解压文件
    echo -e "\n${GREEN}正在解压文件...${NC}"
    if gunzip "$temp_file"; then
        echo -e "${GREEN}✓ 解压成功${NC}"

        # 显示文件信息
        if [ -f "$output_file" ]; then
            FILE_SIZE=$(du -h "$output_file" | cut -f1)
            echo -e "\n${GREEN}文件信息:${NC}"
            echo "  文件名: $output_file"
            echo "  大小: $FILE_SIZE"
            echo -e "\n${GREEN}✓ 数据库已保存到当前目录${NC}"
        fi
    else
        echo -e "${RED}✗ 解压失败${NC}"
        return 1
    fi

    echo ""
}

# 根据用户选择下载
if [ "$DB_TYPE" = "all" ]; then
    echo -e "${BLUE}下载所有数据库...${NC}\n"
    for db_key in "${!DATABASES[@]}"; do
        download_database "${DATABASES[$db_key]}"
    done
elif [ -n "${DATABASES[$DB_TYPE]}" ]; then
    download_database "${DATABASES[$DB_TYPE]}"
else
    echo -e "${RED}错误: 未知的数据库类型 '$DB_TYPE'${NC}"
    echo -e "${YELLOW}请使用: city, country, asn 或 all${NC}"
    exit 1
fi

echo -e "${GREEN}=== 下载完成 ===${NC}\n"

echo -e "${BLUE}使用说明:${NC}"
echo "1. DB-IP 数据库与 MaxMind GeoIP2 格式兼容"
echo "2. 可以使用 MaxMind GeoIP2 库读取数据库"
echo "3. 授权: Creative Commons Attribution 4.0 (CC BY 4.0)"
echo "4. 更多信息请访问: https://db-ip.com"
echo ""
