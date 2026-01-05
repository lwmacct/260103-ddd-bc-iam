// Package usersettings 提供 User Settings Bounded Context。
//
// # Overview
//
// 本模块作为独立 BC 管理用户配置覆盖功能：
//   - 存储用户自定义配置值（覆盖系统默认值）
//   - 依赖 Settings BC 获取配置 Schema 进行校验
//   - 提供 CRUD + 批量操作 + 层级查询 API
//
// # 架构
//
//	pkg/modules/settings/
//	├── domain/       # 领域层：实体、仓储接口
//	├── app/          # 应用层：用例处理器、DTO
//	├── infra/        # 基础设施层：持久化、缓存
//	└── adapters/     # 适配器层：HTTP Handler、路由
//
// # 依赖
//
// 本模块依赖以下外部 BC：
//   - Settings BC: 配置定义和校验
//
// # 使用方式
//
//	fx.New(
//	    // ... 其他模块
//	    usersettings.Module(),
//	    // ...
//	)
package settings

import (
	"go.uber.org/fx"

	"github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/settings/adapters/gin/handler"
	"github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/settings/app"
	"github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/settings/infra/cache"
	"github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/settings/infra/persistence"
)

// Module 返回 User Settings BC 的 Fx 模块
func Module() fx.Option {
	return fx.Module("settings",
		cache.CacheModule,
		persistence.RepositoryModule,
		app.UseCaseModule,
		handler.HandlerModule,
	)
}
