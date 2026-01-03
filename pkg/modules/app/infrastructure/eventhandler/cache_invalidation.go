package eventhandler

import (
	"context"
	"log/slog"

	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/iam/infrastructure/auth"
	"github.com/lwmacct/260101-go-pkg-ddd/pkg/shared/event"
)

// CacheInvalidationHandler 缓存失效处理器
// 处理角色权限变更和用户角色分配事件，自动失效相关缓存
type CacheInvalidationHandler struct {
	permissionCache *auth.PermissionCacheService
	logger          *slog.Logger
}

// NewCacheInvalidationHandler 创建缓存失效处理器
func NewCacheInvalidationHandler(
	permissionCache *auth.PermissionCacheService,
) *CacheInvalidationHandler {
	return &CacheInvalidationHandler{
		permissionCache: permissionCache,
		logger:          slog.Default(),
	}
}

// Handle 处理事件
func (h *CacheInvalidationHandler) Handle(ctx context.Context, e event.Event) error {
	switch evt := e.(type) {
	case *event.UserRoleAssignedEvent:
		return h.handleUserRoleAssigned(ctx, evt)
	case *event.RolePermissionsChangedEvent:
		return h.handleRolePermissionsChanged(ctx, evt)
	case *event.UserDeletedEvent:
		return h.handleUserDeleted(ctx, evt)
	default:
		// 忽略不处理的事件
		return nil
	}
}

// handleUserRoleAssigned 处理用户角色分配事件
// 失效单个用户的权限缓存
func (h *CacheInvalidationHandler) handleUserRoleAssigned(ctx context.Context, evt *event.UserRoleAssignedEvent) error {
	h.logger.Info("invalidating permission cache for user",
		"event", evt.EventName(),
		"user_id", evt.UserID,
	)

	if err := h.permissionCache.InvalidateUser(ctx, evt.UserID); err != nil {
		h.logger.Error("failed to invalidate user permission cache",
			"user_id", evt.UserID,
			"error", err,
		)
		// 缓存失效失败不应该阻塞业务流程，只记录错误
		return nil
	}

	return nil
}

// handleRolePermissionsChanged 处理角色权限变更事件
// 失效拥有该角色的所有用户的权限缓存（批量删除）
func (h *CacheInvalidationHandler) handleRolePermissionsChanged(ctx context.Context, evt *event.RolePermissionsChangedEvent) error {
	h.logger.Info("invalidating permission cache for role",
		"event", evt.EventName(),
		"role_id", evt.RoleID,
	)

	// 批量失效拥有该角色的所有用户缓存（内部使用 Pipeline）
	if err := h.permissionCache.InvalidateUsersWithRole(ctx, evt.RoleID); err != nil {
		h.logger.Error("failed to invalidate users permission cache",
			"role_id", evt.RoleID,
			"error", err,
		)
	}

	return nil
}

// handleUserDeleted 处理用户删除事件
// 清理用户相关缓存
func (h *CacheInvalidationHandler) handleUserDeleted(ctx context.Context, evt *event.UserDeletedEvent) error {
	h.logger.Info("cleaning up cache for deleted user",
		"event", evt.EventName(),
		"user_id", evt.UserID,
	)

	if err := h.permissionCache.InvalidateUser(ctx, evt.UserID); err != nil {
		h.logger.Error("failed to invalidate deleted user cache",
			"user_id", evt.UserID,
			"error", err,
		)
	}

	return nil
}

// Ensure interface is implemented
var _ event.EventHandler = (*CacheInvalidationHandler)(nil)
