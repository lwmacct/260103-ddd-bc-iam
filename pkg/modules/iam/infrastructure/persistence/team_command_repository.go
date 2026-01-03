package persistence

import (
	"context"

	"github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/domain/org"
	"gorm.io/gorm"
)

// teamCommandRepository 团队命令仓储的 GORM 实现
type teamCommandRepository struct {
	db *gorm.DB
}

// NewTeamCommandRepository 创建团队命令仓储实例
func NewTeamCommandRepository(db *gorm.DB) org.TeamCommandRepository {
	return &teamCommandRepository{db: db}
}

// Create 创建团队
func (r *teamCommandRepository) Create(ctx context.Context, team *org.Team) error {
	model := newTeamModelFromEntity(team)
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return err
	}
	team.ID = model.ID
	return nil
}

// Update 更新团队
func (r *teamCommandRepository) Update(ctx context.Context, team *org.Team) error {
	model := newTeamModelFromEntity(team)
	return r.db.WithContext(ctx).Save(model).Error
}

// Delete 删除团队
func (r *teamCommandRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&TeamModel{}, id).Error
}
