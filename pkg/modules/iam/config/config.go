package config

import "time"

// JWT 认证配置
type JWT struct {
	Secret             string        `koanf:"secret"`
	AccessTokenExpiry  time.Duration `koanf:"access-token-expiry"`
	RefreshTokenExpiry time.Duration `koanf:"refresh-token-expiry"`
}

// Auth 认证配置
type Auth struct {
	DevSecret       string `koanf:"dev-secret"`
	TwoFAIssuer     string `koanf:"twofa-issuer"`
	CaptchaRequired bool   `koanf:"captcha-required"`
}

// RedisCache Redis 缓存配置
type RedisCache struct {
	KeyPrefix string `koanf:"key-prefix"`
}

// Config IAM 模块配置
type Config struct {
	JWT        JWT        `koanf:"jwt"`
	Auth       Auth       `koanf:"auth"`
	RedisCache RedisCache `koanf:"redis-cache"`
}
