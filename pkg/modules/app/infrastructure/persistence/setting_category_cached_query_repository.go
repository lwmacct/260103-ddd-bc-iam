package persistence

import (
	"context"
	"log/slog"

	appsetting "github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/app/application/setting"
	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/app/domain/setting"
)

// cachedSettingCategoryQueryRepository 带缓存的 SettingCategory 查询仓储装饰器。
//
// 装饰 [setting.SettingCategoryQueryRepository]，使用 Application 层 SettingsCacheService。
//
// 缓存策略（简化设计）：
//   - 所有查询先从全量缓存过滤，数据量极小（4-10 条），内存过滤开销可忽略
//   - [FindAll]: 缓存命中直接返回，未命中查库后回写全量缓存
//   - [FindByID], [FindByKey], [FindByIDs]: 从全量缓存过滤，未命中回退数据库
//   - 单条查询未命中时不单独回写，等待 FindAll 触发全量回写
type cachedSettingCategoryQueryRepository struct {
	delegate setting.SettingCategoryQueryRepository // 被装饰的原始仓储
	cache    appsetting.SettingsCacheService        // Application 层缓存服务
}

// NewCachedSettingCategoryQueryRepository 创建带缓存的 SettingCategory 查询仓储。
func NewCachedSettingCategoryQueryRepository(
	delegate setting.SettingCategoryQueryRepository,
	cacheService appsetting.SettingsCacheService,
) setting.SettingCategoryQueryRepository {
	return &cachedSettingCategoryQueryRepository{
		delegate: delegate,
		cache:    cacheService,
	}
}

// FindByID 根据 ID 查找配置分类（从全量缓存过滤）。
func (r *cachedSettingCategoryQueryRepository) FindByID(ctx context.Context, id uint) (*setting.SettingCategory, error) {
	// 从全量缓存过滤
	all, err := r.cache.GetAllCategories(ctx)
	if err != nil {
		slog.Warn("cache get all categories failed, fallback to db", "error", err.Error())
	}
	for _, c := range all {
		if c.ID == id {
			return c, nil
		}
	}

	// 缓存未命中，查数据库（不单独回写，等待 FindAll 触发全量回写）
	return r.delegate.FindByID(ctx, id)
}

// FindByKey 根据 Key 查找配置分类（从全量缓存过滤）。
func (r *cachedSettingCategoryQueryRepository) FindByKey(ctx context.Context, key string) (*setting.SettingCategory, error) {
	// 从全量缓存过滤
	all, err := r.cache.GetAllCategories(ctx)
	if err != nil {
		slog.Warn("cache get all categories failed, fallback to db", "error", err.Error())
	}
	for _, c := range all {
		if c.Key == key {
			return c, nil
		}
	}

	// 缓存未命中，查数据库
	return r.delegate.FindByKey(ctx, key)
}

// FindByIDs 批量查找配置分类（从全量缓存过滤）。
func (r *cachedSettingCategoryQueryRepository) FindByIDs(ctx context.Context, ids []uint) ([]*setting.SettingCategory, error) {
	if len(ids) == 0 {
		return []*setting.SettingCategory{}, nil
	}

	// 从全量缓存过滤
	all, err := r.cache.GetAllCategories(ctx)
	if err != nil {
		slog.Warn("cache get all categories failed, fallback to db", "error", err.Error())
	}
	if all != nil {
		// 构建 ID 集合
		idSet := make(map[uint]struct{}, len(ids))
		for _, id := range ids {
			idSet[id] = struct{}{}
		}

		// 过滤匹配的
		result := make([]*setting.SettingCategory, 0, len(ids))
		for _, c := range all {
			if _, ok := idSet[c.ID]; ok {
				result = append(result, c)
			}
		}

		// 全部命中则返回
		if len(result) == len(ids) {
			return result, nil
		}
	}

	// 部分或全部未命中，回退数据库
	return r.delegate.FindByIDs(ctx, ids)
}

// FindAll 查找所有配置分类。
//
// 缓存命中直接返回，未命中查库后回写全量缓存。
func (r *cachedSettingCategoryQueryRepository) FindAll(ctx context.Context) ([]*setting.SettingCategory, error) {
	// 尝试从缓存获取
	cachedList, err := r.cache.GetAllCategories(ctx)
	if err != nil {
		slog.Warn("cache get all categories failed, fallback to db", "error", err.Error())
		return r.delegate.FindAll(ctx)
	}

	// 缓存命中
	if len(cachedList) > 0 {
		return cachedList, nil
	}

	// 缓存未命中，查数据库
	result, err := r.delegate.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	// 回写全量缓存
	if len(result) > 0 {
		if err := r.cache.SetAllCategories(ctx, result); err != nil {
			slog.Warn("cache set all categories failed", "count", len(result), "error", err.Error())
		}
	}

	return result, nil
}

// ExistsByKey 检查 Key 是否已存在（从全量缓存过滤）。
func (r *cachedSettingCategoryQueryRepository) ExistsByKey(ctx context.Context, key string) (bool, error) {
	// 从全量缓存过滤
	all, _ := r.cache.GetAllCategories(ctx)
	for _, c := range all {
		if c.Key == key {
			return true, nil
		}
	}

	// 缓存未命中，查数据库
	return r.delegate.ExistsByKey(ctx, key)
}

var _ setting.SettingCategoryQueryRepository = (*cachedSettingCategoryQueryRepository)(nil)
