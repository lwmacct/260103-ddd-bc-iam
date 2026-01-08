package persistence

import (
	"github.com/lwmacct/260103-ddd-iam-bc/pkg/modules/settings/domain/user"
	"gorm.io/gorm"
)

// Repositories 用户配置仓储聚合
type Repositories struct {
	Command user.CommandRepository
	Query   user.QueryRepository
}

// NewRepositories 创建仓储聚合实例
func NewRepositories(db *gorm.DB) Repositories {
	return Repositories{
		Command: NewCommandRepository(db),
		Query:   NewQueryRepository(db),
	}
}
