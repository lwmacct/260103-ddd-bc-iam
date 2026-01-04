package org

// SetCommand 设置组织配置命令
type SetCommand struct {
	OrgID uint
	Key   string
	Value any
}

// ResetCommand 重置单个配置命令
type ResetCommand struct {
	OrgID uint
	Key   string
}
