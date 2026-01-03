package persistence

import (
	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/iam/domain/org"
	"gorm.io/gorm"
)

// TeamRepositories 聚合团队读写仓储
type TeamRepositories struct {
	Command org.TeamCommandRepository
	Query   org.TeamQueryRepository
}

// NewTeamRepositories 创建团队仓储聚合实例
func NewTeamRepositories(db *gorm.DB) TeamRepositories {
	return TeamRepositories{
		Command: NewTeamCommandRepository(db),
		Query:   NewTeamQueryRepository(db),
	}
}
