// Package handler 实现 IAM 模块的 HTTP 处理器层。
//
// # Overview
//
// 本包是 IAM Bounded Context 的 HTTP 适配器层，负责：
//   - 将 HTTP 请求绑定到 Application 层 UseCase
//   - 请求参数验证和 DTO 转换
//   - 响应格式化和错误处理
//   - 通过 Fx 依赖注入聚合所有 Handler
//
// 主要 Handler：
//   - 认证相关：[auth.LoginHandler]、[auth.RegisterHandler]、[auth.LogoutHandler]
//   - 用户管理：[userprofile.UserProfileHandler]、[adminuser.AdminUserHandler]、[userorg.UserOrgHandler]
//   - 角色权限：[role.RoleHandler]
//   - 组织管理：[org.OrgHandler]、[orgmember.OrgMemberHandler]、[team.TeamHandler]、[teammember.TeamMemberHandler]
//   - PAT 令牌：[pat.PATHandler]
//   - 双因素认证：[twofa.TwoFAHandler]
//   - 审计日志：[audit.AuditHandler]
//   - 验证码：[captcha.CaptchaHandler]
//
// # Usage
//
// Handler 通过 Fx 容器自动注册，无需手动实例化：
//
//	fx.New(
//	    iam.Module(),
//	    // ... 其他模块
//	)
//
// 每个 Handler 对应一个或多个 HTTP 端点，路由定义见 [routes](github.com/lwmacct/260103-ddd-iam-bc/pkg/modules/iam/adapters/gin/routes) 包。
//
// # Thread Safety
//
// 所有 Handler 都是无状态的，仅依赖注入的 UseCase（通过 Fx 管理）。
// Handler 本身是并发安全的，可以安全地在多个 goroutine 中共享。
//
// # 依赖关系
//
// 本包依赖 Application 层的 UseCase（见 [app](github.com/lwmacct/260103-ddd-iam-bc/pkg/modules/iam/app)），
// 不直接访问 Domain 或 Infrastructure 层。
package handler
