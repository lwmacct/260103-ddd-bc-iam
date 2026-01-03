package persistence

import (
	"encoding/json"
	"time"

	"github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/settings/domain/user"
	"gorm.io/datatypes"
)

// UserSettingModel 用户配置的 GORM 模型
type UserSettingModel struct {
	ID         uint           `gorm:"primaryKey"`
	UserID     uint           `gorm:"not null;uniqueIndex:idx_user_settings_user_key"`
	SettingKey string         `gorm:"column:setting_key;size:100;not null;uniqueIndex:idx_user_settings_user_key"`
	Value      datatypes.JSON `gorm:"type:jsonb;not null"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// TableName 指定表名
func (*UserSettingModel) TableName() string {
	return "user_settings"
}

// ToEntity 转换为领域实体
func (m *UserSettingModel) ToEntity() *user.UserSetting {
	if m == nil {
		return nil
	}

	var value any
	if len(m.Value) > 0 {
		_ = json.Unmarshal(m.Value, &value)
	}

	return &user.UserSetting{
		ID:         m.ID,
		UserID:     m.UserID,
		SettingKey: m.SettingKey,
		Value:      value,
		CreatedAt:  m.CreatedAt,
		UpdatedAt:  m.UpdatedAt,
	}
}

// newModelFromEntity 从领域实体创建模型
func newModelFromEntity(e *user.UserSetting) *UserSettingModel {
	if e == nil {
		return nil
	}

	var valueJSON []byte
	if e.Value != nil {
		var err error
		valueJSON, err = json.Marshal(e.Value)
		if err != nil {
			valueJSON = []byte("null")
		}
	} else {
		valueJSON = []byte("null")
	}

	return &UserSettingModel{
		ID:         e.ID,
		UserID:     e.UserID,
		SettingKey: e.SettingKey,
		Value:      datatypes.JSON(valueJSON),
		CreatedAt:  e.CreatedAt,
		UpdatedAt:  e.UpdatedAt,
	}
}

// toEntities 批量转换为领域实体
func toEntities(models []*UserSettingModel) []*user.UserSetting {
	entities := make([]*user.UserSetting, 0, len(models))
	for _, m := range models {
		if m != nil {
			entities = append(entities, m.ToEntity())
		}
	}
	return entities
}
