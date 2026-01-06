package user

import (
	"encoding/json"

	"github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/settings/domain/user"
	setting "github.com/lwmacct/260103-ddd-bc-settings/pkg/modules/settings/app/setting"
	settingdomain "github.com/lwmacct/260103-ddd-bc-settings/pkg/modules/settings/domain/setting"
)

// ToSettingsItemDTO 将配置定义和用户配置合并为扁平结构 DTO
//
// 参数：
//   - def: 配置定义（来自 Settings BC）
//   - us: 用户配置（可为 nil，表示使用默认值）
//   - categoryKey: 分类 Key（用于前端分组）
//
// 返回：
//   - 合并后的 SettingsItemDTO（复用上游类型）
func ToSettingsItemDTO(def *settingdomain.Setting, us *user.UserSetting, categoryKey string) *setting.SettingsItemDTO {
	if def == nil {
		return nil
	}

	group := def.Group
	if group == "" {
		group = "default"
	}

	dto := &setting.SettingsItemDTO{
		Key:            def.Key,
		Category:       categoryKey,
		Group:          group,
		Value:          def.DefaultValue,
		DefaultValue:   def.DefaultValue,
		IsCustomized:   false,
		VisibleAt:      def.VisibleAt,
		ConfigurableAt: def.ConfigurableAt,
		ValueType:      def.ValueType,
		Label:          def.Label,
		Order:          def.Order,
		InputType:      def.InputType,
		Validation:     def.Validation,
		UIConfig:       parseUIConfig(def.UIConfig),
	}

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
		ID:        c.ID,
		Key:       c.Key,
		Label:     c.Label,
		Icon:      c.Icon,
		Order:     c.Order,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
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
func parseUIConfig(raw string) setting.UIConfigDTO {
	if raw == "" {
		return setting.UIConfigDTO{}
	}

	var config struct {
		Hint    string `json:"hint"`
		Options []struct {
			Label string `json:"label"`
			Value string `json:"value"`
		} `json:"options"`
		DependsOn *struct {
			Key      string `json:"key"`
			Value    any    `json:"value"`
			Operator string `json:"operator"`
		} `json:"depends_on"`
	}

	if err := json.Unmarshal([]byte(raw), &config); err != nil {
		return setting.UIConfigDTO{}
	}

	dto := setting.UIConfigDTO{
		Hint: config.Hint,
	}

	// 转换 options
	if len(config.Options) > 0 {
		dto.Options = make([]setting.SelectOptionDTO, len(config.Options))
		for i, opt := range config.Options {
			dto.Options[i] = setting.SelectOptionDTO{
				Label: opt.Label,
				Value: opt.Value,
			}
		}
	}

	// 转换 depends_on
	if config.DependsOn != nil {
		dto.DependsOn = &setting.DependsOnConfigDTO{
			Key:      config.DependsOn.Key,
			Value:    config.DependsOn.Value,
			Operator: config.DependsOn.Operator,
		}
	}

	return dto
}
