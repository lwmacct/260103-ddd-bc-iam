package auth

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"

	domainAuth "github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/domain/auth"
)

// newTestAuthService 创建测试用认证服务。
func newTestAuthService() domainAuth.Service {
	jwtManager := NewJWTManager("test-secret-key-for-testing", time.Hour, 24*time.Hour)
	tokenGenerator := NewTokenGenerator()
	return NewAuthService(jwtManager, tokenGenerator, nil)
}

// newTestAuthServiceWithPolicy 创建带自定义密码策略的测试服务。
func newTestAuthServiceWithPolicy(policy *domainAuth.PasswordPolicy) domainAuth.Service {
	jwtManager := NewJWTManager("test-secret-key-for-testing", time.Hour, 24*time.Hour)
	tokenGenerator := NewTokenGenerator()
	return NewAuthService(jwtManager, tokenGenerator, policy)
}

func TestPasswordPolicy_Validate(t *testing.T) {
	tests := []struct {
		name     string
		policy   *domainAuth.PasswordPolicy
		password string
		want     bool
	}{
		{
			name:     "默认策略 - 符合最小长度",
			policy:   domainAuth.DefaultPasswordPolicy(),
			password: "123456",
			want:     true,
		},
		{
			name:     "默认策略 - 不符合最小长度",
			policy:   domainAuth.DefaultPasswordPolicy(),
			password: "12345",
			want:     false,
		},
		{
			name: "要求大写字母 - 符合",
			policy: &domainAuth.PasswordPolicy{
				MinLength:    6,
				RequireUpper: true,
			},
			password: "Abcdef",
			want:     true,
		},
		{
			name: "要求大写字母 - 不符合",
			policy: &domainAuth.PasswordPolicy{
				MinLength:    6,
				RequireUpper: true,
			},
			password: "abcdef",
			want:     false,
		},
		{
			name: "要求小写字母 - 符合",
			policy: &domainAuth.PasswordPolicy{
				MinLength:    6,
				RequireLower: true,
			},
			password: "ABCDEf",
			want:     true,
		},
		{
			name: "要求小写字母 - 不符合",
			policy: &domainAuth.PasswordPolicy{
				MinLength:    6,
				RequireLower: true,
			},
			password: "ABCDEF",
			want:     false,
		},
		{
			name: "要求数字 - 符合",
			policy: &domainAuth.PasswordPolicy{
				MinLength:     6,
				RequireNumber: true,
			},
			password: "abcde1",
			want:     true,
		},
		{
			name: "要求数字 - 不符合",
			policy: &domainAuth.PasswordPolicy{
				MinLength:     6,
				RequireNumber: true,
			},
			password: "abcdef",
			want:     false,
		},
		{
			name: "要求特殊字符 - 符合",
			policy: &domainAuth.PasswordPolicy{
				MinLength:      6,
				RequireSpecial: true,
			},
			password: "abcde!",
			want:     true,
		},
		{
			name: "要求特殊字符 - 不符合",
			policy: &domainAuth.PasswordPolicy{
				MinLength:      6,
				RequireSpecial: true,
			},
			password: "abcdef",
			want:     false,
		},
		{
			name: "复杂密码策略 - 符合所有要求",
			policy: &domainAuth.PasswordPolicy{
				MinLength:      8,
				RequireUpper:   true,
				RequireLower:   true,
				RequireNumber:  true,
				RequireSpecial: true,
			},
			password: "Abcdef1!",
			want:     true,
		},
		{
			name: "复杂密码策略 - 缺少大写",
			policy: &domainAuth.PasswordPolicy{
				MinLength:      8,
				RequireUpper:   true,
				RequireLower:   true,
				RequireNumber:  true,
				RequireSpecial: true,
			},
			password: "abcdef1!",
			want:     false,
		},
		{
			name: "支持的特殊字符集",
			policy: &domainAuth.PasswordPolicy{
				MinLength:      1,
				RequireSpecial: true,
			},
			password: "!@#$%^&*",
			want:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.policy.Validate(tt.password)
			assert.Equal(t, tt.want, got, "PasswordPolicy.Validate(%q)", tt.password)
		})
	}
}

