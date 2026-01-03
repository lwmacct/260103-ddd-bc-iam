package persistence

import (
	"gorm.io/gorm"

	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/crm/domain/company"
)

// CompanyRepositories 公司仓储聚合。
type CompanyRepositories struct {
	Command company.CommandRepository
	Query   company.QueryRepository
}

// NewCompanyRepositories 创建公司仓储聚合。
func NewCompanyRepositories(db *gorm.DB) CompanyRepositories {
	return CompanyRepositories{
		Command: NewCompanyCommandRepository(db),
		Query:   NewCompanyQueryRepository(db),
	}
}
