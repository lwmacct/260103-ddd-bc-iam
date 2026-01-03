// Package task 提供团队任务管理的 Application 用例处理器。
//
// 本包遵循 CQRS 模式，将命令（写操作）和查询（读操作）分离：
//
// 命令处理器：
//   - [CreateHandler]: 创建任务
//   - [UpdateHandler]: 更新任务
//   - [DeleteHandler]: 删除任务
//
// 查询处理器：
//   - [GetHandler]: 获取任务详情
//   - [ListHandler]: 分页获取任务列表
//
// 多租户设计：
// 所有命令和查询都需要指定 OrgID 和 TeamID 参数，
// 以确保多租户数据隔离。
package task
