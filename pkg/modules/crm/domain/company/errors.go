package company

import "errors"

// 公司相关错误
var (
	// ErrCompanyNotFound 公司不存在
	ErrCompanyNotFound = errors.New("company not found")

	// ErrCompanyNameExists 公司名称已存在
	ErrCompanyNameExists = errors.New("company name already exists")

	// ErrInvalidSize 无效的公司规模
	ErrInvalidSize = errors.New("invalid company size")
)
