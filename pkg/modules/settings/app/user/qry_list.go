package user

import (
	"context"
	"fmt"
	"sort"

	settingdomain "github.com/lwmacct/260103-ddd-bc-settings/pkg/modules/settings/domain/setting"
	"github.com/lwmacct/260103-ddd-iam-bc/pkg/modules/settings/domain/user"
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
// 返回扁平结构的配置列表，每个 item 包含 category 和 group 字段供前端分组
func (h *ListHandler) Handle(ctx context.Context, query ListQuery) ([]SettingsItemDTO, error) {
	// 1. 获取配置定义列表
	defs, err := h.fetchSettings(ctx, query.Category)
	if err != nil {
		return nil, err
	}

	if len(defs) == 0 {
		return []SettingsItemDTO{}, nil
	}

	// 2. 过滤出对普通用户可见的设置
	defs = FilterByVisibleToUser(defs)
	if len(defs) == 0 {
		return []SettingsItemDTO{}, nil
	}

	// 3. 获取用户自定义值
	userSettings, err := h.queryRepo.FindByUser(ctx, query.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to find user settings: %w", err)
	}

	// 4. 构建用户配置映射
	userMap := make(map[string]*user.UserSetting)
	for _, us := range userSettings {
		userMap[us.SettingKey] = us
	}

	// 5. 获取所有分类元数据（用于填充 category key）
	categories, err := h.categoryQueryRepo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch categories: %w", err)
	}

	// 6. 构建 CategoryID -> Key 映射
	categoryKeyByID := make(map[uint]string, len(categories))
	categoryOrderByKey := make(map[string]int, len(categories))
	for _, cat := range categories {
		categoryKeyByID[cat.ID] = cat.Key
		categoryOrderByKey[cat.Key] = cat.Order
	}

	// 7. 转换为扁平 DTO 列表
	result := make([]SettingsItemDTO, 0, len(defs))
	for _, def := range defs {
		categoryKey := categoryKeyByID[def.CategoryID]
		us := userMap[def.Key]
		dto := ToSettingsItemDTO(def, us, categoryKey)
		if dto != nil {
			result = append(result, *dto)
		}
	}

	// 8. 按 Category Order + Group + Setting Order 排序
	sort.Slice(result, func(i, j int) bool {
		catOrderI := categoryOrderByKey[result[i].Category]
		catOrderJ := categoryOrderByKey[result[j].Category]
		if catOrderI != catOrderJ {
			return catOrderI < catOrderJ
		}
		if result[i].Group != result[j].Group {
			if result[i].Group == "default" {
				return false
			}
			if result[j].Group == "default" {
				return true
			}
			return result[i].Group < result[j].Group
		}
		return result[i].Order < result[j].Order
	})

	return result, nil
}

// fetchSettings 根据 Category 获取配置定义列表
func (h *ListHandler) fetchSettings(ctx context.Context, categoryKey string) ([]*settingdomain.Setting, error) {
	if categoryKey == "" {
		defs, err := h.settingQueryRepo.FindVisibleToUser(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to find settings: %w", err)
		}
		return defs, nil
	}

	category, err := h.categoryQueryRepo.FindByKey(ctx, categoryKey)
	if err != nil {
		return nil, fmt.Errorf("failed to find category by key %q: %w", categoryKey, err)
	}
	if category == nil {
		return nil, fmt.Errorf("category not found: %s", categoryKey)
	}

	defs, err := h.settingQueryRepo.FindByCategoryID(ctx, category.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to find settings by category: %w", err)
	}
	return defs, nil
}
