package setting

import (
	"context"
)

// SettingCategoryQueryRepository 配置分类查询仓储接口。
//
// 提供 SettingCategory 的只读查询操作。
type SettingCategoryQueryRepository interface {
	// FindByID 根据 ID 查询分类。
	//
	// 如果未找到，返回 nil 和 nil（无错误）。
	FindByID(ctx context.Context, id uint) (*SettingCategory, error)

	// FindByKey 根据分类键查询。
	//
	// 如果未找到，返回 nil 和 nil（无错误）。
	FindByKey(ctx context.Context, key string) (*SettingCategory, error)

	// FindByIDs 根据 ID 列表批量查询分类，按 Order 升序排列。
	FindByIDs(ctx context.Context, ids []uint) ([]*SettingCategory, error)

	// FindAll 查询所有分类，按 Order 升序排列。
	FindAll(ctx context.Context) ([]*SettingCategory, error)

	// ExistsByKey 检查指定 Key 是否已存在。
	ExistsByKey(ctx context.Context, key string) (bool, error)
}
