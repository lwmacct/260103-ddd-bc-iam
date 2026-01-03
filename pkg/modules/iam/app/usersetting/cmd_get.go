package usersetting

import (
	"context"
	"encoding/json"
	"fmt"

	usersettingdomain "github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/domain/usersetting"
	setting "github.com/lwmacct/260103-ddd-bc-settings/pkg/modules/settings/domain/setting"
)

// GetHandler 获取单个用户设置命令处理器
type GetHandler struct {
	userSettingQueryRepo usersettingdomain.QueryRepository
	settingQueryRepo     setting.QueryRepository
}

// NewGetHandler 创建获取单个用户设置命令处理器
func NewGetHandler(
	userSettingQueryRepo usersettingdomain.QueryRepository,
	settingQueryRepo setting.QueryRepository,
) *GetHandler {
	return &GetHandler{
		userSettingQueryRepo: userSettingQueryRepo,
		settingQueryRepo:     settingQueryRepo,
	}
}

// Handle 处理获取单个用户设置查询
func (h *GetHandler) Handle(ctx context.Context, userID uint, key string) (*UserSettingDTO, error) {
	// 1. 获取 Schema（验证键名存在）
	schema, err := h.settingQueryRepo.FindByKey(ctx, key)
	if err != nil || schema == nil {
		return nil, usersettingdomain.ErrInvalidSettingKey
	}

	// 2. 检查 Schema 是否对用户可见
	if schema.Scope != "user" && schema.Scope != "system" {
		return nil, usersettingdomain.ErrInvalidSettingKey
	}

	// 3. 获取用户自定义值
	userValue, err := h.userSettingQueryRepo.FindByUserAndKey(ctx, userID, key)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user setting: %w", err)
	}

	// 4. 构建 DTO
	// 将 DefaultValue 转换为 JSON 字符串
	defaultValueJSON, err := json.Marshal(schema.DefaultValue)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal default value: %w", err)
	}

	dto := &UserSettingDTO{
		Key:        schema.Key,
		CategoryID: schema.CategoryID,
		ValueType:  schema.ValueType,
		Label:      schema.Label,
	}

	// 如果用户有自定义值，使用用户值；否则使用系统默认值
	if userValue != nil {
		dto.Value = userValue.Value
		dto.IsCustom = true
	} else {
		dto.Value = string(defaultValueJSON)
		dto.IsCustom = false
	}

	return dto, nil
}
