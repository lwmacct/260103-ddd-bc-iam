package user

// CreateCommand 创建用户命令
type CreateCommand struct {
	Username  string
	Email     string
	Password  string
	RealName  string
	Nickname  string
	Phone     string
	Signature string
	Status    *string // 可选：初始状态，默认 "active"
	RoleIDs   []uint  // 可选：创建时分配角色
}

// UpdateCommand 更新用户命令
type UpdateCommand struct {
	UserID    uint
	Username  *string
	Email     *string
	RealName  *string
	Nickname  *string
	Phone     *string
	Signature *string
	Avatar    *string
	Bio       *string
	Status    *string
}

// DeleteCommand 删除用户命令
type DeleteCommand struct {
	UserID uint
}

// BatchCreateCommand 批量创建用户命令
type BatchCreateCommand struct {
	Users []BatchItemDTO
}

// ChangePasswordCommand 修改密码命令
type ChangePasswordCommand struct {
	UserID      uint
	OldPassword string
	NewPassword string
}

// AssignRolesCommand 分配用户角色命令
type AssignRolesCommand struct {
	UserID  uint
	RoleIDs []uint
}
