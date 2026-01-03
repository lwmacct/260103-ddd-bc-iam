package audit

import "errors"

var (
	// ErrAuditLogNotFound 审计日志不存在
	ErrAuditLogNotFound = errors.New("审计日志不存在")

	// ErrInvalidLogID 无效的日志 ID
	ErrInvalidLogID = errors.New("无效的日志 ID")

	// ErrInvalidFilter 无效的过滤条件
	ErrInvalidFilter = errors.New("无效的过滤条件")

	// ErrInvalidDateRange 无效的日期范围
	ErrInvalidDateRange = errors.New("无效的日期范围")

	// ErrInvalidAction 无效的操作类型
	ErrInvalidAction = errors.New("无效的操作类型")

	// ErrInvalidResource 无效的资源类型
	ErrInvalidResource = errors.New("无效的资源类型")
)
