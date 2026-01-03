package role

import (
	"context"
	"fmt"

	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/iam/domain/role"
)

// CreateHandler 创建角色命令处理器
type CreateHandler struct {
	roleCommandRepo role.CommandRepository
	roleQueryRepo   role.QueryRepository
}

// NewCreateHandler 创建角色命令处理器
func NewCreateHandler(
	roleCommandRepo role.CommandRepository,
	roleQueryRepo role.QueryRepository,
) *CreateHandler {
	return &CreateHandler{
		roleCommandRepo: roleCommandRepo,
		roleQueryRepo:   roleQueryRepo,
	}
}

// Handle 处理创建角色命令
func (h *CreateHandler) Handle(ctx context.Context, cmd CreateCommand) (*CreateResultDTO, error) {
	// 1. 验证角色名是否已存在
	exists, err := h.roleQueryRepo.ExistsByName(ctx, cmd.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to check role name existence: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("role name already exists: %s", cmd.Name)
	}

	// 2. 创建角色实体
	newRole := &role.Role{
		Name:        cmd.Name,
		DisplayName: cmd.DisplayName,
		Description: cmd.Description,
		IsSystem:    false, // 用户创建的角色不是系统角色
	}

	// 3. 保存角色
	if err := h.roleCommandRepo.Create(ctx, newRole); err != nil {
		return nil, fmt.Errorf("failed to create role: %w", err)
	}

	return &CreateResultDTO{
		RoleID:      newRole.ID,
		Name:        newRole.Name,
		DisplayName: newRole.DisplayName,
	}, nil
}