func TestAuthService_ValidatePasswordPolicy(t *testing.T) {
	ctx := context.Background()

	t.Run("默认策略 - 有效密码", func(t *testing.T) {
		svc := newTestAuthService()
		err := svc.ValidatePasswordPolicy(ctx, "password123")
		assert.NoError(t, err, "ValidatePasswordPolicy() 应该成功")
	})

	t.Run("默认策略 - 密码太短", func(t *testing.T) {
		svc := newTestAuthService()
		err := svc.ValidatePasswordPolicy(ctx, "12345")
		assert.ErrorIs(t, err, domainAuth.ErrWeakPassword, "应该返回 ErrWeakPassword")
	})

	t.Run("自定义严格策略", func(t *testing.T) {
		strictPolicy := &domainAuth.PasswordPolicy{
			MinLength:      10,
			RequireUpper:   true,
			RequireLower:   true,
			RequireNumber:  true,
			RequireSpecial: true,
		}
		svc := newTestAuthServiceWithPolicy(strictPolicy)

		// 符合策略
		err := svc.ValidatePasswordPolicy(ctx, "MySecure1!")
		require.NoError(t, err, "符合策略的密码应该通过验证")

		// 不符合策略
		err = svc.ValidatePasswordPolicy(ctx, "simple")
		assert.ErrorIs(t, err, domainAuth.ErrWeakPassword, "不符合策略的密码应该返回 ErrWeakPassword")
	})
}

func TestAuthService_GeneratePasswordHash(t *testing.T) {
	ctx := context.Background()
	svc := newTestAuthService()

	t.Run("成功生成密码哈希", func(t *testing.T) {
		password := "mySecurePassword123"
		hash, err := svc.GeneratePasswordHash(ctx, password)

		require.NoError(t, err, "GeneratePasswordHash() 应该成功")
		assert.NotEmpty(t, hash, "GeneratePasswordHash() 不应返回空哈希")

		// 验证哈希格式（bcrypt 哈希以 $2a$ 或 $2b$ 开头）
		assert.True(t, strings.HasPrefix(hash, "$2a$") || strings.HasPrefix(hash, "$2b$"),
			"哈希格式应该是 bcrypt 格式")
	})

	t.Run("相同密码产生不同哈希（盐值）", func(t *testing.T) {
		password := "samePassword"
		hash1, _ := svc.GeneratePasswordHash(ctx, password)
		hash2, _ := svc.GeneratePasswordHash(ctx, password)

		assert.NotEqual(t, hash1, hash2, "相同密码应该产生不同哈希（由于盐值）")
	})

	t.Run("空密码也能生成哈希", func(t *testing.T) {
		hash, err := svc.GeneratePasswordHash(ctx, "")

		require.NoError(t, err, "空密码也应该能生成哈希")
		assert.NotEmpty(t, hash, "空密码也应该产生非空哈希")
	})
}

func TestAuthService_VerifyPassword(t *testing.T) {
	ctx := context.Background()
	svc := newTestAuthService()

	t.Run("正确密码验证成功", func(t *testing.T) {
		password := "correctPassword"
		hash, _ := svc.GeneratePasswordHash(ctx, password)

		err := svc.VerifyPassword(ctx, hash, password)

		assert.NoError(t, err, "正确密码应该验证成功")
	})

	t.Run("错误密码验证失败", func(t *testing.T) {
		password := "correctPassword"
		hash, _ := svc.GeneratePasswordHash(ctx, password)

		err := svc.VerifyPassword(ctx, hash, "wrongPassword")

		assert.ErrorIs(t, err, domainAuth.ErrPasswordMismatch, "错误密码应该返回 ErrPasswordMismatch")
	})

	t.Run("无效哈希验证失败", func(t *testing.T) {
		err := svc.VerifyPassword(ctx, "invalid-hash", "password")

		assert.ErrorIs(t, err, domainAuth.ErrPasswordMismatch, "无效哈希应该返回 ErrPasswordMismatch")
	})

	t.Run("空密码和空哈希", func(t *testing.T) {
		// 空密码的哈希
		hash, _ := svc.GeneratePasswordHash(ctx, "")
		err := svc.VerifyPassword(ctx, hash, "")

		assert.NoError(t, err, "空密码验证应该成功")
	})
}

