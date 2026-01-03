package org

import (
	"context"

	"github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/domain/org"
)

// TeamGetHandler 获取团队查询处理器
type TeamGetHandler struct {
	teamQuery org.TeamQueryRepository
}

// NewTeamGetHandler 创建获取团队查询处理器
func NewTeamGetHandler(teamQuery org.TeamQueryRepository) *TeamGetHandler {
	return &TeamGetHandler{teamQuery: teamQuery}
}

// Handle 处理获取团队查询
func (h *TeamGetHandler) Handle(ctx context.Context, query GetTeamQuery) (*TeamDTO, error) {
	team, err := h.teamQuery.GetByID(ctx, query.TeamID)
	if err != nil {
		return nil, err
	}

	// 验证团队属于指定组织
	if !team.BelongsTo(query.OrgID) {
		return nil, org.ErrTeamNotInOrg
	}

	return ToTeamDTO(team), nil
}
