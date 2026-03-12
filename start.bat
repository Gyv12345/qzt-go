@echo off
chcp 65001 >nul
echo ========================================
echo QZT Go 项目初始化脚本
echo ========================================
echo.

REM 设置颜色
color 0A

echo [1/4] 检查 Go 环境...
where go >nul 2>&1
if %errorlevel% neq 0 (
    echo [错误] 未找到 Go，请确保 Go 已安装并添加到 PATH
    echo 下载地址: https://golang.org/dl/
    pause
    exit /b 1
)
go version
echo.

echo [2/4] 下载依赖...
go mod tidy
if %errorlevel% neq 0 (
    echo [错误] 依赖下载失败
    pause
    exit /b 1
)
echo 依赖下载完成！
echo.

echo [3/4] 检查 PostgreSQL...
echo 请确保已创建数据库: qzt
echo 如果未创建，请运行:
echo   psql -U postgres
echo   CREATE DATABASE qzt WITH ENCODING='UTF8';
echo.
pause

echo [4/4] 启动服务...
echo.
echo ========================================
echo 服务启动中...
echo 访问地址: http://localhost:8080
echo 健康检查: http://localhost:8080/health
echo ========================================
echo.

go run cmd/server/main.go

pause
