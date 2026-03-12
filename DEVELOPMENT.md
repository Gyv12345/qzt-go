# QZT-GO 开发指南

> 为开发者提供的详细开发指南

## 开发环境设置

### 1. 安装 Go

从 https://go.dev/dl/ 下载并安装 Go 1.22 或更高版本。

配置 GOPATH 和 GOROOT：
```bash
go env -w GOPATH=%USERPROFILE%\go
go env -w GOBIN=%GOPATH%\bin
```

### 2. 安装开发工具

```bash
# 安装 Gin 框架开发工具
go install github.com/swaggo/swag/cmd/swag@latest

# 安装代码格式化工具
go install golang.org/x/tools/cmd/goimports@latest

# 安装代码检查工具
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# 安装热重载工具
go install github.com/cosmtrek/air@latest
```

### 3. 配置 VS Code

推荐安装的 VS Code 扩展：
- Go (官方)
- Go Test Explorer
- REST Client (用于测试 API)

## 项目结构说明

### 分层架构

```
Handler (HTTP 层)
    ↓
Service (业务逻辑层)
    ↓
Repository (数据访问层)
    ↓
Database (数据库)
```

### 各层职责

#### Handler 层
- 处理 HTTP 请求和响应
- 参数验证和绑定
- 调用 Service 层
- 返回 JSON 数据

#### Service 层
- 实现业务逻辑
- 数据验证
- 调用 Repository 层
- 事务管理

#### Repository 层
- 数据库 CRUD 操作
- 复杂查询
- 数据映射

## 开发流程

### 1. 添加新功能模块

假设要添加"产品管理"模块：

#### Step 1: 创建领域模型

```go
// internal/domain/product.go
package domain

type Product struct {
    ID          string    `gorm:"primaryKey;size:30" json:"id"`
    Name        string    `gorm:"size:100;not null" json:"name"`
    Price       float64   `gorm:"type:decimal(10,2)" json:"price"`
    Status      int       `gorm:"default:1" json:"status"`
    CreatedAt   time.Time `json:"createdAt"`
    UpdatedAt   time.Time `json:"updatedAt"`
}

func (Product) TableName() string {
    return "products"
}
```

#### Step 2: 创建仓储接口和实现

```go
// internal/repository/product.go
package repository

import (
    "context"
    "github.com/qzt/backend/internal/domain"
    "gorm.io/gorm"
)

type ProductRepository interface {
    Create(ctx context.Context, product *domain.Product) error
    GetByID(ctx context.Context, id string) (*domain.Product, error)
    List(ctx context.Context, params ListParams) ([]domain.Product, int64, error)
    Update(ctx context.Context, product *domain.Product) error
    Delete(ctx context.Context, id string) error
}

type productRepository struct {
    db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
    return &productRepository{db: db}
}

// 实现接口方法...
```

#### Step 3: 创建服务层

```go
// internal/service/product.go
package service

import (
    "context"
    "github.com/qzt/backend/internal/domain"
    "github.com/qzt/backend/internal/repository"
)

type ProductService interface {
    Create(ctx context.Context, req *CreateProductReq) (*domain.Product, error)
    GetByID(ctx context.Context, id string) (*domain.Product, error)
    List(ctx context.Context, params repository.ListParams) (*ListResult, error)
    Update(ctx context.Context, id string, req *UpdateProductReq) (*domain.Product, error)
    Delete(ctx context.Context, id string) error
}

type productService struct {
    repo repository.ProductRepository
}

func NewProductService(repo repository.ProductRepository) ProductService {
    return &productService{repo: repo}
}

// 实现接口方法...
```

#### Step 4: 创建处理器

```go
// internal/handler/product.go
package handler

import (
    "github.com/gin-gonic/gin"
    "github.com/qzt/backend/internal/service"
)

type ProductHandler struct {
    svc service.ProductService
}

func NewProductHandler(svc service.ProductService) *ProductHandler {
    return &ProductHandler{svc: svc}
}

// List, Create, Get, Update, Delete 方法...
```

#### Step 5: 注册路由

在 `cmd/server/main.go` 中添加路由：

```go
// 在需要认证的路由组中添加
api.GET("/products", handlers.Product.List)
api.POST("/products", handlers.Product.Create)
api.GET("/products/:id", handlers.Product.Get)
api.PUT("/products/:id", handlers.Product.Update)
api.DELETE("/products/:id", handlers.Product.Delete)
```

### 2. 添加中间件

```go
// internal/middleware/custom.go
package middleware

import "github.com/gin-gonic/gin"

func CustomMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 前置处理
        c.Next()
        // 后置处理
    }
}
```

在路由中使用：

```go
r.Use(middleware.CustomMiddleware())
```

### 3. 添加配置

在 `config/config.yaml` 中添加配置：

```yaml
custom:
  enabled: true
  value: 123
```

在 `internal/config/config.go` 中添加配置结构：

```go
type Config struct {
    // ... 现有配置
    Custom CustomConfig `mapstructure:"custom"`
}

type CustomConfig struct {
    Enabled bool `mapstructure:"enabled"`
    Value   int  `mapstructure:"value"`
}
```

