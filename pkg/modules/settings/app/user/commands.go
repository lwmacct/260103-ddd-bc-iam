package user

// SetCommand 设置用户配置命令
type SetCommand struct {
	UserID uint
	Key    string
	Value  any
}

// BatchSetCommand 批量设置用户配置命令
type BatchSetCommand struct {
	UserID   uint
	Settings []SettingItemDTO
}

// ResetCommand 重置单个配置命令
type ResetCommand struct {
	UserID uint
	Key    string
}

// ResetAllCommand 重置所有配置命令
type ResetAllCommand struct {
	UserID uint
}
