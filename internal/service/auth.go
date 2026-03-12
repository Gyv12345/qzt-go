// Package service 定义了业务逻辑层
// 实现了核心业务规则和流程编排
package service

import (
	"context"
	"encoding/json"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/qzt/backend/internal/config"
	"github.com/qzt/backend/internal/domain"
	"github.com/qzt/backend/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

// AuthService 认证服务接口
// 定义了用户认证相关的业务方法
type AuthService interface {
	// Login 用户登录
	// 验证用户凭据并返回访问令牌和刷新令牌
	Login(ctx context.Context, username, password string) (*LoginResult, error)
	// Register 用户注册
	// 创建新用户账户
	Register(ctx context.Context, req *RegisterReq) (*domain.User, error)
	// RefreshToken 刷新令牌
	// 使用刷新令牌获取新的访问令牌
	RefreshToken(ctx context.Context, refreshToken string) (*LoginResult, error)
}

// authService 认证服务实现
type authService struct {
	userRepo repository.UserRepository // 用户仓储
	cfg      *config.Config             // 应用配置
}

// NewAuthService 创建认证服务实例
// 参数：
//   - userRepo: 用户仓储接口
//   - cfg: 应用配置对象
// 返回值：
//   - AuthService: 认证服务接口
func NewAuthService(userRepo repository.UserRepository, cfg *config.Config) AuthService {
	return &authService{
		userRepo: userRepo,
		cfg:      cfg,
	}
}

// LoginResult 登录结果
// 包含令牌信息和用户信息
type LoginResult struct {
	Token        string        `json:"token"`         // 访问令牌
	RefreshToken string        `json:"refreshToken"`   // 刷新令牌
	ExpiresAt    int64         `json:"expiresAt"`      // 访问令牌过期时间（Unix 时间戳）
	User         *domain.User `json:"user"`           // 用户信息
}

// RegisterReq 用户注册请求
type RegisterReq struct {
	Username string `json:"username" binding:"required,min=3,max=50"` // 用户名（必填，3-50字符）
	Password string `json:"password" binding:"required,min=6"`       // 密码（必填，至少6字符）
	Email    string `json:"email" binding:"email"`                   // 电子邮箱
	Mobile   string `json:"mobile"`                                  // 手机号码
	RealName string `json:"realName"`                                // 真实姓名
}

// JWTClaims JWT 令牌声明
// 自定义的 JWT Claims，包含用户信息
type JWTClaims struct {
	UserID   string `json:"userId"`   // 用户 ID
	Username string `json:"username"` // 用户名
	jwt.RegisteredClaims                // JWT 标准声明（过期时间、签发时间等）
}

// Login 用户登录
// 验证用户名和密码，生成并返回 JWT 令牌
// 参数：
//   - ctx: 上下文
//   - username: 用户名
//   - password: 密码
// 返回值：
//   - *LoginResult: 登录结果，包含令牌和用户信息
//   - error: 错误信息（如果登录失败）
func (s *authService) Login(ctx context.Context, username, password string) (*LoginResult, error) {
	// ========== 第一步：查询用户 ==========
	// 根据用户名查询用户信息
	user, err := s.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	// ========== 第二步：验证密码 ==========
	// 使用 bcrypt 比对密码哈希值
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	// ========== 第三步：检查用户状态 ==========
	// 确保用户账户未被禁用
	if user.Status != domain.UserStatusActive {
		return nil, ErrUserInactive
	}

	// ========== 第四步：生成访问令牌 ==========
	now := time.Now()
	expiresAt := now.Add(s.cfg.JWT.Expire)

	// 构建 JWT Claims
	claims := &JWTClaims{
		UserID:   user.ID,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	// 使用 HS256 算法签名生成 Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.cfg.JWT.Secret))
	if err != nil {
		return nil, err
	}

	// ========== 第五步：生成刷新令牌 ==========
	// 刷新令牌的有效期比访问令牌长
	refreshClaims := &JWTClaims{
		UserID:   user.ID,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(s.cfg.JWT.RefreshExpire)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(s.cfg.JWT.Secret))
	if err != nil {
		return nil, err
	}

	// ========== 第六步：返回结果 ==========
	return &LoginResult{
		Token:        tokenString,
		RefreshToken: refreshTokenString,
		ExpiresAt:    expiresAt.Unix(),
		User:         user,
	}, nil
}

// Register 用户注册
// 创建新的用户账户
// 参数：
//   - ctx: 上下文
//   - req: 注册请求
// 返回值：
//   - *domain.User: 创建的用户对象
//   - error: 错误信息
func (s *authService) Register(ctx context.Context, req *RegisterReq) (*domain.User, error) {
	// ========== 第一步：检查用户名是否已存在 ==========
	if _, err := s.userRepo.GetByUsername(ctx, req.Username); err == nil {
		return nil, ErrUsernameExists
	}

	// ========== 第二步：加密密码 ==========
	// 使用 bcrypt 算法对密码进行哈希加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// ========== 第三步：创建用户 ==========
	user := &domain.User{
		Username: req.Username,
		Password: string(hashedPassword),
		Email:    req.Email,
		Mobile:   req.Mobile,
		RealName: req.RealName,
		Status:   domain.UserStatusActive, // 默认激活状态
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// RefreshToken 刷新令牌
// 使用刷新令牌获取新的访问令牌和刷新令牌
// 参数：
//   - ctx: 上下文
//   - refreshToken: 刷新令牌
// 返回值：
//   - *LoginResult: 新的登录结果
//   - error: 错误信息
func (s *authService) RefreshToken(ctx context.Context, refreshToken string) (*LoginResult, error) {
	// ========== 第一步：解析刷新令牌 ==========
	claims := &JWTClaims{}
	token, err := jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.cfg.JWT.Secret), nil
	})

	if err != nil || !token.Valid {
		return nil, ErrInvalidToken
	}

	// ========== 第二步：获取用户信息 ==========
	user, err := s.userRepo.GetByID(ctx, claims.UserID)
	if err != nil {
		return nil, err
	}

	// ========== 第三步：检查用户状态 ==========
	if user.Status != domain.UserStatusActive {
		return nil, ErrUserInactive
	}

	// ========== 第四步：生成新的访问令牌 ==========
	now := time.Now()
	expiresAt := now.Add(s.cfg.JWT.Expire)

	newClaims := &JWTClaims{
		UserID:   user.ID,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, newClaims)
	tokenString, err := newToken.SignedString([]byte(s.cfg.JWT.Secret))
	if err != nil {
		return nil, err
	}

	// ========== 第五步：生成新的刷新令牌 ==========
	refreshClaims := &JWTClaims{
		UserID:   user.ID,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(s.cfg.JWT.RefreshExpire)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}
	newRefreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	newRefreshTokenString, err := newRefreshToken.SignedString([]byte(s.cfg.JWT.Secret))
	if err != nil {
		return nil, err
	}

	// ========== 第六步：返回结果 ==========
	return &LoginResult{
		Token:        tokenString,
		RefreshToken: newRefreshTokenString,
		ExpiresAt:    expiresAt.Unix(),
		User:         user,
	}, nil
}
