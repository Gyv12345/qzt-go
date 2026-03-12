#!/usr/bin/env powershell
# 推送到 GitHub 脚本

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "   QZT Go 推送到 GitHub" -ForegroundColor Cyan
Write-Host "========================================`n" -ForegroundColor Cyan

$RepoUrl = Read-Host "请输入 GitHub 仓库地址 (例如: https://github.com/username/qzt-go.git)"

if ($RepoUrl -eq "") {
    Write-Host "❌ 仓库地址不能为空" -ForegroundColor Red
    exit 1
}

Write-Host "`n[1/3] 添加远程仓库..." -ForegroundColor Yellow
& git remote add origin $RepoUrl 2>&1 | Out-Null

Write-Host "[2/3] 推送到 GitHub..." -ForegroundColor Yellow
& git push -u origin master

if ($LASTEXITCODE -eq 0) {
    Write-Host "`n✅ 推送成功！" -ForegroundColor Green
    Write-Host "仓库地址: $RepoUrl" -ForegroundColor Cyan
} else {
    Write-Host "`n❌ 推送失败" -ForegroundColor Red
    Write-Host "请检查仓库地址和网络连接" -ForegroundColor Yellow
}
