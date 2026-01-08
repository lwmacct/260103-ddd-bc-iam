package pat

import "errors"

var (
	// ErrTokenNotFound 令牌不存在
	ErrTokenNotFound = errors.New("token not found")

	// ErrTokenExpired 令牌已过期
	ErrTokenExpired = errors.New("token has expired")

	// ErrTokenDisabled 令牌已禁用
	ErrTokenDisabled = errors.New("token has been disabled")

	// ErrTokenAlreadyDisabled 令牌已处于禁用状态
	ErrTokenAlreadyDisabled = errors.New("token is already disabled")

	// ErrTokenAlreadyEnabled 令牌已处于启用状态
	ErrTokenAlreadyEnabled = errors.New("token is already enabled")

	// ErrInvalidTokenFormat 无效的令牌格式
	ErrInvalidTokenFormat = errors.New("invalid token format")

	// ErrInvalidTokenPrefix 无效的令牌前缀
	ErrInvalidTokenPrefix = errors.New("invalid token prefix")

	// ErrIPNotAllowed IP 不在白名单中
	ErrIPNotAllowed = errors.New("IP address not in whitelist")

	// ErrInsufficientPermissions 权限不足
	ErrInsufficientPermissions = errors.New("insufficient permissions")

	// ErrTokenNameExists 令牌名称已存在
	ErrTokenNameExists = errors.New("token name already exists")

	// ErrMaxTokensReached 已达到最大令牌数量
	ErrMaxTokensReached = errors.New("maximum token limit reached")
)
