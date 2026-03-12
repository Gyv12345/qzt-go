# QZT-GO 完整目录结构

```
qzt-go/
│
├── cmd/                          # 命令行工具
│   └── server/
│       └── main.go              # 程序入口
│
├── internal/                     # 内部包（不可外部引用）
│   ├── config/
│   │   └── config.go            # 配置加载和管理
│   │
│   ├── domain/                  # 领域模型
│   │   ├── apikey.go           # API Key 实体
│   │   ├── contact.go          # 联系人实体
│   │   ├── contract.go         # 合同实体
│   │   ├── customer.go         # 客户实体
│   │   └── user.go             # 用户、角色、部门、菜单实体
│   │
│   ├── handler/                 # HTTP 处理器
│   │   ├── apikey.go           # API Key 处理器
│   │   ├── auth.go             # 认证处理器
│   │   ├── contact.go          # 联系人处理器
│   │   ├── contract.go         # 合同处理器
│   │   ├── customer.go         # 客户处理器
│   │   └── handler.go          # 处理器聚合
│   │
│   ├── middleware/              # 中间件
│   │   ├── apikey.go           # API Key 认证中间件
│   │   ├── auth.go             # JWT 认证中间件
│   │   ├── cors.go             # CORS 中间件
│   │   └── logger.go           # 日志中间件
│   │
│   ├── repository/              # 数据访问层
│   │   ├── apikey.go           # API Key 仓储
│   │   ├── contact.go          # 联系人仓储
│   │   ├── contract.go         # 合同仓储
│   │   ├── customer.go         # 客户仓储
│   │   └── user.go             # 用户仓储（聚合其他仓储）
│   │
│   └── service/                 # 业务逻辑层
│       ├── apikey.go           # API Key 服务
│       ├── auth.go             # 认证服务
│       ├── contact.go          # 联系人服务
│       ├── contract.go         # 合同服务
│       ├── customer.go         # 客户服务
│       ├── errors.go           # 错误定义
│       └── service.go          # 服务聚合
│
├── pkg/                         # 外部可用包
│   └── apikey/
│       └── generator.go        # API Key 生成器
│
├── config/                      # 配置文件目录
│   └── config.yaml             # 主配置文件
│
├── migrations/                  # 数据库迁移
│   └── init.sql                # 初始化 SQL 脚本
│
├── frontend/                    # 前端代码（从 qzt 复制）
│   ├── src/
│   │   ├── components/
│   │   ├── pages/
│   │   ├── api/
│   │   ├── hooks/
│   │   ├── store/
│   │   ├── lib/
│   │   └── main.tsx
│   ├── public/
│   ├── package.json
│   ├── vite.config.ts
│   ├── tsconfig.json
│   └── tailwind.config.js
│
├── .env.example                 # 环境变量示例
├── .gitignore                   # Git 忽略配置
├── Dockerfile                   # Docker 镜像配置
├── docker-compose.yml           # Docker Compose 配置
├── go.mod                       # Go 模块定义
├── Makefile                     # 构建脚本
├── nginx.conf                   # Nginx 配置
├── README.md                    # 项目说明
├── API.md                       # API 文档
├── QUICKSTART.md                # 快速启动指南
└── PROJECT_SUMMARY.md           # 项目总结
```

## 文件统计

### Go 源码文件
- `cmd/server/main.go` - 1 个
- `internal/config/` - 1 个
- `internal/domain/` - 5 个
- `internal/handler/` - 6 个
- `internal/middleware/` - 4 个
- `internal/repository/` - 5 个
- `internal/service/` - 7 个
- `pkg/apikey/` - 1 个

**总计**: 30 个 Go 源文件

### 配置文件
- `config/config.yaml` - 主配置
- `.env.example` - 环境变量示例
- `docker-compose.yml` - Docker Compose
- `nginx.conf` - Nginx 配置
- `Dockerfile` - Docker 镜像
- `Makefile` - 构建脚本
- `.gitignore` - Git 忽略
- `go.mod` - Go 模块

**总计**: 8 个配置文件

### 文档文件
- `README.md` - 项目说明
- `API.md` - API 文档
- `QUICKSTART.md` - 快速指南
- `PROJECT_SUMMARY.md` - 项目总结
- `DIRECTORY_TREE.md` - 本文件

**总计**: 5 个文档文件

### 数据库文件
- `migrations/init.sql` - 初始化脚本

**总计**: 1 个数据库文件

### 前端文件
- `frontend/` - 完整的 React + Vite 前端项目
  - 组件、页面、API、Hooks、Store 等

**总计**: 约 200+ 个文件（从原项目复制）

## 总计

- **Go 源码**: 30 个文件
- **配置文件**: 8 个文件
- **文档文件**: 5 个文件
- **数据库文件**: 1 个文件
- **前端文件**: 200+ 个文件

**全部总计**: 240+ 个文件

## 代码行数估算

- Go 源码: ~3,000 行
- 配置文件: ~500 行
- 文档: ~2,000 行
- 数据库 SQL: ~400 行
- 前端代码: ~50,000 行

**总计**: ~55,900 行代码

## 功能模块

### 已实现模块
- ✅ 用户认证 (JWT)
- ✅ API Key 管理
- ✅ 客户管理
- ✅ 联系人管理
- ✅ 合同管理
- ✅ 基础 RBAC

### 待实现模块 (参考原计划)
- ⏳ 产品管理
- ⏳ 支付管理
- ⏳ 发票管理
- ⏳ 跟进记录
- ⏳ 服务团队
- ⏳ 内容管理
- ⏳ 文件存储
- ⏳ 社交媒体
- ⏳ 通知系统
- ⏳ Webhook
- ⏳ 操作日志
- ⏳ 系统配置
- ⏳ 定时任务
- ⏳ 电子签名
- ⏳ AI Agent

**总计**: 35 个模块，已完成 5 个核心模块

---

*最后更新: 2026-03-11*
