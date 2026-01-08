package team

import (
	"encoding/json"

	"github.com/lwmacct/260103-ddd-iam-bc/pkg/modules/settings/domain/org"
	"github.com/lwmacct/260103-ddd-iam-bc/pkg/modules/settings/domain/team"
	setting "github.com/lwmacct/260103-ddd-settings-bc/pkg/modules/settings/app/setting"
	settingdomain "github.com/lwmacct/260103-ddd-settings-bc/pkg/modules/settings/domain/setting"
)

// ToSettingsItemDTO 将配置定义、团队配置和组织配置合并为扁平结构 DTO
//
// 三级继承优先级：团队 > 组织 > 系统默认值
// InheritedFrom 字段指示当前生效值的来源
func ToSettingsItemDTO(def *settingdomain.Setting, ts *team.TeamSetting, os *org.OrgSetting, categoryKey string) *setting.SettingsItemDTO {
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

	// 三级继承：团队 > 组织 > 系统
	var hasTeam, hasOrg bool
	if ts != nil {
		hasTeam = !ts.IsEmpty()
	}
	if os != nil {
		hasOrg = !os.IsEmpty()
	}

	switch {
	case hasTeam:
		// 团队有自定义值
		dto.Value = ts.Value
		dto.IsCustomized = true
	case hasOrg:
		// 组织有配置值
		dto.Value = os.Value
		dto.IsCustomized = false
	default:
		// 使用系统默认值
		dto.Value = def.DefaultValue
		dto.IsCustomized = false
	}

	return dto
}
