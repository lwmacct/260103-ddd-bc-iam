package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/lwmacct/260103-ddd-iam-bc/pkg/modules/iam/app/audit"
	"github.com/lwmacct/260103-ddd-iam-bc/pkg/modules/iam/domain/auth"
	"github.com/lwmacct/260103-ddd-iam-bc/pkg/modules/iam/domain/user"
)

// RefreshTokenHandler 刷新令牌命令处理器
type RefreshTokenHandler struct {
	userQueryRepo   user.QueryRepository
	authService     auth.Service
	auditLogHandler *audit.CreateHandler
}

// NewRefreshTokenHandler 创建刷新令牌命令处理器
func NewRefreshTokenHandler(
	userQueryRepo user.QueryRepository,
	authService auth.Service,
	auditLogHandler *audit.CreateHandler,
) *RefreshTokenHandler {
	return &RefreshTokenHandler{
		userQueryRepo:   userQueryRepo,
		authService:     authService,
		auditLogHandler: auditLogHandler,
	}
}

// Handle 处理刷新令牌命令
func (h *RefreshTokenHandler) Handle(ctx context.Context, cmd RefreshTokenCommand) (*RefreshTokenResultDTO, error) {
	// 1. 验证 refresh token
	userID, err := h.authService.ValidateRefreshToken(ctx, cmd.RefreshToken)
	if err != nil {
		h.logRefreshEvent(ctx, 0, "", cmd.ClientIP, cmd.UserAgent, "invalid_token", "failure")
		return nil, err
	}

	// 2. 获取用户信息
	u, err := h.userQueryRepo.GetByIDWithRoles(ctx, userID)
	if err != nil {
		h.logRefreshEvent(ctx, userID, "", cmd.ClientIP, cmd.UserAgent, "user_not_found", "failure")
		return nil, auth.ErrUserNotFound
	}

	// 3. 检查用户状态
	if !u.CanLogin() {
		if u.IsBanned() {
			h.logRefreshEvent(ctx, u.ID, u.Username, cmd.ClientIP, cmd.UserAgent, "user_banned", "failure")
			return nil, auth.ErrUserBanned
		}
		h.logRefreshEvent(ctx, u.ID, u.Username, cmd.ClientIP, cmd.UserAgent, "user_inactive", "failure")
		return nil, auth.ErrUserInactive
	}

	// 4. 生成新的访问令牌（新架构：不传递 roles，权限从缓存查询）
	accessToken, expiresAt, err := h.authService.GenerateAccessToken(ctx, u.ID, u.Username)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	// 5. 生成新的刷新令牌
	newRefreshToken, _, err := h.authService.GenerateRefreshToken(ctx, u.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// 6. 记录审计日志
	h.logRefreshEvent(ctx, u.ID, u.Username, cmd.ClientIP, cmd.UserAgent, "token_refreshed", "success")

	// 构建角色列表
	roles := make([]LoginRoleDTO, 0, len(u.Roles))
	for _, r := range u.Roles {
		roles = append(roles, LoginRoleDTO{
			ID:   r.ID,
			Name: r.Name,
		})
	}

	return &RefreshTokenResultDTO{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int(time.Until(expiresAt).Seconds()),
		UserID:       u.ID,
		Username:     u.Username,
		Roles:        roles,
	}, nil
}

// logRefreshEvent 异步记录刷新事件到审计日志
func (h *RefreshTokenHandler) logRefreshEvent(ctx context.Context, userID uint, username, clientIP, userAgent, event, status string) {
	if h.auditLogHandler == nil {
		return
	}
	go func() {
		_ = h.auditLogHandler.Handle(context.WithoutCancel(ctx), audit.CreateCommand{
			UserID:     userID,
			Username:   username,
			Action:     "refresh_token",
			Resource:   "auth",
			ResourceID: "",
			IPAddress:  clientIP,
			UserAgent:  userAgent,
			Details:    fmt.Sprintf(`{"event":"%s"}`, event),
			Status:     status,
		})
	}()
}
