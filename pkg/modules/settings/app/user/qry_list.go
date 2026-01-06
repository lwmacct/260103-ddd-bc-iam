package user

import (
	"context"
	"fmt"

	"github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/settings/domain/user"
	settingdomain "github.com/lwmacct/260103-ddd-bc-settings/pkg/modules/settings/domain/setting"
)

// ListHandler 获取用户配置列表查询处理器
type ListHandler struct {
	settingQueryRepo  settingdomain.QueryRepository
	categoryQueryRepo settingdomain.SettingCategoryQueryRepository
	queryRepo         user.QueryRepository
}

// NewListHandler 创建获取配置列表查询处理器
func NewListHandler(
	settingQueryRepo settingdomain.QueryRepository,
	categoryQueryRepo settingdomain.SettingCategoryQueryRepository,
	queryRepo user.QueryRepository,
) *ListHandler {
	return &ListHandler{
		settingQueryRepo:  settingQueryRepo,
		categoryQueryRepo: categoryQueryRepo,
		queryRepo:         queryRepo,
	}
}

// Handle 处理获取用户配置列表查询
//
// 返回合并后的配置列表：如果用户有自定义值则使用用户值，否则使用默认值
func (h *ListHandler) Handle(ctx context.Context, query ListQuery) ([]*UserSettingDTO, error) {
	// 1. 获取配置定义列表
	var defs []*settingdomain.Setting
	var err error

	if query.Category != "" {
		// 根据 category key 查找分类 ID
		category, catErr := h.categoryQueryRepo.FindByKey(ctx, query.Category)
		if catErr != nil {
			return nil, fmt.Errorf("category not found: %s", query.Category)
		}
		defs, err = h.settingQueryRepo.FindByCategoryID(ctx, category.ID)
	} else {
		defs, err = h.settingQueryRepo.FindVisibleToUser(ctx)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find settings: %w", err)
	}

	if len(defs) == 0 {
		return []*UserSettingDTO{}, nil
	}

	// 2. 获取用户自定义值
	userSettings, err := h.queryRepo.FindByUser(ctx, query.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to find user settings: %w", err)
	}

	// 3. 构建用户配置映射
	userMap := make(map[string]*user.UserSetting)
	for _, us := range userSettings {
		userMap[us.SettingKey] = us
	}

	// 4. 合并返回
	result := make([]*UserSettingDTO, 0, len(defs))
	for _, def := range defs {
		us := userMap[def.Key]
		result = append(result, ToUserSettingDTO(def, us))
	}

	return result, nil
}
