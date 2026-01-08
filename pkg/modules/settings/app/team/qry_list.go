package team

import (
	"context"
	"fmt"
	"sort"

	settingdomain "github.com/lwmacct/260103-ddd-bc-settings/pkg/modules/settings/domain/setting"
	"github.com/lwmacct/260103-ddd-iam-bc/pkg/modules/settings/domain/org"
	"github.com/lwmacct/260103-ddd-iam-bc/pkg/modules/settings/domain/team"
)

// ListHandler 获取团队配置列表查询处理器
type ListHandler struct {
	settingQueryRepo  settingdomain.QueryRepository
	categoryQueryRepo settingdomain.SettingCategoryQueryRepository
	teamQueryRepo     team.QueryRepository
	orgQueryRepo      org.QueryRepository // 用于继承查询
}

// NewListHandler 创建获取配置列表查询处理器
func NewListHandler(
	settingQueryRepo settingdomain.QueryRepository,
	categoryQueryRepo settingdomain.SettingCategoryQueryRepository,
	teamQueryRepo team.QueryRepository,
	orgQueryRepo org.QueryRepository,
) *ListHandler {
	return &ListHandler{
		settingQueryRepo:  settingQueryRepo,
		categoryQueryRepo: categoryQueryRepo,
		teamQueryRepo:     teamQueryRepo,
		orgQueryRepo:      orgQueryRepo,
	}
}

// Handle 处理获取团队配置列表查询
//
// 返回扁平结构的配置列表，每个 item 包含 category 和 group 字段供前端分组
// 优先级：团队配置 > 组织配置 > 系统默认值
func (h *ListHandler) Handle(ctx context.Context, query ListQuery) ([]SettingsItemDTO, error) {
	// 1. 获取配置定义列表
	defs, err := h.fetchSettings(ctx, query.Category)
	if err != nil {
		return nil, err
	}

	if len(defs) == 0 {
		return []SettingsItemDTO{}, nil
	}

	// 2. 确保只返回对 Team 可配置的设置
	defs = FilterByConfigurableAt(defs, settingdomain.ScopeLevelTeam)

	if len(defs) == 0 {
		return []SettingsItemDTO{}, nil
	}

	// 3. 批量获取团队和组织配置
	keys := make([]string, 0, len(defs))
	for _, def := range defs {
		keys = append(keys, def.Key)
	}

	teamSettings, _ := h.teamQueryRepo.FindByTeamAndKeys(ctx, query.TeamID, keys)
	orgSettings, _ := h.orgQueryRepo.FindByOrgAndKeys(ctx, query.OrgID, keys)

	// 4. 构建映射
	teamMap := make(map[string]*team.TeamSetting)
	for _, ts := range teamSettings {
		teamMap[ts.SettingKey] = ts
	}

	orgMap := make(map[string]*org.OrgSetting)
	for _, os := range orgSettings {
		orgMap[os.SettingKey] = os
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
		ts := teamMap[def.Key]
		os := orgMap[def.Key]
		dto := ToSettingsItemDTO(def, ts, os, categoryKey)
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
		defs, err := h.settingQueryRepo.FindByConfigurableAt(ctx, settingdomain.ScopeLevelTeam)
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
