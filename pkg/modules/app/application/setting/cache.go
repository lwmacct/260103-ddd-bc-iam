package setting

import (
	"context"

	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/app/domain/setting"
)

// =========================================================================
// SettingsCacheService - Settings API 响应缓存
// =========================================================================

// SettingsCacheService Settings 响应缓存服务接口。
//
// 缓存 Setting Settings API 的最终响应 DTO，避免重复的数据库查询和构建逻辑。
// 采用 user+category 维度缓存，支持懒加载场景。
//
// # 架构说明
//
// 所有缓存接口统一定义在 Application 层：
//   - 缓存的是 Application 层 DTO（[SettingsCategoryDTO]）
//   - 遵循 DDD 依赖方向：Application 定义接口，Infrastructure 实现
//   - Domain 层不定义缓存接口
//
// # Key 命名规范
//
//   - 用户 Settings：{prefix}schema:user:{userID}:{categoryKey}
//   - 管理员 Settings：{prefix}schema:admin:{categoryKey}
//   - categoryKey 为空时使用 "_all" 表示全量
//
// 默认 TTL：30 分钟
//
// # 缓存失效策略
//
//   - 用户修改自己的设置 → [DeleteUserSettingsAll]
//   - 管理员修改系统设置 → [DeleteByCategoryKey]
//   - Category 结构变更 → [DeleteAll]
//
// 实现位于 [infrastructure/redis.settingsCacheService]。
//
//nolint:interfacebloat // 按职责分组：用户 Settings、管理员 Settings、批量操作、分类列表
type SettingsCacheService interface {
	// =========================================================================
	// 用户 Settings 操作
	// =========================================================================

	// GetUserSettings 获取用户 Settings 缓存。
	//
	// userID: 用户 ID
	// categoryKey: 分类键，为空表示全量 Schema
	//
	// 缓存未命中返回 nil, nil。
	// 缓存数据损坏时自动清除并返回 nil, nil。
	GetUserSettings(ctx context.Context, userID uint, categoryKey string) ([]SettingsCategoryDTO, error)

	// SetUserSettings 设置用户 Settings 缓存。
	//
	// schema 为空时也会缓存（防止缓存穿透）。
	SetUserSettings(ctx context.Context, userID uint, categoryKey string, schema []SettingsCategoryDTO) error

	// DeleteUserSettings 删除用户的指定 category Settings 缓存。
	DeleteUserSettings(ctx context.Context, userID uint, categoryKey string) error

	// DeleteUserSettingsAll 删除用户的所有 Settings 缓存。
	//
	// 使用 SCAN 遍历 {prefix}schema:user:{userID}:* 模式。
	// 用于用户修改设置后失效所有相关缓存。
	DeleteUserSettingsAll(ctx context.Context, userID uint) error

	// =========================================================================
	// 管理员 Settings 操作
	// =========================================================================

	// GetAdminSettings 获取管理员 Settings 缓存。
	//
	// categoryKey: 分类键，为空表示全量 Schema
	//
	// 缓存未命中返回 nil, nil。
	GetAdminSettings(ctx context.Context, categoryKey string) ([]SettingsCategoryDTO, error)

	// SetAdminSettings 设置管理员 Settings 缓存。
	SetAdminSettings(ctx context.Context, categoryKey string, schema []SettingsCategoryDTO) error

	// DeleteAdminSettings 删除管理员的指定 category Settings 缓存。
	DeleteAdminSettings(ctx context.Context, categoryKey string) error

	// DeleteAdminSettingsAll 删除管理员的所有 Settings 缓存。
	//
	// 使用 SCAN 遍历 {prefix}schema:admin:* 模式。
	DeleteAdminSettingsAll(ctx context.Context) error

	// =========================================================================
	// 批量失效操作
	// =========================================================================

	// DeleteByCategoryKey 删除所有用户和管理员的指定 category Settings 缓存。
	//
	// 当系统设置定义变更时调用，使所有相关 Settings 缓存失效。
	// 使用 SCAN 遍历 {prefix}schema:*:{categoryKey} 模式。
	//
	// 注意：这是低频操作（仅管理员修改系统设置时触发），
	// SCAN 遍历的性能影响可接受。
	DeleteByCategoryKey(ctx context.Context, categoryKey string) error

	// DeleteAll 删除所有 Settings 缓存。
	//
	// 当 Category 结构变更或执行数据迁移时调用。
	// 使用 SCAN 遍历 {prefix}schema:* 模式。
	DeleteAll(ctx context.Context) error

	// =========================================================================
	// 分类列表缓存操作
	// =========================================================================

	// GetUserCategories 获取用户分类列表缓存。
	//
	// 缓存的是 scope="user" 的分类元信息列表，系统级数据（不区分用户）。
	// Key 格式：{prefix}schema:categories:user
	//
	// 缓存未命中返回 nil, nil。
	GetUserCategories(ctx context.Context) ([]CategoryMetaDTO, error)

	// SetUserCategories 设置用户分类列表缓存。
	SetUserCategories(ctx context.Context, categories []CategoryMetaDTO) error

	// DeleteUserCategories 删除用户分类列表缓存。
	//
	// 当 Category 结构变更时调用（与 [DeleteAll] 关联）。
	DeleteUserCategories(ctx context.Context) error

	// =========================================================================
	// Category 实体缓存操作（供 Repository 装饰器使用）
	// =========================================================================

	// GetAllCategories 获取所有 SettingCategory 实体缓存。
	//
	// Key 格式：{prefix}schema:categories:all
	// 缓存完整的领域实体（包含 ID, Key, Label, Icon, Order）。
	// 缓存未命中返回 nil, nil。
	GetAllCategories(ctx context.Context) ([]*setting.SettingCategory, error)

	// SetAllCategories 设置所有 SettingCategory 实体缓存。
	SetAllCategories(ctx context.Context, categories []*setting.SettingCategory) error

	// DeleteAllCategories 删除 SettingCategory 实体缓存。
	DeleteAllCategories(ctx context.Context) error
}

