package persistence

import (
	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/iam/domain/pat"
	"gorm.io/gorm"
)

// PATRepositories 聚合个人访问令牌读写仓储
type PATRepositories struct {
	Command pat.CommandRepository
	Query   pat.QueryRepository
}

// NewPATRepositories 创建 PAT 仓储聚合实例
func NewPATRepositories(db *gorm.DB) PATRepositories {
	return PATRepositories{
		Command: NewPATCommandRepository(db),
		Query:   NewPATQueryRepository(db),
	}
}
