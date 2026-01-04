package persistence

import (
	"encoding/json"

	"github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/settings/domain/team"
	"gorm.io/datatypes"
)

// TeamSettingModel 团队配置的 GORM 模型
type TeamSettingModel struct {
	ID         uint           `gorm:"primaryKey"`
	TeamID     uint           `gorm:"not null;uniqueIndex:idx_team_settings_team_key"`
	SettingKey string         `gorm:"column:setting_key;size:100;not null;uniqueIndex:idx_team_settings_team_key"`
	Value      datatypes.JSON `gorm:"type:jsonb;not null"`
	CreatedAt  int64          `gorm:"autoCreateTime:milli"`
	UpdatedAt  int64          `gorm:"autoUpdateTime:milli"`
}

// TableName 指定表名
func (*TeamSettingModel) TableName() string {
	return "team_settings"
}

// ToEntity 转换为领域实体
func (m *TeamSettingModel) ToEntity() *team.TeamSetting {
	if m == nil {
		return nil
	}

	var value any
	if len(m.Value) > 0 {
		_ = json.Unmarshal(m.Value, &value)
	}

	return &team.TeamSetting{
		ID:         m.ID,
		TeamID:     m.TeamID,
		SettingKey: m.SettingKey,
		Value:      value,
	}
}

// newTeamModelFromEntity 从领域实体创建团队配置模型
func newTeamModelFromEntity(e *team.TeamSetting) *TeamSettingModel {
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

	return &TeamSettingModel{
		ID:         e.ID,
		TeamID:     e.TeamID,
		SettingKey: e.SettingKey,
		Value:      datatypes.JSON(valueJSON),
	}
}

// toTeamEntities 批量转换为领域实体
func toTeamEntities(models []*TeamSettingModel) []*team.TeamSetting {
	entities := make([]*team.TeamSetting, 0, len(models))
	for _, m := range models {
		if m != nil {
			entities = append(entities, m.ToEntity())
		}
	}
	return entities
}
