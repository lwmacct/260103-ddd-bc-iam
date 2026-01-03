package container

import (
	internalConfig "github.com/lwmacct/260103-ddd-bc-iam/internal/config"
	"github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/config"
)

// ToIAMConfig 从通用配置中提取 IAM 模块配置。
//
// 这是一个适配器函数，用于平滑迁移。YAML 配置文件仍加载到完整的
// internal/config.Config 中，然后转换到 IAM 专属配置结构。
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
