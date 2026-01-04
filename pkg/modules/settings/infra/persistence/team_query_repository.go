package persistence

import (
	"context"
	"errors"
	"fmt"

	"github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/settings/domain/team"
	"gorm.io/gorm"
)

// teamQueryRepository 团队配置查询仓储的 GORM 实现
type teamQueryRepository struct {
	db *gorm.DB
}

// NewTeamQueryRepository 创建团队配置查询仓储实例
func NewTeamQueryRepository(db *gorm.DB) team.QueryRepository {
	return &teamQueryRepository{db: db}
}

// FindByTeamAndKey 根据团队 ID 和键名查找团队配置
// 如果不存在返回 nil, nil
func (r *teamQueryRepository) FindByTeamAndKey(ctx context.Context, teamID uint, key string) (*team.TeamSetting, error) {
	var model TeamSettingModel
	err := r.db.WithContext(ctx).
		Where("team_id = ? AND setting_key = ?", teamID, key).
		First(&model).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			//nolint:nilnil // 返回 nil, nil 表示未找到，这是预期情况
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find team setting: %w", err)
	}

	return model.ToEntity(), nil
}

// FindByTeam 查找团队的所有自定义配置
func (r *teamQueryRepository) FindByTeam(ctx context.Context, teamID uint) ([]*team.TeamSetting, error) {
	var models []*TeamSettingModel
	err := r.db.WithContext(ctx).
		Where("team_id = ?", teamID).
		Order("setting_key ASC").
		Find(&models).Error

	if err != nil {
		return nil, fmt.Errorf("failed to find team settings: %w", err)
	}

	return toTeamEntities(models), nil
}

// FindByTeamAndKeys 批量查询团队的多个配置
func (r *teamQueryRepository) FindByTeamAndKeys(ctx context.Context, teamID uint, keys []string) ([]*team.TeamSetting, error) {
	if len(keys) == 0 {
		return []*team.TeamSetting{}, nil
	}

	var models []*TeamSettingModel
	err := r.db.WithContext(ctx).
		Where("team_id = ? AND setting_key IN ?", teamID, keys).
		Find(&models).Error

	if err != nil {
		return nil, fmt.Errorf("failed to find team settings by keys: %w", err)
	}

	return toTeamEntities(models), nil
}

// CountByTeam 统计团队的自定义配置数量
func (r *teamQueryRepository) CountByTeam(ctx context.Context, teamID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&TeamSettingModel{}).
		Where("team_id = ?", teamID).
		Count(&count).Error

	if err != nil {
		return 0, fmt.Errorf("failed to count team settings: %w", err)
	}

	return count, nil
}
