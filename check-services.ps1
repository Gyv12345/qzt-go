#!/usr/bin/env powershell
# 检查并启动所有服务

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "   QZT Go 全栈启动检查" -ForegroundColor Cyan
Write-Host "========================================`n" -ForegroundColor Cyan

# 1. 检查 Go
Write-Host "[检查] Go 环境..." -ForegroundColor Yellow
$goPaths = @(
    "C:\Program Files\Go\bin\go.exe",
    "C:\Go\bin\go.exe"
)
$goExe = $null
foreach ($path in $goPaths) {
    if (Test-Path $path) {
        $goExe = $path
        Write-Host "✅ Go 已安装: $path" -ForegroundColor Green
        & $goExe version
        break
    }
}
if (-not $goExe) {
    Write-Host "❌ Go 未找到，请重启 PowerShell 或手动安装" -ForegroundColor Red
}

# 2. 检查 PostgreSQL
Write-Host "`n[检查] PostgreSQL..." -ForegroundColor Yellow
$env:PGPASSWORD = "postgres"
try {
    $result = & psql -U postgres -c "SELECT version();" 2>&1
    if ($LASTEXITCODE -eq 0) {
        Write-Host "✅ PostgreSQL 连接成功" -ForegroundColor Green
    }
} catch {
    Write-Host "❌ PostgreSQL 连接失败: $_" -ForegroundColor Red
}

# 3. 检查 Redis
Write-Host "`n[检查] Redis..." -ForegroundColor Yellow
try {
    $redis = & redis-cli -a 123456 ping 2>&1
    if ($redis -eq "PONG") {
        Write-Host "✅ Redis 连接成功" -ForegroundColor Green
    }
} catch {
    Write-Host "❌ Redis 连接失败: $_" -ForegroundColor Red
}

# 4. 检查端口占用
Write-Host "`n[检查] 端口占用..." -ForegroundColor Yellow
$ports = @(8080, 5173, 3456, 3457, 3458)
foreach ($port in $ports) {
    $connection = netstat -ano | Select-String ":$port "
    if ($connection) {
        Write-Host "⚠️  端口 $port 已被占用" -ForegroundColor Yellow
        $connection | ForEach-Object { Write-Host "  $_" -ForegroundColor Gray }
    } else {
        Write-Host "✅ 端口 $port 可用" -ForegroundColor Green
    }
}

Write-Host "`n========================================`n" -ForegroundColor Cyan

# 启动选项
Write-Host "启动选项:" -ForegroundColor Yellow
Write-Host "  1. 启动后端（需要新的 PowerShell 窗口）"
Write-Host "  2. 前端已在运行: http://localhost:3458/"
Write-Host "  3. 退出"
Write-Host ""

$choice = Read-Host "请选择 (1-3)"
switch ($choice) {
    "1" {
        if ($goExe) {
            Write-Host "`n启动后端..." -ForegroundColor Green
            Set-Location D:\ccproject\qzt-go
            & $goExe run cmd/server/main.go
        } else {
            Write-Host "❌ Go 未找到，无法启动后端" -ForegroundColor Red
        }
    }
    "2" {
        Write-Host "`n前端地址: http://localhost:3458/" -ForegroundColor Green
        Start-Process "http://localhost:3458/"
    }
    "3" {
        Write-Host "退出" -ForegroundColor Yellow
    }
}
