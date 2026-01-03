package user

import (
	"context"
	"fmt"

	"github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/settings/domain/user"
)

// ResetAllHandler 重置所有配置命令处理器
type ResetAllHandler struct {
	cmdRepo user.CommandRepository
}

// NewResetAllHandler 创建重置所有配置命令处理器
func NewResetAllHandler(cmdRepo user.CommandRepository) *ResetAllHandler {
	return &ResetAllHandler{
		cmdRepo: cmdRepo,
	}
}

// Handle 处理重置所有配置命令
//
// 删除用户的所有自定义配置，恢复使用系统默认值
func (h *ResetAllHandler) Handle(ctx context.Context, cmd ResetAllCommand) error {
	if err := h.cmdRepo.DeleteByUser(ctx, cmd.UserID); err != nil {
		return fmt.Errorf("failed to reset all user settings: %w", err)
	}
	return nil
}
