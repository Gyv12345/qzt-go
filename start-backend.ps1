#!/usr/bin/env powershell
# 启动 Go 后端服务
# 需要在新的 PowerShell 窗口中运行

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "   QZT Go 后端服务启动" -ForegroundColor Cyan
Write-Host "========================================`n" -ForegroundColor Cyan

# 查找 Go 安装位置
$goPaths = @(
    "C:\Program Files\Go\bin\go.exe",
    "C:\Go\bin\go.exe",
    "${env:ProgramFiles}\Go\bin\go.exe",
    "${env:ProgramFiles(x86)}\Go\bin\go.exe"
)

$goExe = $null
foreach ($path in $goPaths) {
    if (Test-Path $path) {
        $goExe = $path
        break
    }
}

if (-not $goExe) {
    Write-Host "❌ 未找到 Go 安装" -ForegroundColor Red
    Write-Host @"

Go 已安装但需要重启 PowerShell 才能生效。

请：
1. 关闭所有 PowerShell 窗口
2. 重新打开 PowerShell
3. 运行: cd D:\ccproject\qzt-go; go run cmd/server/main.go

"@ -ForegroundColor Yellow
    Read-Host "按回车键退出"
    exit 1
}

Write-Host "✅ 找到 Go: $goExe" -ForegroundColor Green

# 检查数据库
Write-Host "`n[1/3] 检查数据库..." -ForegroundColor Yellow
$env:PGPASSWORD = "postgres"
try {
    $dbExists = & psql -U postgres -tAc "SELECT 1 FROM pg_database WHERE datname='qzt'" 2>&1
    if ($dbExists -eq "1") {
        Write-Host "✅ 数据库 qzt 已存在" -ForegroundColor Green
    } else {
        Write-Host "创建数据库 qzt..." -ForegroundColor Yellow
        & psql -U postgres -c "CREATE DATABASE qzt WITH ENCODING='UTF8'" 2>&1 | Out-Null
        Write-Host "✅ 数据库创建成功" -ForegroundColor Green
    }
} catch {
    Write-Host "⚠️  无法检查数据库，继续启动..." -ForegroundColor Yellow
}

# 下载依赖
Write-Host "`n[2/3] 下载 Go 依赖..." -ForegroundColor Yellow
Set-Location D:\ccproject\qzt-go
& $goExe mod tidy 2>&1 | Out-Null
Write-Host "✅ 依赖下载完成" -ForegroundColor Green

# 启动服务
Write-Host "`n[3/3] 启动后端服务..." -ForegroundColor Yellow
Write-Host @"

========================================
  服务启动中...
  访问地址: http://localhost:8080
  健康检查: http://localhost:8080/health
========================================

"@ -ForegroundColor Green

& $goExe run cmd/server/main.go
