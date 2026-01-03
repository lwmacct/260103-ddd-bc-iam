package auth

import (
	"context"
	"time"
)

// PasswordPolicy 密码策略值对象。
// 定义密码强度要求，可根据安全需求配置。
type PasswordPolicy struct {
	MinLength      int
	RequireUpper   bool
	RequireLower   bool
	RequireNumber  bool
	RequireSpecial bool
}

// Validate 验证密码是否符合策略
func (p *PasswordPolicy) Validate(password string) bool {
	if len(password) < p.MinLength {
		return false
	}

	var hasUpper, hasLower, hasNumber, hasSpecial bool
	for _, char := range password {
		switch {
		case 'A' <= char && char <= 'Z':
			hasUpper = true
		case 'a' <= char && char <= 'z':
			hasLower = true
		case '0' <= char && char <= '9':
			hasNumber = true
		case char == '!' || char == '@' || char == '#' || char == '$' || char == '%' || char == '^' || char == '&' || char == '*':
			hasSpecial = true
		}
	}

	if p.RequireUpper && !hasUpper {
		return false
	}
	if p.RequireLower && !hasLower {
		return false
	}
	if p.RequireNumber && !hasNumber {
		return false
	}
	if p.RequireSpecial && !hasSpecial {
		return false
	}

	return true
}

// DefaultPasswordPolicy 默认密码策略
func DefaultPasswordPolicy() *PasswordPolicy {
	return &PasswordPolicy{
		MinLength:      6,
		RequireUpper:   false,
		RequireLower:   false,
		RequireNumber:  false,
		RequireSpecial: false,
	}
}

// Service 认证领域服务接口
// 定义密码管理、Token 生成等领域能力
type Service interface {
	// VerifyPassword 验证密码是否正确
	VerifyPassword(ctx context.Context, hashedPassword, plainPassword string) error

	// GeneratePasswordHash 生成密码哈希
	GeneratePasswordHash(ctx context.Context, password string) (string, error)

	// ValidatePasswordPolicy 验证密码是否符合策略
	ValidatePasswordPolicy(ctx context.Context, password string) error

	// GenerateAccessToken 生成访问令牌
	// 新架构：Token 只包含 user_id/username，权限信息从缓存实时查询
	GenerateAccessToken(ctx context.Context, userID uint, username string) (string, time.Time, error)

	// GenerateRefreshToken 生成刷新令牌
	GenerateRefreshToken(ctx context.Context, userID uint) (string, time.Time, error)

	// ValidateAccessToken 验证访问令牌
	ValidateAccessToken(ctx context.Context, token string) (*TokenClaims, error)

	// ValidateRefreshToken 验证刷新令牌
	ValidateRefreshToken(ctx context.Context, token string) (uint, error)

	// GeneratePATToken 生成个人访问令牌
	GeneratePATToken(ctx context.Context) (string, error)

	// HashPATToken 哈希个人访问令牌（用于存储）
	HashPATToken(ctx context.Context, token string) string
}

// TokenClaims Token 声明
type TokenClaims struct {
	UserID   uint     `json:"user_id"`
	Username string   `json:"username"`
	Roles    []string `json:"roles"`
	Exp      int64    `json:"exp"`
}

// IsExpired 检查 Token 是否过期
func (tc *TokenClaims) IsExpired() bool {
	return time.Now().Unix() > tc.Exp
}
