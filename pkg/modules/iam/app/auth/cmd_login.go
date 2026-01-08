package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/lwmacct/260103-ddd-iam-bc/pkg/modules/iam/app/audit"
	"github.com/lwmacct/260103-ddd-iam-bc/pkg/modules/iam/domain/auth"
	"github.com/lwmacct/260103-ddd-iam-bc/pkg/modules/iam/domain/twofa"
	"github.com/lwmacct/260103-ddd-iam-bc/pkg/modules/iam/domain/user"
	"github.com/lwmacct/260103-ddd-shared/pkg/shared/captcha"
)

// LoginHandler 登录命令处理器
type LoginHandler struct {
	userQueryRepo      user.QueryRepository
	captchaCommandRepo captcha.CommandRepository
	twofaQueryRepo     twofa.QueryRepository
	authService        auth.Service
	loginSession       auth.SessionService
	auditLogHandler    *audit.CreateHandler
}

// NewLoginHandler 创建登录命令处理器
func NewLoginHandler(
	userQueryRepo user.QueryRepository,
	captchaCommandRepo captcha.CommandRepository,
	twofaQueryRepo twofa.QueryRepository,
	authService auth.Service,
	loginSession auth.SessionService,
	auditLogHandler *audit.CreateHandler,
) *LoginHandler {
	return &LoginHandler{
		userQueryRepo:      userQueryRepo,
		captchaCommandRepo: captchaCommandRepo,
		twofaQueryRepo:     twofaQueryRepo,
		authService:        authService,
		loginSession:       loginSession,
		auditLogHandler:    auditLogHandler,
	}
}

// Handle 处理登录命令
func (h *LoginHandler) Handle(ctx context.Context, cmd LoginCommand) (*LoginResultDTO, error) {
	// 1. 验证图形验证码
	valid, err := h.captchaCommandRepo.Verify(ctx, cmd.CaptchaID, cmd.Captcha)
	if err != nil {
		return nil, fmt.Errorf("failed to verify captcha: %w", err)
	}
	if !valid {
		h.logLoginEvent(ctx, 0, cmd.Account, cmd.ClientIP, cmd.UserAgent, "invalid_captcha", "failure")
		return nil, auth.ErrInvalidCaptcha
	}

	// 2. 查找用户（支持用户名或邮箱登录）
	var u *user.User
	if u, err = h.userQueryRepo.GetByUsernameWithRoles(ctx, cmd.Account); err != nil {
		// 尝试通过邮箱查找
		if u, err = h.userQueryRepo.GetByEmailWithRoles(ctx, cmd.Account); err != nil {
			h.logLoginEvent(ctx, 0, cmd.Account, cmd.ClientIP, cmd.UserAgent, "user_not_found", "failure")
			return nil, auth.ErrInvalidCredentials
		}
	}

	// 3. 检查用户状态
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

	// 4. 服务账户不能使用密码登录（必须使用 PAT）
	if u.IsServiceAccount() {
		h.logLoginEvent(ctx, u.ID, u.Username, cmd.ClientIP, cmd.UserAgent, "service_account_password_login", "failure")
		return nil, user.ErrServiceAccountPasswordLogin
	}

	// 5. 验证密码
	if err = h.authService.VerifyPassword(ctx, u.Password, cmd.Password); err != nil {
		h.logLoginEvent(ctx, u.ID, u.Username, cmd.ClientIP, cmd.UserAgent, "invalid_password", "failure")
		return nil, auth.ErrInvalidCredentials
	}

	// 6. 检查是否启用 2FA
	tfa, err := h.twofaQueryRepo.FindByUserID(ctx, u.ID)
	if err == nil && tfa != nil && tfa.Enabled {
		// 需要 2FA 验证，生成临时 session token
		sessionToken, sessionErr := h.loginSession.GenerateSessionToken(ctx, u.ID, cmd.Account)
		if sessionErr != nil {
			return nil, fmt.Errorf("failed to generate session token: %w", sessionErr)
		}

		// 2FA 挑战不记录为成功登录，等待 2FA 验证完成后记录
		return &LoginResultDTO{
			Requires2FA:  true,
			SessionToken: sessionToken,
			UserID:       u.ID,
			Username:     u.Username,
		}, nil
	}

	// 7. 生成访问令牌（新架构：不传递 roles，权限从缓存查询）
	accessToken, expiresAt, err := h.authService.GenerateAccessToken(ctx, u.ID, u.Username)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, _, err := h.authService.GenerateRefreshToken(ctx, u.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	expiresIn := int(time.Until(expiresAt).Seconds())

	// 记录登录成功
	h.logLoginEvent(ctx, u.ID, u.Username, cmd.ClientIP, cmd.UserAgent, "login_success", "success")

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
func (h *LoginHandler) logLoginEvent(ctx context.Context, userID uint, username, clientIP, userAgent, event, status string) {
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
