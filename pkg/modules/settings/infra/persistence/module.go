package persistence

import "go.uber.org/fx"

// RepositoryModule 仓储 Fx 模块
var RepositoryModule = fx.Module("user_settings.persistence",
	fx.Provide(NewRepositories),
)
