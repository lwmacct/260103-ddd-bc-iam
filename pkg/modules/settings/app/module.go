package app

import (
	"go.uber.org/fx"

	"github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/settings/app/org"
	"github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/settings/app/team"
	"github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/settings/app/user"
)

// 类型别名 - 向后兼容（供 Handler 层便捷访问）
type (
	UserUseCases = user.UserUseCases
	OrgUseCases  = org.OrgUseCases
	TeamUseCases = team.TeamUseCases
)

// UseCaseModule 用例 Fx 模块 - 聚合所有子模块
var UseCaseModule = fx.Module("settings.usecase",
	user.Module,
	org.Module,
	team.Module,
)