// =========================================================================
// UserSettingCacheService - 用户有效设置值缓存
// =========================================================================

// UserSettingCacheService 用户设置缓存服务接口。
//
// 存储用户的有效设置值（合并后的结果），而非原始 UserSetting 记录。
// 这样查询时无需再合并 Setting.DefaultValue 和 UserSetting.Value。
//
// Key 格式：{prefix}user:{userID}:setting:{key}
// 默认 TTL：30 分钟（比系统设置缓存更长，因为用户设置变更频率较低）
//
// 多实例安全：
//   - 删除操作直接生效，无需跨实例通知（无本地缓存）
//
// 实现位于 [infrastructure/cache.userSettingCacheService]。
type UserSettingCacheService interface {
	// =========================================================================
	// 单条操作
	// =========================================================================

	// Get 获取用户的有效设置值。
	//
	// userID: 用户 ID
	// key: 设置 key（如 "general.theme"）
	//
	// 缓存未命中返回 nil, nil。
	// 缓存数据损坏时自动清除并返回 nil, nil。
	Get(ctx context.Context, userID uint, key string) (*EffectiveUserSetting, error)

	// Set 缓存用户的有效设置值。
	Set(ctx context.Context, userID uint, value *EffectiveUserSetting) error

	// =========================================================================
	// 批量操作
	// =========================================================================

	// GetByKeys 批量获取用户的有效设置值。
	//
	// 使用 MGET 一次网络往返获取多个 key。
	// 返回 key -> *EffectiveUserSetting 映射，未命中的 key 不在结果中。
	GetByKeys(ctx context.Context, userID uint, keys []string) (map[string]*EffectiveUserSetting, error)

	// SetBatch 批量设置用户的有效设置值。
	//
	// 使用 Pipeline 批量写入。
	// values 为空时直接返回 nil。
	SetBatch(ctx context.Context, userID uint, values []*EffectiveUserSetting) error

	// =========================================================================
	// 删除操作
	// =========================================================================

	// Delete 删除用户的指定设置缓存。
	Delete(ctx context.Context, userID uint, key string) error

	// DeleteByKeys 批量删除用户的指定设置缓存。
	DeleteByKeys(ctx context.Context, userID uint, keys []string) error

	// DeleteByUser 删除用户的所有设置缓存。
	//
	// 使用 SCAN 遍历 {prefix}user:{userID}:setting:* 模式。
	// 用于用户重置所有设置或用户删除场景。
	DeleteByUser(ctx context.Context, userID uint) error

	// DeleteBySettingKey 删除所有用户的某个设置缓存。
	//
	// 当系统默认值变更时调用，使所有用户的该设置缓存失效。
	// 使用 SCAN 遍历 {prefix}user:*:setting:{key} 模式。
	//
	// 注意：这是低频操作（仅管理员修改系统设置时触发），
	// SCAN 遍历的性能影响可接受。
	DeleteBySettingKey(ctx context.Context, key string) error

	// DeleteBySettingKeys 批量删除所有用户的多个设置缓存。
	//
	// 当批量修改系统默认值时调用。
	DeleteBySettingKeys(ctx context.Context, keys []string) error
}

