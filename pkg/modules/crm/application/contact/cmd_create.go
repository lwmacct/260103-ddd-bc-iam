package contact

import (
	"context"

	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/crm/domain/contact"
)

// CreateHandler 创建联系人处理器。
type CreateHandler struct {
	cmdRepo   contact.CommandRepository
	queryRepo contact.QueryRepository
}

// NewCreateHandler 创建 CreateHandler。
func NewCreateHandler(cmdRepo contact.CommandRepository, queryRepo contact.QueryRepository) *CreateHandler {
	return &CreateHandler{cmdRepo: cmdRepo, queryRepo: queryRepo}
}

// Handle 处理创建联系人命令。
func (h *CreateHandler) Handle(ctx context.Context, cmd CreateCommand) (*ContactDTO, error) {
	// 检查邮箱是否已存在
	existing, _ := h.queryRepo.GetByEmail(ctx, cmd.Email)
	if existing != nil {
		return nil, contact.ErrEmailAlreadyExists
	}

	entity := &contact.Contact{
		FirstName: cmd.FirstName,
		LastName:  cmd.LastName,
		Email:     cmd.Email,
		Phone:     cmd.Phone,
		Title:     cmd.Title,
		CompanyID: cmd.CompanyID,
		OwnerID:   cmd.OwnerID,
	}

	if err := h.cmdRepo.Create(ctx, entity); err != nil {
		return nil, err
	}

	return ToContactDTO(entity), nil
}
