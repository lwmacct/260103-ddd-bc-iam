package captcha

import "errors"

var (
	// ErrCaptchaNotFound 验证码不存在
	ErrCaptchaNotFound = errors.New("验证码不存在")

	// ErrCaptchaExpired 验证码已过期
	ErrCaptchaExpired = errors.New("验证码已过期")

	// ErrInvalidCaptcha 验证码错误
	ErrInvalidCaptcha = errors.New("验证码错误")

	// ErrCaptchaAlreadyUsed 验证码已被使用
	ErrCaptchaAlreadyUsed = errors.New("验证码已被使用")

	// ErrTooManyAttempts 尝试次数过多
	ErrTooManyAttempts = errors.New("验证码尝试次数过多")

	// ErrCaptchaGenerationFailed 验证码生成失败
	ErrCaptchaGenerationFailed = errors.New("验证码生成失败")
)
