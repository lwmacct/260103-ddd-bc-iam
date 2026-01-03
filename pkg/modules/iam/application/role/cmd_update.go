package role

import (
	"context"
	"errors"
	"fmt"

	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/iam/domain/role"
)

// UpdateHandler 更新角色命令处理器
type UpdateHandler struct {
	roleCommandRepo role.CommandRepository
	roleQueryRepo   role.QueryRepository
}

// NewUpdateHandler 创建更新角色命令处理器
func NewUpdateHandler(
	roleCommandRepo role.CommandRepository,
	roleQueryRepo role.QueryRepository,
) *UpdateHandler {
	return &UpdateHandler{
		roleCommandRepo: roleCommandRepo,
		roleQueryRepo:   roleQueryRepo,
	}
}

// Handle 处理更新角色命令
func (h *UpdateHandler) Handle(ctx context.Context, cmd UpdateCommand) (*RoleDTO, error) {
	// 1. 查找角色
	existingRole, err := h.roleQueryRepo.FindByID(ctx, cmd.RoleID)
	if err != nil {
		return nil, fmt.Errorf("failed to find role: %w", err)
	}
	if existingRole == nil {
		return nil, fmt.Errorf("role not found with id: %d", cmd.RoleID)
	}

	// 2. 检查是否为系统角色（系统角色不可修改）
	if existingRole.IsSystem {
		return nil, errors.New("cannot modify system role")
	}

	// 3. 更新字段
	if cmd.DisplayName != nil {
		existingRole.DisplayName = *cmd.DisplayName
	}
	if cmd.Description != nil {
		existingRole.Description = *cmd.Description
	}

	// 4. 保存更新
	if err := h.roleCommandRepo.Update(ctx, existingRole); err != nil {
		return nil, fmt.Errorf("failed to update role: %w", err)
	}

	return ToRoleDTO(existingRole), nil
}
