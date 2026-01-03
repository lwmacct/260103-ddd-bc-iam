package pat

import (
	"slices"
	"time"
)

// PersonalAccessToken 个人访问令牌实体，用于 API 认证
type PersonalAccessToken struct {
	ID        uint
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time

	UserID      uint   // 所属用户
	Name        string // Token 名称
	Token       string // Token 哈希值（不返回）
	TokenPrefix string // Token 前缀（明文）

	Scopes StringList // 权限范围（full, self, sys）

	ExpiresAt  *time.Time // 过期时间（nil=永久）
	LastUsedAt *time.Time // 最后使用时间
	Status     string     // active, disabled, expired

	IPWhitelist StringList // IP 白名单（可选）
	Description string     // 描述
}

// IsExpired 检查 Token 是否已过期
func (p *PersonalAccessToken) IsExpired() bool {
	if p.ExpiresAt == nil {
		return false
	}
	return p.ExpiresAt.Before(time.Now())
}

// IsActive 检查 Token 是否处于活跃状态（未过期且未禁用）
func (p *PersonalAccessToken) IsActive() bool {
	return p.Status == "active" && !p.IsExpired()
}

// Token 状态常量
const (
	StatusActive   = "active"
	StatusDisabled = "disabled"
	StatusExpired  = "expired"
)

// IsIPAllowed 检查给定 IP 是否允许使用此 Token。
// 如果 IP 白名单为空，则允许所有 IP。
func (p *PersonalAccessToken) IsIPAllowed(ip string) bool {
	if len(p.IPWhitelist) == 0 {
		return true
	}
	return slices.Contains(p.IPWhitelist, ip)
}

// HasScope 检查 Token 是否包含指定 Scope。
func (p *PersonalAccessToken) HasScope(scope string) bool {
	return slices.Contains(p.Scopes, scope)
}

// HasFullScope 检查 Token 是否具有完整权限。
func (p *PersonalAccessToken) HasFullScope() bool {
	return p.HasScope(string(ScopeFull))
}

// Disable 禁用 Token
func (p *PersonalAccessToken) Disable() {
	p.Status = StatusDisabled
}

// Enable 启用 Token
func (p *PersonalAccessToken) Enable() {
	p.Status = StatusActive
}

// MarkExpired 标记 Token 为已过期
func (p *PersonalAccessToken) MarkExpired() {
	p.Status = StatusExpired
}

// IsDisabled 检查 Token 是否被禁用
func (p *PersonalAccessToken) IsDisabled() bool {
	return p.Status == StatusDisabled
}

// CanBeUsed 检查 Token 是否可以被使用（活跃且 IP 允许）
func (p *PersonalAccessToken) CanBeUsed(ip string) bool {
	return p.IsActive() && p.IsIPAllowed(ip)
}

// ToListItem 将 PAT 转换为列表项（不含完整 Token）
func (p *PersonalAccessToken) ToListItem() *TokenListItem {
	return &TokenListItem{
		ID:          p.ID,
		Name:        p.Name,
		TokenPrefix: p.TokenPrefix,
		Scopes:      p.Scopes,
		ExpiresAt:   p.ExpiresAt,
		LastUsedAt:  p.LastUsedAt,
		Status:      p.Status,
		CreatedAt:   p.CreatedAt,
	}
}

// TokenListItem PAT 列表项（不含完整 Token）
type TokenListItem struct {
	ID          uint
	Name        string
	TokenPrefix string // 用于识别
	Scopes      []string
	ExpiresAt   *time.Time
	LastUsedAt  *time.Time
	Status      string
	CreatedAt   time.Time
}
