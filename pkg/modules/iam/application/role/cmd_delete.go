package role

import (
	"context"
	"errors"
	"fmt"

	"github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/domain/role"
)

// DeleteHandler 删除角色命令处理器
type DeleteHandler struct {
	roleCommandRepo role.CommandRepository
	roleQueryRepo   role.QueryRepository
}

// NewDeleteHandler 创建删除角色命令处理器
func NewDeleteHandler(
	roleCommandRepo role.CommandRepository,
	roleQueryRepo role.QueryRepository,
) *DeleteHandler {
	return &DeleteHandler{
		roleCommandRepo: roleCommandRepo,
		roleQueryRepo:   roleQueryRepo,
	}
}

// Handle 处理删除角色命令
func (h *DeleteHandler) Handle(ctx context.Context, cmd DeleteCommand) error {
	// 1. 查找角色
	existingRole, err := h.roleQueryRepo.FindByID(ctx, cmd.RoleID)
	if err != nil {
		return fmt.Errorf("failed to find role: %w", err)
	}
	if existingRole == nil {
		return fmt.Errorf("role not found with id: %d", cmd.RoleID)
	}

	// 2. 检查是否为系统角色（系统角色不可删除）
	if existingRole.IsSystem {
		return errors.New("cannot delete system role")
	}

	// 3. 删除角色
	if err := h.roleCommandRepo.Delete(ctx, cmd.RoleID); err != nil {
		return fmt.Errorf("failed to delete role: %w", err)
	}

	return nil
}
