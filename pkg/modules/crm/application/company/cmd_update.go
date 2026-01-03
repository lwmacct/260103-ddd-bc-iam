package company

import (
	"context"

	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/crm/domain/company"
)

// UpdateHandler 更新公司处理器。
type UpdateHandler struct {
	cmdRepo   company.CommandRepository
	queryRepo company.QueryRepository
}

// NewUpdateHandler 创建 UpdateHandler。
func NewUpdateHandler(cmdRepo company.CommandRepository, queryRepo company.QueryRepository) *UpdateHandler {
	return &UpdateHandler{cmdRepo: cmdRepo, queryRepo: queryRepo}
}

// Handle 处理更新公司命令。
func (h *UpdateHandler) Handle(ctx context.Context, cmd UpdateCommand) (*CompanyDTO, error) {
	c, err := h.queryRepo.GetByID(ctx, cmd.ID)
	if err != nil {
		return nil, err
	}

	// 如果更新名称，检查唯一性
	if cmd.Name != nil && *cmd.Name != c.Name {
		existing, _ := h.queryRepo.GetByName(ctx, *cmd.Name)
		if existing != nil {
			return nil, company.ErrCompanyNameExists
		}
		c.Name = *cmd.Name
	}

	if cmd.Industry != nil {
		c.Industry = *cmd.Industry
	}
	if cmd.Size != nil {
		if *cmd.Size != "" && !company.IsValidSize(*cmd.Size) {
			return nil, company.ErrInvalidSize
		}
		c.Size = *cmd.Size
	}
	if cmd.Website != nil {
		c.Website = *cmd.Website
	}
	if cmd.Address != nil {
		c.Address = *cmd.Address
	}

	if err := h.cmdRepo.Update(ctx, c); err != nil {
		return nil, err
	}

	return ToCompanyDTO(c), nil
}
