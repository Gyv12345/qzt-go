# API 测试文档

## 环境要求

- 服务已启动（http://localhost:8080）
- PowerShell 3.0+ 或 curl

## 快速测试

### 方式一：PowerShell（推荐）
```powershell
cd D:\ccproject\qzt-go\scripts
.\test-api.ps1
```

### 方式二：批处理
```cmd
cd D:\ccproject\qzt-go\scripts
test-api.bat
```

### 方式三：检查数据库
```powershell
cd D:\ccproject\qzt-go\scripts
go run checkdb.go
```

## 测试接口列表

| 序号 | 接口 | 方法 | 说明 |
|------|------|------|------|
| 1 | /health | GET | 健康检查 |
| 2 | /api/auth/register | POST | 用户注册 |
| 3 | /api/auth/login | POST | 用户登录 |
| 4 | /api/customers | POST | 创建客户 |
| 5 | /api/customers | GET | 获取客户列表 |
| 6 | /api/customers/:id | GET | 获取客户详情 |
| 7 | /api/customers/:id | PUT | 更新客户 |
| 8 | /api/contacts | POST | 创建联系人 |
| 9 | /api/contracts | POST | 创建合同 |
| 10 | /api/apikeys | POST | 创建 API Key |
| 11 | /api/v1/customers | GET | API Key 认证测试 |

## 预期结果

✅ 所有测试应返回：
```json
{
  "code": 0,
  "message": "成功",
  "data": { ... }
}
```

❌ 失败时返回：
```json
{
  "code": 10xxx,
  "message": "错误描述",
  "data": null
}
```
