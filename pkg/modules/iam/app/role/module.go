package role

import (
	"go.uber.org/fx"

	"github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/infra/persistence"
	"github.com/lwmacct/260103-ddd-shared/pkg/shared/event"
)

// RoleUseCases 角色管理用例处理器聚合
type RoleUseCases struct {
	Create         *CreateHandler
	Update         *UpdateHandler
	Delete         *DeleteHandler
	SetPermissions *SetPermissionsHandler
	Get            *GetHandler
	List           *ListHandler
}

// Module 注册 Role 子模块依赖
var Module = fx.Module("iam.role",
	fx.Provide(newRoleUseCases),
)

type roleUseCasesParams struct {
	fx.In

	RoleRepos persistence.RoleRepositories
	EventBus  event.EventBus
}

func newRoleUseCases(p roleUseCasesParams) *RoleUseCases {
	return &RoleUseCases{
		Create:         NewCreateHandler(p.RoleRepos.Command, p.RoleRepos.Query),
		Update:         NewUpdateHandler(p.RoleRepos.Command, p.RoleRepos.Query),
		Delete:         NewDeleteHandler(p.RoleRepos.Command, p.RoleRepos.Query),
		SetPermissions: NewSetPermissionsHandler(p.RoleRepos.Command, p.RoleRepos.Query, p.EventBus),
		Get:            NewGetHandler(p.RoleRepos.Query),
		List:           NewListHandler(p.RoleRepos.Query),
	}
}
