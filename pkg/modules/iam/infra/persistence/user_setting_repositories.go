package persistence

import (
	"github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/domain/usersetting"
	"gorm.io/gorm"
)

// UserSettingRepositories 聚合用户设置读写仓储，便于同时注入 Command/Query
type UserSettingRepositories struct {
	Command usersetting.CommandRepository
	Query   usersetting.QueryRepository
}

// NewUserSettingRepositories 初始化用户设置仓储聚合
func NewUserSettingRepositories(db *gorm.DB) UserSettingRepositories {
	return UserSettingRepositories{
		Command: NewUserSettingCommandRepository(db),
		Query:   NewUserSettingQueryRepository(db),
	}
}
