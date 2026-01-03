// Package app 提供核心治理业务模块（App Bounded Context）。
//
// # Overview
//
// 本模块是核心业务的聚合根，包含以下子域：
//   - 组织管理（Organization, Team, Member）
//   - 设置管理（Setting, UserSetting, Category）
//   - 审计日志（AuditLog）
//   - 任务管理（Task）
//   - 统计查询（Stats）
//   - 验证码（Captcha）
//   - 健康检查（Health）
//
// # Architecture
//
// 本模块采用 DDD 四层架构 + CQRS 模式：
//   - domain: 领域模型、仓储接口、领域错误
//   - application: 用例处理器（Command/Query Handler）、DTO
//   - infrastructure: 仓储实现、缓存、事件处理器
//   - transport: HTTP Handler、路由
//
// # 模块化设计
//
// 使用 Uber Fx 实现依赖注入，通过 [Module] 函数导出完整配置：
//
//	fx.New(
//	    fx.Supply(cfg),
//	    platform.Module(),
//	    app.Module(),       // 本模块
//	    iam.Module(),
//	    crm.Module(),
//	)
//
// # 子模块
//
//   - persistence.RepositoryModule: 仓储层（包括缓存装饰器）
//   - application.UseCaseModule: 用例层（业务逻辑编排）
//   - handler.HandlerModule: HTTP 处理器层
//
// # Thread Safety
//
// 所有导出的 Handler 和 UseCase 都是并发安全的。
//
// # 依赖关系
//
//   - 依赖 Platform 层：DB、Redis、Config、EventBus
//   - 依赖 IAM 模块：UserRepositories (UserSetting 用例)、UserUseCases (AdminUser Handler)
//   - 被 IAM/CRM 模块依赖：提供审计、组织、设置等基础能力
package app
