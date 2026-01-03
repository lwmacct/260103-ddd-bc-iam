package persistence

import (
	"github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/user_settings/domain/userset"
	"gorm.io/gorm"
)

// Repositories 用户配置仓储聚合
type Repositories struct {
	Command userset.CommandRepository
	Query   userset.QueryRepository
}

// NewRepositories 创建仓储聚合实例
func NewRepositories(db *gorm.DB) Repositories {
	return Repositories{
		Command: NewCommandRepository(db),
		Query:   NewQueryRepository(db),
	}
}
