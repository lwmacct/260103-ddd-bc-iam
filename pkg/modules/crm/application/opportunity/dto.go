package opportunity

import (
	"time"

	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/crm/domain/opportunity"
)

// OpportunityDTO 商机数据传输对象。
type OpportunityDTO struct {
	ID             uint       `json:"id"`
	Name           string     `json:"name"`
	ContactID      uint       `json:"contact_id"`
	CompanyID      *uint      `json:"company_id,omitempty"`
	LeadID         *uint      `json:"lead_id,omitempty"`
	Stage          string     `json:"stage"`
	Amount         float64    `json:"amount"`
	Probability    int        `json:"probability"`
	WeightedAmount float64    `json:"weighted_amount"`
	ExpectedClose  *time.Time `json:"expected_close,omitempty"`
	OwnerID        uint       `json:"owner_id"`
	Notes          string     `json:"notes,omitempty"`
	ClosedAt       *time.Time `json:"closed_at,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

// CreateOpportunityDTO 创建商机请求。
type CreateOpportunityDTO struct {
	Name          string     `json:"name" binding:"required,max=200"`
	ContactID     uint       `json:"contact_id" binding:"required"`
	CompanyID     *uint      `json:"company_id,omitempty"`
	LeadID        *uint      `json:"lead_id,omitempty"`
	Amount        float64    `json:"amount" binding:"min=0"`
	Probability   int        `json:"probability" binding:"min=0,max=100"`
	ExpectedClose *time.Time `json:"expected_close,omitempty"`
	Notes         string     `json:"notes,omitempty" binding:"max=2000"`
}

// UpdateOpportunityDTO 更新商机请求。
type UpdateOpportunityDTO struct {
	Name          *string    `json:"name,omitempty" binding:"omitempty,max=200"`
	ContactID     *uint      `json:"contact_id,omitempty"`
	CompanyID     *uint      `json:"company_id,omitempty"`
	Amount        *float64   `json:"amount,omitempty" binding:"omitempty,min=0"`
	Probability   *int       `json:"probability,omitempty" binding:"omitempty,min=0,max=100"`
	ExpectedClose *time.Time `json:"expected_close,omitempty"`
	Notes         *string    `json:"notes,omitempty" binding:"omitempty,max=2000"`
}

// AdvanceStageDTO 推进阶段请求。
type AdvanceStageDTO struct {
	Stage string `json:"stage" binding:"required,oneof=proposal negotiation"`
}

// OpportunityListDTO 商机列表结果。
type OpportunityListDTO struct {
	Items []*OpportunityDTO `json:"items"`
	Total int64             `json:"total"`
}

// ToOpportunityDTO 将商机实体转换为 DTO。
func ToOpportunityDTO(opp *opportunity.Opportunity) *OpportunityDTO {
	return &OpportunityDTO{
		ID:             opp.ID,
		Name:           opp.Name,
		ContactID:      opp.ContactID,
		CompanyID:      opp.CompanyID,
		LeadID:         opp.LeadID,
		Stage:          string(opp.Stage),
		Amount:         opp.Amount,
		Probability:    opp.Probability,
		WeightedAmount: opp.WeightedAmount(),
		ExpectedClose:  opp.ExpectedClose,
		OwnerID:        opp.OwnerID,
		Notes:          opp.Notes,
		ClosedAt:       opp.ClosedAt,
		CreatedAt:      opp.CreatedAt,
		UpdatedAt:      opp.UpdatedAt,
	}
}

// ToOpportunityDTOs 将商机实体列表转换为 DTO 列表。
func ToOpportunityDTOs(opps []*opportunity.Opportunity) []*OpportunityDTO {
	dtos := make([]*OpportunityDTO, len(opps))
	for i, opp := range opps {
		dtos[i] = ToOpportunityDTO(opp)
	}
	return dtos
}
