package user

import (
	"context"
	"fmt"

	"github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/settings/domain/user"
	settingdomain "github.com/lwmacct/260103-ddd-bc-settings/pkg/modules/settings/domain/setting"
)

// GetHandler 获取单个配置查询处理器
type GetHandler struct {
	settingQueryRepo  settingdomain.QueryRepository
	categoryQueryRepo settingdomain.SettingCategoryQueryRepository
	queryRepo         user.QueryRepository
}

// NewGetHandler 创建获取配置查询处理器
func NewGetHandler(
	settingQueryRepo settingdomain.QueryRepository,
	categoryQueryRepo settingdomain.SettingCategoryQueryRepository,
	queryRepo user.QueryRepository,
) *GetHandler {
	return &GetHandler{
		settingQueryRepo:  settingQueryRepo,
		categoryQueryRepo: categoryQueryRepo,
		queryRepo:         queryRepo,
	}
}

// Handle 处理获取单个配置查询
//
// 返回合并后的配置：如果用户有自定义值则使用用户值，否则使用默认值
func (h *GetHandler) Handle(ctx context.Context, query GetQuery) (*SettingsItemDTO, error) {
	// 1. 获取配置定义
	def, err := h.settingQueryRepo.FindByKey(ctx, query.Key)
	if err != nil {
		return nil, fmt.Errorf("failed to find setting: %w", err)
	}
	if def == nil {
		return nil, user.ErrInvalidSettingKey
	}

	// 检查是否对用户可见
	if !def.IsVisibleToUser() {
		return nil, user.ErrInvalidSettingKey
	}

	// 2. 获取用户自定义值（可能为 nil）
	us, err := h.queryRepo.FindByUserAndKey(ctx, query.UserID, query.Key)
	if err != nil {
		return nil, fmt.Errorf("failed to find user setting: %w", err)
	}

	// 3. 获取 category key
	category, err := h.categoryQueryRepo.FindByID(ctx, def.CategoryID)
	if err != nil {
		return nil, fmt.Errorf("failed to find category: %w", err)
	}
	categoryKey := ""
	if category != nil {
		categoryKey = category.Key
	}

	// 4. 合并返回
	return ToSettingsItemDTO(def, us, categoryKey), nil
}
