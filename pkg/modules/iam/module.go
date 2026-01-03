// Package iam 提供身份与访问管理业务模块（IAM Bounded Context）。
//
// 本模块包含：
//   - Domain: 用户、角色、权限、认证、双因素认证、PAT
//   - Application: 用例处理器
//   - Infrastructure: JWT、TOTP、仓储实现
//   - Transport: HTTP 适配器、中间件
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
//	    iam.Module(),       // 本模块
//	    crm.Module(),
//	)
package iam

import (
	"go.uber.org/fx"

	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/iam/application"
	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/iam/infrastructure/persistence"
	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/iam/transport/gin/handler"
)

// Module 返回 IAM 模块的完整 Fx 配置。
//
// 依赖注入顺序：
//  1. Platform 层 (DB, Redis, Config, EventBus) - 由外部提供
//  2. App 模块 (AuditUseCases, OrganizationUseCases) - 跨模块依赖
//  3. CacheModule (Platform 层) - 由外部提供
//  4. RepositoryModule (本模块) - 仓储层
//  5. UseCaseModule (本模块) - 用例层
//  6. HandlerModule (本模块) - HTTP 处理器层
//
// 注意：
//   - 认证中间件（Auth, RBAC）在 transport/gin/middleware 包中
//   - 路由注册由顶层容器处理，不在模块内部
func Module() fx.Option {
	return fx.Module("iam",
		// 子模块
		persistence.RepositoryModule,
		application.UseCaseModule,
		handler.HandlerModule,
	)
}
