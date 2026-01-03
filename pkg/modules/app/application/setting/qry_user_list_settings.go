package setting

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/app/domain/setting"
)

// UserListSettingsHandler 获取用户配置 Settings 查询处理器
type UserListSettingsHandler struct {
	settingQueryRepo  setting.QueryRepository
	queryRepo         setting.UserSettingQueryRepository
	categoryQueryRepo setting.SettingCategoryQueryRepository
	settingsCache     SettingsCacheService
}

// NewUserListSettingsHandler 创建 UserListSettingsHandler 实例
func NewUserListSettingsHandler(
	settingQueryRepo setting.QueryRepository,
	queryRepo setting.UserSettingQueryRepository,
	categoryQueryRepo setting.SettingCategoryQueryRepository,
	settingsCache SettingsCacheService,
) *UserListSettingsHandler {
	return &UserListSettingsHandler{
		settingQueryRepo:  settingQueryRepo,
		queryRepo:         queryRepo,
		categoryQueryRepo: categoryQueryRepo,
		settingsCache:     settingsCache,
	}
}

// Handle 处理获取用户配置 Settings 查询
// 返回按 Category → Group → Settings 层级组织的数据，包含用户自定义值
//
// 返回用户可见的设置项：
//   - scope="user" 的设置（用户可编辑）
//   - scope="system" 且 public=true 的设置（只读，用于依赖检查和默认值展示）
//
// 支持 CategoryKey 过滤：
//   - 为空时返回全量用户可见设置（用于总配置页）
//   - 指定 Key 时只返回该分类（用于分散页面的懒加载）
//
// 缓存策略：
//   - 先查缓存，命中直接返回
//   - 未命中时查数据库，同步回写缓存
func (h *UserListSettingsHandler) Handle(ctx context.Context, query UserListSettingsQuery) ([]SettingsCategoryDTO, error) {
	// 1. 查缓存
	if cached, err := h.settingsCache.GetUserSettings(ctx, query.UserID, query.CategoryKey); err == nil && cached != nil {
		return cached, nil
	}

	// 2. 缓存未命中，执行原有逻辑
	result, err := h.fetchAndBuild(ctx, query)
	if err != nil {
		return nil, err
	}

	// 3. 同步回写缓存（仅非空结果，避免缓存无效数据）
	if len(result) > 0 {
		if err := h.settingsCache.SetUserSettings(ctx, query.UserID, query.CategoryKey, result); err != nil {
			slog.Warn("failed to cache user settings", "userID", query.UserID, "categoryKey", query.CategoryKey, "error", err.Error())
		}
	}

	return result, nil
}

// fetchAndBuild 从数据库获取数据并构建 Settings
func (h *UserListSettingsHandler) fetchAndBuild(ctx context.Context, query UserListSettingsQuery) ([]SettingsCategoryDTO, error) {
	// 1. 根据 CategoryKey 决定查询范围（只查询 user scope）
	defs, err := h.fetchUserSettings(ctx, query.CategoryKey)
	if err != nil {
		return nil, err
	}

	// 2. 查找用户的所有自定义配置
	userSettings, err := h.queryRepo.FindByUser(ctx, query.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to find user settings: %w", err)
	}

	// 3. 查询所有分类元数据
	categories, err := h.categoryQueryRepo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch setting categories: %w", err)
	}

	// 4. 构建用户配置映射
	userSettingMap := make(map[string]*setting.UserSetting, len(userSettings))
	for _, us := range userSettings {
		userSettingMap[us.SettingKey] = us
	}

	// 5. 使用共享构建器
	builder := NewSettingsBuilder(categories)
	return builder.Build(defs, userSettingMap, UserSettingMapper), nil
}

// fetchUserSettings 根据 CategoryKey 获取用户可见的设置列表
//
// 返回 scope=user（可编辑）+ scope=system 且 public=true（只读）的设置
func (h *UserListSettingsHandler) fetchUserSettings(ctx context.Context, categoryKey string) ([]*setting.Setting, error) {
	// 全量查询用户可见的设置
	if categoryKey == "" {
		defs, err := h.settingQueryRepo.FindVisibleToUser(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to find visible settings: %w", err)
		}
		return defs, nil
	}

	// 按 Category Key 过滤
	category, err := h.categoryQueryRepo.FindByKey(ctx, categoryKey)
	if err != nil {
		return nil, fmt.Errorf("failed to find category by key %q: %w", categoryKey, err)
	}
	if category == nil {
		return nil, fmt.Errorf("category not found: %s", categoryKey)
	}

	allDefs, err := h.settingQueryRepo.FindByCategoryID(ctx, category.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to find setting definitions: %w", err)
	}

	// 过滤只保留用户可见的设置
	return filterVisibleToUserSettings(allDefs), nil
}

// filterVisibleToUserSettings 过滤只保留用户可见的设置
// 包含 scope=user 和 scope=system 且 public=true
func filterVisibleToUserSettings(settings []*setting.Setting) []*setting.Setting {
	result := make([]*setting.Setting, 0, len(settings))
	for _, s := range settings {
		if s.IsVisibleToUser() {
			result = append(result, s)
		}
	}
	return result
}
