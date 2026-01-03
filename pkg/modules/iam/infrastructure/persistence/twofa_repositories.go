package persistence

import (
	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/iam/domain/twofa"
	"gorm.io/gorm"
)

// TwoFARepositories 聚合两步验证读写仓储
type TwoFARepositories struct {
	Command twofa.CommandRepository
	Query   twofa.QueryRepository
}

// NewTwoFARepositories 创建两步验证仓储聚合实例
func NewTwoFARepositories(db *gorm.DB) TwoFARepositories {
	return TwoFARepositories{
		Command: NewTwoFACommandRepository(db),
		Query:   NewTwoFAQueryRepository(db),
	}
}
