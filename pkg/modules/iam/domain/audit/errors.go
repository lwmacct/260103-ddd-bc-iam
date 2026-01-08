package audit

import "errors"

var (
	// ErrAuditLogNotFound 审计日志不存在
	ErrAuditLogNotFound = errors.New("audit log not found")

	// ErrInvalidLogID 无效的日志 ID
	ErrInvalidLogID = errors.New("invalid log ID")

	// ErrInvalidFilter 无效的过滤条件
	ErrInvalidFilter = errors.New("invalid filter criteria")

	// ErrInvalidDateRange 无效的日期范围
	ErrInvalidDateRange = errors.New("invalid date range")

	// ErrInvalidAction 无效的操作类型
	ErrInvalidAction = errors.New("invalid action type")

	// ErrInvalidResource 无效的资源类型
	ErrInvalidResource = errors.New("invalid resource type")
)
