package persistence

import (
	"gorm.io/gorm"

	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/crm/domain/contact"
)

// ContactRepositories 联系人仓储聚合。
type ContactRepositories struct {
	Command contact.CommandRepository
	Query   contact.QueryRepository
}

// NewContactRepositories 创建联系人仓储聚合。
func NewContactRepositories(db *gorm.DB) ContactRepositories {
	return ContactRepositories{
		Command: NewContactCommandRepository(db),
		Query:   NewContactQueryRepository(db),
	}
}
