#!/bin/bash

# QZT-GO Linux/Mac 快速启动脚本

echo "===================================="
echo "  QZT-GO 快速启动"
echo "===================================="
echo ""

# 检查 Go 是否安装
if ! command -v go &> /dev/null; then
    echo "[错误] Go 未安装，请先安装 Go 1.22+"
    exit 1
fi

echo "[1/5] 检查 Go 版本..."
go version
echo ""

# 检查数据库连接配置
echo "[2/5] 检查数据库配置..."
if [ ! -f "config/config.yaml" ]; then
    echo "[警告] 配置文件不存在，请创建 config/config.yaml"
    exit 1
fi
echo "配置文件存在"
echo ""

# 下载依赖
echo "[3/5] 下载 Go 依赖..."
go mod download
if [ $? -ne 0 ]; then
    echo "[错误] 依赖下载失败"
    exit 1
fi
echo "依赖下载完成"
echo ""

# 编译项目
echo "[4/5] 编译项目..."
go build -o qzt-server cmd/server/main.go
if [ $? -ne 0 ]; then
    echo "[错误] 编译失败"
    exit 1
fi
echo "编译成功"
echo ""

# 启动服务
echo "[5/5] 启动服务..."
echo ""
echo "服务将在 http://localhost:8080 启动"
echo "按 Ctrl+C 停止服务"
echo ""
./qzt-server
