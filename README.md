# QZT Go 项目

> 企智通 CRM 系统 Go 语言后端

## 🚀 快速开始

### 环境要求
- Go 1.22+
- PostgreSQL 12+
- Redis 6+

### 安装运行
```bash
# 1. 克隆项目
git clone https://github.com/yourusername/qzt-go.git
cd qzt-go

# 2. 安装依赖
go mod tidy

# 3. 创建数据库
psql -U postgres
CREATE DATABASE qzt WITH ENCODING='UTF8';

# 4. 配置
cp config/config.yaml config/config.local.yaml
# 编辑 config.local.yaml，更新数据库和 Redis 配置

# 5. 运行
go run cmd/server/main.go
```

### 访问
- 服务地址: http://localhost:8080
- 健康检查: http://localhost:8080/health

## 📚 文档

- [快速启动](QUICKSTART.md)
- [测试指南](docs/TESTING.md)
- [数据库设置](docs/database-setup.md)

## 🛠️ 技术栈

- **Web 框架**: Gin
- **ORM**: GORM
- **数据库**: PostgreSQL
- **缓存**: Redis
- **日志**: Zap
- **配置**: Viper

## 📦 项目结构

```
qzt-go/
├── cmd/server/          # 程序入口
├── internal/            # 内部代码
│   ├── config/         # 配置管理
│   ├── domain/         # 领域模型
│   ├── handler/        # HTTP 处理器
│   ├── middleware/     # 中间件
│   ├── repository/     # 数据访问层
│   ├── service/        # 业务逻辑层
│   └── infrastructure/ # 基础设施
├── pkg/                 # 公共包
│   ├── errors/         # 错误处理
│   └── response/       # 响应格式
├── scripts/             # 脚本工具
└── config/              # 配置文件
```

## ✨ 特性

- ✅ JWT + API Key 双重认证
- ✅ Request ID 链路追踪
- ✅ 统一错误码和响应格式
- ✅ API 限流（令牌桶/滑动窗口）
- ✅ 优雅关闭
- ✅ 数据库索引优化
- ✅ 单元测试
- ✅ Swagger 文档

## 🧪 测试

```powershell
# 完整测试
.\scripts\full-test.ps1

# 仅测试接口
.\scripts\test-api.ps1
```

## 📝 开发计划

- [ ] WebSocket 支持
- [ ] 集成测试
- [ ] CI/CD 配置
- [ ] Docker 部署
- [ ] 性能监控

## 📄 许可证

MIT License

## 👥 联系方式

- 公司: 河南爱编程网络科技有限公司
- 邮箱: shichenyang@devlovecode.com
- 官网: https://devlovecode.com
