package org

import (
	"context"
	"fmt"

	"github.com/lwmacct/260103-ddd-iam-bc/pkg/modules/iam/domain/org"
)

// TeamMemberRemoveHandler 移除团队成员命令处理器
type TeamMemberRemoveHandler struct {
	teamMemberCommand org.TeamMemberCommandRepository
	teamQuery         org.TeamQueryRepository
}

// NewTeamMemberRemoveHandler 创建移除团队成员命令处理器
func NewTeamMemberRemoveHandler(
	teamMemberCommand org.TeamMemberCommandRepository,
	teamQuery org.TeamQueryRepository,
) *TeamMemberRemoveHandler {
	return &TeamMemberRemoveHandler{
		teamMemberCommand: teamMemberCommand,
		teamQuery:         teamQuery,
	}
}

// Handle 处理移除团队成员命令
func (h *TeamMemberRemoveHandler) Handle(ctx context.Context, cmd RemoveTeamMemberCommand) error {
	// 1. 验证团队属于组织
	belongsTo, err := h.teamQuery.BelongsToOrg(ctx, cmd.TeamID, cmd.OrgID)
	if err != nil {
		return fmt.Errorf("failed to check team org: %w", err)
	}
	if !belongsTo {
		return org.ErrTeamNotInOrg
	}

	// 2. 移除成员
	if err := h.teamMemberCommand.Remove(ctx, cmd.TeamID, cmd.UserID); err != nil {
		return fmt.Errorf("failed to remove team member: %w", err)
	}

	return nil
}