func TestAuthService_GenerateAccessToken(t *testing.T) {
	ctx := context.Background()
	svc := newTestAuthService()

	t.Run("成功生成访问令牌", func(t *testing.T) {
		token, expiresAt, err := svc.GenerateAccessToken(ctx, 123, "testuser")

		require.NoError(t, err, "GenerateAccessToken() 应该成功")
		assert.NotEmpty(t, token, "GenerateAccessToken() 不应返回空令牌")
		assert.True(t, expiresAt.After(time.Now()), "过期时间应该在未来")
	})

	t.Run("不同用户生成不同令牌", func(t *testing.T) {
		token1, _, _ := svc.GenerateAccessToken(ctx, 1, "user1")
		token2, _, _ := svc.GenerateAccessToken(ctx, 2, "user2")

		assert.NotEqual(t, token1, token2, "不同用户应该生成不同令牌")
	})
}

func TestAuthService_GenerateRefreshToken(t *testing.T) {
	ctx := context.Background()
	svc := newTestAuthService()

	t.Run("成功生成刷新令牌", func(t *testing.T) {
		token, expiresAt, err := svc.GenerateRefreshToken(ctx, 123)

		require.NoError(t, err, "GenerateRefreshToken() 应该成功")
		assert.NotEmpty(t, token, "GenerateRefreshToken() 不应返回空令牌")
		assert.True(t, expiresAt.After(time.Now()), "过期时间应该在未来")
	})

	t.Run("刷新令牌过期时间比访问令牌长", func(t *testing.T) {
		_, accessExpires, _ := svc.GenerateAccessToken(ctx, 1, "user")
		_, refreshExpires, _ := svc.GenerateRefreshToken(ctx, 1)

		assert.True(t, refreshExpires.After(accessExpires), "刷新令牌过期时间应该比访问令牌长")
	})
}

func TestAuthService_ValidateAccessToken(t *testing.T) {
	ctx := context.Background()
	svc := newTestAuthService()

	t.Run("验证有效令牌", func(t *testing.T) {
		token, _, _ := svc.GenerateAccessToken(ctx, 123, "testuser")

		claims, err := svc.ValidateAccessToken(ctx, token)

		require.NoError(t, err, "ValidateAccessToken() 应该成功")
		assert.Equal(t, uint(123), claims.UserID, "UserID 应该匹配")
		assert.Equal(t, "testuser", claims.Username, "Username 应该匹配")
	})

	t.Run("无效令牌", func(t *testing.T) {
		_, err := svc.ValidateAccessToken(ctx, "invalid.token.here")

		assert.ErrorIs(t, err, domainAuth.ErrInvalidToken, "无效令牌应该返回 ErrInvalidToken")
	})

	t.Run("空令牌", func(t *testing.T) {
		_, err := svc.ValidateAccessToken(ctx, "")

		assert.ErrorIs(t, err, domainAuth.ErrInvalidToken, "空令牌应该返回 ErrInvalidToken")
	})

	t.Run("过期令牌", func(t *testing.T) {
		// 创建一个快速过期的服务
		jwtManager := NewJWTManager("test-secret", time.Nanosecond, time.Hour)
		tokenGenerator := NewTokenGenerator()
		quickExpirySvc := NewAuthService(jwtManager, tokenGenerator, nil)

		token, _, _ := quickExpirySvc.GenerateAccessToken(ctx, 1, "user")
		time.Sleep(time.Millisecond * 10) // 等待令牌过期

		_, err := quickExpirySvc.ValidateAccessToken(ctx, token)
		// 注意：JWT 库在解析时会检查过期，返回解析错误
		// 所以实际返回 ErrInvalidToken（因为解析失败）而不是 ErrTokenExpired
		//
		assert.True(t, errors.Is(err, domainAuth.ErrInvalidToken) || errors.Is(err, domainAuth.ErrTokenExpired),
			"过期令牌应该返回 ErrInvalidToken 或 ErrTokenExpired")
	})
}

func TestAuthService_ValidateRefreshToken(t *testing.T) {
	ctx := context.Background()
	svc := newTestAuthService()

	t.Run("验证有效刷新令牌", func(t *testing.T) {
		token, _, _ := svc.GenerateRefreshToken(ctx, 456)

		userID, err := svc.ValidateRefreshToken(ctx, token)

		require.NoError(t, err, "ValidateRefreshToken() 应该成功")
		assert.Equal(t, uint(456), userID, "UserID 应该匹配")
	})

	t.Run("无效刷新令牌", func(t *testing.T) {
		_, err := svc.ValidateRefreshToken(ctx, "invalid.refresh.token")

		assert.ErrorIs(t, err, domainAuth.ErrInvalidToken, "无效刷新令牌应该返回 ErrInvalidToken")
	})
}

