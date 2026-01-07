// Package routes 定义 IAM 模块的所有 HTTP 路由。
//
// 本包遵循 DDD 架构原则，路由按业务模块组织，通过 fx Module 自动装配。
// 中间件由应用层（internal/container）注入。
//
// # 路由组织
//
//   - [Auth]: 认证模块（注册、登录、令牌刷新）
//   - [Captcha]: 验证码模块
//   - [TwoFA]: 双因素认证模块
//   - [User]: 用户模块（用户资料、用户管理、用户组织视图）
//   - [PAT]: 个人访问令牌模块
//   - [Role]: 角色管理模块
//   - [Audit]: 审计日志模块
//   - [Org]: 组织模块（组织管理、组织成员、团队、团队成员）
//
// # Fx 模块装配
//
//	iam.Module() 包含 routes.RoutesModule，通过 fx 自动注入 Handlers 聚合
package routes
