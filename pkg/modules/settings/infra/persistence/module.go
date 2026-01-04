package persistence

import "go.uber.org/fx"

// RepositoryModule 仓储 Fx 模块
var RepositoryModule = fx.Module("settings.persistence",
	fx.Provide(
		NewRepositories,
		NewOrgRepositories,
		NewTeamRepositories,
	),
)
