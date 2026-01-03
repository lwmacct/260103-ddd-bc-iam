// Package app 提供核心治理业务模块（App Bounded Context）。
//
// 本模块包含：
//   - Domain: 组织、设置、统计、任务、审计
//   - Application: 用例处理器
//   - Infrastructure: 仓储实现
//   - Transport: HTTP 适配器
//
// 模块化设计：
//   - persistence.RepositoryModule: 仓储层（包括缓存装饰器）
//   - application.UseCaseModule: 用例层（业务逻辑编排）
//   - handler.HandlerModule: HTTP 处理器层
//
// 使用方式：
//
//	fx.New(
//	    fx.Supply(cfg),
//	    platform.Module(),
//	    app.Module(),
//	    iam.Module(),
//	    crm.Module(),
//	)
package app

import (
	"go.uber.org/fx"

	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/app/application"
	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/app/infrastructure/persistence"
	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/app/transport/gin/handler"
)

// Module 返回 App 模块的完整 Fx 配置。
//
// 依赖注入顺序：
//  1. Platform 层 (DB, Redis, Config) - 由外部提供
//  2. CacheModule (Platform 层) - 由外部提供
//  3. RepositoryModule (本模块) - 仓储层
//  4. UseCaseModule (本模块) - 用例层
//  5. HandlerModule (本模块) - HTTP 处理器层
//
// 注意：路由注册由顶层容器处理，不在模块内部。
func Module() fx.Option {
	return fx.Module("app",
		// 子模块
		persistence.RepositoryModule,
		application.UseCaseModule,
		handler.HandlerModule,
	)
}
