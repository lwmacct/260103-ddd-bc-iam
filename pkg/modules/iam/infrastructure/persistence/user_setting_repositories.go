package persistence

import (
	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/app/domain/setting"
	"gorm.io/gorm"
)

// UserSettingRepositories 聚合用户配置读写仓储
type UserSettingRepositories struct {
	Command setting.UserSettingCommandRepository
	Query   setting.UserSettingQueryRepository
}

// NewUserSettingRepositories 创建用户配置仓储聚合实例
func NewUserSettingRepositories(db *gorm.DB) UserSettingRepositories {
	return UserSettingRepositories{
		Command: NewUserSettingCommandRepository(db),
		Query:   NewUserSettingQueryRepository(db),
	}
}
