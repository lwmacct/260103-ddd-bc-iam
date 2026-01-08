package org

import (
	"context"
	"fmt"

	"github.com/lwmacct/260103-ddd-iam-bc/pkg/modules/iam/domain/org"
)

// MemberUpdateRoleHandler 更新成员角色命令处理器
type MemberUpdateRoleHandler struct {
	memberCommand org.MemberCommandRepository
	memberQuery   org.MemberQueryRepository
}

// NewMemberUpdateRoleHandler 创建更新成员角色命令处理器
func NewMemberUpdateRoleHandler(
	memberCommand org.MemberCommandRepository,
	memberQuery org.MemberQueryRepository,
) *MemberUpdateRoleHandler {
	return &MemberUpdateRoleHandler{
		memberCommand: memberCommand,
		memberQuery:   memberQuery,
	}
}

// Handle 处理更新成员角色命令
func (h *MemberUpdateRoleHandler) Handle(ctx context.Context, cmd UpdateMemberRoleCommand) error {
	// 1. 验证角色
	newRole := org.MemberRole(cmd.Role)
	if !org.IsValidMemberRole(newRole) {
		return org.ErrInvalidMemberRole
	}

	// 2. 获取当前成员信息
	member, err := h.memberQuery.GetByOrgAndUser(ctx, cmd.OrgID, cmd.UserID)
	if err != nil {
		return err
	}

	// 3. 如果当前是 owner 且要降级，检查是否是最后一个
	if member.IsOwner() && newRole != org.MemberRoleOwner {
		ownerCount, err := h.memberQuery.CountOwners(ctx, cmd.OrgID)
		if err != nil {
			return fmt.Errorf("failed to count owners: %w", err)
		}
		if ownerCount <= 1 {
			return org.ErrCannotDemoteLastOwner
		}
	}

	// 4. 更新角色
	if err := h.memberCommand.UpdateRole(ctx, cmd.OrgID, cmd.UserID, newRole); err != nil {
		return fmt.Errorf("failed to update role: %w", err)
	}

	return nil
}
