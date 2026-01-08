package container

import (
	internalConfig "github.com/lwmacct/260103-ddd-iam-bc/internal/config"
	"github.com/lwmacct/260103-ddd-iam-bc/pkg/modules/iam/config"
)

// ToIAMConfig 从通用配置转换到 IAM 专属配置。
func ToIAMConfig(cfg *internalConfig.Config) config.Config {
	return config.Config{
		JWT: config.JWT{
			Secret:             cfg.JWT.Secret,
			AccessTokenExpiry:  cfg.JWT.AccessTokenExpiry,
			RefreshTokenExpiry: cfg.JWT.RefreshTokenExpiry,
		},
		Auth: config.Auth{
			DevSecret:       cfg.Auth.DevSecret,
			TwoFAIssuer:     cfg.Auth.TwoFAIssuer,
			CaptchaRequired: cfg.Auth.CaptchaRequired,
		},
		Redis: config.Redis{
			KeyPrefix: cfg.Data.RedisKeyPrefix,
		},
	}
}
