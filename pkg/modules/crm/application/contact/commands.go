package contact

// CreateCommand 创建联系人命令。
type CreateCommand struct {
	FirstName string
	LastName  string
	Email     string
	Phone     string
	Title     string
	CompanyID *uint
	OwnerID   uint
}

// UpdateCommand 更新联系人命令。
type UpdateCommand struct {
	ID        uint
	FirstName string
	LastName  string
	Email     string
	Phone     string
	Title     string
	CompanyID *uint
}

// DeleteCommand 删除联系人命令。
type DeleteCommand struct {
	ID uint
}
