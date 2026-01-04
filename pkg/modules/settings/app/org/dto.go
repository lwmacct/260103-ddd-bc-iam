package org

// OrgSettingDTO 组织配置响应（合并视图）
type OrgSettingDTO struct {
	Key          string      `json:"key"`
	Value        any         `json:"value"`         // 实际生效值（组织值 > 默认值）
	DefaultValue any         `json:"default_value"` // 系统默认值
	IsCustomized bool        `json:"is_customized"` // 是否组织自定义
	CategoryID   uint        `json:"category_id"`
	Group        string      `json:"group"`
	ValueType    string      `json:"value_type"`
	Label        string      `json:"label"`
	Order        int         `json:"order"`
	InputType    string      `json:"input_type"`
	Validation   string      `json:"validation,omitempty"`
	UIConfig     UIConfigDTO `json:"ui_config"`
}

// UIConfigDTO UI 配置
type UIConfigDTO struct {
	Hint      string            `json:"hint,omitempty"`
	Options   []SelectOptionDTO `json:"options,omitempty"`
	DependsOn *DependsOnDTO     `json:"depends_on,omitempty"`
}

// SelectOptionDTO 下拉选项
type SelectOptionDTO struct {
	Label string `json:"label"`
	Value any    `json:"value"`
}

// DependsOnDTO 依赖关系配置
type DependsOnDTO struct {
	Key   string `json:"key"`
	Value any    `json:"value"`
}
