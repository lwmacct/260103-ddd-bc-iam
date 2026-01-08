package org

import (
	"encoding/json"

	setting "github.com/lwmacct/260103-ddd-bc-settings/pkg/modules/settings/app/setting"
	settingdomain "github.com/lwmacct/260103-ddd-bc-settings/pkg/modules/settings/domain/setting"
	"github.com/lwmacct/260103-ddd-iam-bc/pkg/modules/settings/domain/org"
)

// ToSettingsItemDTO 将配置定义和组织配置合并为扁平结构 DTO
func ToSettingsItemDTO(def *settingdomain.Setting, os *org.OrgSetting, categoryKey string) *setting.SettingsItemDTO {
	group := def.Group
	if group == "" {
		group = "default"
	}

	dto := &setting.SettingsItemDTO{
		Key:            def.Key,
		Category:       categoryKey,
		Group:          group,
		DefaultValue:   def.DefaultValue,
		VisibleAt:      def.VisibleAt,
		ConfigurableAt: def.ConfigurableAt,
		ValueType:      def.ValueType,
		Label:          def.Label,
		Order:          def.Order,
		InputType:      def.InputType,
		Validation:     def.Validation,
	}

	// 解析 UIConfig
	if def.UIConfig != "" {
		_ = json.Unmarshal([]byte(def.UIConfig), &dto.UIConfig)
	}

	// 如果有组织自定义值，使用组织值
	if os != nil && !os.IsEmpty() {
		dto.Value = os.Value
		dto.IsCustomized = true
	} else {
		dto.Value = def.DefaultValue
		dto.IsCustomized = false
	}

	return dto
}
