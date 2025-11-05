#!/bin/bash

# GeoLite2 数据库下载脚本
# 使用前请先在 MaxMind 注册账号并生成 License Key
# 网址: https://www.maxmind.com/en/geolite2/signup

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}=== GeoLite2 数据库下载工具 ===${NC}\n"

# 检查是否提供了 License Key
if [ -z "$1" ]; then
    echo -e "${RED}错误: 请提供 MaxMind License Key${NC}"
    echo -e "${YELLOW}用法: $0 <YOUR_LICENSE_KEY>${NC}"
    echo -e "${YELLOW}请访问 https://www.maxmind.com/en/geolite2/signup 获取免费的 License Key${NC}"
    exit 1
else
    LICENSE_KEY="$1"
fi

echo -e "${GREEN}使用 License Key: ${LICENSE_KEY:0:10}...${NC}"
DOWNLOAD_URL="https://download.maxmind.com/app/geoip_download?edition_id=GeoLite2-City&license_key=${LICENSE_KEY}&suffix=tar.gz"
TEMP_FILE="GeoLite2-City.tar.gz"

echo -e "${GREEN}正在下载 GeoLite2-City 数据库...${NC}"
echo "下载地址: ${DOWNLOAD_URL}"
echo ""

# 下载数据库
if curl -L -o "$TEMP_FILE" "$DOWNLOAD_URL"; then
    echo -e "\n${GREEN}✓ 下载成功${NC}"
else
    echo -e "\n${RED}✗ 下载失败，请检查 License Key 是否正确${NC}"
    exit 1
fi

# 解压文件
echo -e "\n${GREEN}正在解压文件...${NC}"
tar -xzf "$TEMP_FILE"

# 查找解压后的 .mmdb 文件
MMDB_FILE=$(find . -name "GeoLite2-City.mmdb" -type f | head -n 1)

if [ -n "$MMDB_FILE" ]; then
    # 移动到当前目录
    mv "$MMDB_FILE" ./GeoLite2-City.mmdb
    echo -e "${GREEN}✓ 已解压到: GeoLite2-City.mmdb${NC}"

    # 清理临时文件
    rm -rf "$TEMP_FILE" GeoLite2-City_*

    # 显示文件信息
    echo -e "\n${GREEN}=== 数据库信息 ===${NC}"
    ls -lh GeoLite2-City.mmdb

    echo -e "\n${GREEN}✓ 下载完成！${NC}"
    echo -e "${YELLOW}提示: MaxMind 每周二更新数据库，建议定期运行此脚本更新${NC}"
else
    echo -e "${RED}✗ 未找到 .mmdb 文件${NC}"
    exit 1
fi
