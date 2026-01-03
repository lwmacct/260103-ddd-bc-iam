package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/app/audit"
	"github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/domain/auth"
	"github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/domain/twofa"
	"github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/domain/user"
)

// Login2FAHandler 二次认证登录命令处理器
type Login2FAHandler struct {
	userQueryRepo   user.QueryRepository
	authService     auth.Service
	loginSession    auth.SessionService
	twofaService    twofa.Service
	auditLogHandler *audit.CreateHandler
}

// NewLogin2FAHandler 创建二次认证登录命令处理器
func NewLogin2FAHandler(
	userQueryRepo user.QueryRepository,
	authService auth.Service,
	loginSession auth.SessionService,
	twofaService twofa.Service,
	auditLogHandler *audit.CreateHandler,
) *Login2FAHandler {
	return &Login2FAHandler{
		userQueryRepo:   userQueryRepo,
		authService:     authService,
		loginSession:    loginSession,
		twofaService:    twofaService,
		auditLogHandler: auditLogHandler,
	}
}

// Handle 处理二次认证登录命令
func (h *Login2FAHandler) Handle(ctx context.Context, cmd Login2FACommand) (*LoginResultDTO, error) {
	// 1. 验证 session token（防止 2FA 暴力破解）
	sessionData, err := h.loginSession.VerifySessionToken(ctx, cmd.SessionToken)
	if err != nil {
		h.logLoginEvent(ctx, 0, "", cmd.ClientIP, cmd.UserAgent, "session_expired", "failure")
		return nil, errors.New("session expired or invalid, please login again")
	}

	// 2. 验证 2FA 验证码
	valid, err := h.twofaService.Verify(ctx, sessionData.UserID, cmd.TwoFactorCode)
	if err != nil {
		h.logLoginEvent(ctx, sessionData.UserID, sessionData.Account, cmd.ClientIP, cmd.UserAgent, "2fa_verify_error", "failure")
		return nil, fmt.Errorf("2FA verification failed: %w", err)
	}
	if !valid {
		h.logLoginEvent(ctx, sessionData.UserID, sessionData.Account, cmd.ClientIP, cmd.UserAgent, "2fa_invalid_code", "failure")
		return nil, errors.New("invalid two factor code")
	}

	// 3. 获取用户信息
	u, err := h.userQueryRepo.GetByIDWithRoles(ctx, sessionData.UserID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// 4. 检查用户状态
	if !u.CanLogin() {
		if u.IsBanned() {
			h.logLoginEvent(ctx, u.ID, u.Username, cmd.ClientIP, cmd.UserAgent, "user_banned", "failure")
			return nil, auth.ErrUserBanned
		}
		if u.IsInactive() {
			h.logLoginEvent(ctx, u.ID, u.Username, cmd.ClientIP, cmd.UserAgent, "user_inactive", "failure")
			return nil, auth.ErrUserInactive
		}
	}

	// 5. 生成访问令牌
	accessToken, expiresAt, err := h.authService.GenerateAccessToken(ctx, u.ID, u.Username)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, _, err := h.authService.GenerateRefreshToken(ctx, u.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	expiresIn := int(time.Until(expiresAt).Seconds())

	// 记录 2FA 登录成功
	h.logLoginEvent(ctx, u.ID, u.Username, cmd.ClientIP, cmd.UserAgent, "2fa_login_success", "success")

	// 构建角色列表
	roles := make([]LoginRoleDTO, 0, len(u.Roles))
	for _, r := range u.Roles {
		roles = append(roles, LoginRoleDTO{
			ID:   r.ID,
			Name: r.Name,
		})
	}

	return &LoginResultDTO{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    expiresIn,
		UserID:       u.ID,
		Username:     u.Username,
		Requires2FA:  false,
		Roles:        roles,
	}, nil
}

// logLoginEvent 异步记录登录事件到审计日志
func (h *Login2FAHandler) logLoginEvent(ctx context.Context, userID uint, username, clientIP, userAgent, event, status string) {
	if h.auditLogHandler == nil {
		return
	}
	go func() {
		_ = h.auditLogHandler.Handle(context.WithoutCancel(ctx), audit.CreateCommand{
			UserID:     userID,
			Username:   username,
			Action:     "login",
			Resource:   "auth",
			ResourceID: "",
			IPAddress:  clientIP,
			UserAgent:  userAgent,
			Details:    fmt.Sprintf(`{"event":"%s"}`, event),
			Status:     status,
		})
	}()
}
