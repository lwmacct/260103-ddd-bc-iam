package role

import (
	"context"
	"fmt"

	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/iam/domain/role"
)

// ListHandler 列出角色查询处理器
type ListHandler struct {
	roleQueryRepo role.QueryRepository
}

// NewListHandler 创建列出角色查询处理器
func NewListHandler(roleQueryRepo role.QueryRepository) *ListHandler {
	return &ListHandler{
		roleQueryRepo: roleQueryRepo,
	}
}

// Handle 处理列出角色查询
func (h *ListHandler) Handle(ctx context.Context, query ListQuery) (*ListRolesDTO, error) {
	// 查询角色列表
	roles, total, err := h.roleQueryRepo.List(ctx, query.Page, query.Limit)
	if err != nil {
		return nil, fmt.Errorf("failed to list roles: %w", err)
	}

	// 转换为 DTO
	roleResponses := make([]*RoleDTO, 0, len(roles))
	for i := range roles {
		roleResponses = append(roleResponses, ToRoleDTO(&roles[i]))
	}

	return &ListRolesDTO{
		Roles: roleResponses,
		Total: total,
		Page:  query.Page,
		Limit: query.Limit,
	}, nil
}
