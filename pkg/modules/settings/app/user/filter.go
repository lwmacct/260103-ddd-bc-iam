package user

import (
	settingdomain "github.com/lwmacct/260103-ddd-bc-settings/pkg/modules/settings/domain/setting"
)

// FilterByVisibleToUser 过滤出对普通用户可见的设置
//
// 可见条件：visible_at 为 user 或 public
func FilterByVisibleToUser(defs []*settingdomain.Setting) []*settingdomain.Setting {
	result := make([]*settingdomain.Setting, 0, len(defs))
	for _, def := range defs {
		if def.IsVisibleToUser() {
			result = append(result, def)
		}
	}
	return result
}
