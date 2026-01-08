package persistence

import (
	"github.com/lwmacct/260103-ddd-iam-bc/pkg/modules/iam/domain/org"
	"gorm.io/gorm"
)

// OrganizationRepositories 聚合组织读写仓储（含成员相关）
type OrganizationRepositories struct {
	Command           org.CommandRepository
	Query             org.QueryRepository
	MemberCommand     org.MemberCommandRepository
	MemberQuery       org.MemberQueryRepository
	TeamCommand       org.TeamCommandRepository
	TeamQuery         org.TeamQueryRepository
	TeamMemberCommand org.TeamMemberCommandRepository
	TeamMemberQuery   org.TeamMemberQueryRepository
}

// NewOrganizationRepositories 创建组织仓储聚合实例
func NewOrganizationRepositories(db *gorm.DB) OrganizationRepositories {
	return OrganizationRepositories{
		Command:           NewOrganizationCommandRepository(db),
		Query:             NewOrganizationQueryRepository(db),
		MemberCommand:     NewOrgMemberCommandRepository(db),
		MemberQuery:       NewOrgMemberQueryRepository(db),
		TeamCommand:       NewTeamCommandRepository(db),
		TeamQuery:         NewTeamQueryRepository(db),
		TeamMemberCommand: NewTeamMemberCommandRepository(db),
		TeamMemberQuery:   NewTeamMemberQueryRepository(db),
	}
}
