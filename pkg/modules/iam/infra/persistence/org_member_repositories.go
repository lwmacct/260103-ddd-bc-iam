package persistence

import (
	"github.com/lwmacct/260103-ddd-iam-bc/pkg/modules/iam/domain/org"
	"gorm.io/gorm"
)

// OrgMemberRepositories 聚合组织成员读写仓储
type OrgMemberRepositories struct {
	Command org.MemberCommandRepository
	Query   org.MemberQueryRepository
}

// NewOrgMemberRepositories 创建组织成员仓储聚合实例
func NewOrgMemberRepositories(db *gorm.DB) OrgMemberRepositories {
	return OrgMemberRepositories{
		Command: NewOrgMemberCommandRepository(db),
		Query:   NewOrgMemberQueryRepository(db),
	}
}