func TestAuthService_GeneratePATToken(t *testing.T) {
	ctx := context.Background()
	svc := newTestAuthService()

	t.Run("成功生成 PAT 令牌", func(t *testing.T) {
		token, err := svc.GeneratePATToken(ctx)

		require.NoError(t, err, "GeneratePATToken() 应该成功")
		assert.NotEmpty(t, token, "GeneratePATToken() 不应返回空令牌")
		assert.True(t, strings.HasPrefix(token, "pat_"), "令牌应该以 'pat_' 开头")
	})

	t.Run("每次生成不同的 PAT 令牌", func(t *testing.T) {
		token1, _ := svc.GeneratePATToken(ctx)
		token2, _ := svc.GeneratePATToken(ctx)

		assert.NotEqual(t, token1, token2, "每次生成的 PAT 令牌应该不同")
	})
}

func TestAuthService_HashPATToken(t *testing.T) {
	ctx := context.Background()
	svc := newTestAuthService()

	t.Run("哈希 PAT 令牌", func(t *testing.T) {
		token, _ := svc.GeneratePATToken(ctx)
		hash := svc.HashPATToken(ctx, token)

		assert.NotEmpty(t, hash, "HashPATToken() 不应返回空哈希")
		// SHA-256 哈希应该是 64 个十六进制字符
		assert.Len(t, hash, 64, "哈希长度应该是 64")
	})

	t.Run("相同令牌产生相同哈希", func(t *testing.T) {
		token := "pat_abc12_" + strings.Repeat("x", 32)
		hash1 := svc.HashPATToken(ctx, token)
		hash2 := svc.HashPATToken(ctx, token)

		assert.Equal(t, hash1, hash2, "相同令牌应该产生相同哈希")
	})

	t.Run("不同令牌产生不同哈希", func(t *testing.T) {
		token1, _ := svc.GeneratePATToken(ctx)
		token2, _ := svc.GeneratePATToken(ctx)
		hash1 := svc.HashPATToken(ctx, token1)
		hash2 := svc.HashPATToken(ctx, token2)

		assert.NotEqual(t, hash1, hash2, "不同令牌应该产生不同哈希")
	})
}

func TestTokenClaims_IsExpired(t *testing.T) {
	tests := []struct {
		name string
		exp  int64
		want bool
	}{
		{
			name: "未过期",
			exp:  time.Now().Add(time.Hour).Unix(),
			want: false,
		},
		{
			name: "已过期",
			exp:  time.Now().Add(-time.Hour).Unix(),
			want: true,
		},
		{
			name: "刚好过期",
			exp:  time.Now().Unix() - 1,
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tc := &domainAuth.TokenClaims{Exp: tt.exp}
			got := tc.IsExpired()
			assert.Equal(t, tt.want, got, "TokenClaims.IsExpired()")
		})
	}
}

func BenchmarkAuthService_GeneratePasswordHash(b *testing.B) {
	ctx := context.Background()
	svc := newTestAuthService()

	for b.Loop() {
		_, _ = svc.GeneratePasswordHash(ctx, "testPassword123")
	}
}

func BenchmarkAuthService_VerifyPassword(b *testing.B) {
	ctx := context.Background()
	svc := newTestAuthService()
	hash, _ := bcrypt.GenerateFromPassword([]byte("testPassword123"), bcrypt.DefaultCost)

	for b.Loop() {
		_ = svc.VerifyPassword(ctx, string(hash), "testPassword123")
	}
}

func BenchmarkAuthService_GenerateAccessToken(b *testing.B) {
	ctx := context.Background()
	svc := newTestAuthService()

	for b.Loop() {
		_, _, _ = svc.GenerateAccessToken(ctx, 1, "user")
	}
}

func BenchmarkAuthService_ValidateAccessToken(b *testing.B) {
	ctx := context.Background()
	svc := newTestAuthService()
	token, _, _ := svc.GenerateAccessToken(ctx, 1, "user")

	for b.Loop() {
		_, _ = svc.ValidateAccessToken(ctx, token)
	}
}
