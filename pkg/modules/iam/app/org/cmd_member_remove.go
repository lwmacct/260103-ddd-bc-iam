package org

import (
	"context"
	"fmt"

	"github.com/lwmacct/260103-ddd-iam-bc/pkg/modules/iam/domain/org"
)

// MemberRemoveHandler 移除组织成员命令处理器
type MemberRemoveHandler struct {
	memberCommand org.MemberCommandRepository
	memberQuery   org.MemberQueryRepository
}

// NewMemberRemoveHandler 创建移除成员命令处理器
func NewMemberRemoveHandler(
	memberCommand org.MemberCommandRepository,
	memberQuery org.MemberQueryRepository,
) *MemberRemoveHandler {
	return &MemberRemoveHandler{
		memberCommand: memberCommand,
		memberQuery:   memberQuery,
	}
}

// Handle 处理移除成员命令
func (h *MemberRemoveHandler) Handle(ctx context.Context, cmd RemoveMemberCommand) error {
	// 1. 获取成员信息
	member, err := h.memberQuery.GetByOrgAndUser(ctx, cmd.OrgID, cmd.UserID)
	if err != nil {
		return err
	}

	// 2. 如果是 owner，检查是否是最后一个
	if member.IsOwner() {
		ownerCount, err := h.memberQuery.CountOwners(ctx, cmd.OrgID)
		if err != nil {
			return fmt.Errorf("failed to count owners: %w", err)
		}
		if ownerCount <= 1 {
			return org.ErrCannotRemoveLastOwner
		}
	}

	// 3. 移除成员
	if err := h.memberCommand.Remove(ctx, cmd.OrgID, cmd.UserID); err != nil {
		return fmt.Errorf("failed to remove member: %w", err)
	}

	return nil
}
