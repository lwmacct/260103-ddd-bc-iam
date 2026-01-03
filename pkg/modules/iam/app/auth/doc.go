// Package auth 实现认证相关的应用层用例。
//
// 本包仅提供 Command Handler（认证操作均为写操作）：
//
// # Command（写操作）
//
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
// 认证流程：
//  1. 用户提交凭据（用户名 + 密码）
//  2. 验证凭据有效性（密码哈希比对）
//  3. 检查用户状态（是否激活、是否禁用）
//  4. 生成 JWT Token 并返回
//
// 安全特性：
//   - 密码通过 BCrypt 哈希存储
//   - JWT Token 有过期时间
//   - 支持 Token 刷新机制
//
// 依赖：
//   - [domain/auth.Service]: 认证领域服务接口
//   - [domain/user.QueryRepository]: 用户查询仓储
//
// 依赖注入：所有 Handler 通过 [bootstrap.Container] 注册。
package auth
