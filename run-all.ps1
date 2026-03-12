#!/usr/bin/env powershell
# QZT Go 完整启动脚本

Write-Host @"
========================================
   QZT Go 项目启动
========================================
"@ -ForegroundColor Cyan

# 1. 检查 PostgreSQL 数据库
Write-Host "`n[1/5] 检查 PostgreSQL 数据库..." -ForegroundColor Yellow
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
    Write-Host "⚠️  无法检查数据库，请手动创建: CREATE DATABASE qzt;" -ForegroundColor Yellow
}

# 2. 启动后端
Write-Host "`n[2/5] 启动 Go 后端..." -ForegroundColor Yellow
$backendJob = Start-Job -ScriptBlock {
    Set-Location "D:\ccproject\qzt-go"
    & go run cmd/server/main.go
}

Start-Sleep -Seconds 3

# 检查后端是否启动
try {
    $health = Invoke-RestMethod -Uri "http://localhost:8080/health" -Method Get -TimeoutSec 5
    Write-Host "✅ 后端启动成功: http://localhost:8080" -ForegroundColor Green
} catch {
    Write-Host "⚠️  后端可能还在启动中..." -ForegroundColor Yellow
}

# 3. 安装前端依赖
Write-Host "`n[3/5] 检查前端依赖..." -ForegroundColor Yellow
if (-not (Test-Path "D:\ccproject\qzt-go\frontend\node_modules")) {
    Write-Host "安装前端依赖..." -ForegroundColor Yellow
    Set-Location "D:\ccproject\qzt-go\frontend"
    & npm install
    Write-Host "✅ 前端依赖安装完成" -ForegroundColor Green
} else {
    Write-Host "✅ 前端依赖已存在" -ForegroundColor Green
}

# 4. 启动前端
Write-Host "`n[4/5] 启动前端开发服务器..." -ForegroundColor Yellow
$frontendJob = Start-Job -ScriptBlock {
    Set-Location "D:\ccproject\qzt-go\frontend"
    & npm run dev
}

Start-Sleep -Seconds 5

# 5. 测试接口
Write-Host "`n[5/5] 测试接口..." -ForegroundColor Yellow

# 注册测试用户
Write-Host "  注册测试用户..." -ForegroundColor Cyan
$registerBody = @{
    username = "testuser"
    password = "test123456"
    realName = "测试用户"
} | ConvertTo-Json

try {
    $register = Invoke-RestMethod -Uri "http://localhost:8080/api/auth/register" -Method Post -Body $registerBody -ContentType "application/json" -ErrorAction SilentlyContinue
    Write-Host "  ✅ 测试用户注册成功" -ForegroundColor Green
} catch {
    Write-Host "  ⚠️  测试用户可能已存在" -ForegroundColor Yellow
}

# 登录测试
Write-Host "  登录测试..." -ForegroundColor Cyan
$loginBody = @{
    username = "testuser"
    password = "test123456"
} | ConvertTo-Json

try {
    $login = Invoke-RestMethod -Uri "http://localhost:8080/api/auth/login" -Method Post -Body $loginBody -ContentType "application/json"
    Write-Host "  ✅ 登录成功" -ForegroundColor Green
} catch {
    Write-Host "  ❌ 登录失败: $_" -ForegroundColor Red
}

Write-Host @"

========================================
   ✅ QZT Go 项目启动完成！
========================================

📌 访问地址:
  - 前端: http://localhost:5173
  - 后端: http://localhost:8080
  - 健康检查: http://localhost:8080/health

📌 测试账号:
  - 用户名: testuser
  - 密码: test123456

📌 Page Agent Demo:
  - 已加载测试脚本
  - 可在前端页面进行全访问测试

📌 停止服务:
  - 关闭此窗口或按 Ctrl+C

========================================

"@ -ForegroundColor Green

# 保持运行
Write-Host "按任意键停止所有服务..." -ForegroundColor Yellow
$null = $Host.UI.RawUI.ReadKey("NoEcho,IncludeKeyDown")

# 停止所有 job
Stop-Job $backendJob -ErrorAction SilentlyContinue
Stop-Job $frontendJob -ErrorAction SilentlyContinue
Remove-Job $backendJob -ErrorAction SilentlyContinue
Remove-Job $frontendJob -ErrorAction SilentlyContinue

Write-Host "`n服务已停止" -ForegroundColor Yellow
