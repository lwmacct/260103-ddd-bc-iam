package user

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
//
// 只返回用户可见的分类（Scope 为 user 或 public）。
func (h *ListCategoriesHandler) Handle(ctx context.Context, _ ListCategoriesQuery) ([]*CategoryDTO, error) {
	// 只返回用户可见的分类
	categories, err := h.categoryQueryRepo.FindByVisibleScope(ctx, settingdomain.ScopeLevelUser)
	if err != nil {
		return nil, fmt.Errorf("failed to find visible categories: %w", err)
	}

	return ToCategoryDTOs(categories), nil
}
