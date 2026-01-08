package team

import setting "github.com/lwmacct/260103-ddd-settings-bc/pkg/modules/settings/app/setting"

// SettingsItemDTO 团队配置项响应（复用上游扁平结构类型）
//
// 三级继承优先级：团队 > 组织 > 系统默认值
// InheritedFrom 字段表示当前生效值的来源："team"/"org"/"system"
type SettingsItemDTO = setting.SettingsItemDTO

// ResetResultDTO 配置重置结果
type ResetResultDTO struct {
	Key     string `json:"key"`
	Message string `json:"message"`
}
