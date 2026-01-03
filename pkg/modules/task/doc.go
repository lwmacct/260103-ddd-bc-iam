// Package task 提供团队任务管理的 Bounded Context。
//
// 本 BC 实现完整的四层 DDD 架构：
//   - [Domain]: 任务实体、状态值对象、仓储接口
//   - [Infrastructure]: GORM 数据持久化实现
//   - [Application]: CQRS 用例处理器
//   - [Transport]: HTTP Handler 和路由
//
// 多租户设计：
// 任务与 org/team 模块关联，通过 OrgID + TeamID 实现数据隔离。
//
// 使用方式：
//
//	import taskpkg "github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/task"
//
//	fx.New(
//		// ...
//		taskpkg.Module(),
//	)
package task
