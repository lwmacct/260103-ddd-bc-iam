// Package routes 定义 IAM 模块的 HTTP 路由。
//
// # Overview
//
// 本包负责 IAM Bounded Context 的 HTTP 路由定义和组织：
//   - 定义所有 API 端点路径、HTTP 方法、参数绑定
//   - 集成 Swagger/OpenAPI 注解（用于生成 API 文档）
//   - 定义 OperationID（URN 格式，同时用于权限标识）
//   - 按业务模块分组管理路由（auth、user、role、org 等）
//   - 通过 Fx 依赖注入导出路由集合
//
// 路由分组：
//   - 认证路由：[auth.RegisterRoutes] - 登录、注册、登出
//   - 用户路由：[user.RegisterRoutes] - 用户 Profile、用户管理、用户组织
//   - 角色路由：[role.RegisterRoutes] - 角色 CRUD、权限分配
//   - 组织路由：[org.RegisterRoutes] - 组织、团队、成员管理
//   - PAT 路由：[pat.RegisterRoutes] - 个人访问令牌管理
//   - 双因素认证路由：[twofa.RegisterRoutes] - TOTP 设置和验证
//   - 审计路由：[audit.RegisterRoutes] - 审计日志查询
//   - 验证码路由：[captcha.RegisterRoutes] - CAPTCHA 生成和验证
//
// # Usage
//
// 路由通过 Fx 容器自动注册，使用 `fx.ResultTags` 导出：
//
//	fx.New(
//	    iam.Module(),
//	    // ... 其他模块
//	)
//
// 所有路由遵循以下规范：
//   - 路径以 `/api/` 开头，使用复数名词（如 `/api/users`）
//   - 路径参数使用 OpenAPI 风格 `{id}`（而非 `:id`）
//   - OperationID 使用 URN 格式 `{scope}:{resource}:{action}`
//   - Tags 使用 kebab-case 格式（如 `user-profile`、`org-member`）
//
// # Thread Safety
//
// 路由定义是纯元数据（无状态），编译时确定。
// 路由注册过程由 Fx 管理（应用启动时执行），完成后只读。
// 本包是并发安全的。
//
// # 依赖关系
//
// 本包依赖 Handler 层（见 [handler](github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/adapters/gin/handler)），
// 被 Router 容器（见 [internal/container/router](github.com/lwmacct/260103-ddd-bc-iam/internal/container/router)）消费。
package routes
