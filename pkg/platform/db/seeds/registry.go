package seeds

import "github.com/lwmacct/260103-ddd-bc-iam/pkg/platform/db"

// DefaultSeeders returns the default ordered seeders that bootstrap the system.
// Keep RBAC first because it provisions permissions/roles required by other seeders.
// OrganizationSeeder runs last because it depends on UserSeeder (admin user).
func DefaultSeeders() []db.Seeder {
	return []db.Seeder{
		&RBACSeeder{},
		&UserSeeder{},
		&OrganizationSeeder{},
	}
}
