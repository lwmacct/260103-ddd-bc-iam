package auth

// TokenGenerator 定义 PAT Token 生成器的领域接口。
// 用于生成、哈希和验证个人访问令牌。
//
// 实现：internal/infrastructure/auth/token_generator.go
type TokenGenerator interface {
	// GeneratePAT 生成新的个人访问令牌
	// 返回值：
	//   - plainToken: 明文令牌 (格式: pat_<prefix>_<random>)，仅创建时返回一次
	//   - hashedToken: SHA-256 哈希后的令牌，用于存储
	//   - prefix: 令牌前缀 (格式: pat_<prefix>)，用于快速查找
	//   - error: 生成失败时返回错误
	GeneratePAT() (plainToken, hashedToken, prefix string, err error)

	// HashToken 对明文令牌进行 SHA-256 哈希
	// 用于令牌验证时比对存储的哈希值
	HashToken(plainToken string) string

	// ValidateTokenFormat 验证令牌格式是否正确
	// 期望格式: pat_<5chars>_<32chars>
	ValidateTokenFormat(token string) bool

	// ExtractPrefix 从完整令牌中提取前缀
	// 示例: "pat_2Kj9X_abc123..." -> "pat_2Kj9X"
	ExtractPrefix(token string) (string, error)
}
