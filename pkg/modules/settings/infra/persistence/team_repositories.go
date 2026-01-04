package persistence

import (
	"github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/settings/domain/team"
	"gorm.io/gorm"
)

// TeamRepositories 团队配置仓储聚合
type TeamRepositories struct {
	Command team.CommandRepository
	Query   team.QueryRepository
}

// NewTeamRepositories 创建团队配置仓储聚合实例
func NewTeamRepositories(db *gorm.DB) TeamRepositories {
	return TeamRepositories{
		Command: NewTeamCommandRepository(db),
		Query:   NewTeamQueryRepository(db),
	}
}
