# QZT API 文档

> 企智通 CRM 系统 API 接口文档

## 基础信息

- **Base URL**: `http://localhost:8080`
- **认证方式**: JWT Token / API Key
- **数据格式**: JSON
- **编码**: UTF-8

## 认证

### JWT 认证（内部使用）

在请求头中添加:
```
Authorization: Bearer <token>
```

### API Key 认证（第三方调用）

方式1 - Authorization Header:
```
Authorization: Bearer qzt_live_xxxxxxxxxxxxx
```

方式2 - X-API-Key Header:
```
X-API-Key: qzt_live_xxxxxxxxxxxxx
```

## 通用响应

### 成功响应

```json
{
  "id": "1",
  "name": "示例数据",
  "createdAt": "2026-03-11T10:00:00Z"
}
```

### 错误响应

```json
{
  "error": "错误信息描述"
}
```

### 分页响应

```json
{
  "items": [...],
  "total": 100,
  "page": 1,
  "pageSize": 10
}
```

## API 接口

### 1. 认证模块

#### 1.1 用户登录

**POST** `/api/auth/login`

**请求体**:
```json
{
  "username": "admin",
  "password": "123456"
}
```

**响应**:
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refreshToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "expiresAt": 1234567890,
  "user": {
    "id": "1",
    "username": "admin",
    "realName": "管理员",
    "email": "admin@example.com",
    "mobile": "13800138000"
  }
}
```

**状态码**:
- 200: 成功
- 401: 用户名或密码错误

---

#### 1.2 用户注册

**POST** `/api/auth/register`

**请求体**:
```json
{
  "username": "newuser",
  "password": "123456",
  "email": "user@example.com",
  "mobile": "13800138000",
  "realName": "新用户"
}
```

**响应**:
```json
{
  "id": "2",
  "username": "newuser",
  "email": "user@example.com",
  "mobile": "13800138000",
  "realName": "新用户",
  "status": 1,
  "createdAt": "2026-03-11T10:00:00Z"
}
```

**状态码**:
- 201: 创建成功
- 400: 参数错误或用户名已存在

---

#### 1.3 刷新 Token

**POST** `/api/auth/refresh`

**请求体**:
```json
{
  "refreshToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**响应**: 同登录接口

---

### 2. 客户管理

#### 2.1 获取客户列表

**GET** `/api/customers`

**认证**: 需要 JWT Token

**查询参数**:
- `page` (int): 页码，默认 1
- `pageSize` (int): 每页数量，默认 10
- `keyword` (string): 搜索关键词（公司名称或简称）
- `level` (string): 客户等级 (LEAD/PROSPECT/CUSTOMER/VIP)
- `status` (int): 状态 (1: 启用, 0: 禁用)

**响应**:
```json
{
  "items": [
    {
      "id": "1",
      "name": "河南爱编程网络科技有限公司",
      "shortName": "爱编程",
      "industry": "软件开发",
      "scale": "10-50人",
      "customerLevel": "CUSTOMER",
      "status": 1,
      "followUserId": "1",
      "followUser": {
        "id": "1",
        "realName": "张三"
      },
      "createdAt": "2026-03-11T10:00:00Z"
    }
  ],
  "total": 100,
  "page": 1,
  "pageSize": 10
}
```

---

#### 2.2 创建客户

**POST** `/api/customers`

**认证**: 需要 JWT Token

**请求体**:
```json
{
  "name": "河南爱编程网络科技有限公司",
  "shortName": "爱编程",
  "industry": "软件开发",
  "scale": "10-50人",
  "address": "郑州市金水区",
  "website": "https://www.example.com",
  "sourceChannel": "官网",
  "tags": ["AI", "CRM"],
  "remark": "重要客户",
  "followUserId": "1"
}
```

**响应**:
```json
{
  "id": "1",
  "name": "河南爱编程网络科技有限公司",
  "shortName": "爱编程",
  "industry": "软件开发",
  "customerLevel": "LEAD",
  "status": 1,
  "createdAt": "2026-03-11T10:00:00Z"
}
```

**状态码**:
- 201: 创建成功
- 400: 参数错误

---

#### 2.3 获取客户详情

**GET** `/api/customers/:id`

**认证**: 需要 JWT Token

**响应**:
```json
{
  "id": "1",
  "name": "河南爱编程网络科技有限公司",
  "shortName": "爱编程",
  "industry": "软件开发",
  "scale": "10-50人",
  "address": "郑州市金水区",
  "website": "https://www.example.com",
  "customerLevel": "CUSTOMER",
  "sourceChannel": "官网",
  "followUserId": "1",
  "tags": "[\"AI\",\"CRM\"]",
  "remark": "重要客户",
  "status": 1,
  "followUser": {
    "id": "1",
    "realName": "张三"
  },
  "contacts": [
    {
      "id": "1",
      "name": "李四",
      "position": "技术总监"
    }
  ],
  "createdAt": "2026-03-11T10:00:00Z"
}
```

**状态码**:
- 200: 成功
- 404: 客户不存在

---

#### 2.4 更新客户

**PUT** `/api/customers/:id`

**认证**: 需要 JWT Token

**请求体**: 所有字段可选
```json
{
  "name": "新公司名称",
  "customerLevel": "VIP",
  "status": 1
}
```

**响应**: 同客户详情

---

#### 2.5 删除客户

**DELETE** `/api/customers/:id`

**认证**: 需要 JWT Token

**状态码**:
- 204: 删除成功
- 500: 删除失败

---

### 3. 联系人管理

#### 3.1 获取联系人列表

**GET** `/api/contacts`

**认证**: 需要 JWT Token

**查询参数**:
- `page` (int): 页码
- `pageSize` (int): 每页数量
- `customerId` (string): 客户 ID
- `keyword` (string): 搜索关键词

**响应**:
```json
{
  "items": [
    {
      "id": "1",
      "customerId": "1",
      "name": "李四",
      "position": "技术总监",
      "department": "技术部",
      "phone": "0371-12345678",
      "mobile": "13800138000",
      "email": "lisi@example.com",
      "isPrimary": true,
      "isDecision": true,
      "status": 1,
      "customer": {
        "id": "1",
        "name": "河南爱编程网络科技有限公司"
      }
    }
  ],
  "total": 50,
  "page": 1,
  "pageSize": 10
}
```

---

#### 3.2 创建联系人

**POST** `/api/contacts`

**认证**: 需要 JWT Token

**请求体**:
```json
{
  "customerId": "1",
  "name": "李四",
  "position": "技术总监",
  "department": "技术部",
  "phone": "0371-12345678",
  "mobile": "13800138000",
  "email": "lisi@example.com",
  "wechat": "lisi_wx",
  "isPrimary": true,
  "isDecision": true,
  "gender": "MALE",
  "address": "郑州市",
  "remark": "技术负责人"
}
```

---

#### 3.3 更新联系人

**PUT** `/api/contacts/:id`

**认证**: 需要 JWT Token

---

#### 3.4 删除联系人

**DELETE** `/api/contacts/:id`

**认证**: 需要 JWT Token

---

### 4. 合同管理

#### 4.1 获取合同列表

**GET** `/api/contracts`

**认证**: 需要 JWT Token

**查询参数**:
- `page` (int): 页码
- `pageSize` (int): 每页数量
- `customerId` (string): 客户 ID
- `status` (string): 合同状态 (DRAFT/PENDING/ACTIVE/COMPLETED/CANCELLED)
- `keyword` (string): 搜索关键词

**响应**:
```json
{
  "items": [
    {
      "id": "1",
      "customerId": "1",
      "contractNo": "HT20260311001",
      "title": "CRM 系统开发合同",
      "type": "软件开发",
      "amount": 100000.00,
      "paidAmount": 50000.00,
      "status": "ACTIVE",
      "startDate": "2026-03-01",
      "endDate": "2026-12-31",
      "ownerUserId": "1",
      "customer": {
        "id": "1",
        "name": "河南爱编程网络科技有限公司"
      },
      "owner": {
        "id": "1",
        "realName": "张三"
      }
    }
  ],
  "total": 30,
  "page": 1,
  "pageSize": 10
}
```

---

#### 4.2 创建合同

**POST** `/api/contracts`

**认证**: 需要 JWT Token

**请求体**:
```json
{
  "customerId": "1",
  "title": "CRM 系统开发合同",
  "type": "软件开发",
  "amount": 100000.00,
  "ownerUserId": "1",
  "tags": ["重要", "年度"],
  "remark": "年度框架合同"
}
```

---

#### 4.3 更新合同

**PUT** `/api/contracts/:id`

**认证**: 需要 JWT Token

**请求体示例**:
```json
{
  "status": "ACTIVE",
  "paidAmount": 80000.00
}
```

---

### 5. API Key 管理

#### 5.1 创建 API Key

**POST** `/api/apikeys`

**认证**: 需要 JWT Token

**请求体**:
```json
{
  "name": "OpenClaw 集成",
  "userId": "1",
  "permissions": ["customer:read", "contact:read", "customer:write"],
  "ipWhitelist": ["192.168.1.0/24"],
  "rateLimit": 100
}
```

**响应**:
```json
{
  "apiKey": "qzt_live_abc123def456ghi789jkl012",  // ⚠️ 只显示一次，请保存
  "info": {
    "id": "1",
    "name": "OpenClaw 集成",
    "prefix": "qzt_live_ab",
    "userId": "1",
    "permissions": "[\"customer:read\",\"contact:read\",\"customer:write\"]",
    "rateLimit": 100,
    "status": 1,
    "createdAt": "2026-03-11T10:00:00Z"
  }
}
```

**权限列表**:
- `customer:read`: 读取客户信息
- `customer:write`: 创建/更新/删除客户
- `contact:read`: 读取联系人信息
- `contact:write`: 创建/更新/删除联系人
- `contract:read`: 读取合同信息
- `contract:write`: 创建/更新/删除合同
- `*`: 所有权限

---

#### 5.2 获取 API Key 列表

**GET** `/api/apikeys`

**认证**: 需要 JWT Token

**响应**:
```json
[
  {
    "id": "1",
    "name": "OpenClaw 集成",
    "prefix": "qzt_live_ab",
    "userId": "1",
    "status": 1,
    "lastUsedAt": "2026-03-11T10:00:00Z",
    "createdAt": "2026-03-10T10:00:00Z"
  }
]
```

**注意**: 不返回完整的 API Key，只返回前缀用于识别。

---

#### 5.3 撤销 API Key

**POST** `/api/apikeys/:id/revoke`

**认证**: 需要 JWT Token

**状态码**:
- 204: 撤销成功

---

### 6. API Key 认证接口

所有内部 API 接口都可以通过 API Key 访问，路径前缀为 `/api/v1`。

**示例**:

```bash
# 使用 API Key 获取客户列表
curl -X GET "http://localhost:8080/api/v1/customers" \
  -H "Authorization: Bearer qzt_live_abc123def456ghi789jkl012"

# 使用 API Key 创建客户
curl -X POST "http://localhost:8080/api/v1/customers" \
  -H "Authorization: Bearer qzt_live_abc123def456ghi789jkl012" \
  -H "Content-Type: application/json" \
  -d '{"name":"测试公司"}'
```

## 错误码

| 状态码 | 说明 |
|--------|------|
| 200 | 成功 |
| 201 | 创建成功 |
| 204 | 删除成功（无内容） |
| 400 | 请求参数错误 |
| 401 | 未认证或认证失败 |
| 403 | 权限不足 |
| 404 | 资源不存在 |
| 500 | 服务器内部错误 |

## 数据字典

### 客户等级 (CustomerLevel)
- `LEAD`: 线索
- `PROSPECT`: 潜在客户
- `CUSTOMER`: 客户
- `VIP`: VIP 客户

### 合同状态 (ContractStatus)
- `DRAFT`: 草稿
- `PENDING`: 待审批
- `ACTIVE`: 执行中
- `COMPLETED`: 已完成
- `CANCELLED`: 已取消

### 用户状态 (UserStatus)
- `0`: 禁用
- `1`: 启用

## 最佳实践

### 1. Token 管理
- 登录后保存 `token` 和 `refreshToken`
- Token 过期前使用 `refreshToken` 刷新
- Token 存储在 localStorage 或 cookie 中

### 2. API Key 管理
- 创建后立即保存，系统不会再次显示
- 为不同应用创建不同的 API Key
- 定期轮换 API Key
- 使用 IP 白名单限制访问

### 3. 分页查询
- 默认 page=1, pageSize=10
- 最大 pageSize=100
- 通过 total 字段计算总页数

### 4. 错误处理
- 检查 HTTP 状态码
- 解析 error 字段获取错误信息
- 401 错误需重新登录

## 联系方式

- 技术支持: support@qzt.com
- 问题反馈: https://github.com/qzt/backend/issues
