package persistence

import (
	"go.uber.org/fx"
	"gorm.io/gorm"

	"github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/app/user"
)

// RepositoryModule 提供 IAM 模块的所有仓储实现。
//
// 装饰缓存层的仓储：
//   - User: 缓存查询（GetByIDWithRoles）+ 缓存命令（失效缓存）
var RepositoryModule = fx.Module("iam.repository",
	fx.Provide(
		// 直接使用 persistence 构造函数（无需包装）
		NewAuditRepositories,
		NewOrganizationRepositories,
		NewRoleRepositories,
		NewPATRepositories,
		NewTwoFARepositories,

		// 组织和团队仓储（用于中间件）
		NewOrgMemberRepositories,
		NewTeamRepositories,
		NewTeamMemberRepositories,

		// 带缓存装饰的仓储
		newUserRepositoriesWithCache,
	),
)

// --- 带缓存装饰的仓储构造函数 ---

func newUserRepositoriesWithCache(
	db *gorm.DB,
	userWithRolesCache user.UserWithRolesCacheService,
) UserRepositories {
	rawRepos := NewUserRepositories(db)
	cachedQuery := NewCachedUserQueryRepository(rawRepos.Query, userWithRolesCache)
	cachedCommand := NewCachedUserCommandRepository(rawRepos.Command, userWithRolesCache)
	return UserRepositories{
		Command: cachedCommand,
		Query:   cachedQuery,
	}
}
