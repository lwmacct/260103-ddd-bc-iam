// Package pat 实现个人访问令牌（Personal Access Token）的应用层用例。
//
// 本包提供 CQRS 模式的 Command 和 Query Handler：
//
// # Command（写操作）
//
//   - [command.CreateHandler]: 创建访问令牌
//   - [command.DeleteHandler]: 删除访问令牌
//   - [command.EnableHandler]: 启用访问令牌
//   - [command.DisableHandler]: 禁用访问令牌
//
// # Query（读操作）
//
//   - [query.GetHandler]: 获取令牌详情
//   - [query.ListHandler]: 令牌列表查询
//
// # DTO 与映射
//
// 请求 DTO：
//   - [CreateDTO]: 创建令牌请求（含权限范围、过期时间、IP 白名单）
//
// 响应 DTO：
//   - [TokenDTO]: 令牌信息响应（脱敏）
//   - [CreateResultDTO]: 创建成功响应（含完整令牌，仅返回一次）
//
// 映射函数：
//   - [ToTokenDTO]: PersonalAccessToken -> TokenDTO
//
// 安全特性：
//   - 令牌创建后仅返回一次完整值
//   - 支持权限范围限制
//   - 支持 IP 白名单
//   - 支持过期时间设置
//
// 依赖注入：所有 Handler 通过 [bootstrap.Container] 注册。
package pat
