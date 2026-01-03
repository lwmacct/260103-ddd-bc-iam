package task

import (
	"go.uber.org/fx"

	taskapplication "github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/task/application"
	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/task/infrastructure/persistence"
	taskhandler "github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/task/transport/gin/handler"
)

// Module 返回 task BC 的完整 Fx 模块。
//
// 本函数聚合四层架构，对外提供统一的依赖注入入口：
//   - persistence.RepositoryModule: 仓储实现
//   - application.UseCaseModule: 用例处理器
//   - fx.Provide(NewAllHandlers): HTTP Handler
func Module() fx.Option {
	return fx.Module("task",
		persistence.RepositoryModule,
		taskapplication.UseCaseModule,
		fx.Provide(taskapplication.NewTaskUseCases),
		fx.Provide(taskhandler.NewAllHandlers),
	)
}
