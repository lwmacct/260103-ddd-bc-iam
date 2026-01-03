package persistence

import (
	"gorm.io/gorm"

	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/crm/domain/lead"
)

// LeadRepositories 线索仓储聚合。
type LeadRepositories struct {
	Command lead.CommandRepository
	Query   lead.QueryRepository
}

// NewLeadRepositories 创建线索仓储聚合。
func NewLeadRepositories(db *gorm.DB) LeadRepositories {
	return LeadRepositories{
		Command: NewLeadCommandRepository(db),
		Query:   NewLeadQueryRepository(db),
	}
}
