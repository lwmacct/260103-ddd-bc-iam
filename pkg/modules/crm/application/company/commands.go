package company

// CreateCommand 创建公司命令。
type CreateCommand struct {
	Name     string
	Industry string
	Size     string
	Website  string
	Address  string
	OwnerID  uint
}

// UpdateCommand 更新公司命令。
type UpdateCommand struct {
	ID       uint
	Name     *string
	Industry *string
	Size     *string
	Website  *string
	Address  *string
}

// DeleteCommand 删除公司命令。
type DeleteCommand struct {
	ID uint
}
