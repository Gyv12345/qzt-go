# 数据库设置指南

## PostgreSQL（推荐）

### 1. 创建数据库
```sql
-- 连接 PostgreSQL
psql -U postgres

-- 创建数据库
CREATE DATABASE qzt WITH ENCODING='UTF8';

-- 连接到数据库
\c qzt

-- 验证
SELECT current_database();
```

### 2. 运行服务
```bash
cd D:\ccproject\qzt-go
go mod tidy
go run cmd/server/main.go
```

服务会自动创建表结构（AutoMigrate）。

## Redis

### 1. 确认 Redis 运行
```bash
redis-cli ping
# 应该返回 PONG
```

### 2. 默认配置
- 地址：localhost:6379
- 密码：无
- 数据库：0

## 初始化数据

### 创建管理员账号
```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123","realName":"管理员"}'
```

### 登录测试
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}'
```
