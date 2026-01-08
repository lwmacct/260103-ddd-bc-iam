package org

import (
	"context"
	"fmt"

	"github.com/lwmacct/260103-ddd-iam-bc/pkg/modules/settings/domain/org"
	settingdomain "github.com/lwmacct/260103-ddd-settings-bc/pkg/modules/settings/domain/setting"
)

// SetHandler 设置组织配置命令处理器
type SetHandler struct {
	settingQueryRepo  settingdomain.QueryRepository
	categoryQueryRepo settingdomain.SettingCategoryQueryRepository
	cmdRepo           org.CommandRepository
}

// NewSetHandler 创建设置命令处理器
func NewSetHandler(
	settingQueryRepo settingdomain.QueryRepository,
	categoryQueryRepo settingdomain.SettingCategoryQueryRepository,
	cmdRepo org.CommandRepository,
) *SetHandler {
	return &SetHandler{
		settingQueryRepo:  settingQueryRepo,
		categoryQueryRepo: categoryQueryRepo,
		cmdRepo:           cmdRepo,
	}
}

// Handle 处理设置组织配置命令
//
// 流程：
//  1. 校验配置定义存在（从 Settings BC）
//  2. ValueType 类型校验
//  3. InputType 格式校验（email/url/password 等）
//  4. Upsert 组织配置
func (h *SetHandler) Handle(ctx context.Context, cmd SetCommand) (*SettingsItemDTO, error) {
	// 1. 校验配置定义存在
	def, err := h.settingQueryRepo.FindByKey(ctx, cmd.Key)
	if err != nil {
		return nil, fmt.Errorf("failed to find setting: %w", err)
	}
	if def == nil {
		return nil, org.ErrInvalidSettingKey
	}

	// 2. ValueType 类型校验
	if err := def.ValidateValue(cmd.Value); err != nil {
		return nil, fmt.Errorf("%w: %w", org.ErrInvalidSettingValue, err)
	}

	// 3. InputType 格式校验（email/url/password 等）
	if err := def.ValidateByInputType(cmd.Value); err != nil {
		return nil, fmt.Errorf("%w: %w", org.ErrValidationFailed, err)
	}

	// 4. Upsert 组织配置
	os := org.New(cmd.OrgID, cmd.Key, cmd.Value)
	if err := h.cmdRepo.Upsert(ctx, os); err != nil {
		return nil, fmt.Errorf("failed to save org setting: %w", err)
	}

	// 5. 获取 category key
	category, err := h.categoryQueryRepo.FindByID(ctx, def.CategoryID)
	if err != nil {
		return nil, fmt.Errorf("failed to find category: %w", err)
	}
	categoryKey := ""
	if category != nil {
		categoryKey = category.Key
	}

	return ToSettingsItemDTO(def, os, categoryKey), nil
}
