package usersetting

// UpdateCommand 更新用户设置命令
type UpdateCommand struct {
	UserID uint
	Key    string
	Value  string
}

// DeleteCommand 删除用户设置命令
type DeleteCommand struct {
	UserID uint
	Key    string
}
