package org

import (
	"context"
	"fmt"

	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/iam/domain/org"
)

// MemberListHandler 成员列表查询处理器
type MemberListHandler struct {
	memberQuery org.MemberQueryRepository
}

// NewMemberListHandler 创建成员列表查询处理器
func NewMemberListHandler(memberQuery org.MemberQueryRepository) *MemberListHandler {
	return &MemberListHandler{memberQuery: memberQuery}
}

// MemberListResult 成员列表查询结果
type MemberListResult struct {
	Items []*MemberDTO
	Total int64
}

// Handle 处理成员列表查询
func (h *MemberListHandler) Handle(ctx context.Context, query ListMembersQuery) (*MemberListResult, error) {
	members, err := h.memberQuery.ListByOrgWithUsers(ctx, query.OrgID, query.Offset, query.Limit)
	if err != nil {
		return nil, fmt.Errorf("failed to list members: %w", err)
	}

	total, err := h.memberQuery.CountByOrg(ctx, query.OrgID)
	if err != nil {
		return nil, fmt.Errorf("failed to count members: %w", err)
	}

	return &MemberListResult{
		Items: ToMemberWithUserDTOs(members),
		Total: total,
	}, nil
}
