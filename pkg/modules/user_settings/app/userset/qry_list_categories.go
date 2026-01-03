package userset

import (
	"context"
	"fmt"

	settingdomain "github.com/lwmacct/260103-ddd-bc-settings/pkg/modules/settings/domain/setting"
)

// ListCategoriesHandler 获取分类列表查询处理器
type ListCategoriesHandler struct {
	categoryQueryRepo settingdomain.SettingCategoryQueryRepository
}

// NewListCategoriesHandler 创建获取分类列表查询处理器
func NewListCategoriesHandler(
	categoryQueryRepo settingdomain.SettingCategoryQueryRepository,
) *ListCategoriesHandler {
	return &ListCategoriesHandler{
		categoryQueryRepo: categoryQueryRepo,
	}
}

// Handle 处理获取分类列表查询
func (h *ListCategoriesHandler) Handle(ctx context.Context, _ ListCategoriesQuery) ([]*CategoryDTO, error) {
	categories, err := h.categoryQueryRepo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to find categories: %w", err)
	}

	return ToCategoryDTOs(categories), nil
}
