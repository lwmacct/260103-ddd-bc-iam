package persistence

import (
	"context"
	"errors"
	"fmt"

	"github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/domain/pat"
	"gorm.io/gorm"
)

// patQueryRepository PAT 查询仓储的 GORM 实现
type patQueryRepository struct {
	db *gorm.DB
}

// NewPATQueryRepository 创建 PAT 查询仓储实例
func NewPATQueryRepository(db *gorm.DB) pat.QueryRepository {
	return &patQueryRepository{db: db}
}

// FindByToken 通过令牌哈希查找（用于认证）
func (r *patQueryRepository) FindByToken(ctx context.Context, tokenHash string) (*pat.PersonalAccessToken, error) {
	var model PersonalAccessTokenModel
	err := r.db.WithContext(ctx).
		Where("token = ?", tokenHash).
		First(&model).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("token not found")
		}
		return nil, fmt.Errorf("failed to find PAT by token: %w", err)
	}

	return model.ToEntity(), nil
}

// FindByID 通过 ID 查找令牌
func (r *patQueryRepository) FindByID(ctx context.Context, id uint) (*pat.PersonalAccessToken, error) {
	var model PersonalAccessTokenModel
	err := r.db.WithContext(ctx).
		First(&model, id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("token not found")
		}
		return nil, fmt.Errorf("failed to find PAT by ID: %w", err)
	}

	return model.ToEntity(), nil
}

// FindByPrefix 通过前缀查找令牌
func (r *patQueryRepository) FindByPrefix(ctx context.Context, prefix string) (*pat.PersonalAccessToken, error) {
	var model PersonalAccessTokenModel
	err := r.db.WithContext(ctx).
		Where("token_prefix = ?", prefix).
		First(&model).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("token not found")
		}
		return nil, fmt.Errorf("failed to find PAT by prefix: %w", err)
	}

	return model.ToEntity(), nil
}

// ListByUser 获取指定用户的所有令牌
func (r *patQueryRepository) ListByUser(ctx context.Context, userID uint) ([]*pat.PersonalAccessToken, error) {
	var models []PersonalAccessTokenModel
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&models).Error

	if err != nil {
		return nil, fmt.Errorf("failed to list PATs by user: %w", err)
	}

	return mapPATModelsToEntities(models), nil
}
