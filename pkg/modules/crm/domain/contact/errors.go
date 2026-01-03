package contact

import "errors"

// 联系人相关错误
var (
	// ErrContactNotFound 联系人不存在
	ErrContactNotFound = errors.New("contact not found")

	// ErrEmailAlreadyExists 邮箱已存在
	ErrEmailAlreadyExists = errors.New("email already exists")
)
