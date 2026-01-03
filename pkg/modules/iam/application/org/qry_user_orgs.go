package org

import (
	"context"
	"fmt"

	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/iam/domain/org"
)

// UserOrgsHandler 查询用户加入的组织处理器
type UserOrgsHandler struct {
	memberQuery org.MemberQueryRepository
	orgQuery    org.QueryRepository
}

// NewUserOrgsHandler 创建用户组织查询处理器
func NewUserOrgsHandler(
	memberQuery org.MemberQueryRepository,
	orgQuery org.QueryRepository,
) *UserOrgsHandler {
	return &UserOrgsHandler{
		memberQuery: memberQuery,
		orgQuery:    orgQuery,
	}
}

// Handle 处理查询用户加入的组织
func (h *UserOrgsHandler) Handle(ctx context.Context, query ListUserOrgsQuery) ([]*UserOrgDTO, error) {
	// 1. 获取用户的所有组织成员记录
	memberships, err := h.memberQuery.ListByUser(ctx, query.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to list user memberships: %w", err)
	}

	if len(memberships) == 0 {
		return []*UserOrgDTO{}, nil
	}

	// 2. 获取组织详情并组装结果
	result := make([]*UserOrgDTO, 0, len(memberships))
	for _, m := range memberships {
		org, err := h.orgQuery.GetByID(ctx, m.OrgID)
		if err != nil {
			// 跳过不存在的组织（可能已被删除）
			continue
		}

		result = append(result, &UserOrgDTO{
			OrgDTO:   *ToOrgDTO(org),
			Role:     string(m.Role),
			JoinedAt: m.JoinedAt,
		})
	}

	return result, nil
}
