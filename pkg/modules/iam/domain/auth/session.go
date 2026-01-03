package auth

import (
	"context"
	"time"
)

// SessionService 登录会话服务接口。
//
// 提供 2FA 流程中的临时会话管理：
//   - 用户完成账密验证后生成会话 Token
//   - 用户提交 2FA 验证码时校验会话
//   - 防止暴力破解攻击
//
// 实现位于 infrastructure/auth 包。
type SessionService interface {
	// GenerateSessionToken 生成会话 Token（账密验证通过后）。
	// Token 有效期通常为 5 分钟。
	GenerateSessionToken(ctx context.Context, userID uint, account string) (string, error)

	// VerifySessionToken 验证会话 Token（2FA 验证前）。
	// Token 验证成功后自动删除（一次性使用）。
	// Token 无效或过期返回 error。
	VerifySessionToken(ctx context.Context, token string) (*SessionData, error)
}

// SessionData 会话数据。
type SessionData struct {
	UserID    uint
	Account   string
	CreatedAt time.Time
	ExpireAt  time.Time
}
