package org

import (
	"context"
	"fmt"

	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/iam/domain/org"
)

// UserTeamsHandler 查询用户加入的团队处理器
type UserTeamsHandler struct {
	teamMemberQuery org.TeamMemberQueryRepository
	teamQuery       org.TeamQueryRepository
	orgQuery        org.QueryRepository
}

// NewUserTeamsHandler 创建用户团队查询处理器
func NewUserTeamsHandler(
	teamMemberQuery org.TeamMemberQueryRepository,
	teamQuery org.TeamQueryRepository,
	orgQuery org.QueryRepository,
) *UserTeamsHandler {
	return &UserTeamsHandler{
		teamMemberQuery: teamMemberQuery,
		teamQuery:       teamQuery,
		orgQuery:        orgQuery,
	}
}

// Handle 处理查询用户加入的团队
func (h *UserTeamsHandler) Handle(ctx context.Context, query ListUserTeamsQuery) ([]*UserTeamDTO, error) {
	var memberships []*org.TeamMember
	var err error

	// 1. 获取用户的团队成员记录
	if query.OrgID > 0 {
		// 限定在某个组织内
		memberships, err = h.teamMemberQuery.ListByUserInOrg(ctx, query.UserID, query.OrgID)
	} else {
		// 所有团队
		memberships, err = h.teamMemberQuery.ListByUser(ctx, query.UserID)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to list user team memberships: %w", err)
	}

	if len(memberships) == 0 {
		return []*UserTeamDTO{}, nil
	}

	// 2. 获取团队和组织详情并组装结果
	result := make([]*UserTeamDTO, 0, len(memberships))
	for _, m := range memberships {
		team, err := h.teamQuery.GetByID(ctx, m.TeamID)
		if err != nil {
			// 跳过不存在的团队
			continue
		}

		org, err := h.orgQuery.GetByID(ctx, team.OrgID)
		if err != nil {
			// 跳过不存在的组织
			continue
		}

		result = append(result, &UserTeamDTO{
			TeamDTO:  *ToTeamDTO(team),
			OrgName:  org.Name,
			Role:     string(m.Role),
			JoinedAt: m.JoinedAt,
		})
	}

	return result, nil
}
