package persistence

import (
	"github.com/lwmacct/260103-ddd-iam-bc/pkg/modules/settings/domain/org"
	"gorm.io/gorm"
)

// OrgRepositories 组织配置仓储聚合
type OrgRepositories struct {
	Command org.CommandRepository
	Query   org.QueryRepository
}

// NewOrgRepositories 创建组织配置仓储聚合实例
func NewOrgRepositories(db *gorm.DB) OrgRepositories {
	return OrgRepositories{
		Command: NewOrgCommandRepository(db),
		Query:   NewOrgQueryRepository(db),
	}
}
