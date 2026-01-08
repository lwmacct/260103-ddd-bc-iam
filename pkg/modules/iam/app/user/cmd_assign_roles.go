package user

import (
	"context"

	"github.com/lwmacct/260103-ddd-iam-bc/pkg/modules/iam/domain/user"
	"github.com/lwmacct/260103-ddd-shared/pkg/shared/event"
)

// AssignRolesHandler 负责分配用户角色
type AssignRolesHandler struct {
	userCommandRepo user.CommandRepository
	userQueryRepo   user.QueryRepository
	eventBus        event.EventBus
}

// NewAssignRolesHandler 创建新的分配角色处理器
func NewAssignRolesHandler(
	userCommandRepo user.CommandRepository,
	userQueryRepo user.QueryRepository,
	eventBus event.EventBus,
) *AssignRolesHandler {
	return &AssignRolesHandler{
		userCommandRepo: userCommandRepo,
		userQueryRepo:   userQueryRepo,
		eventBus:        eventBus,
	}
}

// Handle 处理分配角色命令
func (h *AssignRolesHandler) Handle(ctx context.Context, cmd AssignRolesCommand) error {
	// 1. 获取用户
	u, err := h.userQueryRepo.GetByID(ctx, cmd.UserID)
	if err != nil {
		return err
	}

	// 2. root 用户角色保护：root 用户角色不可修改
	if !u.CanModifyRoles() {
		return user.ErrCannotModifyRootRoles
	}

	// 3. 分配角色
	if err := h.userCommandRepo.AssignRoles(ctx, cmd.UserID, cmd.RoleIDs); err != nil {
		return err
	}

	// 4. 发布用户角色分配事件，触发缓存失效
	evt := event.NewUserRoleAssignedEvent(cmd.UserID, cmd.RoleIDs)
	if h.eventBus != nil {
		_ = h.eventBus.Publish(ctx, evt) // 缓存失效失败不阻塞业务
	}

	return nil
}
