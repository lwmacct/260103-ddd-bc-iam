package persistence

import (
	"context"
	"log/slog"
	"time"

	appsetting "github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/app/application/setting"
	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/app/domain/setting"
)

// cachedUserSettingQueryRepository 带缓存的 UserSetting 查询仓储装饰器。
//
// 装饰 [setting.UserSettingQueryRepository]，在查询前检查缓存，未命中再查数据库。
// 采用 Cache-Aside 模式，异步回写缓存。
//
// 缓存策略：
//   - [FindByUser]: 缓存用户的所有自定义配置
//   - [FindByUserAndKey]: 从用户缓存中提取单条，未命中查库
//   - [FindByUserAndKeys]: 从用户缓存中提取，部分未命中查库
type cachedUserSettingQueryRepository struct {
	delegate setting.UserSettingQueryRepository
	cache    appsetting.UserSettingQueryCacheService
}

// NewCachedUserSettingQueryRepository 创建带缓存的 UserSetting 查询仓储。
func NewCachedUserSettingQueryRepository(
	delegate setting.UserSettingQueryRepository,
	cacheService appsetting.UserSettingQueryCacheService,
) setting.UserSettingQueryRepository {
	return &cachedUserSettingQueryRepository{
		delegate: delegate,
		cache:    cacheService,
	}
}

// FindByUserAndKey 根据用户 ID 和 Key 查找用户配置（带缓存）。
func (r *cachedUserSettingQueryRepository) FindByUserAndKey(ctx context.Context, userID uint, key string) (*setting.UserSetting, error) {
	// 1. 尝试从用户全量缓存中获取
	cached, err := r.cache.GetByUser(ctx, userID)
	if err != nil {
		slog.Warn("user setting cache get failed, fallback to db", "userID", userID, "error", err.Error())
	}
	if cached != nil {
		if s, ok := cached[key]; ok {
			return s, nil
		}
		// 缓存存在但 key 不在其中，说明用户没有这个自定义配置
		return nil, nil //nolint:nilnil // cache hit, key not found
	}

	// 2. 缓存未命中，查数据库
	result, err := r.delegate.FindByUserAndKey(ctx, userID, key)
	if err != nil {
		return nil, err
	}

	// 3. 单条查询不触发全量缓存回写，避免频繁单条查询导致缓存不一致
	return result, nil
}

// FindByUser 查找用户的所有自定义配置（带缓存）。
func (r *cachedUserSettingQueryRepository) FindByUser(ctx context.Context, userID uint) ([]*setting.UserSetting, error) {
	// 1. 查缓存
	cached, err := r.cache.GetByUser(ctx, userID)
	if err != nil {
		slog.Warn("user setting cache get failed, fallback to db", "userID", userID, "error", err.Error())
	}
	if cached != nil {
		// 缓存命中，转换为切片返回
		result := make([]*setting.UserSetting, 0, len(cached))
		for _, s := range cached {
			result = append(result, s)
		}
		return result, nil
	}

	// 2. 查数据库
	result, err := r.delegate.FindByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	// 3. 异步回写缓存（包含空结果，防止缓存穿透）
	go func(settings []*setting.UserSetting, uid uint) {
		cacheCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), 3*time.Second)
		defer cancel()
		if err := r.cache.SetByUser(cacheCtx, uid, settings); err != nil {
			slog.Warn("user setting cache set failed", "userID", uid, "error", err.Error())
		}
	}(result, userID)

	return result, nil
}

// FindByUserAndKeys 根据用户 ID 和多个 Key 批量查找用户配置（带缓存）。
func (r *cachedUserSettingQueryRepository) FindByUserAndKeys(ctx context.Context, userID uint, keys []string) ([]*setting.UserSetting, error) {
	if len(keys) == 0 {
		return []*setting.UserSetting{}, nil
	}

	// 1. 尝试从用户全量缓存中获取
	cached, err := r.cache.GetByUser(ctx, userID)
	if err != nil {
		slog.Warn("user setting cache get failed, fallback to db", "userID", userID, "error", err.Error())
	}
	if cached != nil {
		// 从缓存中提取请求的 keys
		result := make([]*setting.UserSetting, 0, len(keys))
		for _, k := range keys {
			if s, ok := cached[k]; ok {
				result = append(result, s)
			}
		}
		return result, nil
	}

	// 2. 缓存未命中，查数据库
	return r.delegate.FindByUserAndKeys(ctx, userID, keys)
}

var _ setting.UserSettingQueryRepository = (*cachedUserSettingQueryRepository)(nil)
