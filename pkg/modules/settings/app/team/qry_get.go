package team

import (
	"context"
	"fmt"

	"github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/settings/domain/org"
	"github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/settings/domain/team"
	settingdomain "github.com/lwmacct/260103-ddd-bc-settings/pkg/modules/settings/domain/setting"
)

// GetHandler 获取单个配置查询处理器
type GetHandler struct {
	settingQueryRepo  settingdomain.QueryRepository
	categoryQueryRepo settingdomain.SettingCategoryQueryRepository
	teamQueryRepo     team.QueryRepository
	orgQueryRepo      org.QueryRepository // 用于继承查询
}

// NewGetHandler 创建获取配置查询处理器
func NewGetHandler(
	settingQueryRepo settingdomain.QueryRepository,
	categoryQueryRepo settingdomain.SettingCategoryQueryRepository,
	teamQueryRepo team.QueryRepository,
	orgQueryRepo org.QueryRepository,
) *GetHandler {
	return &GetHandler{
		settingQueryRepo:  settingQueryRepo,
		categoryQueryRepo: categoryQueryRepo,
		teamQueryRepo:     teamQueryRepo,
		orgQueryRepo:      orgQueryRepo,
	}
}

// Handle 处理获取单个配置查询
//
// 返回合并后的配置，优先级：团队配置 > 组织配置 > 系统默认值
func (h *GetHandler) Handle(ctx context.Context, query GetQuery) (*SettingsItemDTO, error) {
	// 1. 获取配置定义
	def, err := h.settingQueryRepo.FindByKey(ctx, query.Key)
	if err != nil {
		return nil, fmt.Errorf("failed to find setting: %w", err)
	}
	if def == nil {
		return nil, team.ErrInvalidSettingKey
	}

	// 2. 获取团队自定义值（可能为 nil）
	ts, _ := h.teamQueryRepo.FindByTeamAndKey(ctx, query.TeamID, query.Key)

	// 3. 获取组织配置（可能为 nil）
	var os *org.OrgSetting
	if ts == nil || ts.IsEmpty() {
		// 只有团队没有自定义值时才回退到组织配置
		os, _ = h.orgQueryRepo.FindByOrgAndKey(ctx, query.OrgID, query.Key)
	}

	// 4. 获取 category key
	category, err := h.categoryQueryRepo.FindByID(ctx, def.CategoryID)
	if err != nil {
		return nil, fmt.Errorf("failed to find category: %w", err)
	}
	categoryKey := ""
	if category != nil {
		categoryKey = category.Key
	}

	// 5. 合并返回
	return ToSettingsItemDTO(def, ts, os, categoryKey), nil
}
