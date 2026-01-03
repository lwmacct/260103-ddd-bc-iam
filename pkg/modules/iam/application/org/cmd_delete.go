package org

import (
	"context"
	"fmt"

	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/iam/domain/org"
)

// DeleteHandler 删除组织命令处理器
type DeleteHandler struct {
	orgCommandRepo    org.CommandRepository
	orgQueryRepo      org.QueryRepository
	memberQuery       org.MemberQueryRepository
	memberCommandRepo org.MemberCommandRepository
	teamQuery         org.TeamQueryRepository
	teamCommandRepo   org.TeamCommandRepository
	teamMemberQuery   org.TeamMemberQueryRepository
	teamMemberCommand org.TeamMemberCommandRepository
}

// NewDeleteHandler 创建删除组织命令处理器
func NewDeleteHandler(
	orgCommandRepo org.CommandRepository,
	orgQueryRepo org.QueryRepository,
	memberQuery org.MemberQueryRepository,
	memberCommandRepo org.MemberCommandRepository,
	teamQuery org.TeamQueryRepository,
	teamCommandRepo org.TeamCommandRepository,
	teamMemberQuery org.TeamMemberQueryRepository,
	teamMemberCommand org.TeamMemberCommandRepository,
) *DeleteHandler {
	return &DeleteHandler{
		orgCommandRepo:    orgCommandRepo,
		orgQueryRepo:      orgQueryRepo,
		memberQuery:       memberQuery,
		memberCommandRepo: memberCommandRepo,
		teamQuery:         teamQuery,
		teamCommandRepo:   teamCommandRepo,
		teamMemberQuery:   teamMemberQuery,
		teamMemberCommand: teamMemberCommand,
	}
}

// Handle 处理删除组织命令 - 级联删除所有相关数据
// 删除顺序: TeamMembers → Teams → OrgMembers → Org
func (h *DeleteHandler) Handle(ctx context.Context, cmd DeleteOrgCommand) error {
	// 1. 检查组织是否存在
	exists, err := h.orgQueryRepo.Exists(ctx, cmd.OrgID)
	if err != nil {
		return fmt.Errorf("failed to check org existence: %w", err)
	}
	if !exists {
		return org.ErrOrgNotFound
	}

	// 2. 级联删除所有团队及其成员
	teams, err := h.teamQuery.ListByOrg(ctx, cmd.OrgID, 0, 1000)
	if err != nil {
		return fmt.Errorf("failed to list teams: %w", err)
	}
	for _, team := range teams {
		// 删除团队成员
		var teamMembers []*org.TeamMember
		teamMembers, err = h.teamMemberQuery.ListByTeam(ctx, team.ID, 0, 1000)
		if err != nil {
			return fmt.Errorf("failed to list team members: %w", err)
		}
		for _, tm := range teamMembers {
			err = h.teamMemberCommand.Remove(ctx, team.ID, tm.UserID)
			if err != nil {
				return fmt.Errorf("failed to remove team member: %w", err)
			}
		}
		// 删除团队
		err = h.teamCommandRepo.Delete(ctx, team.ID)
		if err != nil {
			return fmt.Errorf("failed to delete team: %w", err)
		}
	}

	// 3. 删除所有组织成员
	members, err := h.memberQuery.ListByOrg(ctx, cmd.OrgID, 0, 1000)
	if err != nil {
		return fmt.Errorf("failed to list members: %w", err)
	}
	for _, member := range members {
		err = h.memberCommandRepo.Remove(ctx, member.OrgID, member.UserID)
		if err != nil {
			return fmt.Errorf("failed to remove member: %w", err)
		}
	}

	// 4. 删除组织
	err = h.orgCommandRepo.Delete(ctx, cmd.OrgID)
	if err != nil {
		return fmt.Errorf("failed to delete organization: %w", err)
	}

	return nil
}
