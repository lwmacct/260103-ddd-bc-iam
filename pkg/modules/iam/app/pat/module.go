package pat

import (
	"go.uber.org/fx"
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
	fx.Provide(
		NewCreateHandler,
		NewDeleteHandler,
		NewDisableHandler,
		NewEnableHandler,
		NewGetHandler,
		NewListHandler,
		newPATUseCases,
	),
)

type patUseCasesParams struct {
	fx.In

	Create  *CreateHandler
	Delete  *DeleteHandler
	Disable *DisableHandler
	Enable  *EnableHandler
	Get     *GetHandler
	List    *ListHandler
}

func newPATUseCases(p patUseCasesParams) *PATUseCases {
	return &PATUseCases{
		Create:  p.Create,
		Delete:  p.Delete,
		Disable: p.Disable,
		Enable:  p.Enable,
		Get:     p.Get,
		List:    p.List,
	}
}