// EffectiveUserSetting 用户有效设置值（缓存数据结构）。
//
// 存储合并后的实际生效值，避免查询时再次合并 Setting + UserSetting。
// 包含 UI 渲染所需的元数据，前端可直接使用。
type EffectiveUserSetting struct {
	// Key 设置键
	Key string `json:"key"`

	// Value 实际生效值（用户值或系统默认值）
	Value any `json:"value"`

	// DefaultValue 系统默认值（用于判断是否自定义、重置功能）
	DefaultValue any `json:"default_value"`

	// IsCustomized 是否被用户自定义
	// true: Value 来自 UserSetting
	// false: Value 等于 DefaultValue
	IsCustomized bool `json:"is_customized"`

	// =========================================================================
	// UI 元数据（透传给前端）
	// =========================================================================

	// ValueType 值类型：string, number, boolean, json
	ValueType string `json:"value_type"`

	// CategoryID 所属分类 ID
	CategoryID uint `json:"category_id"`

	// Group 分类内子分组
	Group string `json:"group"`

	// Label 显示标签
	Label string `json:"label"`

	// UIConfig UI 配置（JSON 字符串）
	UIConfig string `json:"ui_config,omitempty"`

	// Order 排序权重
	Order int `json:"order"`
}

// =========================================================================
// UserSettingQueryCacheService - 用户设置查询缓存（Repository 层）
// =========================================================================

// UserSettingQueryCacheService 用户设置查询缓存接口。
//
// 存储原始 UserSetting 记录，用于 Repository 层减少数据库查询。
// 采用用户维度全量缓存策略：一次查询缓存用户所有自定义配置。
//
// 与 [UserSettingCacheService] 的区别：
//   - 本接口：存储原始 UserSetting（Domain 实体），Repository 层使用
//   - UserSettingCacheService：存储合并后的 EffectiveUserSetting（DTO），Application 层使用
//
// Key 格式：{prefix}usersetting:user:{userID}
// 默认 TTL：30 分钟
type UserSettingQueryCacheService interface {
	// GetByUser 获取用户的所有自定义配置缓存。
	// 返回 map[settingKey]*UserSetting，未命中返回 nil, nil。
	// 缓存数据损坏时自动清除并返回 nil, nil。
	GetByUser(ctx context.Context, userID uint) (map[string]*setting.UserSetting, error)

	// SetByUser 设置用户的所有自定义配置缓存。
	// settings 为空时缓存空结果（防止缓存穿透）。
	SetByUser(ctx context.Context, userID uint, settings []*setting.UserSetting) error

	// DeleteByUser 删除用户的所有配置缓存。
	DeleteByUser(ctx context.Context, userID uint) error
}
