package task

import "errors"

var (
	// ErrTaskNotFound 任务不存在。
	ErrTaskNotFound = errors.New("任务不存在")
	// ErrInvalidStatusTransition 无效的状态转换。
	ErrInvalidStatusTransition = errors.New("无效的状态转换")
)
