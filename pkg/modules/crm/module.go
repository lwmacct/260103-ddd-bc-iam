// Package crm 提供客户关系管理业务模块（CRM Bounded Context）。
//
// 本模块包含：
//   - Domain: 公司、联系人、线索、商机
//   - Application: 用例处理器
//   - Infrastructure: 仓储实现
//   - Transport: HTTP 适配器
//
// 模块化设计：
//   - persistence.RepositoryModule: 仓储层（无缓存装饰器）
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
//	    crm.Module(),       // 本模块
//	)
package crm

import (
	"go.uber.org/fx"

	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/crm/application"
	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/crm/infrastructure/persistence"
	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/crm/transport/gin/handler"
)

// Module 返回 CRM 模块的完整 Fx 配置。
//
// 依赖注入顺序：
//  1. Platform 层 (DB, Redis, Config) - 由外部提供
//  2. RepositoryModule (本模块) - 仓储层
//  3. UseCaseModule (本模块) - 用例层
//  4. HandlerModule (本模块) - HTTP 处理器层
//
// 注意：
//   - CRM 模块不使用缓存装饰器（业务特性决定）
//   - 路由注册由顶层容器处理，不在模块内部
func Module() fx.Option {
	return fx.Module("crm",
		// 子模块
		persistence.RepositoryModule,
		application.UseCaseModule,
		handler.HandlerModule,
	)
}
