package setting

import (
	"context"
	"fmt"

	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/app/domain/setting"
)

// UserResetHandler 重置用户配置命令处理器
type UserResetHandler struct {
	cmdRepo setting.UserSettingCommandRepository
}

// NewUserResetHandler 创建 UserResetHandler 实例
func NewUserResetHandler(cmdRepo setting.UserSettingCommandRepository) *UserResetHandler {
	return &UserResetHandler{cmdRepo: cmdRepo}
}

// Handle 处理重置用户配置命令（删除用户自定义值，恢复为系统默认值）
func (h *UserResetHandler) Handle(ctx context.Context, cmd UserResetCommand) error {
	if err := h.cmdRepo.Delete(ctx, cmd.UserID, cmd.Key); err != nil {
		return fmt.Errorf("failed to reset user setting: %w", err)
	}
	return nil
}
