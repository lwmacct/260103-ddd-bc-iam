package usersetting

import (
	"context"
	"fmt"

	usersettingdomain "github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/domain/usersetting"
	setting "github.com/lwmacct/260103-ddd-bc-settings/pkg/modules/settings/domain/setting"
)

// UpdateHandler 更新用户设置命令处理器
type UpdateHandler struct {
	userSettingCommandRepo usersettingdomain.CommandRepository
	settingQueryRepo       setting.QueryRepository
}

// NewUpdateHandler 创建更新用户设置命令处理器
func NewUpdateHandler(
	userSettingCommandRepo usersettingdomain.CommandRepository,
	settingQueryRepo setting.QueryRepository,
) *UpdateHandler {
	return &UpdateHandler{
		userSettingCommandRepo: userSettingCommandRepo,
		settingQueryRepo:       settingQueryRepo,
	}
}

// Handle 处理更新用户设置命令
func (h *UpdateHandler) Handle(ctx context.Context, cmd UpdateCommand) error {
	// 1. 验证 Key 是否在 Schema 中存在
	schema, err := h.settingQueryRepo.FindByKey(ctx, cmd.Key)
	if err != nil || schema == nil {
		return usersettingdomain.ErrInvalidSettingKey
	}

	// 2. 验证 Schema 是否允许用户自定义（scope=user or system）
	if schema.Scope != "user" && schema.Scope != "system" {
		return usersettingdomain.ErrInvalidSettingKey
	}

	// 3. 验证值格式（根据 Schema.Validation）
	// TODO: 实现 JSON Logic 验证（Phase 4.6）
	// if schema.Validation != "" {
	//     if err := validateJSONLogic(cmd.Value, schema.Validation); err != nil {
	//         return usersetting.ErrInvalidSettingValue
	//     }
	// }

	// 4. Upsert 用户值
	if err := h.userSettingCommandRepo.Upsert(ctx, cmd.UserID, cmd.Key, cmd.Value); err != nil {
		return fmt.Errorf("failed to upsert user setting: %w", err)
	}

	return nil
}
