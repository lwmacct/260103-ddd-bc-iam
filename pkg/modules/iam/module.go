// Package iam 提供身份与访问管理业务模块（IAM Bounded Context）。
//
// 本模块包含：
//   - Domain: 用户、角色、权限、认证、双因素认证、PAT
//   - App: 用例处理器
//   - Infrastructure: JWT、TOTP、仓储实现、缓存服务、种子数据
//   - Adapters: HTTP 适配器、中间件、测试工具
//
// 模块化设计：
//   - persistence.RepositoryModule: 仓储层（包括缓存装饰器）
//   - cache.CacheModule: IAM 专属缓存服务
//   - app.UseCaseModule: 用例层（业务逻辑编排）
//   - handler.HandlerModule: HTTP 处理器层
//   - routes.RoutesModule: 路由层（通过 fx 自动装配）
//
// 注意：User Settings 功能已迁移到独立的 User Settings BC。
//
// 使用方式：
//
//	fx.New(
//	    fx.Supply(cfg),
//	    platform.Module(),
//	    iam.Module(),           // 本模块（完全自治）
//	    usersettings.Module(),  // User Settings BC（独立模块）
//	    crm.Module(),
//	)
package iam

import (
	"go.uber.org/fx"

	"github.com/lwmacct/260103-ddd-iam-bc/pkg/modules/iam/adapters/gin/handler"
	"github.com/lwmacct/260103-ddd-iam-bc/pkg/modules/iam/adapters/gin/routes"
	"github.com/lwmacct/260103-ddd-iam-bc/pkg/modules/iam/app"
	"github.com/lwmacct/260103-ddd-iam-bc/pkg/modules/iam/infra/cache"
	"github.com/lwmacct/260103-ddd-iam-bc/pkg/modules/iam/infra/persistence"
)

// Module 返回 IAM 模块的完整 Fx 配置。
//
// 依赖注入顺序：
//  1. Platform 层 (DB, Redis, Config, EventBus) - 由外部提供
//  2. App 模块 (AuditUseCases, OrganizationUseCases) - 跨模块依赖
//  3. CacheModule (本模块) - IAM 专属缓存服务
//  4. RepositoryModule (本模块) - 仓储层
//  5. UseCaseModule (本模块) - 用例层
//  6. HandlerModule (本模块) - HTTP 处理器层
//  7. RoutesModule (本模块) - 路由层（自动注入 Handlers 聚合）
//
// 注意：
//   - 认证中间件（Auth, RBAC）在 adapters/gin/middleware 包中
//   - User Settings 功能已迁移到独立的 User Settings BC
func Module() fx.Option {
	return fx.Module("iam",
		// 子模块
		cache.CacheModule,
		persistence.RepositoryModule,
		app.UseCaseModule,
		handler.HandlerModule,
		routes.RoutesModule,
	)
}
