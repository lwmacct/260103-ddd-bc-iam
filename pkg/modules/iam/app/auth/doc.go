// Package auth 实现认证相关的应用层用例。
//
// # Overview
//
// 本包提供认证相关的 Command Handler（认证操作均为写操作）：
//   - [command.LoginHandler]: 用户登录（返回 JWT Token）
//   - [command.RegisterHandler]: 用户注册
//   - [command.RefreshTokenHandler]: 刷新访问令牌
//
// # DTO
//
// 请求 DTO：
//   - [LoginDTO]: 登录请求（用户名/密码）
//   - [RegisterDTO]: 注册请求
//   - [RefreshTokenDTO]: 刷新令牌请求
//
// 响应 DTO：
//   - [LoginResponse]: 登录成功响应（含 Token 和过期时间）
//   - [RegisterResponse]: 注册成功响应
//
// # Usage
//
//	// 用户登录
//	loginCmd := auth.LoginCommand{
//	    Username: "alice",
//	    Password: "password123",
//	}
//	result, err := loginHandler.Handle(ctx, loginCmd)
//
//	// 用户注册
//	registerCmd := auth.RegisterCommand{
//	    Username: "bob",
//	    Email:    "bob@example.com",
//	    Password: "securepassword",
//	}
//	result, err := registerHandler.Handle(ctx, registerCmd)
//
// 认证流程：
//  1. 用户提交凭据（用户名 + 密码）
//  2. 验证凭据有效性（密码哈希比对）
//  3. 检查用户状态（是否激活、是否禁用）
//  4. 生成 JWT Token 并返回
//
// # Thread Safety
//
// 所有 Handler 都是无状态的，仅依赖注入的 Repository（通过 Fx 管理）。
// Handler 本身是并发安全的，可以安全地在多个 goroutine 中调用。
//
// # 安全特性
//
//   - 密码通过 BCrypt 哈希存储
//   - JWT Token 有过期时间
//   - 支持 Token 刷新机制
//
// # 依赖关系
//
// 本包依赖 Domain 层的 [domain/auth.Service] 和 [domain/user.QueryRepository]。
package auth
