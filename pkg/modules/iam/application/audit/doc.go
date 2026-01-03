// Package auditlog 实现审计日志的应用层用例。
//
// 本包仅提供 Query Handler（审计日志只读）：
//
// # Query（读操作）
//
//   - [query.GetLogHandler]: 获取审计日志详情
//   - [query.ListLogsHandler]: 审计日志列表查询（支持多维度筛选）
//
// # DTO 与映射
//
// 查询参数：
//   - [ListLogsQuery]: 日志查询参数（用户、操作类型、时间范围等）
//
// 响应 DTO：
//   - [AuditLogResponse]: 审计日志响应
//   - [AuditLogListResponse]: 日志列表响应（含分页）
//
// 映射函数：
//   - [ToAuditLogResponse]: AuditLog -> AuditLogResponse
//
// 审计日志特性：
//   - 记录所有敏感操作（创建、更新、删除）
//   - 包含操作者、IP 地址、时间戳
//   - 包含操作前后的数据快照（可选）
//   - 不可修改、不可删除（只读）
//
// 日志写入：
// 审计日志的写入由中间件自动完成，不在本包处理。
// 参见 [adapters/http/middleware.AuditMiddleware]。
//
// 依赖注入：所有 Handler 通过 [bootstrap.Container] 注册。
package audit
