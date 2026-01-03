package contact

import (
	"context"

	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/crm/domain/contact"
)

// UpdateHandler 更新联系人处理器。
type UpdateHandler struct {
	cmdRepo   contact.CommandRepository
	queryRepo contact.QueryRepository
}

// NewUpdateHandler 创建 UpdateHandler。
func NewUpdateHandler(cmdRepo contact.CommandRepository, queryRepo contact.QueryRepository) *UpdateHandler {
	return &UpdateHandler{cmdRepo: cmdRepo, queryRepo: queryRepo}
}

// Handle 处理更新联系人命令。
func (h *UpdateHandler) Handle(ctx context.Context, cmd UpdateCommand) (*ContactDTO, error) {
	entity, err := h.queryRepo.GetByID(ctx, cmd.ID)
	if err != nil {
		return nil, err
	}

	// 如果邮箱变更，检查是否冲突
	if cmd.Email != "" && cmd.Email != entity.Email {
		existing, _ := h.queryRepo.GetByEmail(ctx, cmd.Email)
		if existing != nil && existing.ID != cmd.ID {
			return nil, contact.ErrEmailAlreadyExists
		}
		entity.Email = cmd.Email
	}

	// 更新字段
	if cmd.FirstName != "" {
		entity.FirstName = cmd.FirstName
	}
	if cmd.LastName != "" {
		entity.LastName = cmd.LastName
	}
	if cmd.Phone != "" {
		entity.Phone = cmd.Phone
	}
	if cmd.Title != "" {
		entity.Title = cmd.Title
	}
	entity.CompanyID = cmd.CompanyID

	if err := h.cmdRepo.Update(ctx, entity); err != nil {
		return nil, err
	}

	return ToContactDTO(entity), nil
}
