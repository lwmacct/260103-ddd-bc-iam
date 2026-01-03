package usersetting

import (
	"context"
	"fmt"

	usersettingdomain "github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/domain/usersetting"
)

// DeleteHandler 删除用户设置命令处理器
type DeleteHandler struct {
	userSettingQueryRepo   usersettingdomain.QueryRepository
	userSettingCommandRepo usersettingdomain.CommandRepository
}

// NewDeleteHandler 创建删除用户设置命令处理器
func NewDeleteHandler(
	userSettingQueryRepo usersettingdomain.QueryRepository,
	userSettingCommandRepo usersettingdomain.CommandRepository,
) *DeleteHandler {
	return &DeleteHandler{
		userSettingQueryRepo:   userSettingQueryRepo,
		userSettingCommandRepo: userSettingCommandRepo,
	}
}

// Handle 处理删除用户设置命令（软删除）
func (h *DeleteHandler) Handle(ctx context.Context, cmd DeleteCommand) error {
	// 1. 查找用户设置
	userSetting, err := h.userSettingQueryRepo.FindByUserAndKey(ctx, cmd.UserID, cmd.Key)
	if err != nil {
		return fmt.Errorf("failed to find user setting: %w", err)
	}

	// 2. 如果不存在，直接返回成功（幂等）
	if userSetting == nil {
		return nil
	}

	// 3. 软删除（恢复系统默认值）
	if err := h.userSettingCommandRepo.Delete(ctx, userSetting.ID); err != nil {
		return fmt.Errorf("failed to delete user setting: %w", err)
	}

	return nil
}
