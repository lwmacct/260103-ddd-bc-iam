package lead

import "errors"

// 线索相关错误
var (
	// ErrLeadNotFound 线索不存在
	ErrLeadNotFound = errors.New("lead not found")
)

// 状态转换相关错误
var (
	// ErrCannotContact 无法转换到已联系状态
	ErrCannotContact = errors.New("cannot contact: lead must be in 'new' status")

	// ErrCannotQualify 无法转换到已确认状态
	ErrCannotQualify = errors.New("cannot qualify: lead must be in 'contacted' status")

	// ErrCannotConvert 无法转化
	ErrCannotConvert = errors.New("cannot convert: lead must be in 'qualified' status")

	// ErrCannotLose 无法标记为丢失
	ErrCannotLose = errors.New("cannot lose: lead is already closed")

	// ErrAlreadyClosed 线索已关闭
	ErrAlreadyClosed = errors.New("lead is already closed")
)
