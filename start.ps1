# PowerShell 初始化脚本

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "QZT Go 项目初始化" -ForegroundColor Cyan
Write-Host "========================================`n" -ForegroundColor Cyan

# 检查 Go
Write-Host "[1/4] 检查 Go 环境..." -ForegroundColor Yellow
try {
    $goVersion = go version
    Write-Host $goVersion -ForegroundColor Green
} catch {
    Write-Host "[错误] 未找到 Go，请先安装" -ForegroundColor Red
    Write-Host "下载: https://golang.org/dl/" -ForegroundColor Gray
    exit 1
}

# 下载依赖
Write-Host "`n[2/4] 下载依赖..." -ForegroundColor Yellow
go mod tidy
if ($LASTEXITCODE -ne 0) {
    Write-Host "[错误] 依赖下载失败" -ForegroundColor Red
    exit 1
}
Write-Host "依赖下载完成！" -ForegroundColor Green

# 提示创建数据库
Write-Host "`n[3/4] 检查 PostgreSQL..." -ForegroundColor Yellow
Write-Host @"

请确保已创建数据库 'qzt'。如果未创建，请执行：

  psql -U postgres
  CREATE DATABASE qzt WITH ENCODING='UTF8';

按任意键继续...
"@ -ForegroundColor Gray
$null = $Host.UI.RawUI.ReadKey("NoEcho,IncludeKeyDown")

# 启动服务
Write-Host "`n[4/4] 启动服务..." -ForegroundColor Yellow
Write-Host @"
========================================
服务启动中...
访问地址: http://localhost:8080
健康检查: http://localhost:8080/health
========================================
"@ -ForegroundColor Cyan

go run cmd/server/main.go
