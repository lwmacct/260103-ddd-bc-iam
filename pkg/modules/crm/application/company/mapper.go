package company

import "github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/crm/domain/company"

// ToCompanyDTO 将公司实体转换为 DTO。
func ToCompanyDTO(c *company.Company) *CompanyDTO {
	if c == nil {
		return nil
	}
	return &CompanyDTO{
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

// ToCompanyDTOs 将公司实体列表转换为 DTO 列表。
func ToCompanyDTOs(companies []*company.Company) []*CompanyDTO {
	dtos := make([]*CompanyDTO, len(companies))
	for i, c := range companies {
		dtos[i] = ToCompanyDTO(c)
	}
	return dtos
}
