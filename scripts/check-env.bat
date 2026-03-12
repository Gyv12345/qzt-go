@echo off
chcp 65001 >nul
echo ========================================
echo    QZT 环境检查
echo ========================================
echo.

echo [检查] Go 环境
where go >nul 2>&1
if %errorlevel% neq 0 (
    echo ❌ Go 未安装
    echo.
    echo 请先安装 Go:
    echo 1. 访问 https://golang.org/dl/
    echo 2. 下载并安装 Go 1.22+
    echo 3. 重启命令行
    echo 4. 重新运行此脚本
    pause
    exit /b 1
)
go version
echo ✅ Go 已安装
echo.

echo [检查] PostgreSQL
psql --version >nul 2>&1
if %errorlevel% neq 0 (
    echo ⚠️  psql 不在 PATH 中，但可能已安装
) else (
    echo ✅ PostgreSQL 已安装
)
echo.

echo [检查] Redis
redis-cli -a 123456 ping >nul 2>&1
if %errorlevel% neq 0 (
    echo ⚠️  Redis 未运行或无法连接
    echo    请确保 Redis 运行在 localhost:6379，密码: 123456
) else (
    echo ✅ Redis 连接成功
)
echo.

echo ========================================
echo    环境检查完成
echo ========================================
echo.
echo 如需运行完整测试，请使用 PowerShell:
echo   .\scripts\full-test.ps1
echo.
pause
