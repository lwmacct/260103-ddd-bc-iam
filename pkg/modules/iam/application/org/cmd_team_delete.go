package org

import (
	"context"
	"fmt"

	"github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/domain/org"
)

// TeamDeleteHandler 删除团队命令处理器
type TeamDeleteHandler struct {
	teamCommand       org.TeamCommandRepository
	teamQuery         org.TeamQueryRepository
	teamMemberQuery   org.TeamMemberQueryRepository
	teamMemberCommand org.TeamMemberCommandRepository
}

// NewTeamDeleteHandler 创建删除团队命令处理器
func NewTeamDeleteHandler(
	teamCommand org.TeamCommandRepository,
	teamQuery org.TeamQueryRepository,
	teamMemberQuery org.TeamMemberQueryRepository,
	teamMemberCommand org.TeamMemberCommandRepository,
) *TeamDeleteHandler {
	return &TeamDeleteHandler{
		teamCommand:       teamCommand,
		teamQuery:         teamQuery,
		teamMemberQuery:   teamMemberQuery,
		teamMemberCommand: teamMemberCommand,
	}
}

// Handle 处理删除团队命令 - 级联删除所有团队成员
func (h *TeamDeleteHandler) Handle(ctx context.Context, cmd DeleteTeamCommand) error {
	// 1. 获取团队
	team, err := h.teamQuery.GetByID(ctx, cmd.TeamID)
	if err != nil {
		return err
	}

	// 2. 验证团队属于指定组织
	if !team.BelongsTo(cmd.OrgID) {
		return org.ErrTeamNotInOrg
	}

	// 3. 级联删除所有团队成员
	members, err := h.teamMemberQuery.ListByTeam(ctx, cmd.TeamID, 0, 1000)
	if err != nil {
		return fmt.Errorf("failed to list team members: %w", err)
	}
	for _, member := range members {
		if err := h.teamMemberCommand.Remove(ctx, cmd.TeamID, member.UserID); err != nil {
			return fmt.Errorf("failed to remove team member: %w", err)
		}
	}

	// 4. 删除团队
	if err := h.teamCommand.Delete(ctx, cmd.TeamID); err != nil {
		return fmt.Errorf("failed to delete team: %w", err)
	}

	return nil
}
