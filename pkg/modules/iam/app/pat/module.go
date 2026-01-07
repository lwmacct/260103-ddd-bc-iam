package pat

import (
	"go.uber.org/fx"

	authInfra "github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/infra/auth"
	"github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/infra/persistence"
)

// PATUseCases 个人访问令牌用例处理器聚合
type PATUseCases struct {
	Create  *CreateHandler
	Delete  *DeleteHandler
	Disable *DisableHandler
	Enable  *EnableHandler
	Get     *GetHandler
	List    *ListHandler
}

// Module 注册 PAT 子模块依赖
var Module = fx.Module("iam.pat",
	fx.Provide(newPATUseCases),
)

type patUseCasesParams struct {
	fx.In

	PATRepos  persistence.PATRepositories
	UserRepos persistence.UserRepositories
	TokenGen  *authInfra.TokenGenerator
}

func newPATUseCases(p patUseCasesParams) *PATUseCases {
	return &PATUseCases{
		Create:  NewCreateHandler(p.PATRepos.Command, p.UserRepos.Query, p.TokenGen),
		Delete:  NewDeleteHandler(p.PATRepos.Command, p.PATRepos.Query),
		Disable: NewDisableHandler(p.PATRepos.Command, p.PATRepos.Query),
		Enable:  NewEnableHandler(p.PATRepos.Command, p.PATRepos.Query),
		Get:     NewGetHandler(p.PATRepos.Query),
		List:    NewListHandler(p.PATRepos.Query),
	}
}
