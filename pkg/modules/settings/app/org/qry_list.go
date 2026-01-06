package org

import (
	"context"
	"fmt"

	"github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/settings/domain/org"
	settingdomain "github.com/lwmacct/260103-ddd-bc-settings/pkg/modules/settings/domain/setting"
)

// ListHandler 获取组织配置列表查询处理器
type ListHandler struct {
	settingQueryRepo  settingdomain.QueryRepository
	categoryQueryRepo settingdomain.SettingCategoryQueryRepository
	queryRepo         org.QueryRepository
}

// NewListHandler 创建获取配置列表查询处理器
func NewListHandler(
	settingQueryRepo settingdomain.QueryRepository,
	categoryQueryRepo settingdomain.SettingCategoryQueryRepository,
	queryRepo org.QueryRepository,
) *ListHandler {
	return &ListHandler{
		settingQueryRepo:  settingQueryRepo,
		categoryQueryRepo: categoryQueryRepo,
		queryRepo:         queryRepo,
	}
}

// Handle 处理获取组织配置列表查询
//
// 返回合并后的配置列表：如果组织有自定义值则使用组织值，否则使用默认值
func (h *ListHandler) Handle(ctx context.Context, query ListQuery) ([]*OrgSettingDTO, error) {
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
		// 获取所有用户可配置的设置
		defs, err = h.settingQueryRepo.FindByVisibleAt(ctx, settingdomain.ScopeLevelUser)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find settings: %w", err)
	}

	if len(defs) == 0 {
		return []*OrgSettingDTO{}, nil
	}

	// 2. 获取组织自定义值
	orgSettings, err := h.queryRepo.FindByOrg(ctx, query.OrgID)
	if err != nil {
		return nil, fmt.Errorf("failed to find org settings: %w", err)
	}

	// 3. 构建组织配置映射
	orgMap := make(map[string]*org.OrgSetting)
	for _, os := range orgSettings {
		orgMap[os.SettingKey] = os
	}

	// 4. 合并返回
	result := make([]*OrgSettingDTO, 0, len(defs))
	for _, def := range defs {
		os := orgMap[def.Key]
		result = append(result, ToOrgSettingDTO(def, os))
	}

	return result, nil
}
