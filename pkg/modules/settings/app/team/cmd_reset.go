package team

import (
	"context"
	"fmt"

	"github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/settings/domain/team"
)

// ResetHandler 重置单个配置命令处理器
type ResetHandler struct {
	cmdRepo team.CommandRepository
}

// NewResetHandler 创建重置命令处理器
func NewResetHandler(cmdRepo team.CommandRepository) *ResetHandler {
	return &ResetHandler{
		cmdRepo: cmdRepo,
	}
}

// Handle 处理重置单个配置命令
//
// 重置配置即删除团队自定义值，恢复使用组织配置或系统默认值
func (h *ResetHandler) Handle(ctx context.Context, cmd ResetCommand) error {
	if err := h.cmdRepo.Delete(ctx, cmd.TeamID, cmd.Key); err != nil {
		return fmt.Errorf("failed to reset team setting: %w", err)
	}
	return nil
}
