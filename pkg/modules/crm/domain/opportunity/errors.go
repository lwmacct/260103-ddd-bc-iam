package opportunity

import "errors"

// 商机相关错误。
var (
	// ErrOpportunityNotFound 商机不存在。
	ErrOpportunityNotFound = errors.New("opportunity not found")

	// ErrContactRequired 必须关联联系人。
	ErrContactRequired = errors.New("contact is required")
)

// 阶段转换相关错误。
var (
	// ErrInvalidStageTransition 无效的阶段转换。
	ErrInvalidStageTransition = errors.New("invalid stage transition")

	// ErrAlreadyClosed 商机已关闭。
	ErrAlreadyClosed = errors.New("opportunity already closed")

	// ErrCannotClose 不能关闭商机。
	ErrCannotClose = errors.New("cannot close opportunity, must be in negotiation stage")
)
