@echo off
chcp 65001 >nul
setlocal enabledelayedexpansion

echo ========================================
echo QZT API 测试脚本
echo ========================================
echo.

set BASE_URL=http://localhost:8080

echo [1] 健康检查
curl -s %BASE_URL%/health
echo.
echo.

echo [2] 用户注册
curl -s -X POST %BASE_URL%/api/auth/register ^
  -H "Content-Type: application/json" ^
  -d "{\"username\":\"testuser\",\"password\":\"test123456\",\"realName\":\"测试用户\"}"
echo.
echo.

echo [3] 用户登录（获取 Token）
for /f "tokens=*" %%i in ('curl -s -X POST %BASE_URL%/api/auth/login -H "Content-Type: application/json" -d "{\"username\":\"testuser\",\"password\":\"test123456\"} ^| findstr /r "token"') do (
    set TOKEN_LINE=%%i
)
echo %TOKEN_LINE%
echo.

echo 请手动从上面的响应中复制 token，然后继续测试...
pause

echo.
echo [4] 创建客户（需要 Token）
set /p TOKEN=请输入 Token: 
echo.
curl -s -X POST %BASE_URL%/api/customers ^
  -H "Content-Type: application/json" ^
  -H "Authorization: Bearer %TOKEN%" ^
  -d "{\"name\":\"测试公司A\",\"shortName\":\"测试A\",\"industry\":\"IT\",\"scale\":\"100-500人\"}"
echo.
echo.

echo [5] 获取客户列表
curl -s %BASE_URL%/api/customers?page=1^&pageSize=10 ^
  -H "Authorization: Bearer %TOKEN%"
echo.
echo.

echo ========================================
echo 测试完成！
echo ========================================
pause
