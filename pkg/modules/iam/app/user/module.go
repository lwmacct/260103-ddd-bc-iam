package user

import (
	"go.uber.org/fx"
)

// UserUseCases 用户管理用例处理器聚合
type UserUseCases struct {
	Create         *CreateHandler
	Update         *UpdateHandler
	Delete         *DeleteHandler
	AssignRoles    *AssignRolesHandler
	ChangePassword *ChangePasswordHandler
	BatchCreate    *BatchCreateHandler
	Get            *GetHandler
	List           *ListHandler
}

// Module 注册 User 子模块依赖
var Module = fx.Module("iam.user",
	fx.Provide(
		NewCreateHandler,
		NewUpdateHandler,
		NewDeleteHandler,
		NewAssignRolesHandler,
		NewChangePasswordHandler,
		NewBatchCreateHandler,
		NewGetHandler,
		NewListHandler,
		newUserUseCases,
	),
)

type userUseCasesParams struct {
	fx.In

	Create         *CreateHandler
	Update         *UpdateHandler
	Delete         *DeleteHandler
	AssignRoles    *AssignRolesHandler
	ChangePassword *ChangePasswordHandler
	BatchCreate    *BatchCreateHandler
	Get            *GetHandler
	List           *ListHandler
}

func newUserUseCases(p userUseCasesParams) *UserUseCases {
	return &UserUseCases{
		Create:         p.Create,
		Update:         p.Update,
		Delete:         p.Delete,
		AssignRoles:    p.AssignRoles,
		ChangePassword: p.ChangePassword,
		BatchCreate:    p.BatchCreate,
		Get:            p.Get,
		List:           p.List,
	}
}
