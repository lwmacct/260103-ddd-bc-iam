package role

import (
	"context"
	"errors"
	"fmt"

	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/iam/domain/role"
	"github.com/lwmacct/260101-go-pkg-ddd/pkg/shared/event"
)

// SetPermissionsHandler 设置权限命令处理器
// 新 RBAC 模型：接受 Permission 模式（OperationPattern + ResourcePattern）
type SetPermissionsHandler struct {
	roleCommandRepo role.CommandRepository
	roleQueryRepo   role.QueryRepository
	eventBus        event.EventBus
}

// NewSetPermissionsHandler 创建设置权限命令处理器
func NewSetPermissionsHandler(
	roleCommandRepo role.CommandRepository,
	roleQueryRepo role.QueryRepository,
	eventBus event.EventBus,
) *SetPermissionsHandler {
	return &SetPermissionsHandler{
		roleCommandRepo: roleCommandRepo,
		roleQueryRepo:   roleQueryRepo,
		eventBus:        eventBus,
	}
}

// Handle 处理设置权限命令
// 新 RBAC 模型：直接设置 Permission 模式，无需验证 PermissionID
func (h *SetPermissionsHandler) Handle(ctx context.Context, cmd SetPermissionsCommand) error {
	// 1. 验证角色是否存在
	exists, err := h.roleQueryRepo.Exists(ctx, cmd.RoleID)
	if err != nil {
		return fmt.Errorf("failed to check role existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("role not found with id: %d", cmd.RoleID)
	}

	// 2. 验证权限模式有效性（基本格式验证）
	for _, perm := range cmd.Permissions {
		if perm.OperationPattern == "" {
			return errors.New("operation_pattern cannot be empty")
		}
		// ResourcePattern 允许为空，默认为 "*"
	}

	// 3. 设置权限
	if err := h.roleCommandRepo.SetPermissions(ctx, cmd.RoleID, cmd.Permissions); err != nil {
		return fmt.Errorf("failed to set permissions: %w", err)
	}

	// 4. 发布角色权限变更事件，触发缓存失效
	// 使用空 PermissionIDs 因为新模型不再使用 ID
	evt := event.NewRolePermissionsChangedEvent(cmd.RoleID, nil)
	if h.eventBus != nil {
		_ = h.eventBus.Publish(ctx, evt) // 缓存失效失败不阻塞业务
	}

	return nil
}
