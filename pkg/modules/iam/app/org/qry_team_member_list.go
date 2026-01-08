package org

import (
	"context"
	"fmt"

	"github.com/lwmacct/260103-ddd-iam-bc/pkg/modules/iam/domain/org"
)

// TeamMemberListHandler 团队成员列表查询处理器
type TeamMemberListHandler struct {
	teamMemberQuery org.TeamMemberQueryRepository
	teamQuery       org.TeamQueryRepository
}

// NewTeamMemberListHandler 创建团队成员列表查询处理器
func NewTeamMemberListHandler(
	teamMemberQuery org.TeamMemberQueryRepository,
	teamQuery org.TeamQueryRepository,
) *TeamMemberListHandler {
	return &TeamMemberListHandler{
		teamMemberQuery: teamMemberQuery,
		teamQuery:       teamQuery,
	}
}

// TeamMemberListResult 团队成员列表查询结果
type TeamMemberListResult struct {
	Items []*TeamMemberDTO
	Total int64
}

// Handle 处理团队成员列表查询
func (h *TeamMemberListHandler) Handle(ctx context.Context, query ListTeamMembersQuery) (*TeamMemberListResult, error) {
	// 1. 验证团队属于组织
	belongsTo, err := h.teamQuery.BelongsToOrg(ctx, query.TeamID, query.OrgID)
	if err != nil {
		return nil, fmt.Errorf("failed to check team org: %w", err)
	}
	if !belongsTo {
		return nil, org.ErrTeamNotInOrg
	}

	// 2. 查询成员列表
	members, err := h.teamMemberQuery.ListByTeam(ctx, query.TeamID, query.Offset, query.Limit)
	if err != nil {
		return nil, fmt.Errorf("failed to list team members: %w", err)
	}

	total, err := h.teamMemberQuery.CountByTeam(ctx, query.TeamID)
	if err != nil {
		return nil, fmt.Errorf("failed to count team members: %w", err)
	}

	return &TeamMemberListResult{
		Items: ToTeamMemberDTOs(members),
		Total: total,
	}, nil
}
