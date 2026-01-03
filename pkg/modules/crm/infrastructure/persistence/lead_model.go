package persistence

import (
	"time"

	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/crm/domain/lead"
)

// LeadModel GORM 线索模型。
type LeadModel struct {
	ID          uint   `gorm:"primaryKey"`
	Title       string `gorm:"size:200;not null"`
	ContactID   *uint  `gorm:"index"`
	CompanyName string `gorm:"size:200"`
	Source      string `gorm:"size:50"`
	Status      string `gorm:"size:20;index;not null;default:'new'"`
	Score       int    `gorm:"default:0"`
	OwnerID     uint   `gorm:"index;not null"`
	Notes       string `gorm:"type:text"`
	ConvertedAt *time.Time
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}

// TableName 返回表名。
func (LeadModel) TableName() string {
	return "leads"
}

// toLeadModel 将线索实体转换为 GORM 模型。
func toLeadModel(l *lead.Lead) *LeadModel {
	return &LeadModel{
		ID:          l.ID,
		Title:       l.Title,
		ContactID:   l.ContactID,
		CompanyName: l.CompanyName,
		Source:      l.Source,
		Status:      string(l.Status),
		Score:       l.Score,
		OwnerID:     l.OwnerID,
		Notes:       l.Notes,
		ConvertedAt: l.ConvertedAt,
		CreatedAt:   l.CreatedAt,
		UpdatedAt:   l.UpdatedAt,
	}
}

// toLeadEntity 将 GORM 模型转换为线索实体。
func toLeadEntity(m *LeadModel) *lead.Lead {
	return &lead.Lead{
		ID:          m.ID,
		Title:       m.Title,
		ContactID:   m.ContactID,
		CompanyName: m.CompanyName,
		Source:      m.Source,
		Status:      lead.Status(m.Status),
		Score:       m.Score,
		OwnerID:     m.OwnerID,
		Notes:       m.Notes,
		ConvertedAt: m.ConvertedAt,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

// toLeadEntities 将 GORM 模型列表转换为线索实体列表。
func toLeadEntities(models []*LeadModel) []*lead.Lead {
	entities := make([]*lead.Lead, len(models))
	for i, m := range models {
		entities[i] = toLeadEntity(m)
	}
	return entities
}
