package pat

import "errors"

var (
	// ErrTokenNotFound 令牌不存在
	ErrTokenNotFound = errors.New("令牌不存在")

	// ErrTokenExpired 令牌已过期
	ErrTokenExpired = errors.New("令牌已过期")

	// ErrTokenDisabled 令牌已禁用
	ErrTokenDisabled = errors.New("令牌已禁用")

	// ErrTokenAlreadyDisabled 令牌已处于禁用状态
	ErrTokenAlreadyDisabled = errors.New("令牌已处于禁用状态")

	// ErrTokenAlreadyEnabled 令牌已处于启用状态
	ErrTokenAlreadyEnabled = errors.New("令牌已处于启用状态")

	// ErrInvalidTokenFormat 无效的令牌格式
	ErrInvalidTokenFormat = errors.New("无效的令牌格式")

	// ErrInvalidTokenPrefix 无效的令牌前缀
	ErrInvalidTokenPrefix = errors.New("无效的令牌前缀")

	// ErrIPNotAllowed IP 不在白名单中
	ErrIPNotAllowed = errors.New("IP 地址不在白名单中")

	// ErrInsufficientPermissions 权限不足
	ErrInsufficientPermissions = errors.New("权限不足")

	// ErrTokenNameExists 令牌名称已存在
	ErrTokenNameExists = errors.New("令牌名称已存在")

	// ErrMaxTokensReached 已达到最大令牌数量
	ErrMaxTokensReached = errors.New("已达到最大令牌数量")
)
