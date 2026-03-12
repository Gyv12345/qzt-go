# QZT Go 快速启动指南

## 方式一：使用脚本（推荐）

### Windows
双击运行 `start.bat` 或在 PowerShell 中运行：
```powershell
.\start.ps1
```

## 方式二：手动启动

### 1. 创建数据库
打开 pgAdmin 或命令行，执行：
```sql
CREATE DATABASE qzt WITH ENCODING='UTF8';
```

### 2. 下载依赖
```bash
go mod tidy
```

### 3. 启动服务
```bash
go run cmd/server/main.go
```

### 4. 测试
访问 http://localhost:8080/health

## 环境要求

- ✅ Go 1.22+
- ✅ PostgreSQL（默认配置）
- ✅ Redis（密码: 123456）

## 常见问题

### Q: 数据库连接失败
A: 检查 PostgreSQL 是否运行，确认数据库 `qzt` 已创建

### Q: Redis 连接失败
A: 检查 Redis 是否运行，确认密码是 `123456`

### Q: 端口被占用
A: 修改 `config/config.yaml` 中的 `server.port`

## 测试账号

注册管理员：
```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d "{\"username\":\"admin\",\"password\":\"admin123\",\"realName\":\"管理员\"}"
```

登录：
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d "{\"username\":\"admin\",\"password\":\"admin123\"}"
```
