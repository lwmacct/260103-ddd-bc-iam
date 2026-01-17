package user

import (
	"context"
	"fmt"

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

// Handle 处理分配角色命令，返回完整用户信息（包含角色）
func (h *AssignRolesHandler) Handle(ctx context.Context, cmd AssignRolesCommand) (*UserWithRolesDTO, error) {
	// 1. 获取用户
	u, err := h.userQueryRepo.GetByID(ctx, cmd.UserID)
	if err != nil {
		return nil, err
	}

	// 2. root 用户角色保护：root 用户角色不可修改
	if !u.CanModifyRoles() {
		return nil, user.ErrCannotModifyRootRoles
	}

	// 3. 分配角色
	if err := h.userCommandRepo.AssignRoles(ctx, cmd.UserID, cmd.RoleIDs); err != nil {
		return nil, err
	}

	// 4. 发布用户角色分配事件，触发缓存失效
	evt := event.NewUserRoleAssignedEvent(cmd.UserID, cmd.RoleIDs)
	if h.eventBus != nil {
		_ = h.eventBus.Publish(ctx, evt) // 缓存失效失败不阻塞业务
	}

	// 5. 获取更新后的完整用户信息（包含角色）
	updatedUser, err := h.userQueryRepo.GetByIDWithRoles(ctx, cmd.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get updated user: %w", err)
	}

	// 6. 转换为 DTO
	return ToUserWithRolesDTO(updatedUser), nil
}
