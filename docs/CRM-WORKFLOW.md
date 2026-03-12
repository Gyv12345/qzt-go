# 完整的 CRM 销售流程功能

## ✅ 已实现功能

### 1. 商机管理（Opportunity）

#### 商机阶段
```
线索 (LEAD) → 资格审查 (QUALIFY) → 方案报价 (PROPOSAL) → 商务谈判 (NEGOTIATE) → 签订合同 (CONTRACT) → 赢单 (CLOSED_WON)
```

#### 每个阶段的成功概率
- 线索：10%
- 资格审查：20%
- 方案报价：40%
- 商务谈判：60%
- 签订合同：80%
- 赢单：100%

#### API 接口
```
GET    /api/opportunities                    商机列表
POST   /api/opportunities                    创建商机
GET    /api/opportunities/:id                商机详情
PUT    /api/opportunities/:id                更新商机
DELETE /api/opportunities/:id                删除商机
POST   /api/opportunities/:id/next-stage     推进到下一阶段
POST   /api/opportunities/:id/prev-stage     回退到上一阶段
POST   /api/opportunities/:id/convert        转化为合同
GET    /api/opportunities/stats/stage        阶段统计（销售漏斗）
```

#### 使用示例
```bash
# 1. 创建商机
curl -X POST http://localhost:8080/api/opportunities \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "某公司软件采购",
    "customerId": "cust_123",
    "amount": 100000,
    "expectedCloseAt": "2026-04-30",
    "description": "大型企业软件采购项目"
  }'

# 2. 推进到下一阶段
curl -X POST http://localhost:8080/api/opportunities/opp_123/next-stage \
  -H "Authorization: Bearer <token>"

# 3. 转化为合同
curl -X POST http://localhost:8080/api/opportunities/opp_123/convert \
  -H "Authorization: Bearer <token>"

# 4. 查看销售漏斗
curl http://localhost:8080/api/opportunities/stats/stage \
  -H "Authorization: Bearer <token>"
```

---

### 2. 回款管理（Payment）

#### API 接口
```
GET    /api/payments                    回款列表
POST   /api/payments                    创建回款
POST   /api/payments/:id/confirm        确认回款
```

#### 使用示例
```bash
# 1. 创建回款记录
curl -X POST http://localhost:8080/api/payments \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "contractId": "contract_123",
    "customerId": "cust_123",
    "amount": 50000,
    "paymentDate": "2026-03-15",
    "paymentType": "银行转账",
    "remark": "第一期款项"
  }'

# 2. 确认回款
curl -X POST http://localhost:8080/api/payments/payment_123/confirm \
  -H "Authorization: Bearer <token>"
```

---

### 3. 客户跟进记录（FollowRecord）

#### API 接口
```
GET    /api/follow-records                    跟进记录列表
POST   /api/follow-records                    创建跟进记录
GET    /api/follow-records/:id                跟进记录详情
PUT    /api/follow-records/:id                更新跟进记录
DELETE /api/follow-records/:id                删除跟进记录
```

#### 使用示例
```bash
# 创建跟进记录
curl -X POST http://localhost:8080/api/follow-records \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "customerId": "cust_123",
    "type": "电话",
    "content": "客户对产品很感兴趣，计划下周安排演示",
    "nextFollowAt": "2026-03-20",
    "status": "FOLLOWING"
  }'
```

---

## 🎯 完整销售流程

```
1. 客户线索
   ↓
2. 创建跟进记录（电话/邮件/拜访）
   ↓
3. 创建商机（LEAD 阶段）
   ↓
4. 推进阶段（QUALIFY → PROPOSAL → NEGOTIATE → CONTRACT）
   ↓
5. 转化为合同
   ↓
6. 创建回款记录
   ↓
7. 确认回款
   ↓
8. 商机关闭（CLOSED_WON）
```

---

## 📊 数据模型

### 商机（Opportunity）
```json
{
  "id": "opp_123",
  "name": "某公司软件采购",
  "customerId": "cust_123",
  "userId": "user_456",
  "amount": 100000,
  "stage": "PROPOSAL",
  "probability": 40,
  "expectedCloseAt": "2026-04-30",
  "description": "...",
  "status": "ACTIVE",
  "createdAt": "2026-03-12T10:00:00Z"
}
```

### 回款（Payment）
```json
{
  "id": "payment_123",
  "contractId": "contract_123",
  "customerId": "cust_123",
  "amount": 50000,
  "paymentDate": "2026-03-15",
  "paymentType": "银行转账",
  "status": "CONFIRMED",
  "remark": "第一期款项"
}
```

---

## 🚀 测试流程

1. **创建客户**
```bash
curl -X POST http://localhost:8080/api/customers \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"name":"测试公司","industry":"IT"}'
```

2. **创建商机**
```bash
curl -X POST http://localhost:8080/api/opportunities \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "name":"测试商机",
    "customerId":"<customer_id>",
    "amount":100000,
    "expectedCloseAt":"2026-04-30"
  }'
```

3. **推进阶段**（重复执行直到 CONTRACT）
```bash
curl -X POST http://localhost:8080/api/opportunities/<id>/next-stage \
  -H "Authorization: Bearer <token>"
```

4. **转化为合同**
```bash
curl -X POST http://localhost:8080/api/opportunities/<id>/convert \
  -H "Authorization: Bearer <token>"
```

5. **创建回款**
```bash
curl -X POST http://localhost:8080/api/payments \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "contractId":"<contract_id>",
    "customerId":"<customer_id>",
    "amount":50000,
    "paymentDate":"2026-03-15"
  }'
```

---

## ✅ 功能完整性

- ✅ 客户管理（CRUD）
- ✅ 联系人管理（CRUD）
- ✅ 合同管理（CRUD）
- ✅ 跟进记录（CRUD）
- ✅ 商机管理（CRUD + 阶段流转）
- ✅ 回款管理（CRUD + 确认）
- ✅ 销售漏斗统计
- ✅ RBAC 权限系统
- ✅ 权限中间件

---

## 📈 下一步开发

- [ ] 合同审批流程
- [ ] 发票管理
- [ ] 产品管理
- [ ] 报表统计
- [ ] 数据导入导出
- [ ] 消息通知
