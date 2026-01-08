package twofa

import "errors"

var (
	// ErrTwoFANotFound 双因素认证配置不存在
	ErrTwoFANotFound = errors.New("two-factor authentication configuration not found")

	// ErrTwoFAAlreadyEnabled 双因素认证已启用
	ErrTwoFAAlreadyEnabled = errors.New("two-factor authentication already enabled")

	// ErrTwoFANotEnabled 双因素认证未启用
	ErrTwoFANotEnabled = errors.New("two-factor authentication not enabled")

	// ErrInvalidTOTPCode 无效的 TOTP 验证码
	ErrInvalidTOTPCode = errors.New("invalid TOTP verification code")

	// ErrInvalidRecoveryCode 无效的恢复码
	ErrInvalidRecoveryCode = errors.New("invalid recovery code")

	// ErrNoRecoveryCodesLeft 没有剩余的恢复码
	ErrNoRecoveryCodesLeft = errors.New("no recovery codes left")

	// ErrSecretNotSet 密钥未设置
	ErrSecretNotSet = errors.New("TOTP secret not set")

	// ErrSetupNotComplete 设置未完成
	ErrSetupNotComplete = errors.New("two-factor authentication setup not completed")

	// ErrRecoveryCodeAlreadyUsed 恢复码已被使用
	ErrRecoveryCodeAlreadyUsed = errors.New("recovery code already used")
)
