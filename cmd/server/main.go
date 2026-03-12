// Package main 是 QZT 后端服务的主程序入口包
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/qzt/backend/internal/config"
	"github.com/qzt/backend/internal/domain"
	"github.com/qzt/backend/internal/handler"
	"github.com/qzt/backend/internal/middleware"
	"github.com/qzt/backend/internal/repository"
	"github.com/qzt/backend/internal/service"
)

// main 主程序入口函数
func main() {
	// ========== 加载配置 ==========
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// ========== 初始化日志 ==========
	logger, err := initLogger(cfg)
	if err != nil {
		log.Fatalf("初始化日志失败: %v", err)
	}
	defer logger.Sync()

	// ========== 连接数据库 ==========
	db, err := gorm.Open(postgres.Open(cfg.Database.DSN()), &gorm.Config{})
	if err != nil {
		logger.Fatal("数据库连接失败", zap.Error(err))
	}

	// 配置连接池
	sqlDB, err := db.DB()
	if err != nil {
		logger.Fatal("获取数据库连接失败", zap.Error(err))
	}
	sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// ========== 数据库迁移 ==========
	if err := db.AutoMigrate(
		&domain.User{},
		&domain.Role{},
		&domain.Department{},
		&domain.Menu{},
		&domain.Permission{},
		&domain.Customer{},
		&domain.Contact{},
		&domain.Contract{},
		&domain.APIKey{},
		&domain.FollowRecord{},
		&domain.Opportunity{},
		&domain.Payment{},
		&domain.Product{},
		&domain.Invoice{},
		&domain.UserRole{},
		&domain.RolePermission{},
	); err != nil {
		logger.Fatal("数据库迁移失败", zap.Error(err))
	}

	// ========== 创建数据库索引 ==========
	if err := createIndexes(db); err != nil {
		logger.Fatal("创建索引失败", zap.Error(err))
	}

	// ========== 依赖注入 ==========
	repos := repository.NewRepository(db)
	services := service.NewService(repos, cfg)
	handlers := handler.NewHandler(services)

	// ========== 配置 Gin ==========
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	
	// 全局中间件
	r.Use(gin.Recovery())
	r.Use(middleware.RequestID())  // Request ID 追踪
	r.Use(middleware.Logger(logger))
	r.Use(middleware.CORS(&cfg.CORS))
	
	// API 限流（每秒 100 请求，最大 200）
	limiter := middleware.NewTokenBucketLimiter(100, 200)
	r.Use(middleware.RateLimit(limiter))

	// ========== 健康检查 ==========
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":    "ok",
			"timestamp": time.Now().Unix(),
			"requestId": middleware.GetRequestID(c),
		})
	})

	// ========== 注册路由 ==========
	setupRoutes(r, handlers, cfg, repos)

	// ========== 启动服务器（优雅关闭）==========
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	// 启动 HTTP 服务器（非阻塞）
	go func() {
		logger.Info("服务器启动", zap.String("addr", addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("服务器启动失败", zap.Error(err))
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("正在关闭服务器...")

	// 优雅关闭（最多等待 30 秒）
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("服务器强制关闭", zap.Error(err))
	}

	logger.Info("服务器已关闭")
}

// initLogger 初始化日志记录器
func initLogger(cfg *config.Config) (*zap.Logger, error) {
	if cfg.Server.Mode == "release" {
		return zap.NewProduction()
	}
	return zap.NewDevelopment()
}

// setupRoutes 配置路由
func setupRoutes(r *gin.Engine, handlers *handler.Handler, cfg *config.Config, repos *repository.Repository) {
	// 公开路由（无需认证）
	auth := r.Group("/api/auth")
	{
		auth.POST("/login", handlers.Auth.Login)
		auth.POST("/register", handlers.Auth.Register)
		auth.POST("/refresh", handlers.Auth.RefreshToken)
	}

	// 需要认证的路由
	authProtected := r.Group("/api/auth")
	authProtected.Use(middleware.JWTAuth(cfg))
	{
		authProtected.GET("/me", handlers.Auth.GetMe)
	}

	// RBAC 权限管理
	rbac := r.Group("/api")
	rbac.Use(middleware.JWTAuth(cfg))
	{
		// 角色管理
		rbac.POST("/roles", handlers.RBAC.CreateRole)
		rbac.GET("/roles", handlers.RBAC.ListRoles)
		rbac.POST("/roles/:id/permissions", handlers.RBAC.AssignPermissions)
		
		// 权限管理
		rbac.POST("/permissions", handlers.RBAC.CreatePermission)
		rbac.GET("/permissions", handlers.RBAC.ListPermissions)
		
		// 用户角色分配
		rbac.POST("/users/:id/roles", handlers.RBAC.AssignRolesToUser)
		rbac.GET("/users/:id/permissions", handlers.RBAC.GetUserPermissions)
	}

	// 跟进记录
	follow := r.Group("/api")
	follow.Use(middleware.JWTAuth(cfg))
	{
		follow.GET("/follow-records", handlers.FollowRecord.List)
		follow.POST("/follow-records", handlers.FollowRecord.Create)
		follow.GET("/follow-records/:id", handlers.FollowRecord.Get)
		follow.PUT("/follow-records/:id", handlers.FollowRecord.Update)
		follow.DELETE("/follow-records/:id", handlers.FollowRecord.Delete)
	}

	// 商机管理
	opp := r.Group("/api")
	opp.Use(middleware.JWTAuth(cfg))
	{
		opp.GET("/opportunities", handlers.Opportunity.List)
		opp.POST("/opportunities", handlers.Opportunity.Create)
		opp.GET("/opportunities/:id", handlers.Opportunity.Get)
		opp.PUT("/opportunities/:id", handlers.Opportunity.Update)
		opp.DELETE("/opportunities/:id", handlers.Opportunity.Delete)
		opp.POST("/opportunities/:id/next-stage", handlers.Opportunity.MoveToNextStage)
		opp.POST("/opportunities/:id/prev-stage", handlers.Opportunity.MoveToPrevStage)
		opp.POST("/opportunities/:id/convert", handlers.Opportunity.ConvertToContract)
		opp.GET("/opportunities/stats/stage", handlers.Opportunity.GetStageStats)
	}

	// 回款管理
	payment := r.Group("/api")
	payment.Use(middleware.JWTAuth(cfg))
	{
		payment.GET("/payments", handlers.Payment.List)
		payment.POST("/payments", handlers.Payment.Create)
		payment.POST("/payments/:id/confirm", handlers.Payment.Confirm)
	}

	// 报表统计
	report := r.Group("/api/reports")
	report.Use(middleware.JWTAuth(cfg))
	{
		report.GET("/dashboard", handlers.Report.GetDashboard)
		report.GET("/sales-funnel", handlers.Report.GetSalesFunnel)
		report.GET("/performance", handlers.Report.GetPerformance)
		report.GET("/payments", handlers.Report.GetPaymentStats)
		report.GET("/customers", handlers.Report.GetCustomerStats)
	}

	// 内部 API（JWT 认证）
	api := r.Group("/api")
	api.Use(middleware.JWTAuth(cfg))
	{
		// 客户管理
		api.GET("/customers", handlers.Customer.List)
		api.POST("/customers", handlers.Customer.Create)
		api.GET("/customers/:id", handlers.Customer.Get)
		api.PUT("/customers/:id", handlers.Customer.Update)
		api.DELETE("/customers/:id", handlers.Customer.Delete)

		// 联系人管理
		api.GET("/contacts", handlers.Contact.List)
		api.POST("/contacts", handlers.Contact.Create)
		api.GET("/contacts/:id", handlers.Contact.Get)
		api.PUT("/contacts/:id", handlers.Contact.Update)
		api.DELETE("/contacts/:id", handlers.Contact.Delete)

		// 合同管理
		api.GET("/contracts", handlers.Contract.List)
		api.POST("/contracts", handlers.Contract.Create)
		api.GET("/contracts/:id", handlers.Contract.Get)
		api.PUT("/contracts/:id", handlers.Contract.Update)
		api.DELETE("/contracts/:id", handlers.Contract.Delete)

		// API Key 管理
		api.GET("/apikeys", handlers.APIKey.List)
		api.POST("/apikeys", handlers.APIKey.Create)
		api.POST("/apikeys/:id/revoke", handlers.APIKey.Revoke)
	}

	// 外部 API（API Key 认证）
	external := r.Group("/api/v1")
	apiKeyMiddleware := middleware.NewAPIKeyMiddleware(repos.APIKey)
	external.Use(apiKeyMiddleware.Handler())
	
	// API Key 限流（使用配置的限流值）
	apiKeyLimiter := middleware.NewSlidingWindowLimiter(cfg.APIKey.DefaultRateLimit, time.Minute)
	external.Use(middleware.RateLimitByAPIKey(apiKeyLimiter))
	{
		// 客户管理
		external.GET("/customers", handlers.Customer.List)
		external.POST("/customers", handlers.Customer.Create)
		external.GET("/customers/:id", handlers.Customer.Get)
		external.PUT("/customers/:id", handlers.Customer.Update)
		external.DELETE("/customers/:id", handlers.Customer.Delete)

		// 联系人管理
		external.GET("/contacts", handlers.Contact.List)
		external.POST("/contacts", handlers.Contact.Create)
		external.GET("/contacts/:id", handlers.Contact.Get)
		external.PUT("/contacts/:id", handlers.Contact.Update)
		external.DELETE("/contacts/:id", handlers.Contact.Delete)

		// 合同管理
		external.GET("/contracts", handlers.Contract.List)
		external.POST("/contracts", handlers.Contract.Create)
		external.GET("/contracts/:id", handlers.Contract.Get)
		external.PUT("/contracts/:id", handlers.Contract.Update)
		external.DELETE("/contracts/:id", handlers.Contract.Delete)
	}
}

// createIndexes 创建数据库索引
func createIndexes(db *gorm.DB) error {
	// 客户表索引
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_customers_name ON customers(name)").Error; err != nil {
		return err
	}
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_customers_level ON customers(customer_level)").Error; err != nil {
		return err
	}
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_customers_status ON customers(status)").Error; err != nil {
		return err
	}
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_customers_follow_user ON customers(follow_user_id)").Error; err != nil {
		return err
	}
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_customers_created_at ON customers(created_at)").Error; err != nil {
		return err
	}

	// 联系人表索引
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_contacts_customer_id ON contacts(customer_id)").Error; err != nil {
		return err
	}

	// 合同表索引
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_contracts_customer_id ON contracts(customer_id)").Error; err != nil {
		return err
	}
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_contracts_status ON contracts(status)").Error; err != nil {
		return err
	}

	// API Key 表索引
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_apikeys_prefix ON api_keys(prefix)").Error; err != nil {
		return err
	}
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_apikeys_user_id ON api_keys(user_id)").Error; err != nil {
		return err
	}

	return nil
}
