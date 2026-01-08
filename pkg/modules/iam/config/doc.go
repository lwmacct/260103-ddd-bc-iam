// Package config 提供 IAM 模块的配置定义。
//
// # Overview
//
// 本包定义 IAM 模块所需的所有配置项，实现模块配置自治：
//   - [JWT]: 认证令牌配置（密钥、过期时间、签发者）
//   - [Auth]: 认证策略配置（密码强度、双因素认证）
//   - [RedisCache]: 缓存键前缀配置
//
// # Usage
//
// 配置通过依赖注入由 internal/container 提供：
//
//	// 配置从环境变量加载
//	cfg := &config.Config{
//	    JWT: &config.JWT{
//	        Secret:     os.Getenv("JWT_SECRET"),
//	        ExpiresIn:  24 * time.Hour,
//	        Issuer:     "iam-service",
//	    },
//	    // ...
//	}
//
//	// 通过 Fx 注入到需要的服务
//	fx.Provide(func(jwt *config.JWT) *auth.Service {
//	    return auth.NewService(jwt)
//	})
//
// # Thread Safety
//
// 配置结构体在应用启动时初始化，之后只读。
// 所有配置字段都是并发安全的（ immutable 或基本类型）。
//
// # 依赖关系
//
// 配置通过依赖注入由 internal/container 提供，IAM 模块不直接依赖 internal/config。
package config
