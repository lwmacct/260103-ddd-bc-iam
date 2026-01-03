package lead

// CreateCommand 创建线索命令。
type CreateCommand struct {
	Title       string
	ContactID   *uint
	CompanyName string
	Source      string
	Score       int
	OwnerID     uint
	Notes       string
}

// UpdateCommand 更新线索命令。
type UpdateCommand struct {
	ID          uint
	Title       *string
	ContactID   *uint
	CompanyName *string
	Source      *string
	Score       *int
	Notes       *string
}

// DeleteCommand 删除线索命令。
type DeleteCommand struct {
	ID uint
}

// ContactCommand 转换到已联系状态命令。
type ContactCommand struct {
	ID uint
}

// QualifyCommand 转换到已确认状态命令。
type QualifyCommand struct {
	ID uint
}

// ConvertCommand 转化为商机命令。
type ConvertCommand struct {
	ID uint
}

// LoseCommand 标记为丢失命令。
type LoseCommand struct {
	ID uint
}
