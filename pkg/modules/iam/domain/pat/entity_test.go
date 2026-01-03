package pat

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/iam/domain/role"
)

// newTestPAT 创建测试用 PAT。
func newTestPAT(status string, expiresAt *time.Time) *PersonalAccessToken {
	return &PersonalAccessToken{
		ID:          1,
		UserID:      100,
		Name:        "Test Token",
		Token:       "hashed_token_value",
		TokenPrefix: "ghp_xxxx",
		Scopes:      StringList{string(ScopeFull)},
		Status:      status,
		ExpiresAt:   expiresAt,
		CreatedAt:   time.Now(),
	}
}

func TestPersonalAccessToken_IsExpired(t *testing.T) {
	now := time.Now()
	pastTime := now.Add(-24 * time.Hour)
	futureTime := now.Add(24 * time.Hour)

	tests := []struct {
		name      string
		expiresAt *time.Time
		want      bool
	}{
		{
			name:      "无过期时间 - 永久有效",
			expiresAt: nil,
			want:      false,
		},
		{
			name:      "已过期",
			expiresAt: &pastTime,
			want:      true,
		},
		{
			name:      "未过期",
			expiresAt: &futureTime,
			want:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pat := newTestPAT("active", tt.expiresAt)
			got := pat.IsExpired()
			assert.Equal(t, tt.want, got, "PersonalAccessToken.IsExpired()")
		})
	}
}

func TestPersonalAccessToken_IsActive(t *testing.T) {
	now := time.Now()
	pastTime := now.Add(-24 * time.Hour)
	futureTime := now.Add(24 * time.Hour)

	tests := []struct {
		name      string
		status    string
		expiresAt *time.Time
		want      bool
	}{
		{
			name:      "活跃且未过期",
			status:    "active",
			expiresAt: &futureTime,
			want:      true,
		},
		{
			name:      "活跃且永久有效",
			status:    "active",
			expiresAt: nil,
			want:      true,
		},
		{
			name:      "活跃但已过期",
			status:    "active",
			expiresAt: &pastTime,
			want:      false,
		},
		{
			name:      "已禁用",
			status:    "disabled",
			expiresAt: &futureTime,
			want:      false,
		},
		{
			name:      "已过期状态",
			status:    "expired",
			expiresAt: nil,
			want:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pat := newTestPAT(tt.status, tt.expiresAt)
			got := pat.IsActive()
			assert.Equal(t, tt.want, got, "PersonalAccessToken.IsActive()")
		})
	}
}

func TestPersonalAccessToken_ToListItem(t *testing.T) {
	now := time.Now()
	expiresAt := now.Add(30 * 24 * time.Hour)
	lastUsedAt := now.Add(-1 * time.Hour)

	pat := &PersonalAccessToken{
		ID:          42,
		Name:        "CI Token",
		TokenPrefix: "ghp_1234",
		Scopes:      StringList{string(ScopeSelf), string(ScopeSys)},
		ExpiresAt:   &expiresAt,
		LastUsedAt:  &lastUsedAt,
		Status:      "active",
		CreatedAt:   now,
	}

	item := pat.ToListItem()

	assert.Equal(t, pat.ID, item.ID, "ID 应该相等")
	assert.Equal(t, pat.Name, item.Name, "Name 应该相等")
	assert.Equal(t, pat.TokenPrefix, item.TokenPrefix, "TokenPrefix 应该相等")
	assert.Equal(t, []string(pat.Scopes), item.Scopes, "Scopes 应该相等")
	assert.Equal(t, pat.ExpiresAt, item.ExpiresAt, "ExpiresAt 应该相等")
	assert.Equal(t, pat.LastUsedAt, item.LastUsedAt, "LastUsedAt 应该相等")
	assert.Equal(t, pat.Status, item.Status, "Status 应该相等")
	assert.Equal(t, pat.CreatedAt, item.CreatedAt, "CreatedAt 应该相等")
}

func TestPersonalAccessToken_ToListItem_NilFields(t *testing.T) {
	pat := &PersonalAccessToken{
		ID:          1,
		Name:        "Simple Token",
		TokenPrefix: "ghp_xxxx",
		Scopes:      nil,
		ExpiresAt:   nil,
		LastUsedAt:  nil,
		Status:      "active",
	}

	item := pat.ToListItem()

	assert.Nil(t, item.ExpiresAt, "ExpiresAt 应该为 nil")
	assert.Nil(t, item.LastUsedAt, "LastUsedAt 应该为 nil")
	assert.Nil(t, item.Scopes, "Scopes 应该为 nil")
}

