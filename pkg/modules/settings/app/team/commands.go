package team

// SetCommand 设置团队配置命令
type SetCommand struct {
	TeamID uint
	Key    string
	Value  any
}

// ResetCommand 重置单个配置命令
type ResetCommand struct {
	TeamID uint
	Key    string
}
