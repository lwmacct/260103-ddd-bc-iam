// Package auth 提供认证服务的基础设施实现。
//
// 本包实现 [domain/auth.Service] 和 [domain/auth.TokenGenerator] 接口。
//
// # 核心组件
//
// JWT 管理：
//   - [JWTManager]: JWT Token 生成与验证
//   - 支持访问令牌和刷新令牌
//   - 可配置的密钥和过期时间
//
// 认证服务：
//   - [authServiceImpl]: 实现 domain/auth.Service 接口
//   - BCrypt 密码哈希（成本因子 10）
//   - 密码策略验证
//   - Token 生成与验证
//
// 令牌生成：
//   - [TokenGenerator]: 实现 domain/auth.TokenGenerator 接口
//   - 安全随机令牌生成
//   - 用于 PAT 等场景
//
// # 辅助服务
//
// 会话管理：
//   - [LoginSession]: 登录会话管理（可选）
//
// PAT 认证：
//   - [PATService]: PAT 认证服务（供中间件使用）
//   - ValidateToken/ValidateTokenWithIP: 验证 PAT
//   - DeleteAllUserTokens: 安全撤销（如密码重置）
//   - CleanupExpiredTokens: 系统维护
//   - 注意：PAT 的 CRUD 由 Application 层负责
//
// 权限缓存：
//   - [PermissionCacheService]: 用户权限缓存服务
//   - 使用 Redis 缓存用户权限
//   - 减少数据库查询
//
// # 配置项
//
//   - JWT 密钥：通过 config.JWTSecret 配置
//   - 访问令牌过期时间：通过 config.JWTExpireHours 配置
//   - 刷新令牌过期时间：通过 config.RefreshTokenExpireHours 配置
//   - BCrypt 成本因子：默认 10
//
// # 依赖
//
//   - Redis：权限缓存（[PermissionCacheService]）
//   - GORM：用户查询（验证用户存在性）
//
// # 使用示例
//
// 认证服务通过依赖注入在 bootstrap 层初始化，供 Application 层 Handler 使用。
package auth
