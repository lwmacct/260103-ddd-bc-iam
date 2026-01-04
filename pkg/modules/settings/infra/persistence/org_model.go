package persistence

import (
	"encoding/json"

	"github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/settings/domain/org"
	"gorm.io/datatypes"
)

// OrgSettingModel 组织配置的 GORM 模型
type OrgSettingModel struct {
	ID         uint           `gorm:"primaryKey"`
	OrgID      uint           `gorm:"not null;uniqueIndex:idx_org_settings_org_key"`
	SettingKey string         `gorm:"column:setting_key;size:100;not null;uniqueIndex:idx_org_settings_org_key"`
	Value      datatypes.JSON `gorm:"type:jsonb;not null"`
	CreatedAt  int64          `gorm:"autoCreateTime:milli"`
	UpdatedAt  int64          `gorm:"autoUpdateTime:milli"`
}

// TableName 指定表名
func (*OrgSettingModel) TableName() string {
	return "org_settings"
}

// ToEntity 转换为领域实体
func (m *OrgSettingModel) ToEntity() *org.OrgSetting {
	if m == nil {
		return nil
	}

	var value any
	if len(m.Value) > 0 {
		_ = json.Unmarshal(m.Value, &value)
	}

	return &org.OrgSetting{
		ID:         m.ID,
		OrgID:      m.OrgID,
		SettingKey: m.SettingKey,
		Value:      value,
	}
}

// newOrgModelFromEntity 从领域实体创建组织配置模型
func newOrgModelFromEntity(e *org.OrgSetting) *OrgSettingModel {
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

	return &OrgSettingModel{
		ID:         e.ID,
		OrgID:      e.OrgID,
		SettingKey: e.SettingKey,
		Value:      datatypes.JSON(valueJSON),
	}
}

// toOrgEntities 批量转换为领域实体
func toOrgEntities(models []*OrgSettingModel) []*org.OrgSetting {
	entities := make([]*org.OrgSetting, 0, len(models))
	for _, m := range models {
		if m != nil {
			entities = append(entities, m.ToEntity())
		}
	}
	return entities
}
