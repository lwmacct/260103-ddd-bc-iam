package setting

import (
	"context"
	"errors"
	"fmt"

	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/app/domain/setting"
)

// UserSetHandler 设置用户配置命令处理器
type UserSetHandler struct {
	settingQueryRepo setting.QueryRepository
	cmdRepo          setting.UserSettingCommandRepository
	validator        setting.Validator
}

// NewUserSetHandler 创建 UserSetHandler 实例
func NewUserSetHandler(
	settingQueryRepo setting.QueryRepository,
	cmdRepo setting.UserSettingCommandRepository,
	validator setting.Validator,
) *UserSetHandler {
	return &UserSetHandler{
		settingQueryRepo: settingQueryRepo,
		cmdRepo:          cmdRepo,
		validator:        validator,
	}
}

// Handle 处理设置用户配置命令
func (h *UserSetHandler) Handle(ctx context.Context, cmd UserSetCommand) (*UserSettingDTO, error) {
	// 1. 验证配置定义是否存在
	def, err := h.settingQueryRepo.FindByKey(ctx, cmd.Key)
	if err != nil {
		return nil, fmt.Errorf("failed to find setting definition: %w", err)
	}
	if def == nil {
		return nil, errors.New("setting key does not exist")
	}

	// 2. 基于 ValueType 的类型校验
	if err := def.ValidateValue(cmd.Value); err != nil {
		return nil, fmt.Errorf("value validation failed: %w", err)
	}

	// 3. 基于 InputType 的格式校验（email/url/password 等）
	if err := def.ValidateByInputType(cmd.Value); err != nil {
		return nil, err
	}

	// 4. 自定义 Validation 规则校验（JSON Logic）
	if err := h.validateWithRule(ctx, def, cmd.Value); err != nil {
		return nil, err
	}

	// 5. Upsert 用户配置
	us := &setting.UserSetting{
		UserID:     cmd.UserID,
		SettingKey: cmd.Key,
		Value:      cmd.Value,
	}
	if err := h.cmdRepo.Upsert(ctx, us); err != nil {
		return nil, fmt.Errorf("failed to set user setting: %w", err)
	}

	return ToUserSettingDTO(def, us), nil
}

// validateWithRule 使用自定义规则验证配置值
func (h *UserSetHandler) validateWithRule(ctx context.Context, def *setting.Setting, value any) error {
	if h.validator == nil || def.Validation == "" {
		return nil
	}

	// 获取所有设置用于跨字段验证
	allSettings, _ := h.getAllSettingsMap(ctx)

	vctx := &setting.ValidationContext{
		Key:         def.Key,
		Value:       value,
		Rule:        def.Validation, // 直接使用实体字段
		AllSettings: allSettings,
	}

	result, err := h.validator.Validate(ctx, vctx)
	if err != nil {
		return fmt.Errorf("validation error: %w", err)
	}
	if !result.Valid {
		return fmt.Errorf("%w: %s", setting.ErrValidationFailed, result.Message)
	}

	return nil
}

// getAllSettingsMap 获取所有配置的 key -> value 映射
func (h *UserSetHandler) getAllSettingsMap(ctx context.Context) (map[string]any, error) {
	defs, err := h.settingQueryRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	result := make(map[string]any, len(defs))
	for _, d := range defs {
		result[d.Key] = d.DefaultValue
	}
	return result, nil
}
