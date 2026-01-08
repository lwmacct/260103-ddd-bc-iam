// Package handler 实现 Settings 模块的 HTTP 处理器层。
//
// # Overview
//
// 本包是 Settings Bounded Context 的 HTTP 适配器层，负责：
//   - 将 HTTP 请求绑定到 Application 层 UseCase
//   - 请求参数验证和 DTO 转换
//   - 响应格式化和错误处理
//   - 通过 Fx 依赖注入聚合所有 Handler
//
// 主要 Handler：
//   - 用户配置：[handler.SetHandler]、[handler.BatchSetHandler]、[handler.GetHandler]、[handler.ListHandler]、[handler.ResetHandler]、[handler.ResetAllHandler]
//   - 组织配置：[orgsettings.SetHandler]、[orgsettings.GetHandler]、[orgsettings.ListHandler]、[orgsettings.ResetHandler]
//   - 团队配置：[teamsettings.SetHandler]、[teamsettings.GetHandler]、[teamsettings.ListHandler]、[teamsettings.ResetHandler]
//
// # Usage
//
// Handler 通过 Fx 容器自动注册，无需手动实例化：
//
//	fx.New(
//	    settings.Module(),
//	    // ... 其他模块
//	)
//
// 每个 Handler 对应一个或多个 HTTP 端点，路由定义见 [routes](github.com/lwmacct/260103-ddd-iam-bc/pkg/modules/settings/adapters/gin/routes) 包。
//
// # Thread Safety
//
// 所有 Handler 都是无状态的，仅依赖注入的 UseCase（通过 Fx 管理）。
// Handler 本身是并发安全的，可以安全地在多个 goroutine 中共享。
//
// # 依赖关系
//
// 本包依赖 Application 层的 UseCase（见 [app](github.com/lwmacct/260103-ddd-iam-bc/pkg/modules/settings/app)），
// 不直接访问 Domain 或 Infrastructure 层。
package handler
