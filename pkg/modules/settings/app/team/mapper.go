package team

import (
	"encoding/json"

	"github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/settings/domain/org"
	"github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/settings/domain/team"
	settingdomain "github.com/lwmacct/260103-ddd-bc-settings/pkg/modules/settings/domain/setting"
)

// ToTeamSettingDTO 将配置定义、团队配置和组织配置合并为 DTO
// 优先级：团队 > 组织 > 系统默认
func ToTeamSettingDTO(def *settingdomain.Setting, ts *team.TeamSetting, os *org.OrgSetting) *TeamSettingDTO {
	dto := &TeamSettingDTO{
		Key:            def.Key,
		DefaultValue:   def.DefaultValue,
		CategoryID:     def.CategoryID,
		Group:          def.Group,
		ValueType:      def.ValueType,
		Label:          def.Label,
		Order:          def.Order,
		InputType:      def.InputType,
		Validation:     def.Validation,
		VisibleAt:      def.VisibleAt,
		ConfigurableAt: def.ConfigurableAt,
		IsTeamDefault:  def.IsTeamDefaultForUser(), // 是否为团队可为用户设置的默认值
	}

	// 解析 UIConfig
	if def.UIConfig != "" {
		_ = json.Unmarshal([]byte(def.UIConfig), &dto.UIConfig)
	}

	// 三级继承：团队 > 组织 > 系统
	// 优先级判断：使用嵌套的三元表达式风格实现 switch
	var hasTeam, hasOrg bool
	if ts != nil {
		hasTeam = !ts.IsEmpty()
	}
	if os != nil {
		hasOrg = !os.IsEmpty()
	}

	// 根据优先级选择源
	switch {
	case hasTeam:
		// 团队有自定义值
		dto.Value = ts.Value
		dto.IsCustomized = true
		dto.InheritedFrom = "team"
	case hasOrg:
		// 组织有配置值
		dto.Value = os.Value
		dto.OrgValue = os.Value
		dto.IsCustomized = false
		dto.InheritedFrom = "org"
	default:
		// 使用系统默认值
		dto.Value = def.DefaultValue
		dto.IsCustomized = false
		dto.InheritedFrom = "system"
	}

	return dto
}