## 代码规范

### 1. 命名规范

- **包名**: 小写单词，不使用下划线或驼峰
  ```go
  package customer // ✅
  package customer_service // ❌
  ```

- **变量名**: 驼峰命名
  ```go
  var customerName string // ✅
  var customer_name string // ❌
  ```

- **常量名**: 大写单词，下划线分隔
  ```go
  const MAX_RETRY_COUNT = 3 // ✅
  const maxRetryCount = 3 // ❌
  ```

- **接口名**: 方法名 + er 后缀
  ```go
  type CustomerRepository interface {} // ✅
  type CustomerRepo interface {} // ❌
  ```

### 2. 错误处理

```go
// ✅ 好的做法
result, err := someFunction()
if err != nil {
    return nil, fmt.Errorf("操作失败: %w", err)
}

// ❌ 不好的做法
result, err := someFunction()
if err != nil {
    log.Println(err) // 不要直接忽略错误
}
```

### 3. Context 使用

```go
// ✅ 传递 context
func (s *service) GetByID(ctx context.Context, id string) (*Customer, error) {
    return s.repo.GetByID(ctx, id)
}

// ❌ 不使用 context
func (s *service) GetByID(id string) (*Customer, error) {
    return s.repo.GetByID(id)
}
```

### 4. 日志记录

```go
// ✅ 使用结构化日志
logger.Info("创建客户",
    zap.String("customerId", customer.ID),
    zap.String("customerName", customer.Name),
)

// ❌ 使用普通字符串
logger.Info("创建客户成功: " + customer.Name)
```

## 测试

### 1. 单元测试

```go
// internal/service/customer_test.go
package service_test

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestCustomerService_Create(t *testing.T) {
    // 准备测试数据
    req := &CreateCustomerReq{
        Name: "测试公司",
    }

    // 执行测试
    result, err := svc.Create(context.Background(), req)

    // 验证结果
    assert.NoError(t, err)
    assert.NotNil(t, result)
    assert.Equal(t, "测试公司", result.Name)
}
```

### 2. 运行测试

```bash
# 运行所有测试
go test ./...

# 运行特定包的测试
go test ./internal/service

# 运行特定测试函数
go test ./internal/service -run TestCustomerService_Create

# 查看覆盖率
go test ./... -cover
```

## 调试

### 1. 使用 Delve 调试器

```bash
# 安装 Delve
go install github.com/go-delve/delve/cmd/dlv@latest

# 调试程序
dlv debug cmd/server/main.go
```

### 2. VS Code 调试配置

创建 `.vscode/launch.json`:

```json
{
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Launch",
            "type": "go",
            "request": "launch",
            "mode": "auto",
            "program": "${workspaceFolder}/cmd/server",
            "env": {
                "GIN_MODE": "debug"
            },
            "args": []
        }
    ]
}
```

## 性能优化

### 1. 数据库查询优化

```go
// ❌ N+1 查询问题
customers, _ := repo.List()
for _, customer := range customers {
    contacts, _ := contactRepo.GetByCustomerID(customer.ID) // 每次循环都查询
}

// ✅ 使用预加载
customers, _ := repo.ListWithPreload() // 一次性查询
```

### 2. 使用缓存

```go
// 使用 Redis 缓存
func (s *service) GetByID(ctx context.Context, id string) (*Customer, error) {
    // 先查缓存
    cached := s.cache.Get("customer:" + id)
    if cached != nil {
        return cached, nil
    }

    // 查数据库
    customer, err := s.repo.GetByID(ctx, id)
    if err != nil {
        return nil, err
    }

    // 写入缓存
    s.cache.Set("customer:" + id, customer, 5*time.Minute)

    return customer, nil
}
```

## 常见问题

### Q1: 如何处理数据库事务？

```go
// 开启事务
tx := db.Begin()
defer func() {
    if r := recover(); r != nil {
        tx.Rollback()
    }
}()

// 执行操作
if err := tx.Create(&customer).Error; err != nil {
    tx.Rollback()
    return err
}

if err := tx.Create(&contact).Error; err != nil {
    tx.Rollback()
    return err
}

// 提交事务
tx.Commit()
```

### Q2: 如何实现软删除？

```go
// 在模型中添加删除时间字段
type Customer struct {
    DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt"`
}

// GORM 会自动处理软删除
db.Delete(&customer) // 实际上更新 deleted_at 字段
```

### Q3: 如何实现分页？

```go
type ListParams struct {
    Page     int
    PageSize int
}

func (r *repository) List(ctx context.Context, params ListParams) ([]Customer, int64, error) {
    var customers []Customer
    var total int64

    // 查询总数
    db.Model(&Customer{}).Count(&total)

    // 分页查询
    offset := (params.Page - 1) * params.PageSize
    db.Offset(offset).Limit(params.PageSize).Find(&customers)

    return customers, total, nil
}
```

## 参考资源

- [Go 官方文档](https://golang.org/doc/)
- [Gin 文档](https://gin-gonic.com/docs/)
- [GORM 文档](https://gorm.io/docs/)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Effective Go](https://golang.org/doc/effective_go.html)

---

**最后更新**: 2026-03-11
