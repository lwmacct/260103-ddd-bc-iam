package org

import "context"

// CommandRepository 组织配置命令仓储接口
// 负责所有修改状态的操作
type CommandRepository interface {
	// Upsert 创建或更新组织配置（基于 org_id + setting_key 唯一约束）
	Upsert(ctx context.Context, setting *OrgSetting) error

	// Delete 删除指定组织的指定配置
	Delete(ctx context.Context, orgID uint, key string) error

	// DeleteByOrg 删除指定组织的所有配置
	DeleteByOrg(ctx context.Context, orgID uint) error
}

// QueryRepository 组织配置查询仓储接口
// 负责所有只读查询操作
type QueryRepository interface {
	// FindByOrgAndKey 根据组织 ID 和键名查找组织配置
	// 如果不存在返回 nil, nil
	FindByOrgAndKey(ctx context.Context, orgID uint, key string) (*OrgSetting, error)

	// FindByOrg 查找组织的所有自定义配置
	FindByOrg(ctx context.Context, orgID uint) ([]*OrgSetting, error)

	// FindByOrgAndKeys 批量查询组织的多个配置
	FindByOrgAndKeys(ctx context.Context, orgID uint, keys []string) ([]*OrgSetting, error)

	// CountByOrg 统计组织的自定义配置数量
	CountByOrg(ctx context.Context, orgID uint) (int64, error)
}
