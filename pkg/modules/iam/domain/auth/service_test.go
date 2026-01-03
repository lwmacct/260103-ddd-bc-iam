package auth

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPasswordPolicy_Validate(t *testing.T) {
	tests := []struct {
		name     string
		policy   *PasswordPolicy
		password string
		want     bool
	}{
		{
			name: "满足所有要求",
			policy: &PasswordPolicy{
				MinLength:      8,
				RequireUpper:   true,
				RequireLower:   true,
				RequireNumber:  true,
				RequireSpecial: true,
			},
			password: "Abc123!@",
			want:     true,
		},
		{
			name: "密码太短",
			policy: &PasswordPolicy{
				MinLength: 8,
			},
			password: "abc123",
			want:     false,
		},
		{
			name: "缺少大写字母",
			policy: &PasswordPolicy{
				MinLength:    6,
				RequireUpper: true,
			},
			password: "abc123",
			want:     false,
		},
		{
			name: "缺少小写字母",
			policy: &PasswordPolicy{
				MinLength:    6,
				RequireLower: true,
			},
			password: "ABC123",
			want:     false,
		},
		{
			name: "缺少数字",
			policy: &PasswordPolicy{
				MinLength:     6,
				RequireNumber: true,
			},
			password: "abcdef",
			want:     false,
		},
		{
			name: "缺少特殊字符",
			policy: &PasswordPolicy{
				MinLength:      6,
				RequireSpecial: true,
			},
			password: "abc123",
			want:     false,
		},
		{
			name: "包含特殊字符 !",
			policy: &PasswordPolicy{
				MinLength:      6,
				RequireSpecial: true,
			},
			password: "abc12!",
			want:     true,
		},
		{
			name: "包含特殊字符 @",
			policy: &PasswordPolicy{
				MinLength:      6,
				RequireSpecial: true,
			},
			password: "abc12@",
			want:     true,
		},
		{
			name: "包含特殊字符 #",
			policy: &PasswordPolicy{
				MinLength:      6,
				RequireSpecial: true,
			},
			password: "abc12#",
			want:     true,
		},
		{
			name: "包含特殊字符 $",
			policy: &PasswordPolicy{
				MinLength:      6,
				RequireSpecial: true,
			},
			password: "abc12$",
			want:     true,
		},
		{
			name: "包含特殊字符 %",
			policy: &PasswordPolicy{
				MinLength:      6,
				RequireSpecial: true,
			},
			password: "abc12%",
			want:     true,
		},
		{
			name: "包含特殊字符 ^",
			policy: &PasswordPolicy{
				MinLength:      6,
				RequireSpecial: true,
			},
			password: "abc12^",
			want:     true,
		},
		{
			name: "包含特殊字符 &",
			policy: &PasswordPolicy{
				MinLength:      6,
				RequireSpecial: true,
			},
			password: "abc12&",
			want:     true,
		},
		{
			name: "包含特殊字符 *",
			policy: &PasswordPolicy{
				MinLength:      6,
				RequireSpecial: true,
			},
			password: "abc12*",
			want:     true,
		},
		{
			name:     "默认策略 - 简单密码",
			policy:   DefaultPasswordPolicy(),
			password: "123456",
			want:     true,
		},
		{
			name:     "默认策略 - 密码太短",
			policy:   DefaultPasswordPolicy(),
			password: "12345",
			want:     false,
		},
		{
			name: "空密码",
			policy: &PasswordPolicy{
				MinLength: 1,
			},
			password: "",
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.policy.Validate(tt.password)
			assert.Equal(t, tt.want, got, "PasswordPolicy.Validate(%q)", tt.password)
		})
	}
}

func TestDefaultPasswordPolicy(t *testing.T) {
	policy := DefaultPasswordPolicy()

	assert.Equal(t, 6, policy.MinLength, "默认最小长度应该是 6")
	assert.False(t, policy.RequireUpper, "默认不要求大写")
	assert.False(t, policy.RequireLower, "默认不要求小写")
	assert.False(t, policy.RequireNumber, "默认不要求数字")
	assert.False(t, policy.RequireSpecial, "默认不要求特殊字符")
}

func TestTokenClaims_IsExpired(t *testing.T) {
	tests := []struct {
		name string
		exp  int64
		want bool
	}{
		{
			name: "已过期",
			exp:  time.Now().Add(-1 * time.Hour).Unix(),
			want: true,
		},
		{
			name: "未过期",
			exp:  time.Now().Add(1 * time.Hour).Unix(),
			want: false,
		},
		{
			name: "刚刚过期",
			exp:  time.Now().Add(-1 * time.Second).Unix(),
			want: true,
		},
		{
			name: "零值过期时间",
			exp:  0,
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims := &TokenClaims{
				UserID:   1,
				Username: "testuser",
				Exp:      tt.exp,
			}
			got := claims.IsExpired()
			assert.Equal(t, tt.want, got, "TokenClaims.IsExpired()")
		})
	}
}

func TestTokenClaims_Fields(t *testing.T) {
	claims := &TokenClaims{
		UserID:   42,
		Username: "john_doe",
		Roles:    []string{"admin", "editor"},
		Exp:      time.Now().Add(1 * time.Hour).Unix(),
	}

	assert.Equal(t, uint(42), claims.UserID)
	assert.Equal(t, "john_doe", claims.Username)
	assert.Equal(t, []string{"admin", "editor"}, claims.Roles)
	assert.False(t, claims.IsExpired())
}
