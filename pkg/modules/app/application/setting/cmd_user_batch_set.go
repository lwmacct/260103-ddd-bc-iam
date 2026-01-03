package setting

import (
	"context"
	"fmt"
	"maps"

	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/app/domain/setting"
)

// UserBatchSetHandler 批量设置用户配置命令处理器
type UserBatchSetHandler struct {
	settingQueryRepo setting.QueryRepository
	cmdRepo          setting.UserSettingCommandRepository
	validator        setting.Validator
}

// NewUserBatchSetHandler 创建 UserBatchSetHandler 实例
func NewUserBatchSetHandler(
	settingQueryRepo setting.QueryRepository,
	cmdRepo setting.UserSettingCommandRepository,
	validator setting.Validator,
) *UserBatchSetHandler {
	return &UserBatchSetHandler{
		settingQueryRepo: settingQueryRepo,
		cmdRepo:          cmdRepo,
		validator:        validator,
	}
}

// Handle 处理批量设置用户配置命令
func (h *UserBatchSetHandler) Handle(ctx context.Context, cmd UserBatchSetCommand) error {
	if len(cmd.Settings) == 0 {
		return nil
	}

	// 1. 提取所有 keys
	keys := make([]string, len(cmd.Settings))
	keyValueMap := make(map[string]any, len(cmd.Settings))
	for i, item := range cmd.Settings {
		keys[i] = item.Key
		keyValueMap[item.Key] = item.Value
	}

	// 2. 批量查询所有配置定义
	defs, err := h.settingQueryRepo.FindByKeys(ctx, keys)
	if err != nil {
		return fmt.Errorf("failed to find setting definitions: %w", err)
	}

	// 3. 构建 key -> def 映射并验证所有 key 存在
	defMap := make(map[string]*setting.Setting, len(defs))
	for _, d := range defs {
		defMap[d.Key] = d
	}

	for _, key := range keys {
		if _, ok := defMap[key]; !ok {
			return fmt.Errorf("setting key %s does not exist", key)
		}
	}

	// 4. 执行验证
	if err := h.validateSettings(ctx, defMap, keyValueMap); err != nil {
		return err
	}

	// 5. 构建用户配置并批量 Upsert
	userSettings := make([]*setting.UserSetting, 0, len(cmd.Settings))
	for _, item := range cmd.Settings {
		userSettings = append(userSettings, &setting.UserSetting{
			UserID:     cmd.UserID,
			SettingKey: item.Key,
			Value:      item.Value,
		})
	}

	if err := h.cmdRepo.BatchUpsert(ctx, userSettings); err != nil {
		return fmt.Errorf("failed to batch set user settings: %w", err)
	}

	return nil
}

// validateSettings 批量验证配置值
func (h *UserBatchSetHandler) validateSettings(
	ctx context.Context,
	defMap map[string]*setting.Setting,
	keyValueMap map[string]any,
) error {
	if h.validator == nil {
		return nil
	}

	// 获取所有设置用于跨字段验证
	allSettings, _ := h.getAllSettingsMap(ctx, keyValueMap)

	// 构建验证上下文
	validationItems := make([]*setting.ValidationContext, 0, len(keyValueMap))
	for key, value := range keyValueMap {
		def := defMap[key]
		rule := extractValidationRule(def.Validation)
		if rule == "" {
			continue
		}
		validationItems = append(validationItems, &setting.ValidationContext{
			Key:         key,
			Value:       value,
			Rule:        rule,
			AllSettings: allSettings,
		})
	}

	if len(validationItems) == 0 {
		return nil
	}

	errors, err := h.validator.ValidateBatch(ctx, validationItems)
	if err != nil {
		return fmt.Errorf("validation error: %w", err)
	}
	if len(errors) > 0 {
		for key, msg := range errors {
			return fmt.Errorf("%w: %s - %s", setting.ErrValidationFailed, key, msg)
		}
	}

	return nil
}

// getAllSettingsMap 获取所有配置的 key -> value 映射（合并待更新值）
func (h *UserBatchSetHandler) getAllSettingsMap(ctx context.Context, pendingUpdates map[string]any) (map[string]any, error) {
	defs, err := h.settingQueryRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	result := make(map[string]any, len(defs))
	for _, d := range defs {
		result[d.Key] = d.DefaultValue
	}

	// 合并待更新的值
	maps.Copy(result, pendingUpdates)

	return result, nil
}
