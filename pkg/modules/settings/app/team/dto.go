package team

import setting "github.com/lwmacct/260103-ddd-bc-settings/pkg/modules/settings/app/setting"

// TeamSettingDTO 团队配置响应（合并视图）
type TeamSettingDTO struct {
	Key            string              `json:"key"`
	Value          any                 `json:"value"`           // 实际生效值（团队 > 组织 > 默认值）
	DefaultValue   any                 `json:"default_value"`   // 系统默认值
	OrgValue       any                 `json:"org_value"`       // 组织配置值
	IsCustomized   bool                `json:"is_customized"`   // 是否团队自定义
	IsTeamDefault  bool                `json:"is_team_default"` // 是否为团队可配置的默认值（User 可覆盖）
	InheritedFrom  string              `json:"inherited_from"`  // 继承来源: "team" | "org" | "system"
	VisibleAt      string              `json:"visible_at"`      // 最小可见级别
	ConfigurableAt string              `json:"configurable_at"` // 最大可配置级别
	CategoryID     uint                `json:"category_id"`
	Group          string              `json:"group"`
	ValueType      string              `json:"value_type"`
	Label          string              `json:"label"`
	Order          int                 `json:"order"`
	InputType      string              `json:"input_type"`
	Validation     any                 `json:"validation,omitempty"`
	UIConfig       setting.UIConfigDTO `json:"ui_config"`
}
