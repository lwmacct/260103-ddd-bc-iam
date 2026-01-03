package userset

import (
	"encoding/json"

	"github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/user_settings/domain/userset"
	settingdomain "github.com/lwmacct/260103-ddd-bc-settings/pkg/modules/settings/domain/setting"
)

// ToUserSettingDTO 将配置定义和用户配置合并为 DTO
//
// 参数：
//   - def: 配置定义（来自 Settings BC）
//   - us: 用户配置（可为 nil，表示使用默认值）
//
// 返回：
//   - 合并后的 UserSettingDTO
func ToUserSettingDTO(def *settingdomain.Setting, us *userset.UserSetting) *UserSettingDTO {
	if def == nil {
		return nil
	}

	dto := &UserSettingDTO{
		Key:          def.Key,
		Value:        def.DefaultValue, // 默认使用系统默认值
		DefaultValue: def.DefaultValue,
		IsCustomized: false,
		CategoryID:   def.CategoryID,
		Group:        def.Group,
		ValueType:    def.ValueType,
		Label:        def.Label,
		Order:        def.Order,
		InputType:    def.InputType,
		Validation:   def.Validation,
		UIConfig:     parseUIConfig(def.UIConfig),
	}

	// 如果用户有自定义值，使用用户值
	if us != nil {
		dto.Value = us.Value
		dto.IsCustomized = true
	}

	return dto
}

// ToCategoryDTO 将分类实体转换为 DTO
func ToCategoryDTO(c *settingdomain.SettingCategory) *CategoryDTO {
	if c == nil {
		return nil
	}
	return &CategoryDTO{
		ID:    c.ID,
		Key:   c.Key,
		Label: c.Label,
		Icon:  c.Icon,
		Order: c.Order,
	}
}

// ToCategoryDTOs 批量转换分类实体为 DTO
func ToCategoryDTOs(categories []*settingdomain.SettingCategory) []*CategoryDTO {
	dtos := make([]*CategoryDTO, 0, len(categories))
	for _, c := range categories {
		if c != nil {
			dtos = append(dtos, ToCategoryDTO(c))
		}
	}
	return dtos
}

// parseUIConfig 解析 UIConfig JSON 字符串
func parseUIConfig(raw string) UIConfigDTO {
	if raw == "" {
		return UIConfigDTO{}
	}

	var config struct {
		Hint    string `json:"hint"`
		Options []struct {
			Label string `json:"label"`
			Value any    `json:"value"`
		} `json:"options"`
		DependsOn *struct {
			Key   string `json:"key"`
			Value any    `json:"value"`
		} `json:"depends_on"`
	}

	if err := json.Unmarshal([]byte(raw), &config); err != nil {
		return UIConfigDTO{}
	}

	dto := UIConfigDTO{
		Hint: config.Hint,
	}

	// 转换 options
	if len(config.Options) > 0 {
		dto.Options = make([]SelectOptionDTO, len(config.Options))
		for i, opt := range config.Options {
			dto.Options[i] = SelectOptionDTO{
				Label: opt.Label,
				Value: opt.Value,
			}
		}
	}

	// 转换 depends_on
	if config.DependsOn != nil {
		dto.DependsOn = &DependsOnDTO{
			Key:   config.DependsOn.Key,
			Value: config.DependsOn.Value,
		}
	}

	return dto
}
