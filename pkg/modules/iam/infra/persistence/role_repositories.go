package persistence

import (
	"github.com/lwmacct/260103-ddd-iam-bc/pkg/modules/iam/domain/role"
	"gorm.io/gorm"
)

// RoleRepositories 聚合角色读写仓储
type RoleRepositories struct {
	Command role.CommandRepository
	Query   role.QueryRepository
}

// NewRoleRepositories 创建角色仓储聚合实例
func NewRoleRepositories(db *gorm.DB) RoleRepositories {
	return RoleRepositories{
		Command: NewRoleCommandRepository(db),
		Query:   NewRoleQueryRepository(db),
	}
}
