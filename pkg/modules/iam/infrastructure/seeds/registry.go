package seeds

import "github.com/lwmacct/260103-ddd-shared/pkg/platform/db"

// DefaultSeeders 返回 IAM 模块的所有种子数据。
//
// 执行顺序：
//  1. RBACSeeder - 角色/权限（其他种子可能依赖）
//  2. UserSeeder - 管理员用户
//  3. OrganizationSeeder - 默认组织（依赖管理员用户）
func DefaultSeeders() []db.Seeder {
	return []db.Seeder{
		&RBACSeeder{},
		&UserSeeder{},
		&OrganizationSeeder{},
	}
}