func TestPersonalAccessToken_EdgeCases(t *testing.T) {
	t.Run("边界时间 - 刚刚过期", func(t *testing.T) {
		// 1 纳秒前过期
		justExpired := time.Now().Add(-1 * time.Nanosecond)
		pat := newTestPAT("active", &justExpired)

		assert.True(t, pat.IsExpired(), "刚刚过期的 Token 应该返回 true")
		assert.False(t, pat.IsActive(), "刚刚过期的 Token 不应该是活跃的")
	})

	t.Run("零值时间", func(t *testing.T) {
		zeroTime := time.Time{}
		pat := newTestPAT("active", &zeroTime)

		// 零值时间是过去的，应该过期
		assert.True(t, pat.IsExpired(), "零值时间应该被视为过期")
	})
}

func TestPersonalAccessToken_IsIPAllowed(t *testing.T) {
	tests := []struct {
		name        string
		ipWhitelist StringList
		ip          string
		want        bool
	}{
		{
			name:        "空白名单 - 允许所有",
			ipWhitelist: nil,
			ip:          "192.168.1.1",
			want:        true,
		},
		{
			name:        "空数组白名单 - 允许所有",
			ipWhitelist: StringList{},
			ip:          "192.168.1.1",
			want:        true,
		},
		{
			name:        "IP 在白名单中",
			ipWhitelist: StringList{"192.168.1.1", "10.0.0.1"},
			ip:          "192.168.1.1",
			want:        true,
		},
		{
			name:        "IP 不在白名单中",
			ipWhitelist: StringList{"192.168.1.1", "10.0.0.1"},
			ip:          "172.16.0.1",
			want:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pat := &PersonalAccessToken{IPWhitelist: tt.ipWhitelist}
			assert.Equal(t, tt.want, pat.IsIPAllowed(tt.ip))
		})
	}
}

func TestPersonalAccessToken_HasScope(t *testing.T) {
	pat := &PersonalAccessToken{
		Scopes: StringList{string(ScopeFull), string(ScopeSelf)},
	}

	tests := []struct {
		name  string
		scope string
		want  bool
	}{
		{"存在的 Scope - full", string(ScopeFull), true},
		{"存在的 Scope - self", string(ScopeSelf), true},
		{"不存在的 Scope - sys", string(ScopeSys), false},
		{"空 Scope", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, pat.HasScope(tt.scope))
		})
	}
}

func TestPersonalAccessToken_HasFullScope(t *testing.T) {
	tests := []struct {
		name   string
		scopes StringList
		want   bool
	}{
		{"有 full scope", StringList{string(ScopeFull)}, true},
		{"多个 scope 包含 full", StringList{string(ScopeSelf), string(ScopeFull)}, true},
		{"没有 full scope", StringList{string(ScopeSelf), string(ScopeSys)}, false},
		{"空 scopes", nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pat := &PersonalAccessToken{Scopes: tt.scopes}
			assert.Equal(t, tt.want, pat.HasFullScope())
		})
	}
}

func TestPersonalAccessToken_DisableEnable(t *testing.T) {
	t.Run("禁用 Token", func(t *testing.T) {
		pat := newTestPAT("active", nil)
		assert.False(t, pat.IsDisabled())

		pat.Disable()
		assert.True(t, pat.IsDisabled())
		assert.Equal(t, StatusDisabled, pat.Status)
		assert.False(t, pat.IsActive())
	})

	t.Run("启用 Token", func(t *testing.T) {
		pat := newTestPAT("disabled", nil)
		assert.True(t, pat.IsDisabled())

		pat.Enable()
		assert.False(t, pat.IsDisabled())
		assert.Equal(t, StatusActive, pat.Status)
		assert.True(t, pat.IsActive())
	})

	t.Run("标记过期", func(t *testing.T) {
		pat := newTestPAT("active", nil)
		pat.MarkExpired()
		assert.Equal(t, StatusExpired, pat.Status)
		assert.False(t, pat.IsActive())
	})
}

