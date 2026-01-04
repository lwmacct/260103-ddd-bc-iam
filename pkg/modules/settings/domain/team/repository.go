package team

import "context"

// CommandRepository 团队配置命令仓储接口
// 负责所有修改状态的操作
type CommandRepository interface {
	// Upsert 创建或更新团队配置（基于 team_id + setting_key 唯一约束）
	Upsert(ctx context.Context, setting *TeamSetting) error

	// Delete 删除指定团队的指定配置
	Delete(ctx context.Context, teamID uint, key string) error

	// DeleteByTeam 删除指定团队的所有配置
	DeleteByTeam(ctx context.Context, teamID uint) error
}

// QueryRepository 团队配置查询仓储接口
// 负责所有只读查询操作
type QueryRepository interface {
	// FindByTeamAndKey 根据团队 ID 和键名查找团队配置
	// 如果不存在返回 nil, nil
	FindByTeamAndKey(ctx context.Context, teamID uint, key string) (*TeamSetting, error)

	// FindByTeam 查找团队的所有自定义配置
	FindByTeam(ctx context.Context, teamID uint) ([]*TeamSetting, error)

	// FindByTeamAndKeys 批量查询团队的多个配置
	FindByTeamAndKeys(ctx context.Context, teamID uint, keys []string) ([]*TeamSetting, error)

	// CountByTeam 统计团队的自定义配置数量
	CountByTeam(ctx context.Context, teamID uint) (int64, error)
}
