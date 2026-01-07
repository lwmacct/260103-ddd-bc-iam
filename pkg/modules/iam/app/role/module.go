package role

import (
	"go.uber.org/fx"
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
	fx.Provide(
		NewCreateHandler,
		NewUpdateHandler,
		NewDeleteHandler,
		NewSetPermissionsHandler,
		NewGetHandler,
		NewListHandler,
		newRoleUseCases,
	),
)

type roleUseCasesParams struct {
	fx.In

	Create         *CreateHandler
	Update         *UpdateHandler
	Delete         *DeleteHandler
	SetPermissions *SetPermissionsHandler
	Get            *GetHandler
	List           *ListHandler
}

func newRoleUseCases(p roleUseCasesParams) *RoleUseCases {
	return &RoleUseCases{
		Create:         p.Create,
		Update:         p.Update,
		Delete:         p.Delete,
		SetPermissions: p.SetPermissions,
		Get:            p.Get,
		List:           p.List,
	}
}
