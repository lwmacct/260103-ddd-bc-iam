package org

import (
	"encoding/json"

	"github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/settings/domain/org"
	settingdomain "github.com/lwmacct/260103-ddd-bc-settings/pkg/modules/settings/domain/setting"
)

// ToOrgSettingDTO 将配置定义和组织配置合并为 DTO
func ToOrgSettingDTO(def *settingdomain.Setting, os *org.OrgSetting) *OrgSettingDTO {
	dto := &OrgSettingDTO{
		Key:          def.Key,
		DefaultValue: def.DefaultValue,
		CategoryID:   def.CategoryID,
		Group:        def.Group,
		ValueType:    def.ValueType,
		Label:        def.Label,
		Order:        def.Order,
		InputType:    def.InputType,
		Validation:   def.Validation,
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
