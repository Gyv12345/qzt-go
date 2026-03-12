# QZT Go 前后端匹配状态

## ✅ 已修复的问题

### 1. API 路径匹配
- ✅ Vite 代理配置已更新：`/api` → `http://localhost:8080`
- ✅ 后端路由：`/api/auth/login`, `/api/customers` 等

### 2. 新增接口
- ✅ `GET /api/auth/me` - 获取当前用户信息（已添加）

### 3. 响应格式
- 前端期望：`{ success, data, message }`
- 后端返回：`{ code, message, data }`
- ⚠️ 需要在前端适配或后端修改

## 📊 API 接口对照表

| 前端路径 | 后端路由 | 状态 |
|---------|---------|------|
| POST /api/auth/login | POST /api/auth/login | ✅ |
| GET /api/auth/me | GET /api/auth/me | ✅ |
| POST /api/auth/register | POST /api/auth/register | ✅ |
| GET /api/customers | GET /api/customers | ✅ |
| POST /api/customers | POST /api/customers | ✅ |
| GET /api/customers/:id | GET /api/customers/:id | ✅ |
| PUT /api/customers/:id | PUT /api/customers/:id | ✅ |
| DELETE /api/customers/:id | DELETE /api/customers/:id | ✅ |
| GET /api/contacts | GET /api/contacts | ✅ |
| POST /api/contacts | POST /api/contacts | ✅ |
| GET /api/contracts | GET /api/contracts | ✅ |
| POST /api/contracts | POST /api/contracts | ✅ |

## 🚀 启动测试

### 1. 启动后端
```powershell
# 新开一个 PowerShell 窗口
cd D:\ccproject\qzt-go
go run cmd/server/main.go
```

### 2. 启动前端
```powershell
cd D:\ccproject\qzt-go\frontend
pnpm dev
```

### 3. 测试登录
```bash
# 注册
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"test123456","realName":"测试用户"}'

# 登录
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"test123456"}'

# 获取用户信息（需要 token）
curl http://localhost:8080/api/auth/me \
  -H "Authorization: Bearer <token>"
```

## ⚠️ 待处理

1. **响应格式统一** - 前端期望 `success` 字段，后端返回 `code`
2. **完整业务逻辑** - 部分高级功能（批量操作、高级搜索等）需要完善
3. **权限控制** - 细粒度权限需要实现

## ✅ 已完成

1. ✅ Vite 代理配置修复
2. ✅ 添加 GET /api/auth/me 接口
3. ✅ 前后端路由匹配
4. ✅ PostgreSQL + Redis 配置
5. ✅ JWT 认证
6. ✅ 基础 CRUD 操作
