package org

import setting "github.com/lwmacct/260103-ddd-bc-settings/pkg/modules/settings/app/setting"

// OrgSettingDTO 组织配置响应（合并视图）
type OrgSettingDTO struct {
	Key            string              `json:"key"`
	Value          any                 `json:"value"`           // 实际生效值（组织值 > 默认值）
	DefaultValue   any                 `json:"default_value"`   // 系统默认值
	IsCustomized   bool                `json:"is_customized"`   // 是否组织自定义
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
