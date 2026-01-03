package persistence

import (
	"time"

	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/crm/domain/opportunity"
)

// OpportunityModel GORM 商机模型。
type OpportunityModel struct {
	ID            uint    `gorm:"primaryKey"`
	Name          string  `gorm:"size:200;not null"`
	ContactID     uint    `gorm:"index;not null"`
	CompanyID     *uint   `gorm:"index"`
	LeadID        *uint   `gorm:"index"`
	Stage         string  `gorm:"size:20;index;not null;default:'prospecting'"`
	Amount        float64 `gorm:"default:0"`
	Probability   int     `gorm:"default:0"`
	ExpectedClose *time.Time
	OwnerID       uint   `gorm:"index;not null"`
	Notes         string `gorm:"type:text"`
	ClosedAt      *time.Time
	CreatedAt     time.Time `gorm:"autoCreateTime"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime"`
}

// TableName 返回表名。
func (OpportunityModel) TableName() string {
	return "opportunities"
}

// toOpportunityModel 将商机实体转换为 GORM 模型。
func toOpportunityModel(o *opportunity.Opportunity) *OpportunityModel {
	return &OpportunityModel{
		ID:            o.ID,
		Name:          o.Name,
		ContactID:     o.ContactID,
		CompanyID:     o.CompanyID,
		LeadID:        o.LeadID,
		Stage:         string(o.Stage),
		Amount:        o.Amount,
		Probability:   o.Probability,
		ExpectedClose: o.ExpectedClose,
		OwnerID:       o.OwnerID,
		Notes:         o.Notes,
		ClosedAt:      o.ClosedAt,
		CreatedAt:     o.CreatedAt,
		UpdatedAt:     o.UpdatedAt,
	}
}

// toOpportunityEntity 将 GORM 模型转换为商机实体。
func toOpportunityEntity(m *OpportunityModel) *opportunity.Opportunity {
	return &opportunity.Opportunity{
		ID:            m.ID,
		Name:          m.Name,
		ContactID:     m.ContactID,
		CompanyID:     m.CompanyID,
		LeadID:        m.LeadID,
		Stage:         opportunity.Stage(m.Stage),
		Amount:        m.Amount,
		Probability:   m.Probability,
		ExpectedClose: m.ExpectedClose,
		OwnerID:       m.OwnerID,
		Notes:         m.Notes,
		ClosedAt:      m.ClosedAt,
		CreatedAt:     m.CreatedAt,
		UpdatedAt:     m.UpdatedAt,
	}
}

// toOpportunityEntities 将 GORM 模型列表转换为商机实体列表。
func toOpportunityEntities(models []*OpportunityModel) []*opportunity.Opportunity {
	entities := make([]*opportunity.Opportunity, len(models))
	for i, m := range models {
		entities[i] = toOpportunityEntity(m)
	}
	return entities
}
