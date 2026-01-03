package org

import (
	"context"
	"fmt"

	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/iam/domain/org"
)

// TeamUpdateHandler 更新团队命令处理器
type TeamUpdateHandler struct {
	teamCommand org.TeamCommandRepository
	teamQuery   org.TeamQueryRepository
}

// NewTeamUpdateHandler 创建更新团队命令处理器
func NewTeamUpdateHandler(
	teamCommand org.TeamCommandRepository,
	teamQuery org.TeamQueryRepository,
) *TeamUpdateHandler {
	return &TeamUpdateHandler{
		teamCommand: teamCommand,
		teamQuery:   teamQuery,
	}
}

// Handle 处理更新团队命令
func (h *TeamUpdateHandler) Handle(ctx context.Context, cmd UpdateTeamCommand) (*TeamDTO, error) {
	// 1. 获取现有团队
	team, err := h.teamQuery.GetByID(ctx, cmd.TeamID)
	if err != nil {
		return nil, err
	}

	// 2. 验证团队属于指定组织
	if !team.BelongsTo(cmd.OrgID) {
		return nil, org.ErrTeamNotInOrg
	}

	// 3. 更新字段
	if cmd.DisplayName != nil {
		team.DisplayName = *cmd.DisplayName
	}
	if cmd.Description != nil {
		team.Description = *cmd.Description
	}
	if cmd.Avatar != nil {
		team.Avatar = *cmd.Avatar
	}

	// 4. 保存更新
	if err := h.teamCommand.Update(ctx, team); err != nil {
		return nil, fmt.Errorf("failed to update team: %w", err)
	}

	return ToTeamDTO(team), nil
}
