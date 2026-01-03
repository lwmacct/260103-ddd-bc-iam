package role

import "context"

// ============================================================================
// Command Repository
// ============================================================================

// CommandRepository 角色命令仓储接口（写操作）
type CommandRepository interface {
	// Create creates a new role
	Create(ctx context.Context, role *Role) error

	// Update updates a role
	Update(ctx context.Context, role *Role) error

	// Delete deletes a role (soft delete)
	Delete(ctx context.Context, id uint) error

	// SetPermissions sets permissions for a role (replaces all existing permissions)
	SetPermissions(ctx context.Context, roleID uint, permissions []Permission) error
}

// ============================================================================
// Query Repository
// ============================================================================

// QueryRepository 角色查询仓储接口（读操作）
type QueryRepository interface {
	// FindByID finds a role by ID
	FindByID(ctx context.Context, id uint) (*Role, error)

	// FindByName finds a role by name
	FindByName(ctx context.Context, name string) (*Role, error)

	// FindByIDWithPermissions finds a role by ID with permissions
	FindByIDWithPermissions(ctx context.Context, id uint) (*Role, error)

	// List returns all roles with pagination
	List(ctx context.Context, page, limit int) ([]Role, int64, error)

	// GetPermissions retrieves all permissions for a role
	GetPermissions(ctx context.Context, roleID uint) ([]Permission, error)

	// Exists checks if a role exists by ID
	Exists(ctx context.Context, id uint) (bool, error)

	// ExistsByName checks if a role exists by name
	ExistsByName(ctx context.Context, name string) (bool, error)
}
