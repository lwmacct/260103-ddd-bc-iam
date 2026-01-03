package usersetting

import (
	"encoding/json"
	"time"
)

// UserSetting 用户设置实体
// 表示用户对系统配置的自定义值，与 Settings BC 的 Setting.Schema 配合使用
type UserSetting struct {
	ID        uint       `json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`

	UserID   uint   `json:"user_id"`  // 所属用户
	Key      string `json:"key"`      // 设置键名（关联 setting.Setting.Key）
	Value    string `json:"value"`    // 用户自定义值（JSON 字符串）
	Category string `json:"category"` // 分类（冗余字段，便于查询和过滤）
}

// ValueOf 解析 JSON 值到指定类型
//
// 示例：
//
//	var enabled bool
//	if err := setting.ValueOf(&enabled); err != nil {
//	    return err
//	}
func (us *UserSetting) ValueOf(target any) error {
	if us == nil {
		return ErrUserSettingNotFound
	}
	return json.Unmarshal([]byte(us.Value), target)
}

// SetValue 从任意类型设置 JSON 值
//
// 示例：
//
//	if err := setting.SetValue(true); err != nil {
//	    return err
//	}
func (us *UserSetting) SetValue(value any) error {
	if us == nil {
		return ErrUserSettingNotFound
	}
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	us.Value = string(data)
	return nil
}
