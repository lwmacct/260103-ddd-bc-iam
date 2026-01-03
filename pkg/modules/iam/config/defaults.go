package config

import "time"

// DefaultConfig 返回 IAM 模块的默认配置。
func DefaultConfig() Config {
	return Config{
		JWT: JWT{
			Secret:             "change-me-in-production",
			AccessTokenExpiry:  15 * time.Minute,
			RefreshTokenExpiry: 7 * 24 * time.Hour,
		},
		Auth: Auth{
			DevSecret:       "dev-secret-change-me",
			TwoFAIssuer:     "Go-DDD-Package-Lib",
			CaptchaRequired: true,
		},
		RedisCache: RedisCache{
			KeyPrefix: "app:",
		},
	}
}
