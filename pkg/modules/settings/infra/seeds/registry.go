package seeds

import "github.com/lwmacct/260103-ddd-shared/pkg/platform/db"

// DefaultSeeders 返回 Settings 模块的所有种子数据。
//
// 执行顺序：
//  1. UserSettingSeeder - 用户配置值（依赖用户已存在）
func DefaultSeeders() []db.Seeder {
	return []db.Seeder{
		&UserSettingSeeder{},
	}
}
