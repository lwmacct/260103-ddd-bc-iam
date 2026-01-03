package persistence

import (
	"context"
	"log/slog"

	appsetting "github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/app/application/setting"
	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/app/domain/setting"
)

// cachedUserSettingCommandRepository 带缓存失效的 UserSetting 命令仓储装饰器。
//
// 装饰 [setting.UserSettingCommandRepository]，在写操作后自动失效相关缓存。
// 同时失效三层缓存：
//   - UserSettingQueryCacheService: 存储原始 UserSetting（Repository 层使用）
//   - UserSettingCacheService: 存储合并后的 EffectiveUserSetting（Application 层使用）
//   - SettingsCacheService: 存储 Schema API 响应（Application 层使用）
type cachedUserSettingCommandRepository struct {
	delegate       setting.UserSettingCommandRepository
	queryCache     appsetting.UserSettingQueryCacheService // Repository 层缓存
	effectiveCache appsetting.UserSettingCacheService      // Application 层缓存
	settingsCache  appsetting.SettingsCacheService         // Schema 响应缓存（Application 层接口）
}

// NewCachedUserSettingCommandRepository 创建带缓存失效的 UserSetting 命令仓储。
func NewCachedUserSettingCommandRepository(
	delegate setting.UserSettingCommandRepository,
	queryCache appsetting.UserSettingQueryCacheService,
	effectiveCache appsetting.UserSettingCacheService,
	settingsCache appsetting.SettingsCacheService,
) setting.UserSettingCommandRepository {
	return &cachedUserSettingCommandRepository{
		delegate:       delegate,
		queryCache:     queryCache,
		effectiveCache: effectiveCache,
		settingsCache:  settingsCache,
	}
}

// Upsert 插入或更新用户配置，并失效缓存。
func (r *cachedUserSettingCommandRepository) Upsert(ctx context.Context, us *setting.UserSetting) error {
	if err := r.delegate.Upsert(ctx, us); err != nil {
		return err
	}

	// 失效 Application 层缓存（EffectiveUserSetting）
	if err := r.effectiveCache.Delete(ctx, us.UserID, us.SettingKey); err != nil {
		slog.Warn("failed to invalidate effective user setting cache after upsert",
			"userID", us.UserID, "key", us.SettingKey, "error", err.Error())
	}

	// 失效 Repository 层查询缓存（原始 UserSetting）
	if err := r.queryCache.DeleteByUser(ctx, us.UserID); err != nil {
		slog.Warn("failed to invalidate user setting query cache after upsert",
			"userID", us.UserID, "error", err.Error())
	}

	// 失效 Schema 响应缓存
	if err := r.settingsCache.DeleteUserSettingsAll(ctx, us.UserID); err != nil {
		slog.Warn("failed to invalidate user schema cache after upsert",
			"userID", us.UserID, "error", err.Error())
	}

	return nil
}

// Delete 删除用户配置，并失效缓存。
func (r *cachedUserSettingCommandRepository) Delete(ctx context.Context, userID uint, key string) error {
	if err := r.delegate.Delete(ctx, userID, key); err != nil {
		return err
	}

	// 失效 Application 层缓存
	if err := r.effectiveCache.Delete(ctx, userID, key); err != nil {
		slog.Warn("failed to invalidate effective user setting cache after delete",
			"userID", userID, "key", key, "error", err.Error())
	}

	// 失效 Repository 层查询缓存
	if err := r.queryCache.DeleteByUser(ctx, userID); err != nil {
		slog.Warn("failed to invalidate user setting query cache after delete",
			"userID", userID, "error", err.Error())
	}

	// 失效 Schema 响应缓存
	if err := r.settingsCache.DeleteUserSettingsAll(ctx, userID); err != nil {
		slog.Warn("failed to invalidate user schema cache after delete",
			"userID", userID, "error", err.Error())
	}

	return nil
}

// DeleteByUser 删除用户的所有配置，并失效缓存。
func (r *cachedUserSettingCommandRepository) DeleteByUser(ctx context.Context, userID uint) error {
	if err := r.delegate.DeleteByUser(ctx, userID); err != nil {
		return err
	}

	// 失效 Application 层缓存
	if err := r.effectiveCache.DeleteByUser(ctx, userID); err != nil {
		slog.Warn("failed to invalidate all effective user setting cache after delete by user",
			"userID", userID, "error", err.Error())
	}

	// 失效 Repository 层查询缓存
	if err := r.queryCache.DeleteByUser(ctx, userID); err != nil {
		slog.Warn("failed to invalidate user setting query cache after delete by user",
			"userID", userID, "error", err.Error())
	}

	// 失效 Schema 响应缓存
	if err := r.settingsCache.DeleteUserSettingsAll(ctx, userID); err != nil {
		slog.Warn("failed to invalidate user schema cache after delete by user",
			"userID", userID, "error", err.Error())
	}

	return nil
}

// BatchUpsert 批量插入或更新用户配置，并失效缓存。
func (r *cachedUserSettingCommandRepository) BatchUpsert(ctx context.Context, settings []*setting.UserSetting) error {
	if err := r.delegate.BatchUpsert(ctx, settings); err != nil {
		return err
	}

	// 收集需要失效的 key
	if len(settings) > 0 {
		userID := settings[0].UserID
		keys := make([]string, 0, len(settings))
		for _, s := range settings {
			keys = append(keys, s.SettingKey)
		}

		// 失效 Application 层缓存
		if err := r.effectiveCache.DeleteByKeys(ctx, userID, keys); err != nil {
			slog.Warn("failed to invalidate effective user setting cache after batch upsert",
				"userID", userID, "count", len(keys), "error", err.Error())
		}

		// 失效 Repository 层查询缓存
		if err := r.queryCache.DeleteByUser(ctx, userID); err != nil {
			slog.Warn("failed to invalidate user setting query cache after batch upsert",
				"userID", userID, "error", err.Error())
		}

		// 失效 Schema 响应缓存
		if err := r.settingsCache.DeleteUserSettingsAll(ctx, userID); err != nil {
			slog.Warn("failed to invalidate user schema cache after batch upsert",
				"userID", userID, "error", err.Error())
		}
	}

	return nil
}

var _ setting.UserSettingCommandRepository = (*cachedUserSettingCommandRepository)(nil)
