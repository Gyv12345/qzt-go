// Package config 提供应用程序配置管理功能
//
// 该包负责从配置文件（config.yaml）和环境变量中加载应用程序配置，
// 并提供类型安全的配置访问接口。
//
// 主要功能：
//   - 加载 YAML 格式的配置文件
//   - 支持环境变量覆盖配置
//   - 提供数据库、服务器、JWT、API Key、Redis、日志和 CORS 等配置
//
// 配置结构：
//   - ServerConfig: HTTP 服务器配置（端口、运行模式）
//   - DatabaseConfig: 数据库连接配置
//   - JWTConfig: JWT 令牌配置
//   - APIKeyConfig: API Key 配置
//   - RedisConfig: Redis 缓存配置
//   - LogConfig: 日志配置
//   - CORSConfig: 跨域资源共享配置
//
// 使用示例：
//   cfg, err := config.Load()
//   if err != nil {
//       log.Fatal(err)
//   }
//   fmt.Println(cfg.Server.Port)
package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

// Config 应用程序总配置结构
//
// 包含应用程序运行所需的所有配置项，从 YAML 配置文件中加载。
// 使用 mapstructure 标签映射配置文件中的字段。
//
// 配置项：
//   - Server: HTTP 服务器配置
//   - Database: 数据库连接配置
//   - JWT: JWT 令牌配置
//   - APIKey: API Key 认证配置
//   - Redis: Redis 缓存配置
//   - Log: 日志输出配置
//   - CORS: 跨域资源共享配置
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`   // HTTP 服务器配置
	Database DatabaseConfig `mapstructure:"database"` // 数据库配置
	JWT      JWTConfig      `mapstructure:"jwt"`      // JWT 令牌配置
	APIKey   APIKeyConfig   `mapstructure:"apikey"`   // API Key 配置
	Redis    RedisConfig    `mapstructure:"redis"`    // Redis 缓存配置
	Log      LogConfig      `mapstructure:"log"`      // 日志配置
	CORS     CORSConfig     `mapstructure:"cors"`     // CORS 配置
}

// ServerConfig HTTP 服务器配置
//
// 定义 HTTP 服务器的运行参数，包括监听端口和运行模式。
//
// 配置示例（YAML）：
//   server:
//     port: 8080
//     mode: release  # release 或 debug
type ServerConfig struct {
	Port int    `mapstructure:"port"` // HTTP 服务器监听端口（默认：8080）
	Mode string `mapstructure:"mode"` // 运行模式：release（生产）或 debug（开发）
}

// DatabaseConfig 数据库连接配置
//
// 定义 PostgreSQL 数据库的连接参数和连接池设置。
// 支持 GORM 的标准配置项。
//
// 配置示例（YAML）：
//   database:
//     driver: postgres
//     host: localhost
//     port: 5432
//     username: postgres
//     password: password
//     database: qzt
//     max_idle_conns: 10
//     max_open_conns: 100
type DatabaseConfig struct {
	Driver       string `mapstructure:"driver"`         // 数据库驱动类型（如：postgres）
	Host         string `mapstructure:"host"`           // 数据库主机地址
	Port         int    `mapstructure:"port"`           // 数据库端口号
	Username     string `mapstructure:"username"`       // 数据库用户名
	Password     string `mapstructure:"password"`       // 数据库密码
	Database     string `mapstructure:"database"`       // 数据库名称
	MaxIdleConns int    `mapstructure:"max_idle_conns"` // 空闲连接池中的最大连接数
	MaxOpenConns int    `mapstructure:"max_open_conns"` // 数据库的最大打开连接数
}

// DSN 生成数据库连接字符串（Data Source Name）
//
// 根据 DatabaseConfig 中的配置生成 PostgreSQL 连接字符串。
// 格式：host=x port=x user=x password=x dbname=x sslmode=disable
//
// 返回值：
//   - string: PostgreSQL DSN 连接字符串
//
// 示例输出：
//   host=localhost port=5432 user=postgres password=password dbname=qzt sslmode=disable
func (d DatabaseConfig) DSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Shanghai",
		d.Host, d.Port, d.Username, d.Password, d.Database)
}

// JWTConfig JWT 令牌配置
//
// 定义 JWT（JSON Web Token）认证的密钥和过期时间设置。
// 用于用户身份验证和会话管理。
//
// 配置示例（YAML）：
//   jwt:
//     secret: your-secret-key
//     expire: 24h
//     refresh_expire: 168h
type JWTConfig struct {
	Secret         string        `mapstructure:"secret"`         // JWT 签名密钥（必须保密）
	Expire         time.Duration `mapstructure:"expire"`         // 访问令牌过期时间（如：24h）
	RefreshExpire  time.Duration `mapstructure:"refresh_expire"` // 刷新令牌过期时间（如：168h）
}

// APIKeyConfig API Key 认证配置
//
// 定义 API Key 认证机制的配置，用于外部系统调用 API。
//
// 配置示例（YAML）：
//   apikey:
//     secret: your-api-key-secret
//     default_rate_limit: 1000
type APIKeyConfig struct {
	Secret          string `mapstructure:"secret"`           // API Key 签名密钥
	DefaultRateLimit int   `mapstructure:"default_rate_limit"` // 默认速率限制（请求数/分钟）
}

// RedisConfig Redis 缓存配置
//
// 定义 Redis 连接参数和连接池设置。
// 用于缓存、会话存储和分布式锁等功能。
//
// 配置示例（YAML）：
//   redis:
//     host: localhost
//     port: 6379
//     password: ""
//     db: 0
//     pool_size: 10
type RedisConfig struct {
	Host     string `mapstructure:"host"`      // Redis 服务器地址
	Port     int    `mapstructure:"port"`      // Redis 服务器端口
	Password string `mapstructure:"password"`  // Redis 密码（无密码则为空）
	DB       int    `mapstructure:"db"`        // Redis 数据库编号
	PoolSize int    `mapstructure:"pool_size"` // 连接池大小
}

// LogConfig 日志配置
//
// 定义日志输出的级别、格式和目标。
// 支持控制台输出和文件输出。
//
// 配置示例（YAML）：
//   log:
//     level: info
//     format: json
//     output: stdout
//     file: logs/app.log
type LogConfig struct {
	Level  string `mapstructure:"level"`  // 日志级别：debug、info、warn、error
	Format string `mapstructure:"format"` // 日志格式：json 或 console
	Output string `mapstructure:"output"` // 输出目标：stdout、stderr 或文件路径
	File   string `mapstructure:"file"`   // 日志文件路径（如果输出到文件）
}

// CORSConfig 跨域资源共享配置
//
// 定义 CORS（Cross-Origin Resource Sharing）策略，
// 控制哪些来源的请求可以访问 API。
//
// 配置示例（YAML）：
//   cors:
//     allow_origins:
//       - http://localhost:3000
//       - https://example.com
//     allow_methods:
//       - GET
//       - POST
//       - PUT
//       - DELETE
//     allow_headers:
//       - Content-Type
//       - Authorization
//     expose_headers:
//       - Content-Length
//     allow_credentials: true
//     max_age: 86400
type CORSConfig struct {
	AllowOrigins     []string `mapstructure:"allow_origins"`     // 允许的来源域名列表
	AllowMethods     []string `mapstructure:"allow_methods"`     // 允许的 HTTP 方法列表
	AllowHeaders     []string `mapstructure:"allow_headers"`     // 允许的请求头列表
	ExposeHeaders    []string `mapstructure:"expose_headers"`    // 暴露给客户端的响应头列表
	AllowCredentials bool     `mapstructure:"allow_credentials"` // 是否允许携带凭证（Cookie）
	MaxAge           int      `mapstructure:"max_age"`           // 预检请求的缓存时间（秒）
}

// Load 从默认配置文件路径加载配置
//
// 从 "config/config.yaml" 加载应用程序配置。
// 这是 LoadWithPath 的便捷封装。
//
// 返回值：
//   - *Config: 加载的配置对象
//   - error: 加载或解析过程中的错误
//
// 示例：
//   cfg, err := config.Load()
//   if err != nil {
//       log.Fatalf("加载配置失败: %v", err)
//   }
func Load() (*Config, error) {
	return LoadWithPath("config/config.yaml")
}

// LoadWithPath 从指定路径加载配置文件
//
// 从指定的文件路径加载 YAML 格式的配置文件，
// 并支持环境变量覆盖配置项。
//
// 加载流程：
//   1. 读取指定路径的 YAML 配置文件
//   2. 解析 YAML 内容到 Config 结构体
//   3. 允许环境变量覆盖配置项
//
// 环境变量覆盖规则：
//   - 环境变量名称格式：CONFIG_SECTION_FIELD（如：DATABASE_HOST）
//   - 会自动覆盖配置文件中的同名配置
//
// 参数：
//   - path: 配置文件的路径（相对或绝对路径）
//
// 返回值：
//   - *Config: 加载的配置对象
//   - error: 加载或解析过程中的错误
//     - fmt.Errorf("读取配置文件失败: %w", err): 文件不存在或无法读取
//     - fmt.Errorf("解析配置失败: %w", err): YAML 格式错误或字段类型不匹配
//
// 示例：
//   cfg, err := config.LoadWithPath("/etc/qzt/config.yaml")
//   if err != nil {
//       log.Fatalf("加载配置失败: %v", err)
//   }
func LoadWithPath(path string) (*Config, error) {
	// 创建 Viper 实例
	v := viper.New()
	
	// 设置配置文件路径和类型
	v.SetConfigFile(path)
	v.SetConfigType("yaml")
	
	// 允许环境变量覆盖配置文件
	// 环境变量格式：DATABASE_HOST 会覆盖 database.host
	v.AutomaticEnv()
	
	// 读取配置文件
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}
	
	// 将配置解析到结构体
	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("解析配置失败: %w", err)
	}
	
	return &cfg, nil
}
