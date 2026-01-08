package user

import (
	"time"

	setting "github.com/lwmacct/260103-ddd-bc-settings/pkg/modules/settings/app/setting"
)

// SettingsItemDTO 用户配置项响应（复用上游扁平结构类型）
type SettingsItemDTO = setting.SettingsItemDTO

// CategoryDTO 配置分类响应
type CategoryDTO struct {
	ID        uint      `json:"id"`
	Key       string    `json:"key"`
	Label     string    `json:"label"`
	Icon      string    `json:"icon"`
	Order     int       `json:"order"`
	Scope     string    `json:"scope"` // 可见性级别：system | org | team | user | public
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// SettingItemDTO 配置项（用于批量设置）
type SettingItemDTO struct {
	Key   string `json:"key"`
	Value any    `json:"value"`
}

// ResetResultDTO 配置重置结果
type ResetResultDTO struct {
	Key     string `json:"key"`
	Message string `json:"message"`
}
