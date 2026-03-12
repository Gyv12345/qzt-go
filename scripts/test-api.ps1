#!/usr/bin/env powershell
# API 测试脚本

$BaseUrl = "http://localhost:8080"
$Headers = @{
    "Content-Type" = "application/json"
}

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "QZT API 接口测试" -ForegroundColor Cyan
Write-Host "========================================`n" -ForegroundColor Cyan

# 1. 健康检查
Write-Host "[1] 健康检查" -ForegroundColor Yellow
try {
    $response = Invoke-RestMethod -Uri "$BaseUrl/health" -Method Get
    Write-Host "✅ 服务运行中" -ForegroundColor Green
    Write-Host ($response | ConvertTo-Json) -ForegroundColor Gray
} catch {
    Write-Host "❌ 服务未启动或无法访问" -ForegroundColor Red
    Write-Host "错误: $_" -ForegroundColor Red
    exit 1
}

Write-Host "`n[2] 用户注册" -ForegroundColor Yellow
$registerBody = @{
    username = "testuser"
    password = "test123456"
    realName = "测试用户"
} | ConvertTo-Json

try {
    $response = Invoke-RestMethod -Uri "$BaseUrl/api/auth/register" -Method Post -Body $registerBody -Headers $Headers
    Write-Host "✅ 注册成功" -ForegroundColor Green
    Write-Host ($response | ConvertTo-Json) -ForegroundColor Gray
} catch {
    if ($_.Exception.Response.StatusCode -eq 500) {
        Write-Host "⚠️  用户可能已存在，继续登录测试" -ForegroundColor Yellow
    } else {
        Write-Host "❌ 注册失败: $_" -ForegroundColor Red
    }
}

Write-Host "`n[3] 用户登录" -ForegroundColor Yellow
$loginBody = @{
    username = "testuser"
    password = "test123456"
} | ConvertTo-Json

try {
    $response = Invoke-RestMethod -Uri "$BaseUrl/api/auth/login" -Method Post -Body $loginBody -Headers $Headers
    Write-Host "✅ 登录成功" -ForegroundColor Green
    $token = $response.data.token
    Write-Host "Token: $token" -ForegroundColor Gray
    
    # 更新 Headers
    $Headers["Authorization"] = "Bearer $token"
} catch {
    Write-Host "❌ 登录失败: $_" -ForegroundColor Red
    exit 1
}

Write-Host "`n[4] 创建客户" -ForegroundColor Yellow
$customerBody = @{
    name = "测试公司A"
    shortName = "测试A"
    industry = "IT"
    scale = "100-500人"
    address = "北京市海淀区"
    website = "https://example.com"
    sourceChannel = "官网"
    tags = @("VIP", "重要客户")
    remark = "这是一个测试客户"
} | ConvertTo-Json -Depth 3

try {
    $response = Invoke-RestMethod -Uri "$BaseUrl/api/customers" -Method Post -Body $customerBody -Headers $Headers
    Write-Host "✅ 创建成功" -ForegroundColor Green
    Write-Host ($response | ConvertTo-Json) -ForegroundColor Gray
    $customerId = $response.data.id
} catch {
    Write-Host "❌ 创建失败: $_" -ForegroundColor Red
}

Write-Host "`n[5] 获取客户列表" -ForegroundColor Yellow
try {
    $response = Invoke-RestMethod -Uri "$BaseUrl/api/customers?page=1&pageSize=10" -Method Get -Headers $Headers
    Write-Host "✅ 查询成功" -ForegroundColor Green
    Write-Host "总数: $($response.data.total)" -ForegroundColor Gray
    Write-Host ($response.data.items | ConvertTo-Json -Depth 3) -ForegroundColor Gray
} catch {
    Write-Host "❌ 查询失败: $_" -ForegroundColor Red
}

