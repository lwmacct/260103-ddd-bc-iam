package lead

import (
	"time"

	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/crm/domain/lead"
)

// LeadDTO 线索数据传输对象。
type LeadDTO struct {
	ID          uint       `json:"id"`
	Title       string     `json:"title"`
	ContactID   *uint      `json:"contact_id"`
	CompanyName string     `json:"company_name"`
	Source      string     `json:"source"`
	Status      string     `json:"status"`
	Score       int        `json:"score"`
	OwnerID     uint       `json:"owner_id"`
	Notes       string     `json:"notes"`
	ConvertedAt *time.Time `json:"converted_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`

	// 状态转换可用性
	CanContact bool `json:"can_contact"`
	CanQualify bool `json:"can_qualify"`
	CanConvert bool `json:"can_convert"`
	CanLose    bool `json:"can_lose"`
}

// CreateLeadDTO 创建线索请求 DTO。
type CreateLeadDTO struct {
	Title       string `json:"title" binding:"required,min=1,max=200"`
	ContactID   *uint  `json:"contact_id"`
	CompanyName string `json:"company_name" binding:"max=200"`
	Source      string `json:"source" binding:"omitempty,oneof=website referral campaign other"`
	Score       int    `json:"score" binding:"min=0,max=100"`
	Notes       string `json:"notes" binding:"max=2000"`
}

// UpdateLeadDTO 更新线索请求 DTO。
type UpdateLeadDTO struct {
	Title       *string `json:"title" binding:"omitempty,min=1,max=200"`
	ContactID   *uint   `json:"contact_id"`
	CompanyName *string `json:"company_name" binding:"omitempty,max=200"`
	Source      *string `json:"source" binding:"omitempty,oneof=website referral campaign other"`
	Score       *int    `json:"score" binding:"omitempty,min=0,max=100"`
	Notes       *string `json:"notes" binding:"omitempty,max=2000"`
}

// ListResultDTO 线索列表结果 DTO。
type ListResultDTO struct {
	Items []*LeadDTO
	Total int64
}

// StatusSummaryDTO 状态统计 DTO。
type StatusSummaryDTO struct {
	Status string `json:"status"`
	Count  int64  `json:"count"`
}

// ToLeadDTO 将线索实体转换为 DTO。
func ToLeadDTO(l *lead.Lead) *LeadDTO {
	if l == nil {
		return nil
	}
	return &LeadDTO{
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
		CanContact:  l.CanContact(),
		CanQualify:  l.CanQualify(),
		CanConvert:  l.CanConvert(),
		CanLose:     l.CanLose(),
	}
}

// ToLeadDTOs 将线索实体列表转换为 DTO 列表。
func ToLeadDTOs(leads []*lead.Lead) []*LeadDTO {
	dtos := make([]*LeadDTO, len(leads))
	for i, l := range leads {
		dtos[i] = ToLeadDTO(l)
	}
	return dtos
}
