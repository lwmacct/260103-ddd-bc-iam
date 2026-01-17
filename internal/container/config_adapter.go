package container

import (
	"github.com/lwmacct/260103-ddd-iam-bc/internal/config"
	iamconfig "github.com/lwmacct/260103-ddd-iam-bc/pkg/modules/iam/config"
)

// ToIAMConfig 从通用配置转换到 IAM 专属配置。
func ToIAMConfig(cfg *config.Config) iamconfig.Config {
	return iamconfig.Config{
		JWT: iamconfig.JWT{
			Secret:             cfg.JWT.Secret,
			AccessTokenExpiry:  cfg.JWT.AccessTokenExpiry,
			RefreshTokenExpiry: cfg.JWT.RefreshTokenExpiry,
		},
		Auth: iamconfig.Auth{
			DevSecret:       cfg.Auth.DevSecret,
			TwoFAIssuer:     cfg.Auth.TwoFAIssuer,
			CaptchaRequired: cfg.Auth.CaptchaRequired,
		},
	}
}
