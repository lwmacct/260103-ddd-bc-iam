package contact

import "github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/crm/domain/contact"

// ToContactDTO 将实体转换为 DTO。
func ToContactDTO(entity *contact.Contact) *ContactDTO {
	if entity == nil {
		return nil
	}
	return &ContactDTO{
		ID:        entity.ID,
		FirstName: entity.FirstName,
		LastName:  entity.LastName,
		FullName:  entity.FullName(),
		Email:     entity.Email,
		Phone:     entity.Phone,
		Title:     entity.Title,
		CompanyID: entity.CompanyID,
		OwnerID:   entity.OwnerID,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
	}
}

// ToContactDTOs 将实体列表转换为 DTO 列表。
func ToContactDTOs(entities []*contact.Contact) []*ContactDTO {
	dtos := make([]*ContactDTO, len(entities))
	for i, entity := range entities {
		dtos[i] = ToContactDTO(entity)
	}
	return dtos
}
