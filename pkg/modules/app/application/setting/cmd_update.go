package setting

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/app/domain/setting"
)

// UpdateHandler 更新配置命令处理器
type UpdateHandler struct {
	commandRepo   setting.CommandRepository
	queryRepo     setting.QueryRepository
	validator     setting.Validator
	settingsCache SettingsCacheService
}

// NewUpdateHandler 创建 UpdateHandler 实例
func NewUpdateHandler(
	commandRepo setting.CommandRepository,
	queryRepo setting.QueryRepository,
	validator setting.Validator,
	settingsCache SettingsCacheService,
) *UpdateHandler {
	return &UpdateHandler{
		commandRepo:   commandRepo,
		queryRepo:     queryRepo,
		validator:     validator,
		settingsCache: settingsCache,
	}
}

// Handle 处理更新配置命令
func (h *UpdateHandler) Handle(ctx context.Context, cmd UpdateCommand) (*SettingDTO, error) {
	// 1. 查询配置定义
	def, err := h.queryRepo.FindByKey(ctx, cmd.Key)
	if err != nil {
		return nil, fmt.Errorf("failed to find setting: %w", err)
	}
	if def == nil {
		return nil, errors.New("setting not found")
	}

	// 2. 基于 ValueType 的类型校验
	if err := def.ValidateValue(cmd.DefaultValue); err != nil {
		return nil, fmt.Errorf("value validation failed: %w", err)
	}

	// 3. 基于 InputType 的格式校验（email/url/password 等）
	if err := def.ValidateByInputType(cmd.DefaultValue); err != nil {
		return nil, err
	}

	// 4. 自定义 Validation 规则校验（JSON Logic）
	if h.validator != nil && def.Validation != "" {
		// 获取所有设置用于跨字段验证
		allSettings, _ := h.getAllSettingsMap(ctx)

		vctx := &setting.ValidationContext{
			Key:         cmd.Key,
			Value:       cmd.DefaultValue, // 直接使用 any 类型，无需转换
			Rule:        def.Validation,   // 直接使用实体字段
			AllSettings: allSettings,
		}

		result, validateErr := h.validator.Validate(ctx, vctx)
		if validateErr != nil {
			return nil, fmt.Errorf("validation error: %w", validateErr)
		}
		if !result.Valid {
			return nil, fmt.Errorf("%w: %s", setting.ErrValidationFailed, result.Message)
		}
	}

	// 5. 更新字段
	def.DefaultValue = cmd.DefaultValue
	if cmd.Label != "" {
		def.Label = cmd.Label
	}
	if cmd.UIConfig != "" {
		def.UIConfig = cmd.UIConfig
	}
	if cmd.Order != 0 {
		def.Order = cmd.Order
	}

	// 6. 保存更新
	if err := h.commandRepo.Update(ctx, def); err != nil {
		return nil, fmt.Errorf("failed to update setting: %w", err)
	}

	// 7. 失效 Settings 缓存
	if h.settingsCache != nil {
		if err := h.settingsCache.DeleteAdminSettingsAll(ctx); err != nil {
			slog.Warn("admin settings cache invalidation failed", "key", cmd.Key, "error", err.Error())
		}
	}

	return ToSettingDTO(def), nil
}

// getAllSettingsMap 获取所有配置的 key -> value 映射
func (h *UpdateHandler) getAllSettingsMap(ctx context.Context) (map[string]any, error) {
	defs, err := h.queryRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	result := make(map[string]any, len(defs))
	for _, d := range defs {
		result[d.Key] = d.DefaultValue // 直接使用 any 类型
	}
	return result, nil
}
