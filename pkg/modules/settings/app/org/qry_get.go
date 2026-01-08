package org

import (
	"context"
	"fmt"

	"github.com/lwmacct/260103-ddd-iam-bc/pkg/modules/settings/domain/org"
	settingdomain "github.com/lwmacct/260103-ddd-settings-bc/pkg/modules/settings/domain/setting"
)

// GetHandler 获取单个配置查询处理器
type GetHandler struct {
	settingQueryRepo  settingdomain.QueryRepository
	categoryQueryRepo settingdomain.SettingCategoryQueryRepository
	queryRepo         org.QueryRepository
}

// NewGetHandler 创建获取配置查询处理器
func NewGetHandler(
	settingQueryRepo settingdomain.QueryRepository,
	categoryQueryRepo settingdomain.SettingCategoryQueryRepository,
	queryRepo org.QueryRepository,
) *GetHandler {
	return &GetHandler{
		settingQueryRepo:  settingQueryRepo,
		categoryQueryRepo: categoryQueryRepo,
		queryRepo:         queryRepo,
	}
}

// Handle 处理获取单个配置查询
//
// 返回合并后的配置：如果组织有自定义值则使用组织值，否则使用默认值
func (h *GetHandler) Handle(ctx context.Context, query GetQuery) (*SettingsItemDTO, error) {
	// 1. 获取配置定义
	def, err := h.settingQueryRepo.FindByKey(ctx, query.Key)
	if err != nil {
		return nil, fmt.Errorf("failed to find setting: %w", err)
	}
	if def == nil {
		return nil, org.ErrInvalidSettingKey
	}

	// 2. 获取组织自定义值（可能为 nil）
	os, err := h.queryRepo.FindByOrgAndKey(ctx, query.OrgID, query.Key)
	if err != nil {
		return nil, fmt.Errorf("failed to find org setting: %w", err)
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
	return ToSettingsItemDTO(def, os, categoryKey), nil
}
