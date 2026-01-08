package user

import (
	"go.uber.org/fx"

	authDomain "github.com/lwmacct/260103-ddd-iam-bc/pkg/modules/iam/domain/auth"
	"github.com/lwmacct/260103-ddd-iam-bc/pkg/modules/iam/infra/persistence"
	"github.com/lwmacct/260103-ddd-shared/pkg/shared/event"
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
	fx.Provide(newUserUseCases),
)

type userUseCasesParams struct {
	fx.In

	UserRepos persistence.UserRepositories
	AuthSvc   authDomain.Service
	EventBus  event.EventBus
}

func newUserUseCases(p userUseCasesParams) *UserUseCases {
	return &UserUseCases{
		Create:         NewCreateHandler(p.UserRepos.Command, p.UserRepos.Query, p.AuthSvc),
		Update:         NewUpdateHandler(p.UserRepos.Command, p.UserRepos.Query),
		Delete:         NewDeleteHandler(p.UserRepos.Command, p.UserRepos.Query, p.EventBus),
		AssignRoles:    NewAssignRolesHandler(p.UserRepos.Command, p.UserRepos.Query, p.EventBus),
		ChangePassword: NewChangePasswordHandler(p.UserRepos.Command, p.UserRepos.Query, p.AuthSvc),
		BatchCreate:    NewBatchCreateHandler(p.UserRepos.Command, p.UserRepos.Query, p.AuthSvc),
		Get:            NewGetHandler(p.UserRepos.Query),
		List:           NewListHandler(p.UserRepos.Query),
	}
}
