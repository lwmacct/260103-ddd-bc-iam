package persistence

import (
	"gorm.io/gorm"

	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/crm/domain/opportunity"
)

// OpportunityRepositories 商机仓储聚合。
type OpportunityRepositories struct {
	Command opportunity.CommandRepository
	Query   opportunity.QueryRepository
}

// NewOpportunityRepositories 创建商机仓储聚合。
func NewOpportunityRepositories(db *gorm.DB) OpportunityRepositories {
	return OpportunityRepositories{
		Command: NewOpportunityCommandRepository(db),
		Query:   NewOpportunityQueryRepository(db),
	}
}
