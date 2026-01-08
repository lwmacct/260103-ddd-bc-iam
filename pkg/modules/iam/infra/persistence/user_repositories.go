package persistence

import (
	"github.com/lwmacct/260103-ddd-iam-bc/pkg/modules/iam/domain/user"
	"gorm.io/gorm"
)

// UserRepositories 聚合用户读写仓储，方便在容器中同时获取
type UserRepositories struct {
	Command user.CommandRepository
	Query   user.QueryRepository
}

// NewUserRepositories 创建聚合实例，同时初始化 Command/Query 仓储
func NewUserRepositories(db *gorm.DB) UserRepositories {
	return UserRepositories{
		Command: NewUserCommandRepository(db),
		Query:   NewUserQueryRepository(db),
	}
}
