package setting

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/app/domain/setting"
)

// BatchUpdateHandler 批量更新配置命令处理器
type BatchUpdateHandler struct {
	commandRepo   setting.CommandRepository
	queryRepo     setting.QueryRepository
	validator     setting.Validator
	settingsCache SettingsCacheService
}

// NewBatchUpdateHandler 创建 BatchUpdateHandler 实例
func NewBatchUpdateHandler(
	commandRepo setting.CommandRepository,
	queryRepo setting.QueryRepository,
	validator setting.Validator,
	settingsCache SettingsCacheService,
) *BatchUpdateHandler {
	return &BatchUpdateHandler{
		commandRepo:   commandRepo,
		queryRepo:     queryRepo,
		validator:     validator,
		settingsCache: settingsCache,
	}
}

// Handle 处理批量更新配置命令
func (h *BatchUpdateHandler) Handle(ctx context.Context, cmd BatchUpdateCommand) error {
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

	// 2. 批量查询所有现有配置定义
	existingDefs, err := h.queryRepo.FindByKeys(ctx, keys)
	if err != nil {
		return fmt.Errorf("failed to find settings: %w", err)
	}

	// 3. 构建 key -> def 映射
	existingMap := make(map[string]*setting.Setting, len(existingDefs))
	for _, d := range existingDefs {
		existingMap[d.Key] = d
	}

	// 4. 验证所有 key 存在并更新值
	settings := make([]*setting.Setting, 0, len(cmd.Settings))
	for _, key := range keys {
		existing, ok := existingMap[key]
		if !ok {
			return fmt.Errorf("setting key %s does not exist", key)
		}
		existing.DefaultValue = keyValueMap[key]
		settings = append(settings, existing)
	}

	// 5. 执行验证（如果有验证器）
	if err := h.validateSettings(ctx, settings); err != nil {
		return err
	}

	// 6. 批量更新
	if err := h.commandRepo.BatchUpsert(ctx, settings); err != nil {
		return fmt.Errorf("failed to batch update settings: %w", err)
	}

	// 7. 失效 Settings 缓存
	if h.settingsCache != nil {
		if err := h.settingsCache.DeleteAdminSettingsAll(ctx); err != nil {
			slog.Warn("admin settings cache invalidation failed after batch update", "error", err.Error())
		}
	}

	return nil
}

// validateSettings 执行配置验证
func (h *BatchUpdateHandler) validateSettings(ctx context.Context, settings []*setting.Setting) error {
	if h.validator == nil {
		return nil
	}

	// 获取所有设置用于跨字段验证
	allSettings, _ := h.getAllSettingsMap(ctx, settings)

	// 构建验证上下文
	validationItems := make([]*setting.ValidationContext, 0)
	for _, s := range settings {
		rule := extractValidationRule(s.Validation)
		if rule == "" {
			continue
		}
		validationItems = append(validationItems, &setting.ValidationContext{
			Key:         s.Key,
			Value:       s.DefaultValue,
			Rule:        rule,
			AllSettings: allSettings,
		})
	}

	if len(validationItems) == 0 {
		return nil
	}

	errors, validateErr := h.validator.ValidateBatch(ctx, validationItems)
	if validateErr != nil {
		return fmt.Errorf("validation error: %w", validateErr)
	}
	if len(errors) > 0 {
		for key, msg := range errors {
			return fmt.Errorf("%w: %s - %s", setting.ErrValidationFailed, key, msg)
		}
	}

	return nil
}

// getAllSettingsMap 获取所有配置的 key -> value 映射
func (h *BatchUpdateHandler) getAllSettingsMap(ctx context.Context, pendingUpdates []*setting.Setting) (map[string]any, error) {
	settings, err := h.queryRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	result := make(map[string]any, len(settings))
	for _, s := range settings {
		result[s.Key] = s.DefaultValue
	}

	// 合并待更新的值
	for _, s := range pendingUpdates {
		result[s.Key] = s.DefaultValue
	}

	return result, nil
}