func TestPersonalAccessToken_CanBeUsed(t *testing.T) {
	futureTime := time.Now().Add(24 * time.Hour)

	tests := []struct {
		name        string
		status      string
		expiresAt   *time.Time
		ipWhitelist StringList
		ip          string
		want        bool
	}{
		{
			name:        "活跃且 IP 允许",
			status:      "active",
			expiresAt:   &futureTime,
			ipWhitelist: nil,
			ip:          "192.168.1.1",
			want:        true,
		},
		{
			name:        "活跃但 IP 不允许",
			status:      "active",
			expiresAt:   &futureTime,
			ipWhitelist: StringList{"10.0.0.1"},
			ip:          "192.168.1.1",
			want:        false,
		},
		{
			name:        "禁用状态",
			status:      "disabled",
			expiresAt:   &futureTime,
			ipWhitelist: nil,
			ip:          "192.168.1.1",
			want:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pat := &PersonalAccessToken{
				Status:      tt.status,
				ExpiresAt:   tt.expiresAt,
				IPWhitelist: tt.ipWhitelist,
			}
			assert.Equal(t, tt.want, pat.CanBeUsed(tt.ip))
		})
	}
}

// TestFilterByScopes_SuperAdminRestriction 验证超级权限在限制性 Scope 下被正确过滤。
//
// 安全要求：PAT 权限不得超过用户本身权限，且 Scope 限制必须生效。
func TestFilterByScopes_SuperAdminRestriction(t *testing.T) {
	superPerms := []role.Permission{
		{OperationPattern: "*:*:*", ResourcePattern: "*:*:*"},
	}

	t.Run("full scope 应保留超级权限", func(t *testing.T) {
		result := FilterByScopes([]string{"full"}, superPerms)
		assert.Len(t, result, 1, "full scope 应返回全部权限")
		assert.Equal(t, "*:*:*", result[0].OperationPattern)
	})

	t.Run("self scope 不应保留超级权限", func(t *testing.T) {
		result := FilterByScopes([]string{"self"}, superPerms)
		assert.Empty(t, result, "self scope 不应继承 *:*:* 超级权限")
	})

	t.Run("sys scope 不应保留超级权限", func(t *testing.T) {
		result := FilterByScopes([]string{"sys"}, superPerms)
		assert.Empty(t, result, "sys scope 不应继承 *:*:* 超级权限")
	})

	t.Run("self+sys scope 不应保留超级权限", func(t *testing.T) {
		result := FilterByScopes([]string{"self", "sys"}, superPerms)
		assert.Empty(t, result, "self+sys scope 不应继承 *:*:* 超级权限")
	})
}

// TestFilterByScopes_MixedPermissions 验证混合权限的过滤逻辑。
func TestFilterByScopes_MixedPermissions(t *testing.T) {
	mixedPerms := []role.Permission{
		{OperationPattern: "*:*:*", ResourcePattern: "*:*:*"},
		{OperationPattern: "self:profile:read", ResourcePattern: "self:user:@me"},
		{OperationPattern: "self:profile:update", ResourcePattern: "self:user:@me"},
		{OperationPattern: "sys:users:list", ResourcePattern: "*:*:*"},
		{OperationPattern: "sys:users:create", ResourcePattern: "*:*:*"},
	}

	t.Run("self scope 只保留 self 权限", func(t *testing.T) {
		result := FilterByScopes([]string{"self"}, mixedPerms)
		assert.Len(t, result, 2, "应只保留 2 个 self 权限")
		for _, p := range result {
			assert.Equal(t, "self:", p.OperationPattern[:5], "所有权限应以 self: 开头")
		}
	})

	t.Run("sys scope 只保留 sys 权限", func(t *testing.T) {
		result := FilterByScopes([]string{"sys"}, mixedPerms)
		assert.Len(t, result, 2, "应只保留 2 个 sys 权限")
		for _, p := range result {
			assert.Equal(t, "sys:", p.OperationPattern[:4], "所有权限应以 sys: 开头")
		}
	})

	t.Run("self+sys scope 取并集", func(t *testing.T) {
		result := FilterByScopes([]string{"self", "sys"}, mixedPerms)
		assert.Len(t, result, 4, "应保留 4 个权限（2 self + 2 sys）")
	})

	t.Run("full scope 返回全部", func(t *testing.T) {
		result := FilterByScopes([]string{"full"}, mixedPerms)
		assert.Len(t, result, 5, "full scope 应返回全部 5 个权限")
	})
}
