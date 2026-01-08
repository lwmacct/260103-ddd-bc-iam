package team

import (
	"context"
	"fmt"

	"github.com/lwmacct/260103-ddd-iam-bc/pkg/modules/settings/domain/team"
	settingdomain "github.com/lwmacct/260103-ddd-settings-bc/pkg/modules/settings/domain/setting"
)

// SetHandler 设置团队配置命令处理器
type SetHandler struct {
	settingQueryRepo  settingdomain.QueryRepository
	categoryQueryRepo settingdomain.SettingCategoryQueryRepository
	cmdRepo           team.CommandRepository
}

// NewSetHandler 创建设置命令处理器
func NewSetHandler(
	settingQueryRepo settingdomain.QueryRepository,
	categoryQueryRepo settingdomain.SettingCategoryQueryRepository,
	cmdRepo team.CommandRepository,
) *SetHandler {
	return &SetHandler{
		settingQueryRepo:  settingQueryRepo,
		categoryQueryRepo: categoryQueryRepo,
		cmdRepo:           cmdRepo,
	}
}

// Handle 处理设置团队配置命令
//
// 流程：
//  1. 校验配置定义存在（从 Settings BC）
//  2. 检查是否对 Team 可见
//  3. 检查是否允许 Team 配置
//  4. ValueType 类型校验
//  5. InputType 格式校验（email/url/password 等）
//  6. Upsert 团队配置
func (h *SetHandler) Handle(ctx context.Context, cmd SetCommand) (*SettingsItemDTO, error) {
	// 1. 校验配置定义存在
	def, err := h.settingQueryRepo.FindByKey(ctx, cmd.Key)
	if err != nil {
		return nil, fmt.Errorf("failed to find setting: %w", err)
	}
	if def == nil {
		return nil, team.ErrInvalidSettingKey
	}

	// 2. 检查是否对 Team 可见
	if !IsVisibleToTeam(def) {
		return nil, fmt.Errorf("%w: %s is not visible at team level", team.ErrSettingNotVisibleAtTeam, cmd.Key)
	}

	// 3. 检查是否允许 Team 配置
	if !IsConfigurableByTeam(def) {
		return nil, fmt.Errorf("%w: %s is not configurable at team level", team.ErrSettingNotConfigurableAtTeam, cmd.Key)
	}

	// 4. ValueType 类型校验
	if err := def.ValidateValue(cmd.Value); err != nil {
		return nil, fmt.Errorf("%w: %w", team.ErrInvalidSettingValue, err)
	}

	// 5. InputType 格式校验（email/url/password 等）
	if err := def.ValidateByInputType(cmd.Value); err != nil {
		return nil, fmt.Errorf("%w: %w", team.ErrValidationFailed, err)
	}

	// 6. Upsert 团队配置
	ts := team.New(cmd.TeamID, cmd.Key, cmd.Value)
	if err := h.cmdRepo.Upsert(ctx, ts); err != nil {
		return nil, fmt.Errorf("failed to save team setting: %w", err)
	}

	// 7. 获取 category key
	category, err := h.categoryQueryRepo.FindByID(ctx, def.CategoryID)
	if err != nil {
		return nil, fmt.Errorf("failed to find category: %w", err)
	}
	categoryKey := ""
	if category != nil {
		categoryKey = category.Key
	}

	// 返回时只包含团队自定义值
	return ToSettingsItemDTO(def, ts, nil, categoryKey), nil
}
