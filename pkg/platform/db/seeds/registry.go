package seeds

import "github.com/lwmacct/260101-go-pkg-ddd/pkg/platform/db"

// DefaultSeeders returns the default ordered seeders that bootstrap the system.
// Keep RBAC first because it provisions permissions/roles required by other seeders.
// SettingCategorySeeder must run before SettingSeeder to ensure categories exist.
// OrganizationSeeder runs last because it depends on UserSeeder (admin user).
func DefaultSeeders() []db.Seeder {
	return []db.Seeder{
		&RBACSeeder{},
		&UserSeeder{},
		&SettingCategorySeeder{},
		&SettingSeeder{},
		&OrganizationSeeder{},
	}
}
