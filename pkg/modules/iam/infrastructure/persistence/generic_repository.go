package persistence

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

// Model 定义 GORM Model 必须实现的接口，用于转换回 Domain Entity。
// 泛型参数 E 是对应的 Domain Entity 类型。
type Model[E any] interface {
	// ToEntity 将 GORM Model 转换为 Domain Entity
	ToEntity() *E
}

// EntityToModel 定义从 Domain Entity 创建 GORM Model 的函数类型。
type EntityToModel[E any, M any] func(*E) M

// NotFoundError 定义"未找到"错误的接口，用于返回 Domain 层定义的错误。
type NotFoundError interface {
	error
	IsNotFound() bool
}

// GenericCommandRepository 提供泛型的写操作仓储实现。
// 泛型参数：
//   - E: Domain Entity 类型
//   - M: GORM Model 类型（必须实现 Model[E] 接口）
type GenericCommandRepository[E any, M Model[E]] struct {
	db            *gorm.DB
	entityToModel EntityToModel[E, M]
}

// NewGenericCommandRepository 创建泛型写仓储实例。
func NewGenericCommandRepository[E any, M Model[E]](
	db *gorm.DB,
	entityToModel EntityToModel[E, M],
) *GenericCommandRepository[E, M] {
	return &GenericCommandRepository[E, M]{
		db:            db,
		entityToModel: entityToModel,
	}
}

// Create 创建实体。创建成功后，entity 会被更新为包含数据库生成的字段（如 ID）。
func (r *GenericCommandRepository[E, M]) Create(ctx context.Context, entity *E) error {
	model := r.entityToModel(entity)
	if err := r.db.WithContext(ctx).Create(&model).Error; err != nil {
		return fmt.Errorf("failed to create: %w", err)
	}
	// 将数据库返回的值（如自增 ID、默认值）同步回 entity
	if saved := model.ToEntity(); saved != nil {
		*entity = *saved
	}
	return nil
}

// Update 更新实体。使用 GORM 的 Save 方法进行全量更新。
func (r *GenericCommandRepository[E, M]) Update(ctx context.Context, entity *E) error {
	model := r.entityToModel(entity)
	if err := r.db.WithContext(ctx).Save(&model).Error; err != nil {
		return fmt.Errorf("failed to update: %w", err)
	}
	// 同步更新后的值（如 UpdatedAt）
	if saved := model.ToEntity(); saved != nil {
		*entity = *saved
	}
	return nil
}

// Delete 删除实体（软删除，如果 Model 包含 DeletedAt 字段）。
func (r *GenericCommandRepository[E, M]) Delete(ctx context.Context, id uint) error {
	var model M
	if err := r.db.WithContext(ctx).Delete(&model, id).Error; err != nil {
		return fmt.Errorf("failed to delete: %w", err)
	}
	return nil
}

// DB 返回底层的 GORM DB 实例，用于子类实现特殊方法。
func (r *GenericCommandRepository[E, M]) DB() *gorm.DB {
	return r.db
}

// GenericQueryRepository 提供泛型的读操作仓储实现。
// 泛型参数：
//   - E: Domain Entity 类型
//   - M: GORM Model 类型（必须实现 Model[E] 接口）
type GenericQueryRepository[E any, M Model[E]] struct {
	db          *gorm.DB
	notFoundErr error // 可选的 Domain 层"未找到"错误
}

// NewGenericQueryRepository 创建泛型读仓储实例。
// notFoundErr 参数可选，用于在记录不存在时返回 Domain 层定义的错误。
func NewGenericQueryRepository[E any, M Model[E]](
	db *gorm.DB,
	notFoundErr error,
) *GenericQueryRepository[E, M] {
	return &GenericQueryRepository[E, M]{
		db:          db,
		notFoundErr: notFoundErr,
	}
}

// GetByID 根据 ID 获取单个实体。
// 如果记录不存在且设置了 notFoundErr，返回该错误；否则返回 nil。
func (r *GenericQueryRepository[E, M]) GetByID(ctx context.Context, id uint) (*E, error) {
	var model M
	if err := r.db.WithContext(ctx).First(&model, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, r.notFoundErr
		}
		return nil, fmt.Errorf("failed to get by id: %w", err)
	}
	return model.ToEntity(), nil
}

// List 获取实体列表（分页）。
func (r *GenericQueryRepository[E, M]) List(ctx context.Context, offset, limit int) ([]*E, error) {
	var models []M
	if err := r.db.WithContext(ctx).Offset(offset).Limit(limit).Find(&models).Error; err != nil {
		return nil, fmt.Errorf("failed to list: %w", err)
	}
	return r.modelsToEntities(models), nil
}

// Count 统计实体总数。
func (r *GenericQueryRepository[E, M]) Count(ctx context.Context) (int64, error) {
	var model M
	var count int64
	if err := r.db.WithContext(ctx).Model(&model).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count: %w", err)
	}
	return count, nil
}

// Exists 检查指定 ID 的实体是否存在。
func (r *GenericQueryRepository[E, M]) Exists(ctx context.Context, id uint) (bool, error) {
	var model M
	var count int64
	if err := r.db.WithContext(ctx).Model(&model).Where("id = ?", id).Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to check existence: %w", err)
	}
	return count > 0, nil
}

// DB 返回底层的 GORM DB 实例，用于子类实现特殊查询。
func (r *GenericQueryRepository[E, M]) DB() *gorm.DB {
	return r.db
}

// NotFoundErr 返回配置的"未找到"错误。
func (r *GenericQueryRepository[E, M]) NotFoundErr() error {
	return r.notFoundErr
}

// modelsToEntities 将 Model 切片转换为 Entity 切片。
func (r *GenericQueryRepository[E, M]) modelsToEntities(models []M) []*E {
	if len(models) == 0 {
		return nil
	}
	entities := make([]*E, 0, len(models))
	for i := range models {
		if entity := models[i].ToEntity(); entity != nil {
			entities = append(entities, entity)
		}
	}
	return entities
}

// ModelsToEntities 将 Model 切片转换为 Entity 切片（导出版本，供子类使用）。
func (r *GenericQueryRepository[E, M]) ModelsToEntities(models []M) []*E {
	return r.modelsToEntities(models)
}
