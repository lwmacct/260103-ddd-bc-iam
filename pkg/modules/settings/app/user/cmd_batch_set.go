package user

import (
	"context"
	"fmt"

	settingdomain "github.com/lwmacct/260103-ddd-bc-settings/pkg/modules/settings/domain/setting"
	"github.com/lwmacct/260103-ddd-iam-bc/pkg/modules/settings/domain/user"
)

// BatchSetHandler 批量设置用户配置命令处理器
type BatchSetHandler struct {
	settingQueryRepo  settingdomain.QueryRepository
	categoryQueryRepo settingdomain.SettingCategoryQueryRepository
	cmdRepo           user.CommandRepository
}

// NewBatchSetHandler 创建批量设置命令处理器
func NewBatchSetHandler(
	settingQueryRepo settingdomain.QueryRepository,
	categoryQueryRepo settingdomain.SettingCategoryQueryRepository,
	cmdRepo user.CommandRepository,
) *BatchSetHandler {
	return &BatchSetHandler{
		settingQueryRepo:  settingQueryRepo,
		categoryQueryRepo: categoryQueryRepo,
		cmdRepo:           cmdRepo,
	}
}

// Handle 处理批量设置用户配置命令
//
// 流程：
//  1. 批量获取配置定义
//  2. 逐个校验值类型和格式
//  3. 批量 Upsert
func (h *BatchSetHandler) Handle(ctx context.Context, cmd BatchSetCommand) ([]SettingsItemDTO, error) {
	if len(cmd.Settings) == 0 {
		return []SettingsItemDTO{}, nil
	}

	// 1. 提取所有 keys
	keys := make([]string, len(cmd.Settings))
	for i, s := range cmd.Settings {
		keys[i] = s.Key
	}

	// 2. 批量获取配置定义
	defs, err := h.settingQueryRepo.FindByKeys(ctx, keys)
	if err != nil {
		return nil, fmt.Errorf("failed to find settings: %w", err)
	}

	// 构建定义映射
	defMap := make(map[string]*settingdomain.Setting)
	for _, def := range defs {
		defMap[def.Key] = def
	}

	// 3. 校验并构建用户配置列表
	userSettings := make([]*user.UserSetting, 0, len(cmd.Settings))
	for _, item := range cmd.Settings {
		def, ok := defMap[item.Key]
		if !ok {
			return nil, fmt.Errorf("%w: %s", user.ErrInvalidSettingKey, item.Key)
		}

		// 检查是否为用户可配置的配置
		if !def.IsUserScope() {
			return nil, fmt.Errorf("%w: %s is not a user-configurable setting", user.ErrInvalidSettingKey, item.Key)
		}

		// 类型校验
		if err := def.ValidateValue(item.Value); err != nil {
			return nil, fmt.Errorf("%w: %s - %w", user.ErrInvalidSettingValue, item.Key, err)
		}

		// 格式校验
		if err := def.ValidateByInputType(item.Value); err != nil {
			return nil, fmt.Errorf("%w: %s - %w", user.ErrValidationFailed, item.Key, err)
		}

		userSettings = append(userSettings, user.New(cmd.UserID, item.Key, item.Value))
	}

	// 4. 批量 Upsert
	if err := h.cmdRepo.BatchUpsert(ctx, userSettings); err != nil {
		return nil, fmt.Errorf("failed to batch save user settings: %w", err)
	}

	// 5. 获取所有分类元数据（用于填充 category key）
	categories, err := h.categoryQueryRepo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch categories: %w", err)
	}

	categoryKeyByID := make(map[uint]string, len(categories))
	for _, cat := range categories {
		categoryKeyByID[cat.ID] = cat.Key
	}

	// 6. 构建返回 DTO
	result := make([]SettingsItemDTO, len(userSettings))
	for i, us := range userSettings {
		def := defMap[us.SettingKey]
		categoryKey := categoryKeyByID[def.CategoryID]
		result[i] = *ToSettingsItemDTO(def, us, categoryKey)
	}

	return result, nil
}
