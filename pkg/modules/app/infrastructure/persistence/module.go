package persistence

import (
	"go.uber.org/fx"
	"gorm.io/gorm"

	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/app/application/setting"
	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/app/domain/stats"
	infrastats "github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/app/infrastructure/stats"
	"github.com/lwmacct/260101-go-pkg-ddd/pkg/shared/captcha"
	infracaptcha "github.com/lwmacct/260101-go-pkg-ddd/pkg/shared/captcha/infrastructure"
)

// RepositoryModule 提供 App 模块的所有仓储实现。
//
// 装饰缓存层的仓储：
//   - Setting: 缓存查询 + 命令，支持多级失效
var RepositoryModule = fx.Module("app.repository",
	fx.Provide(
		// 带缓存装饰的仓储
		newSettingRepositoriesWithCache,

		// 特殊仓储
		newCaptchaRepository,
		newStatsQueryRepository,
	),
)

// --- 带缓存装饰的仓储构造函数 ---

// settingRepositoriesParams 聚合 Setting 仓储所需的缓存服务。
type settingRepositoriesParams struct {
	fx.In

	DB               *gorm.DB
	UserSettingCache setting.UserSettingCacheService
	SettingsCache    setting.SettingsCacheService
}

// newSettingRepositoriesWithCache 创建带缓存装饰的 Setting 仓储。
//
// 简化设计：
//   - Query 直接使用原始仓储，不再缓存（由 Application 层 Settings 缓存覆盖）
//   - Command 装饰器只负责写操作后失效下游缓存（Settings + UserSetting）
//   - CategoryQuery 使用 SettingsCacheService 的 Category 缓存方法
//   - CategoryCommand 直接使用原始仓储，缓存失效在 Handler 层处理
func newSettingRepositoriesWithCache(p settingRepositoriesParams) SettingRepositories {
	rawRepos := NewSettingRepositories(p.DB)

	// 查询直接使用原始仓储，不再缓存
	// 写操作装饰器：失效 Settings + UserSetting 缓存
	wrappedCommand := NewSettingCommandWithCacheInvalidation(
		rawRepos.Command,
		p.UserSettingCache,
		p.SettingsCache,
	)

	// Category 查询使用 SettingsCacheService（合并后的 Application 层缓存）
	cachedCategoryQuery := NewCachedSettingCategoryQueryRepository(
		rawRepos.CategoryQuery,
		p.SettingsCache,
	)

	// Category 命令直接使用原始仓储，缓存失效在 Handler 层统一处理
	return SettingRepositories{
		Command:         wrappedCommand,
		Query:           rawRepos.Query,
		CategoryQuery:   cachedCategoryQuery,
		CategoryCommand: rawRepos.CategoryCommand,
	}
}

// --- 特殊仓储 ---

// CaptchaRepositoryResult 从单个仓储提供 Command 和 Query 两个接口。
type CaptchaRepositoryResult struct {
	fx.Out

	Command captcha.CommandRepository
	Query   captcha.QueryRepository
}

// newCaptchaRepository 创建验证码仓储
func newCaptchaRepository() CaptchaRepositoryResult {
	repo := infracaptcha.NewMemoryRepository()
	return CaptchaRepositoryResult{
		Command: repo,
		Query:   repo,
	}
}

// newStatsQueryRepository 创建统计查询仓储
func newStatsQueryRepository(db *gorm.DB) stats.QueryRepository {
	return infrastats.NewQueryRepository(db)
}
