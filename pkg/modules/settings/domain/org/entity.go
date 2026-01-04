package org

import (
	"encoding/json"
	"time"
)

// OrgSetting 组织配置实体
// 存储组织对系统配置项的自定义覆盖值
type OrgSetting struct {
	ID         uint      `json:"id"`
	OrgID      uint      `json:"org_id"`
	SettingKey string    `json:"setting_key"`
	Value      any       `json:"value"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// New 创建新的组织配置实体
func New(orgID uint, key string, value any) *OrgSetting {
	return &OrgSetting{
		OrgID:      orgID,
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
func (os *OrgSetting) ValueOf(target any) error {
	if os == nil {
		return ErrOrgSettingNotFound
	}

	// 如果 Value 已经是目标类型，直接赋值
	data, err := json.Marshal(os.Value)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, target)
}

// SetValue 设置配置值
func (os *OrgSetting) SetValue(value any) {
	if os != nil {
		os.Value = value
	}
}

// IsEmpty 检查配置值是否为空
func (os *OrgSetting) IsEmpty() bool {
	return os == nil || os.Value == nil
}
