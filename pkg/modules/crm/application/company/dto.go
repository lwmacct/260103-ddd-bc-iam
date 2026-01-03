package company

import "time"

// CompanyDTO 公司数据传输对象。
type CompanyDTO struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Industry  string    `json:"industry"`
	Size      string    `json:"size"`
	Website   string    `json:"website"`
	Address   string    `json:"address"`
	OwnerID   uint      `json:"owner_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateCompanyDTO 创建公司请求 DTO。
type CreateCompanyDTO struct {
	Name     string `json:"name" binding:"required,min=1,max=200"`
	Industry string `json:"industry" binding:"max=100"`
	Size     string `json:"size" binding:"omitempty,oneof=small medium large enterprise"`
	Website  string `json:"website" binding:"omitempty,url,max=500"`
	Address  string `json:"address" binding:"max=500"`
}

// UpdateCompanyDTO 更新公司请求 DTO。
type UpdateCompanyDTO struct {
	Name     *string `json:"name" binding:"omitempty,min=1,max=200"`
	Industry *string `json:"industry" binding:"omitempty,max=100"`
	Size     *string `json:"size" binding:"omitempty,oneof=small medium large enterprise"`
	Website  *string `json:"website" binding:"omitempty,url,max=500"`
	Address  *string `json:"address" binding:"omitempty,max=500"`
}

// ListResultDTO 公司列表结果 DTO。
type ListResultDTO struct {
	Items []*CompanyDTO
	Total int64
}
