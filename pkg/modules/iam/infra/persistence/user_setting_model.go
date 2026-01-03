package persistence

import (
	"strings"
	"time"

	"github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/domain/usersetting"
)

// UserSettingModel 定义用户设置的 GORM 实体
//

type UserSettingModel struct {
	ID        uint       `gorm:"primarykey"`
	CreatedAt time.Time  `gorm:"autoCreateTime"`
	UpdatedAt time.Time  `gorm:"autoUpdateTime"`
	DeletedAt *time.Time `gorm:"index"`
	UserID    uint       `gorm:"not null;index:idx_user_settings_user_key,unique"`
	Key       string     `gorm:"size:255;not null;index:idx_user_settings_user_key,unique"`
	Value     string     `gorm:"type:text;not null"`
	Category  string     `gorm:"size:100;not null;index"`
}

// TableName 指定用户设置表名
func (UserSettingModel) TableName() string {
	return "user_settings"
}

// newUserSettingModelFromEntity 从领域实体创建 GORM 模型
func newUserSettingModelFromEntity(entity *usersetting.UserSetting) *UserSettingModel {
	if entity == nil {
		return nil
	}
	return &UserSettingModel{
		ID:        entity.ID,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
		DeletedAt: entity.DeletedAt,
		UserID:    entity.UserID,
		Key:       entity.Key,
		Value:     entity.Value,
		Category:  entity.Category,
	}
}

// newUserSettingEntityFromModel 从 GORM 模型创建领域实体
func newUserSettingEntityFromModel(model *UserSettingModel) *usersetting.UserSetting {
	if model == nil {
		return nil
	}
	return &usersetting.UserSetting{
		ID:        model.ID,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
		DeletedAt: model.DeletedAt,
		UserID:    model.UserID,
		Key:       model.Key,
		Value:     model.Value,
		Category:  model.Category,
	}
}

// newUserSettingEntitiesFromModels 从 GORM 模型切片创建领域实体切片
func newUserSettingEntitiesFromModels(models []*UserSettingModel) []*usersetting.UserSetting {
	entities := make([]*usersetting.UserSetting, 0, len(models))
	for _, model := range models {
		if model != nil {
			entities = append(entities, newUserSettingEntityFromModel(model))
		}
	}
	return entities
}

// extractCategoryFromKey 从键名中提取分类
// 例如：theme.dark_mode -> theme, notification.email_enabled -> notification
func extractCategoryFromKey(key string) string {
	parts := strings.Split(key, ".")
	if len(parts) > 0 {
		return parts[0]
	}
	return "general"
}
