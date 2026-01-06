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

// IsConfigurableByTeam 检查设置定义是否支持团队级别配置（configurable_at >= team）
func IsConfigurableByTeam(def *settingdomain.Setting) bool {
	return def.IsConfigurableAtScope(settingdomain.ScopeLevelTeam)
}

// IsVisibleToTeam 检查设置是否对团队管理员可见
//
// 团队管理员可以看到满足以下任一条件的设置：
// 1. visible_at >= team (team 或 user 级别可见)
// 2. configurable_at >= team (支持团队级别配置的设置，即使 visible_at=user)
//
// 规则 2 的原因：对于 "Team 默认值" 场景（如 visible_at=user, configurable_at=team），
// 团队管理员需要能看到设置才能为成员配置默认值。
func IsVisibleToTeam(def *settingdomain.Setting) bool {
	// 如果设置支持团队级别配置，则应该对团队管理员可见
	if def.IsConfigurableAtScope(settingdomain.ScopeLevelTeam) {
		return true
	}
	// 否则检查可见性
	return def.IsVisibleAtScope(settingdomain.ScopeLevelTeam)
}
