package company

import (
	"context"

	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/crm/domain/company"
)

// CreateHandler 创建公司处理器。
type CreateHandler struct {
	cmdRepo   company.CommandRepository
	queryRepo company.QueryRepository
}

// NewCreateHandler 创建 CreateHandler。
func NewCreateHandler(cmdRepo company.CommandRepository, queryRepo company.QueryRepository) *CreateHandler {
	return &CreateHandler{cmdRepo: cmdRepo, queryRepo: queryRepo}
}

// Handle 处理创建公司命令。
func (h *CreateHandler) Handle(ctx context.Context, cmd CreateCommand) (*CompanyDTO, error) {
	// 验证公司名称唯一性
	existing, _ := h.queryRepo.GetByName(ctx, cmd.Name)
	if existing != nil {
		return nil, company.ErrCompanyNameExists
	}

	// 验证规模值
	if cmd.Size != "" && !company.IsValidSize(cmd.Size) {
		return nil, company.ErrInvalidSize
	}

	c := &company.Company{
		Name:     cmd.Name,
		Industry: cmd.Industry,
		Size:     cmd.Size,
		Website:  cmd.Website,
		Address:  cmd.Address,
		OwnerID:  cmd.OwnerID,
	}

	if err := h.cmdRepo.Create(ctx, c); err != nil {
		return nil, err
	}

	return ToCompanyDTO(c), nil
}
