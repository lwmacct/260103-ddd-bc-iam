package role

import (
	"context"

	"github.com/lwmacct/260103-ddd-iam-bc/pkg/modules/iam/domain/role"
)

// GetHandler 获取角色查询处理器
type GetHandler struct {
	roleQueryRepo role.QueryRepository
}

// NewGetHandler 创建获取角色查询处理器
func NewGetHandler(roleQueryRepo role.QueryRepository) *GetHandler {
	return &GetHandler{
		roleQueryRepo: roleQueryRepo,
	}
}

// Handle 处理获取角色查询
func (h *GetHandler) Handle(ctx context.Context, query GetQuery) (*RoleDTO, error) {
	// 查询角色（包含权限）
	roleEntity, err := h.roleQueryRepo.FindByIDWithPermissions(ctx, query.RoleID)
	if err != nil {
		return nil, err // 直接返回 repository 错误（保留 domain 错误类型）
	}
	if roleEntity == nil {
		return nil, role.ErrRoleNotFound // 返回 domain 错误
	}

	// 转换为 DTO
	return ToRoleDTO(roleEntity), nil
}
