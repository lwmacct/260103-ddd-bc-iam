package twofa

import "errors"

var (
	// ErrTwoFANotFound 双因素认证配置不存在
	ErrTwoFANotFound = errors.New("双因素认证配置不存在")

	// ErrTwoFAAlreadyEnabled 双因素认证已启用
	ErrTwoFAAlreadyEnabled = errors.New("双因素认证已启用")

	// ErrTwoFANotEnabled 双因素认证未启用
	ErrTwoFANotEnabled = errors.New("双因素认证未启用")

	// ErrInvalidTOTPCode 无效的 TOTP 验证码
	ErrInvalidTOTPCode = errors.New("无效的 TOTP 验证码")

	// ErrInvalidRecoveryCode 无效的恢复码
	ErrInvalidRecoveryCode = errors.New("无效的恢复码")

	// ErrNoRecoveryCodesLeft 没有剩余的恢复码
	ErrNoRecoveryCodesLeft = errors.New("没有剩余的恢复码")

	// ErrSecretNotSet 密钥未设置
	ErrSecretNotSet = errors.New("TOTP 密钥未设置")

	// ErrSetupNotComplete 设置未完成
	ErrSetupNotComplete = errors.New("双因素认证设置未完成")

	// ErrRecoveryCodeAlreadyUsed 恢复码已被使用
	ErrRecoveryCodeAlreadyUsed = errors.New("恢复码已被使用")
)
