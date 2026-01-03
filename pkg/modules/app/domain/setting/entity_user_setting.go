package setting

import "time"

// UserSetting 用户配置实体。
// 存储用户对配置项的自定义值，覆盖系统默认值。
//
// Value 字段（JSONB）直接存储原生 JSON 值，类型应与对应配置项的 ValueType 一致。
//
// 设计说明：
//   - 用户配置采用覆盖模式：只存储用户修改过的配置
//   - 删除 UserSetting 记录即可恢复默认值
//   - 通过 (UserID, SettingKey) 唯一约束保证每个用户每个配置只有一条记录
type UserSetting struct {
	ID         uint   `json:"id"`          // 唯一标识
	UserID     uint   `json:"user_id"`     // 用户 ID，关联 users 表
	SettingKey string `json:"setting_key"` // 配置键
	Value      any    `json:"value"`       // 用户自定义值（JSONB 原生值）

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
