package userset

import (
	"context"
	"fmt"

	"github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/user_settings/domain/userset"
)

// ResetHandler 重置单个配置命令处理器
type ResetHandler struct {
	cmdRepo userset.CommandRepository
}

// NewResetHandler 创建重置命令处理器
func NewResetHandler(cmdRepo userset.CommandRepository) *ResetHandler {
	return &ResetHandler{
		cmdRepo: cmdRepo,
	}
}

// Handle 处理重置单个配置命令
//
// 重置配置即删除用户自定义值，恢复使用系统默认值
func (h *ResetHandler) Handle(ctx context.Context, cmd ResetCommand) error {
	if err := h.cmdRepo.Delete(ctx, cmd.UserID, cmd.Key); err != nil {
		return fmt.Errorf("failed to reset user setting: %w", err)
	}
	return nil
}
