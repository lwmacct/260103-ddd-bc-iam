package persistence

import (
	"github.com/lwmacct/260103-ddd-iam-bc/pkg/modules/iam/domain/org"
	"gorm.io/gorm"
)

// TeamMemberRepositories 聚合团队成员读写仓储
type TeamMemberRepositories struct {
	Command org.TeamMemberCommandRepository
	Query   org.TeamMemberQueryRepository
}

// NewTeamMemberRepositories 创建团队成员仓储聚合实例
func NewTeamMemberRepositories(db *gorm.DB) TeamMemberRepositories {
	return TeamMemberRepositories{
		Command: NewTeamMemberCommandRepository(db),
		Query:   NewTeamMemberQueryRepository(db),
	}
}
