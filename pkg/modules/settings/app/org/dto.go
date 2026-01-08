package org

import setting "github.com/lwmacct/260103-ddd-bc-settings/pkg/modules/settings/app/setting"

// SettingsItemDTO 组织配置项响应（复用上游扁平结构类型）
type SettingsItemDTO = setting.SettingsItemDTO

// ResetResultDTO 配置重置结果
type ResetResultDTO struct {
	Key     string `json:"key"`
	Message string `json:"message"`
}
