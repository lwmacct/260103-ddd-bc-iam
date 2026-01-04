package team

import (
	"encoding/json"
	"time"
)

// TeamSetting 团队配置实体
// 存储团队对系统配置项的自定义覆盖值
type TeamSetting struct {
	ID         uint      `json:"id"`
	TeamID     uint      `json:"team_id"`
	SettingKey string    `json:"setting_key"`
	Value      any       `json:"value"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// New 创建新的团队配置实体
func New(teamID uint, key string, value any) *TeamSetting {
	return &TeamSetting{
		TeamID:     teamID,
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
func (ts *TeamSetting) ValueOf(target any) error {
	if ts == nil {
		return ErrTeamSettingNotFound
	}

	// 如果 Value 已经是目标类型，直接赋值
	data, err := json.Marshal(ts.Value)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, target)
}

// SetValue 设置配置值
func (ts *TeamSetting) SetValue(value any) {
	if ts != nil {
		ts.Value = value
	}
}

// IsEmpty 检查配置值是否为空
func (ts *TeamSetting) IsEmpty() bool {
	return ts == nil || ts.Value == nil
}
