// Package config 提供 IAM 模块的配置定义。
//
// # IAM 模块配置独立管理
//
// 本包定义 IAM 模块所需的所有配置项，实现模块配置自治：
//   - JWT：认证令牌配置
//   - Auth：认证策略配置
//   - RedisCache：缓存键前缀配置
//
// 配置通过依赖注入由 internal/container 提供，IAM 模块不直接依赖 internal/config。
package config
