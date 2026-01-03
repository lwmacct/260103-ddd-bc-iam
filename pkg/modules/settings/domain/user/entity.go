package user

import (
	"encoding/json"
	"time"
)

// UserSetting 用户配置实体
// 存储用户对系统配置项的自定义覆盖值
type UserSetting struct {
	ID         uint      `json:"id"`
	UserID     uint      `json:"user_id"`
	SettingKey string    `json:"setting_key"`
	Value      any       `json:"value"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// New 创建新的用户配置实体
func New(userID uint, key string, value any) *UserSetting {
	return &UserSetting{
		UserID:     userID,
		SettingKey: key,
		Value:      value,
	}
}

// ValueOf 将 Value 解析到指定类型
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

	// 如果 Value 已经是目标类型，直接赋值
	data, err := json.Marshal(us.Value)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, target)
}

// SetValue 设置配置值
func (us *UserSetting) SetValue(value any) {
	if us != nil {
		us.Value = value
	}
}

// IsEmpty 检查配置值是否为空
func (us *UserSetting) IsEmpty() bool {
	return us == nil || us.Value == nil
}
