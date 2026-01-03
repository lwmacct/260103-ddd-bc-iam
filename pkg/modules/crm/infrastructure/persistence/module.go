package persistence

import (
	"go.uber.org/fx"
)

// RepositoryModule 提供 CRM 模块的所有仓储实现。
//
// CRM 模块不使用缓存装饰器（业务特性决定）。
var RepositoryModule = fx.Module("crm.repository",
	fx.Provide(
		// 直接使用 persistence 构造函数
		NewContactRepositories,
		NewCompanyRepositories,
		NewLeadRepositories,
		NewOpportunityRepositories,
	),
)
