package team

import (
	"context"
	"fmt"

	"github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/settings/domain/org"
	"github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/settings/domain/team"
	settingdomain "github.com/lwmacct/260103-ddd-bc-settings/pkg/modules/settings/domain/setting"
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
// 返回合并后的配置列表，优先级：团队配置 > 组织配置 > 系统默认值
func (h *ListHandler) Handle(ctx context.Context, query ListQuery) ([]*TeamSettingDTO, error) {
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
		// 获取所有 Team 可配置的设置
		defs, err = h.settingQueryRepo.FindByConfigurableAt(ctx, settingdomain.ScopeLevelTeam)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find settings: %w", err)
	}

	if len(defs) == 0 {
		return []*TeamSettingDTO{}, nil
	}

	// 2. 确保只返回对 Team 可见的设置
	defs = FilterByConfigurableAt(defs, settingdomain.ScopeLevelTeam)

	if len(defs) == 0 {
		return []*TeamSettingDTO{}, nil
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

	// 5. 合并返回
	result := make([]*TeamSettingDTO, 0, len(defs))
	for _, def := range defs {
		ts := teamMap[def.Key]
		os := orgMap[def.Key]
		result = append(result, ToTeamSettingDTO(def, ts, os))
	}

	return result, nil
}
