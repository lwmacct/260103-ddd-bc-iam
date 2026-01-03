package setting

import (
	"context"
	"fmt"

	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/app/domain/setting"
)

// UserResetAllHandler 重置用户所有配置命令处理器
type UserResetAllHandler struct {
	cmdRepo setting.UserSettingCommandRepository
}

// NewUserResetAllHandler 创建 UserResetAllHandler 实例
func NewUserResetAllHandler(cmdRepo setting.UserSettingCommandRepository) *UserResetAllHandler {
	return &UserResetAllHandler{cmdRepo: cmdRepo}
}

// Handle 处理重置用户所有配置命令
func (h *UserResetAllHandler) Handle(ctx context.Context, cmd UserResetAllCommand) error {
	if err := h.cmdRepo.DeleteByUser(ctx, cmd.UserID); err != nil {
		return fmt.Errorf("failed to reset all user settings: %w", err)
	}
	return nil
}
