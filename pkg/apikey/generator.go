// Package apikey 提供 API Key 生成和验证功能
//
// 该包实现了安全的 API Key 生成机制，用于第三方系统调用本系统 API 时的身份验证。
//
// 主要功能：
//   - 生成安全的随机 API Key
//   - 支持生产和测试环境的 Key 前缀区分
//   - 提供 Key 格式验证
//   - 提供 Key 哈希计算（用于安全存储）
//
// 安全特性：
//   - 使用加密安全的随机数生成器（crypto/rand）
//   - 不存储原始 Key，只存储 SHA256 哈希值
//   - Key 包含前缀便于识别和调试
//
// Key 格式：
//   - 生产环境：qzt_live_xxxxxxxxxxxxxxxxxxxxxxxx
//   - 测试环境：qzt_test_xxxxxxxxxxxxxxxxxxxxxxxx
//
// 使用示例：
//   // 生成新的 API Key
//   rawKey, prefix, hash, err := apikey.Generate(true)
//   if err != nil {
//       log.Fatal(err)
//   }
//   fmt.Println("API Key:", rawKey)        // 仅显示一次，请妥善保存
//   fmt.Println("Prefix:", prefix)          // 用于识别
//   fmt.Println("Hash:", hash)              // 用于存储到数据库
//
//   // 验证 Key 格式
//   if apikey.Validate(userKey) {
//       // 验证通过，查询数据库
//       keyHash := apikey.Hash(userKey)
//       // 使用 keyHash 查询数据库
//   }
package apikey

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
)

// API Key 前缀常量定义
const (
	prefixLive = "qzt_live_" // 生产环境 Key 前缀
	prefixTest = "qzt_test_" // 测试环境 Key 前缀
	keyLength  = 32          // 随机部分的字节长度
)

// Generate 生成 API Key
//
// 生成一个安全的随机 API Key，包含前缀和随机部分。
// 同时计算 Key 的 SHA256 哈希值，用于安全存储。
//
// 参数：
//   - isProduction: 是否为生产环境
//     - true: 使用 "qzt_live_" 前缀
//     - false: 使用 "qzt_test_" 前缀
//
// 返回值：
//   - string: 原始 API Key（完整 Key，仅此一次显示，请妥善保存）
//   - string: Key 前缀（用于识别和显示，如 "qzt_live_abcd"）
//   - string: Key 哈希值（SHA256，用于存储到数据库）
//   - error: 生成过程中的错误（通常是随机数生成失败）
//
// 安全说明：
//   - 原始 Key 仅在生成时返回一次，请立即保存
//   - 数据库中只存储哈希值，不存储原始 Key
//   - 使用 crypto/rand 生成加密安全的随机数
//
// 示例：
//   rawKey, prefix, hash, err := apikey.Generate(true)
//   // rawKey: "qzt_live_a1b2c3d4e5f6g7h8i9j0"
//   // prefix: "qzt_live_a1b2"
//   // hash:   "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
func Generate(isProduction bool) (string, string, string, error) {
	// 根据环境选择前缀
	prefix := prefixTest
	if isProduction {
		prefix = prefixLive
	}

	// 生成加密安全的随机字节
	bytes := make([]byte, keyLength)
	if _, err := rand.Read(bytes); err != nil {
		return "", "", "", err
	}

	// 将随机字节转换为十六进制字符串
	randomPart := hex.EncodeToString(bytes)

	// 构建完整的 API Key（前缀 + 随机部分的前 24 个字符）
	apiKey := prefix + randomPart[:24] // qzt_live_xxxxxxxxxxxxxxxxxxxxxxxx

	// 提取 Key 前缀用于显示（前 12 个字符）
	keyPrefix := apiKey[:12] // qzt_live_abcd

	// 计算 SHA256 哈希值用于安全存储
	hash := sha256.Sum256([]byte(apiKey))
	keyHash := hex.EncodeToString(hash[:])

	return apiKey, keyPrefix, keyHash, nil
}

// Validate 验证 API Key 格式
//
// 检查给定的字符串是否符合 API Key 的格式要求。
// 仅验证格式，不验证 Key 是否存在或有效。
//
// 验证规则：
//   - 长度必须至少 20 个字符
//   - 必须以 "qzt_live_" 或 "qzt_test_" 开头
//
// 参数：
//   - key: 待验证的 API Key 字符串
//
// 返回值：
//   - bool: true 表示格式正确，false 表示格式错误
//
// 示例：
//   if apikey.Validate(userKey) {
//       // 格式正确，继续验证
//       hash := apikey.Hash(userKey)
//       // 查询数据库验证 hash
//   } else {
//       // 格式错误，拒绝访问
//   }
func Validate(key string) bool {
	// 检查最小长度
	if len(key) < 20 {
		return false
	}
	// 检查前缀是否正确
	return key[:9] == prefixLive || key[:9] == prefixTest
}

// Hash 计算 API Key 的哈希值
//
// 使用 SHA256 算法计算 API Key 的哈希值。
// 用于在数据库中查询和比对 Key。
//
// 参数：
//   - key: 原始 API Key 字符串
//
// 返回值：
//   - string: SHA256 哈希值（十六进制字符串）
//
// 使用场景：
//   - 用户提交 Key 后，计算哈希值
//   - 使用哈希值在数据库中查询 Key 信息
//
// 示例：
//   userKey := "qzt_live_a1b2c3d4e5f6g7h8i9j0"
//   hash := apikey.Hash(userKey)
//   // 使用 hash 查询数据库
//   keyInfo, err := repo.GetByHash(ctx, hash)
func Hash(key string) string {
	// 计算 SHA256 哈希
	hash := sha256.Sum256([]byte(key))
	// 转换为十六进制字符串
	return hex.EncodeToString(hash[:])
}
