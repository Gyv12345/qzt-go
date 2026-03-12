// Package service 定义了业务逻辑层
// 实现了核心业务规则和流程编排
package service

import (
	"context"

	"github.com/qzt/backend/internal/config"
	"github.com/qzt/backend/internal/domain"
	"github.com/qzt/backend/internal/repository"
	"github.com/qzt/backend/pkg/apikey"
)

// APIKeyService API Key 服务接口
// 定义了 API Key 管理相关的业务方法
type APIKeyService interface {
	// Create 创建 API Key
	Create(ctx context.Context, req *CreateAPIKeyReq) (*CreateAPIKeyResult, error)
	// GetByID 根据 ID 获取 API Key 信息
	GetByID(ctx context.Context, id string) (*domain.APIKey, error)
	// List 获取指定用户的所有 API Key 列表
	List(ctx context.Context, userID string) ([]domain.APIKey, error)
	// Revoke 撤销（禁用）API Key
	Revoke(ctx context.Context, id string) error
}

// apiKeyService API Key 服务实现
type apiKeyService struct {
	repo repository.APIKeyRepository // API Key 仓储
	cfg  *config.Config              // 应用配置
}

// NewAPIKeyService 创建 API Key 服务实例
// 参数：
//   - repo: API Key 仓储接口
//   - cfg: 应用配置对象
// 返回值：
//   - APIKeyService: API Key 服务接口
func NewAPIKeyService(repo repository.APIKeyRepository, cfg *config.Config) APIKeyService {
	return &apiKeyService{
		repo: repo,
		cfg:  cfg,
	}
}

// CreateAPIKeyReq 创建 API Key 请求
type CreateAPIKeyReq struct {
	Name        string   `json:"name" binding:"required"`   // API Key 名称（便于识别）
	UserID      string   `json:"userId" binding:"required"` // 所属用户 ID（必填）
	Permissions []string `json:"permissions"`               // 权限列表（如 ["customer:read", "contact:read"]）
	IPWhitelist []string `json:"ipWhitelist"`               // IP 白名单（如 ["192.168.1.0/24"]）
	RateLimit   int      `json:"rateLimit"`                 // 速率限制（每分钟最大请求数）
}

// CreateAPIKeyResult 创建 API Key 结果
// 包含原始密钥（仅显示一次）和 API Key 信息
type CreateAPIKeyResult struct {
	APIKey string         `json:"apiKey"` // 原始 key，只显示一次，请妥善保存
	Info   *domain.APIKey `json:"info"`   // API Key 信息（不含原始密钥）
}

// Create 创建 API Key
// 生成新的 API Key，包括密钥生成、哈希计算和权限配置
// 参数：
//   - ctx: 上下文
//   - req: 创建 API Key 请求
// 返回值：
//   - *CreateAPIKeyResult: 创建结果（包含原始密钥和 API Key 信息）
//   - error: 错误信息
func (s *apiKeyService) Create(ctx context.Context, req *CreateAPIKeyReq) (*CreateAPIKeyResult, error) {
	// ========== 第一步：生成 API Key ==========
	// 生成原始密钥、前缀和哈希值
	rawKey, prefix, keyHash, err := apikey.Generate(true)
	if err != nil {
		return nil, err
	}

	// ========== 第二步：序列化配置信息 ==========
	// 将权限列表和 IP 白名单序列化为 JSON 字符串
	permissionsJSON, _ := serializeStringArray(req.Permissions)
	ipWhitelistJSON, _ := serializeStringArray(req.IPWhitelist)

	// ========== 第三步：设置速率限制 ==========
	// 如果未指定速率限制，使用默认配置
	rateLimit := req.RateLimit
	if rateLimit == 0 {
		rateLimit = s.cfg.APIKey.DefaultRateLimit
	}

	// ========== 第四步：构建 API Key 对象 ==========
	apiKey := &domain.APIKey{
		Name:        req.Name,
		KeyHash:     keyHash,              // 存储哈希值，不存储原始密钥
		Prefix:      prefix,               // 存储前缀用于识别
		UserID:      req.UserID,
		Permissions: permissionsJSON,
		IPWhitelist: ipWhitelistJSON,
		RateLimit:   rateLimit,
		Status:      1, // 默认启用
	}

	// ========== 第五步：保存到数据库 ==========
	if err := s.repo.Create(ctx, apiKey); err != nil {
		return nil, err
	}

	// ========== 第六步：返回结果 ==========
	return &CreateAPIKeyResult{
		APIKey: rawKey, // 原始密钥，仅此一次显示
		Info:   apiKey,
	}, nil
}

// GetByID 根据 ID 获取 API Key 信息
// 参数：
//   - ctx: 上下文
//   - id: API Key ID
// 返回值：
//   - *domain.APIKey: API Key 对象
//   - error: 错误信息
func (s *apiKeyService) GetByID(ctx context.Context, id string) (*domain.APIKey, error) {
	return s.repo.GetByID(ctx, id)
}

// List 获取指定用户的所有 API Key 列表
// 参数：
//   - ctx: 上下文
//   - userID: 用户 ID
// 返回值：
//   - []domain.APIKey: API Key 列表
//   - error: 错误信息
func (s *apiKeyService) List(ctx context.Context, userID string) ([]domain.APIKey, error) {
	return s.repo.List(ctx, userID)
}

// Revoke 撤销（禁用）API Key
// 将 API Key 状态设置为禁用，使其无法继续使用
// 参数：
//   - ctx: 上下文
//   - id: API Key ID
// 返回值：
//   - error: 错误信息
func (s *apiKeyService) Revoke(ctx context.Context, id string) error {
	// 查询 API Key 信息
	key, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// 将状态设置为禁用
	key.Status = 0
	
	// 更新到数据库
	return s.repo.Update(ctx, key)
}

// serializeStringArray 将字符串数组序列化为 JSON 字符串
// 简单的 JSON 序列化实现，用于权限列表和 IP 白名单
// 参数：
//   - arr: 字符串数组
// 返回值：
//   - string: JSON 字符串
//   - error: 错误信息
func serializeStringArray(arr []string) (string, error) {
	// 如果数组为空，返回空数组 JSON
	if len(arr) == 0 {
		return "[]", nil
	}
	
	// 简单的 JSON 序列化
	result := "["
	for i, s := range arr {
		if i > 0 {
			result += ","
		}
		result += `"` + s + `"`
	}
	result += "]"
	return result, nil
}
