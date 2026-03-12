#!/usr/bin/env powershell
# 快速启动 QZT 全栈应用

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "   QZT Go 全栈应用启动" -ForegroundColor Cyan
Write-Host "========================================` -ForegroundColor Cyan

# 检查并停止占用端口的进程
Write-Host "`n[1/3] 清理端口占用..." -ForegroundColor Yellow
$port8080 = Get-NetTCPConnection -LocalPort 8080 -ErrorAction SilentlyContinue
if ($port8080) {
    Write-Host "停止占用端口 8080 的进程..." -ForegroundColor Yellow
    Stop-Process -Id $port8080.OwningProcess -Force -ErrorAction SilentlyContinue
    Start-Sleep -Seconds 2
}

# 启动后端
Write-Host "`n[2/3] 启动 Go 后端..." -ForegroundColor Yellow
$backendScript = @"
$env:Path = "C:\Program Files\Go\bin;$C:\Go\bin;$$env:Path
Set-Location D:\ccproject\qzt-go

Write-Host "下载依赖..." -ForegroundColor Gray
& go mod tidy 2>&1 | Out-Null

Write-Host "启动后端..." -ForegroundColor Gray
Start-Process go -FilePath go run cmd/server/main.go -WindowStyle Hidden -ErrorAction SilentlyContinue

Write-Host "等待后端启动..." -ForegroundColor Yellow
Start-Sleep -Seconds 5

# 检查后端
Write-Host "`n[2/3] 测试后端连接..." -ForegroundColor Yellow
$attempts = 0
$maxAttempts = 3
try {
    $health = Invoke-RestMethod -Uri "http://localhost:8080/health" -Method Get -TimeoutSec 5
    Write-Host "✅ 后端服务运行正常!" -ForegroundColor Green
    Write-Host ($health | ConvertTo-Json) -ForegroundColor Green
} catch {
    Write-Host "⚠️ 尝试 $i 失败，重试..." -ForegroundColor Yellow
    Start-Sleep -Seconds 2
}

# 启动前端
Write-Host "`n[2/3] 检查前端状态..." -ForegroundColor Yellow
try {
    $frontend = Invoke-RestMethod -Uri "http://localhost:3458" -Method Get -TimeoutSec 5
    Write-Host "✅ 前端服务运行正常!" -ForegroundColor Green
} catch {
    Write-Host "⚠️ 前端未响应" -ForegroundColor Yellow
}

# 创建测试用户
Write-Host "`n[1/3] 创建测试用户..." -ForegroundColor Yellow
$registerBody = @{
    username = "testuser"
    password = "test123456"
    realName = "测试用户"
} | ConvertTo-Json

try {
    $register = Invoke-RestMethod -Uri "http://localhost:8080/api/auth/register" -Method Post -Body $registerBody -ContentType "application/json" -TimeoutSec 5
    Write-Host "✅ 测试用户注册成功!" -ForegroundColor Green
} catch {
    Write-Host "⚠️ 测试用户可能已存在" -ForegroundColor Yellow
}

# 登录测试
Write-Host "`n[1/3] 登录测试..." -ForegroundColor Yellow
$loginBody = @{
    username = "testuser"
    password = "test123456"
} | ConvertTo-Json

try {
    $login = Invoke-RestMethod -Uri "http://localhost:8080/api/auth/login" -Method Post -Body $loginBody -ContentType "application/json" -TimeoutSec 5
    Write-Host "✅ 登录成功!" -ForegroundColor Green
    $token = $login.data.token
    Write-Host "Token 已获取" -ForegroundColor Gray
} catch {
    Write-Host "❌ 登录失败: $_" -ForegroundColor Red
    exit 1
}

Write-Host @"

========================================
   ✅ QZT Go 全栈应用启动完成！ 
========================================

📌 服务地址:
  - 前端: http://localhost:3458
  - 后端: http://localhost:8080
  - 健康检查: http://localhost:8080/health

📌 测试账号:
  - 用户名: testuser
  - 密码: test123456

📌 Page Agent Demo
  - 已加载测试脚本
  - 可在前端页面进行全访问测试

📌 停止服务
  - 关闭此窗口

========================================

"@ -ForegroundColor Green

# 保持窗口打开
Read-Host "`n按任意键停止服务..." -ForegroundColor Yellow
$null = $Host.UI.RawUI.ReadKey("NoEcho,IncludeKeyDown")
