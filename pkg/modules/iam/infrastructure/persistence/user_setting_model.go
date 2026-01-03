package persistence

import (
	"encoding/json"
	"time"

	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/app/domain/setting"
	"gorm.io/datatypes"
)

// UserSettingModel 用户配置的 GORM 实体
//
//nolint:recvcheck // TableName uses value receiver per GORM convention
type UserSettingModel struct {
	ID         uint           `gorm:"primaryKey"`
	UserID     uint           `gorm:"not null;uniqueIndex:idx_user_setting_key"` // 复合唯一索引
	SettingKey string         `gorm:"size:100;not null;uniqueIndex:idx_user_setting_key"`
	Value      datatypes.JSON `gorm:"type:jsonb;not null"` // JSONB 原生值

	CreatedAt time.Time
	UpdatedAt time.Time
}

// TableName 指定用户配置表名
func (UserSettingModel) TableName() string {
	return "user_settings"
}

func newUserSettingModelFromEntity(entity *setting.UserSetting) *UserSettingModel {
	if entity == nil {
		return nil
	}

	// 将 any 类型的 Value 序列化为 JSON
	var valueJSON []byte
	if entity.Value != nil {
		valueJSON, _ = json.Marshal(entity.Value) //nolint:errchkjson // Value 是任意 JSONB 值，序列化不会失败
	} else {
		valueJSON = []byte("null")
	}

	return &UserSettingModel{
		ID:         entity.ID,
		UserID:     entity.UserID,
		SettingKey: entity.SettingKey,
		Value:      datatypes.JSON(valueJSON),
		CreatedAt:  entity.CreatedAt,
		UpdatedAt:  entity.UpdatedAt,
	}
}

// ToEntity 将 GORM Model 转换为 Domain Entity
func (m *UserSettingModel) ToEntity() *setting.UserSetting {
	if m == nil {
		return nil
	}

	// 将 JSON 反序列化为 any 类型
	var value any
	_ = json.Unmarshal(m.Value, &value)

	return &setting.UserSetting{
		ID:         m.ID,
		UserID:     m.UserID,
		SettingKey: m.SettingKey,
		Value:      value,
		CreatedAt:  m.CreatedAt,
		UpdatedAt:  m.UpdatedAt,
	}
}

func mapUserSettingModelsToEntities(models []UserSettingModel) []*setting.UserSetting {
	if len(models) == 0 {
		return nil
	}

	settings := make([]*setting.UserSetting, 0, len(models))
	for i := range models {
		if entity := models[i].ToEntity(); entity != nil {
			settings = append(settings, entity)
		}
	}

	return settings
}
