package persistence

import (
	"time"

	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/crm/domain/company"
)

// CompanyModel GORM 公司模型。
type CompanyModel struct {
	ID        uint      `gorm:"primaryKey"`
	Name      string    `gorm:"uniqueIndex;size:200;not null"`
	Industry  string    `gorm:"size:100"`
	Size      string    `gorm:"size:20"`
	Website   string    `gorm:"size:500"`
	Address   string    `gorm:"size:500"`
	OwnerID   uint      `gorm:"index;not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

// TableName 返回表名。
func (CompanyModel) TableName() string {
	return "companies"
}

// toCompanyModel 将公司实体转换为 GORM 模型。
func toCompanyModel(c *company.Company) *CompanyModel {
	return &CompanyModel{
		ID:        c.ID,
		Name:      c.Name,
		Industry:  c.Industry,
		Size:      c.Size,
		Website:   c.Website,
		Address:   c.Address,
		OwnerID:   c.OwnerID,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}

// toCompanyEntity 将 GORM 模型转换为公司实体。
func toCompanyEntity(m *CompanyModel) *company.Company {
	return &company.Company{
		ID:        m.ID,
		Name:      m.Name,
		Industry:  m.Industry,
		Size:      m.Size,
		Website:   m.Website,
		Address:   m.Address,
		OwnerID:   m.OwnerID,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

// toCompanyEntities 将 GORM 模型列表转换为公司实体列表。
func toCompanyEntities(models []*CompanyModel) []*company.Company {
	entities := make([]*company.Company, len(models))
	for i, m := range models {
		entities[i] = toCompanyEntity(m)
	}
	return entities
}
