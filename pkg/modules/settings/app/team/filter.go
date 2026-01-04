package team

import (
	settingdomain "github.com/lwmacct/260103-ddd-bc-settings/pkg/modules/settings/domain/setting"
)

// FilterByConfigurableAt 过滤出在指定级别可配置的设置
func FilterByConfigurableAt(defs []*settingdomain.Setting, level settingdomain.ScopeLevel) []*settingdomain.Setting {
	result := make([]*settingdomain.Setting, 0, len(defs))
	for _, def := range defs {
		if def.IsConfigurableAtScope(level) {
			result = append(result, def)
		}
	}
	return result
}

// IsConfigurableByTeam 检查设置是否允许 Team 配置
func IsConfigurableByTeam(def *settingdomain.Setting) bool {
	return def.IsConfigurableAtScope(settingdomain.ScopeLevelTeam)
}

// IsVisibleToTeam 检查设置是否对 Team 可见
//
// Team 可以看到满足以下任一条件的设置：
// 1. visible_at >= team (team 或 user 级别可见)
// 2. configurable_at >= team (Team 可配置的设置，即使 visible_at=user)
//
// 规则 2 的原因：对于 "Team 默认值" 场景（如 visible_at=user, configurable_at=team），
// Team 需要能看到设置才能为用户配置默认值。
func IsVisibleToTeam(def *settingdomain.Setting) bool {
	// 如果设置允许 Team 配置，则 Team 应该能看到
	if def.IsConfigurableAtScope(settingdomain.ScopeLevelTeam) {
		return true
	}
	// 否则检查可见性
	return def.IsVisibleAtScope(settingdomain.ScopeLevelTeam)
}
