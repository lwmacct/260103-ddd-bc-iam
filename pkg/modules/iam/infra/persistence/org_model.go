package persistence

import (
	"time"

	"github.com/lwmacct/260103-ddd-iam-bc/pkg/modules/iam/domain/org"
	"gorm.io/gorm"
)

// OrgModel 定义组织的 GORM 持久化模型
//
//nolint:recvcheck // TableName uses value receiver per GORM convention
type OrgModel struct {
	ID          uint   `gorm:"primaryKey"`
	Name        string `gorm:"uniqueIndex;size:50;not null"`
	DisplayName string `gorm:"size:100;not null"`
	Description string `gorm:"type:text"`
	Avatar      string `gorm:"size:255"`
	Status      string `gorm:"size:20;default:'active';not null"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	// 聚合关系
	Teams   []TeamModel      `gorm:"foreignKey:OrgID"`
	Members []OrgMemberModel `gorm:"foreignKey:OrgID"`
}

// TableName 指定组织表名
func (OrgModel) TableName() string {
	return "orgs"
}

func newOrgModelFromEntity(entity *org.Org) *OrgModel {
	if entity == nil {
		return nil
	}

	model := &OrgModel{
		ID:          entity.ID,
		CreatedAt:   entity.CreatedAt,
		UpdatedAt:   entity.UpdatedAt,
		Name:        entity.Name,
		DisplayName: entity.DisplayName,
		Description: entity.Description,
		Avatar:      entity.Avatar,
		Status:      entity.Status,
	}

	if entity.DeletedAt != nil {
		model.DeletedAt = gorm.DeletedAt{Time: *entity.DeletedAt, Valid: true}
	}

	return model
}

// ToEntity 将 GORM Model 转换为 Domain Entity
func (m *OrgModel) ToEntity() *org.Org {
	if m == nil {
		return nil
	}

	entity := &org.Org{
		ID:          m.ID,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
		Name:        m.Name,
		DisplayName: m.DisplayName,
		Description: m.Description,
		Avatar:      m.Avatar,
		Status:      m.Status,
		Teams:       mapTeamModelsToEntities(m.Teams),
		Members:     mapOrgMemberModelsToEntities(m.Members),
	}

	if m.DeletedAt.Valid {
		t := m.DeletedAt.Time
		entity.DeletedAt = &t
	}

	return entity
}

func mapOrgModelsToEntities(models []OrgModel) []*org.Org {
	if len(models) == 0 {
		return nil
	}

	orgs := make([]*org.Org, 0, len(models))
	for i := range models {
		if entity := models[i].ToEntity(); entity != nil {
			orgs = append(orgs, entity)
		}
	}
	return orgs
}
