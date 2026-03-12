# QZT 测试报告

## 📋 测试清单

### 环境检查
- [ ] Go 1.22+ 已安装
- [ ] PostgreSQL 已安装并运行
- [ ] Redis 已安装并运行（密码: 123456）

### 数据库准备
- [ ] 数据库 `qzt` 已创建

### 运行测试

#### 方式一：完整自动化测试（推荐）
```powershell
cd D:\ccproject\qzt-go
.\scripts\full-test.ps1
```

**测试内容：**
1. ✅ 环境检查（Go/PostgreSQL/Redis）
2. ✅ 依赖下载
3. ✅ 数据库创建
4. ✅ 服务启动
5. ✅ 11 个接口自动化测试

#### 方式二：分步测试
```powershell
# 1. 检查环境
.\scripts\check-env.bat

# 2. 启动服务
.\start.ps1

# 3. 测试接口（新终端）
.\scripts\test-api.ps1
```

## 📊 预期测试结果

所有接口应返回：
```json
{
  "code": 0,
  "message": "成功",
  "data": { ... }
}
```

## ⚠️ 常见问题

### 1. Go 未找到
**解决：** 安装 Go 并添加到 PATH
- 下载：https://golang.org/dl/
- 安装后重启 PowerShell

### 2. PostgreSQL 连接失败
**解决：**
```sql
-- 创建数据库
CREATE DATABASE qzt WITH ENCODING='UTF8';

-- 检查连接
psql -U postgres -d qzt
```

### 3. Redis 连接失败
**解决：**
```bash
# 检查 Redis 是否运行
redis-cli ping

# 如果有密码
redis-cli -a 123456 ping
```

### 4. 端口被占用
**解决：** 修改 `config/config.yaml` 中的端口
```yaml
server:
  port: 8081  # 改为其他端口
```

## 📝 测试数据

测试脚本会自动创建：
- 用户：testuser_[随机]/test123456
- 客户：测试公司 [随机]
- 联系人：张三
- 合同：HT-[随机]
- API Key：测试 Key

## ✅ 成功标志

看到以下输出表示测试成功：
```
========================================
   ✅ 所有测试完成！
========================================
```
