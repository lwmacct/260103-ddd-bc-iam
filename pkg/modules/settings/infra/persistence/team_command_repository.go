package persistence

import (
	"context"
	"fmt"

	"github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/settings/domain/team"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// teamCommandRepository 团队配置命令仓储的 GORM 实现
type teamCommandRepository struct {
	db *gorm.DB
}

// NewTeamCommandRepository 创建团队配置命令仓储实例
func NewTeamCommandRepository(db *gorm.DB) team.CommandRepository {
	return &teamCommandRepository{db: db}
}

// Upsert 创建或更新团队配置（基于 team_id + setting_key 唯一约束）
func (r *teamCommandRepository) Upsert(ctx context.Context, setting *team.TeamSetting) error {
	model := newTeamModelFromEntity(setting)
	err := r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "team_id"}, {Name: "setting_key"}},
		DoUpdates: clause.AssignmentColumns([]string{"value", "updated_at"}),
	}).Create(model).Error

	if err != nil {
		return fmt.Errorf("failed to upsert team setting: %w", err)
	}

	// 回写生成的 ID
	setting.ID = model.ID
	return nil
}

// Delete 删除指定团队的指定配置
func (r *teamCommandRepository) Delete(ctx context.Context, teamID uint, key string) error {
	result := r.db.WithContext(ctx).
		Where("team_id = ? AND setting_key = ?", teamID, key).
		Delete(&TeamSettingModel{})

	if result.Error != nil {
		return fmt.Errorf("failed to delete team setting: %w", result.Error)
	}
	return nil
}

// DeleteByTeam 删除指定团队的所有配置
func (r *teamCommandRepository) DeleteByTeam(ctx context.Context, teamID uint) error {
	result := r.db.WithContext(ctx).
		Where("team_id = ?", teamID).
		Delete(&TeamSettingModel{})

	if result.Error != nil {
		return fmt.Errorf("failed to delete team settings: %w", result.Error)
	}
	return nil
}
