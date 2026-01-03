package contact

import "time"

// ContactDTO 联系人数据传输对象。
type ContactDTO struct {
	ID        uint      `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	FullName  string    `json:"full_name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Title     string    `json:"title"`
	CompanyID *uint     `json:"company_id,omitempty"`
	OwnerID   uint      `json:"owner_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateContactDTO 创建联系人请求。
type CreateContactDTO struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Phone     string `json:"phone"`
	Title     string `json:"title"`
	CompanyID *uint  `json:"company_id"`
}

// UpdateContactDTO 更新联系人请求。
type UpdateContactDTO struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email" binding:"omitempty,email"`
	Phone     string `json:"phone"`
	Title     string `json:"title"`
	CompanyID *uint  `json:"company_id"`
}

// ListResultDTO 联系人列表结果。
type ListResultDTO struct {
	Items []*ContactDTO `json:"items"`
	Total int64         `json:"total"`
}
