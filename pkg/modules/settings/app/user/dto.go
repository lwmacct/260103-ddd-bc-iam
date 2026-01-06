package user

import (
	"time"

	setting "github.com/lwmacct/260103-ddd-bc-settings/pkg/modules/settings/app/setting"
)

// UserSettingDTO 用户配置响应（合并视图）
type UserSettingDTO struct {
	Key          string              `json:"key"`
	Value        any                 `json:"value"`         // 实际生效值（用户值 > 默认值）
	DefaultValue any                 `json:"default_value"` // 系统默认值
	IsCustomized bool                `json:"is_customized"` // 是否用户自定义
	CategoryID   uint                `json:"category_id"`
	Group        string              `json:"group"`
	ValueType    string              `json:"value_type"`
	Label        string              `json:"label"`
	Order        int                 `json:"order"`
	InputType    string              `json:"input_type"`
	Validation   any                 `json:"validation,omitempty"`
	UIConfig     setting.UIConfigDTO `json:"ui_config"`
}

// CategoryDTO 配置分类响应
type CategoryDTO struct {
	ID        uint      `json:"id"`
	Key       string    `json:"key"`
	Label     string    `json:"label"`
	Icon      string    `json:"icon"`
	Order     int       `json:"order"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// SettingsCategoryDTO 层级结构响应
type SettingsCategoryDTO struct {
	Category string             `json:"category"`
	Label    string             `json:"label"`
	Icon     string             `json:"icon"`
	Groups   []SettingsGroupDTO `json:"groups"`
}

// SettingsGroupDTO 配置分组
type SettingsGroupDTO struct {
	Name     string           `json:"name"`
	Settings []UserSettingDTO `json:"settings"`
}

// SettingItemDTO 配置项（用于批量设置）
type SettingItemDTO struct {
	Key   string `json:"key"`
	Value any    `json:"value"`
}