Write-Host "`n[6] 获取客户详情" -ForegroundColor Yellow
try {
    $response = Invoke-RestMethod -Uri "$BaseUrl/api/customers/$customerId" -Method Get -Headers $Headers
    Write-Host "✅ 查询成功" -ForegroundColor Green
    Write-Host ($response | ConvertTo-Json -Depth 3) -ForegroundColor Gray
} catch {
    Write-Host "❌ 查询失败: $_" -ForegroundColor Red
}

Write-Host "`n[7] 更新客户" -ForegroundColor Yellow
$updateBody = @{
    name = "测试公司B（已更新）"
    remark = "更新后的备注"
} | ConvertTo-Json

try {
    $response = Invoke-RestMethod -Uri "$BaseUrl/api/customers/$customerId" -Method Put -Body $updateBody -Headers $Headers
    Write-Host "✅ 更新成功" -ForegroundColor Green
    Write-Host ($response | ConvertTo-Json -Depth 3) -ForegroundColor Gray
} catch {
    Write-Host "❌ 更新失败: $_" -ForegroundColor Red
}

Write-Host "`n[8] 创建联系人" -ForegroundColor Yellow
$contactBody = @{
    name = "张三"
    position = "技术总监"
    phone = "13800138000"
    email = "zhangsan@example.com"
    customerId = $customerId
    isPrimary = $true
} | ConvertTo-Json

try {
    $response = Invoke-RestMethod -Uri "$BaseUrl/api/contacts" -Method Post -Body $contactBody -Headers $Headers
    Write-Host "✅ 创建成功" -ForegroundColor Green
    Write-Host ($response | ConvertTo-Json) -ForegroundColor Gray
    $contactId = $response.data.id
} catch {
    Write-Host "❌ 创建失败: $_" -ForegroundColor Red
}

Write-Host "`n[9] 创建合同" -ForegroundColor Yellow
$contractBody = @{
    contractNo = "HT-2026-001"
    name = "软件开发合同"
    customerId = $customerId
    amount = 100000
    status = "PENDING"
    startDate = "2026-03-01"
    endDate = "2026-06-30"
} | ConvertTo-Json

try {
    $response = Invoke-RestMethod -Uri "$BaseUrl/api/contracts" -Method Post -Body $contractBody -Headers $Headers
    Write-Host "✅ 创建成功" -ForegroundColor Green
    Write-Host ($response | ConvertTo-Json) -ForegroundColor Gray
    $contractId = $response.data.id
} catch {
    Write-Host "❌ 创建失败: $_" -ForegroundColor Red
}

Write-Host "`n[10] 创建 API Key" -ForegroundColor Yellow
$apiKeyBody = @{
    name = "测试 API Key"
    permissions = @("customer:read", "customer:write")
    rateLimit = 100
} | ConvertTo-Json

try {
    $response = Invoke-RestMethod -Uri "$BaseUrl/api/apikeys" -Method Post -Body $apiKeyBody -Headers $Headers
    Write-Host "✅ 创建成功" -ForegroundColor Green
    Write-Host "API Key: $($response.data.apiKey)" -ForegroundColor Yellow
    Write-Host "⚠️  请保存 API Key，只显示一次！" -ForegroundColor Yellow
    $apiKey = $response.data.apiKey
} catch {
    Write-Host "❌ 创建失败: $_" -ForegroundColor Red
}

Write-Host "`n[11] 使用 API Key 访问（外部接口）" -ForegroundColor Yellow
$apiKeyHeaders = @{
    "Content-Type" = "application/json"
    "X-API-Key" = $apiKey
}

try {
    $response = Invoke-RestMethod -Uri "$BaseUrl/api/v1/customers?page=1&pageSize=5" -Method Get -Headers $apiKeyHeaders
    Write-Host "✅ API Key 认证成功" -ForegroundColor Green
    Write-Host "查询到 $($response.data.total) 个客户" -ForegroundColor Gray
} catch {
    Write-Host "❌ API Key 认证失败: $_" -ForegroundColor Red
}

Write-Host "`n========================================" -ForegroundColor Cyan
Write-Host "测试完成！" -ForegroundColor Green
Write-Host "========================================" -ForegroundColor Cyan
