package persistence

import (
	"time"

	"gorm.io/gorm"

	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/crm/domain/contact"
)

// ContactModel 联系人数据模型。
type ContactModel struct {
	ID        uint           `gorm:"primaryKey"`
	FirstName string         `gorm:"size:100;not null"`
	LastName  string         `gorm:"size:100;not null"`
	Email     string         `gorm:"size:255;uniqueIndex;not null"`
	Phone     string         `gorm:"size:50"`
	Title     string         `gorm:"size:100"`
	CompanyID *uint          `gorm:"index"`
	OwnerID   uint           `gorm:"index;not null"`
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// TableName 返回表名。
func (m *ContactModel) TableName() string {
	return "contacts"
}

// newContactModelFromEntity 从实体创建模型。
func newContactModelFromEntity(entity *contact.Contact) *ContactModel {
	return &ContactModel{
		ID:        entity.ID,
		FirstName: entity.FirstName,
		LastName:  entity.LastName,
		Email:     entity.Email,
		Phone:     entity.Phone,
		Title:     entity.Title,
		CompanyID: entity.CompanyID,
		OwnerID:   entity.OwnerID,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
	}
}

// ToEntity 转换为实体。
func (m *ContactModel) ToEntity() *contact.Contact {
	return &contact.Contact{
		ID:        m.ID,
		FirstName: m.FirstName,
		LastName:  m.LastName,
		Email:     m.Email,
		Phone:     m.Phone,
		Title:     m.Title,
		CompanyID: m.CompanyID,
		OwnerID:   m.OwnerID,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

// mapContactModelsToEntities 批量转换模型为实体。
func mapContactModelsToEntities(models []ContactModel) []*contact.Contact {
	entities := make([]*contact.Contact, len(models))
	for i, m := range models {
		entities[i] = m.ToEntity()
	}
	return entities
}
