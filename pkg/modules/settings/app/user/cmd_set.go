package user

import (
	"context"
	"fmt"

	"github.com/lwmacct/260103-ddd-iam-bc/pkg/modules/settings/domain/user"
	settingdomain "github.com/lwmacct/260103-ddd-settings-bc/pkg/modules/settings/domain/setting"
)

// SetHandler 设置用户配置命令处理器
type SetHandler struct {
	settingQueryRepo  settingdomain.QueryRepository // 跨 BC 依赖：Settings BC
	categoryQueryRepo settingdomain.SettingCategoryQueryRepository
	cmdRepo           user.CommandRepository
}

// NewSetHandler 创建设置命令处理器
func NewSetHandler(
	settingQueryRepo settingdomain.QueryRepository,
	categoryQueryRepo settingdomain.SettingCategoryQueryRepository,
	cmdRepo user.CommandRepository,
) *SetHandler {
	return &SetHandler{
		settingQueryRepo:  settingQueryRepo,
		categoryQueryRepo: categoryQueryRepo,
		cmdRepo:           cmdRepo,
	}
}

// Handle 处理设置用户配置命令
//
// 流程：
//  1. 校验配置定义存在（从 Settings BC）
//  2. ValueType 类型校验
//  3. InputType 格式校验（email/url/password 等）
//  4. Upsert 用户配置
func (h *SetHandler) Handle(ctx context.Context, cmd SetCommand) (*SettingsItemDTO, error) {
	// 1. 校验配置定义存在
	def, err := h.settingQueryRepo.FindByKey(ctx, cmd.Key)
	if err != nil {
		return nil, fmt.Errorf("failed to find setting: %w", err)
	}
	if def == nil {
		return nil, user.ErrInvalidSettingKey
	}

	// 2. 检查是否为用户可配置的配置
	if !def.IsUserScope() {
		return nil, fmt.Errorf("%w: %s is not a user-configurable setting", user.ErrInvalidSettingKey, cmd.Key)
	}

	// 3. ValueType 类型校验
	if err := def.ValidateValue(cmd.Value); err != nil {
		return nil, fmt.Errorf("%w: %w", user.ErrInvalidSettingValue, err)
	}

	// 4. InputType 格式校验（email/url/password 等）
	if err := def.ValidateByInputType(cmd.Value); err != nil {
		return nil, fmt.Errorf("%w: %w", user.ErrValidationFailed, err)
	}

	// 5. Upsert 用户配置
	us := user.New(cmd.UserID, cmd.Key, cmd.Value)
	if err := h.cmdRepo.Upsert(ctx, us); err != nil {
		return nil, fmt.Errorf("failed to save user setting: %w", err)
	}

	// 6. 获取 category key
	category, err := h.categoryQueryRepo.FindByID(ctx, def.CategoryID)
	if err != nil {
		return nil, fmt.Errorf("failed to find category: %w", err)
	}
	categoryKey := ""
	if category != nil {
		categoryKey = category.Key
	}

	return ToSettingsItemDTO(def, us, categoryKey), nil
}
