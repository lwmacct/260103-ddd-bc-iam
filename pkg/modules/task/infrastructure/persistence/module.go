package persistence

import "go.uber.org/fx"

// RepositoryModule 提供 task 领域的仓储实现。
var RepositoryModule = fx.Module("task.repository",
	fx.Provide(NewTaskRepositories),
)
