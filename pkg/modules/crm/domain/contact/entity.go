package contact

import "time"

// Contact 联系人实体。
type Contact struct {
	ID        uint
	FirstName string
	LastName  string
	Email     string // 唯一约束
	Phone     string
	Title     string // 职位
	CompanyID *uint  // 可选关联公司
	OwnerID   uint   // 负责人（用户 ID）
	CreatedAt time.Time
	UpdatedAt time.Time
}

// FullName 返回联系人全名。
func (c *Contact) FullName() string {
	if c.FirstName == "" {
		return c.LastName
	}
	if c.LastName == "" {
		return c.FirstName
	}
	return c.FirstName + " " + c.LastName
}
