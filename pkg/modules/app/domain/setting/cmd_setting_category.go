package setting

import (
	"context"
)

// SettingCategoryCommandRepository 配置分类写仓储接口。
//
// 提供 SettingCategory 的写操作（Create/Update/Delete）。
type SettingCategoryCommandRepository interface {
	// Create 创建配置分类。
	//
	// 成功后会回写生成的 ID 到 category.ID。
	Create(ctx context.Context, category *SettingCategory) error

	// Update 更新配置分类。
	//
	// 仅更新 Label、Icon、Order 字段，Key 不可修改。
	Update(ctx context.Context, category *SettingCategory) error

	// Delete 删除配置分类。
	//
	// 调用前应检查是否有关联的 Setting 记录。
	Delete(ctx context.Context, id uint) error
}
