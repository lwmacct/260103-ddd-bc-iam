package auth

import (
	"context"
	"fmt"
	"time"

	domainAuth "github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/iam/domain/auth"
	"golang.org/x/crypto/bcrypt"
)

// authServiceImpl 认证服务实现
type authServiceImpl struct {
	jwtManager     *JWTManager
	tokenGenerator *TokenGenerator
	passwordPolicy *domainAuth.PasswordPolicy
}

// NewAuthService 创建认证服务实例
func NewAuthService(
	jwtManager *JWTManager,
	tokenGenerator *TokenGenerator,
	passwordPolicy *domainAuth.PasswordPolicy,
) domainAuth.Service {
	if passwordPolicy == nil {
		passwordPolicy = domainAuth.DefaultPasswordPolicy()
	}
	return &authServiceImpl{
		jwtManager:     jwtManager,
		tokenGenerator: tokenGenerator,
		passwordPolicy: passwordPolicy,
	}
}

// VerifyPassword 验证密码是否正确
func (s *authServiceImpl) VerifyPassword(ctx context.Context, hashedPassword, plainPassword string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword)); err != nil {
		return domainAuth.ErrPasswordMismatch
	}
	return nil
}

// GeneratePasswordHash 生成密码哈希
func (s *authServiceImpl) GeneratePasswordHash(ctx context.Context, password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hashedBytes), nil
}

// ValidatePasswordPolicy 验证密码是否符合策略
func (s *authServiceImpl) ValidatePasswordPolicy(ctx context.Context, password string) error {
	if !s.passwordPolicy.Validate(password) {
		return domainAuth.ErrWeakPassword
	}
	return nil
}

// GenerateAccessToken 生成访问令牌
// 新架构：Token 只包含 user_id/username，权限信息从缓存实时查询
func (s *authServiceImpl) GenerateAccessToken(ctx context.Context, userID uint, username string) (string, time.Time, error) {
	token, err := s.jwtManager.GenerateAccessToken(userID, username, "")
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to generate access token: %w", err)
	}

	expiresAt := time.Now().Add(s.jwtManager.accessTokenDuration)
	return token, expiresAt, nil
}

// GenerateRefreshToken 生成刷新令牌
func (s *authServiceImpl) GenerateRefreshToken(ctx context.Context, userID uint) (string, time.Time, error) {
	token, err := s.jwtManager.GenerateRefreshToken(userID)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	expiresAt := time.Now().Add(s.jwtManager.refreshTokenDuration)
	return token, expiresAt, nil
}

// ValidateAccessToken 验证访问令牌
func (s *authServiceImpl) ValidateAccessToken(ctx context.Context, token string) (*domainAuth.TokenClaims, error) {
	claims, err := s.jwtManager.ValidateToken(token)
	if err != nil {
		return nil, domainAuth.ErrInvalidToken
	}

	if time.Unix(claims.ExpiresAt.Unix(), 0).Before(time.Now()) {
		return nil, domainAuth.ErrTokenExpired
	}

	return &domainAuth.TokenClaims{
		UserID:   claims.UserID,
		Username: claims.Username,
		Roles:    nil, // 角色信息从缓存中获取，不存储在 JWT 中
		Exp:      claims.ExpiresAt.Unix(),
	}, nil
}

// ValidateRefreshToken 验证刷新令牌
func (s *authServiceImpl) ValidateRefreshToken(ctx context.Context, token string) (uint, error) {
	claims, err := s.jwtManager.ValidateToken(token)
	if err != nil {
		return 0, domainAuth.ErrInvalidToken
	}

	if time.Unix(claims.ExpiresAt.Unix(), 0).Before(time.Now()) {
		return 0, domainAuth.ErrTokenExpired
	}

	return claims.UserID, nil
}

// GeneratePATToken 生成个人访问令牌
func (s *authServiceImpl) GeneratePATToken(ctx context.Context) (string, error) {
	plainToken, _, _, err := s.tokenGenerator.GeneratePAT()
	if err != nil {
		return "", fmt.Errorf("failed to generate PAT: %w", err)
	}
	return plainToken, nil
}

// HashPATToken 哈希个人访问令牌（用于存储）
func (s *authServiceImpl) HashPATToken(ctx context.Context, token string) string {
	return s.tokenGenerator.HashToken(token)
}
