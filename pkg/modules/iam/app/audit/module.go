package audit

import (
	"go.uber.org/fx"
)

// AuditUseCases 审计日志用例处理器聚合
type AuditUseCases struct {
	CreateLog *CreateHandler
	Get       *GetHandler
	List      *ListHandler
}

// Module 注册 Audit 子模块依赖
var Module = fx.Module("iam.audit",
	fx.Provide(
		NewCreateHandler,
		NewGetHandler,
		NewListHandler,
		newAuditUseCases,
	),
)

type auditUseCasesParams struct {
	fx.In

	CreateLog *CreateHandler
	Get       *GetHandler
	List      *ListHandler
}

func newAuditUseCases(p auditUseCasesParams) *AuditUseCases {
	return &AuditUseCases{
		CreateLog: p.CreateLog,
		Get:       p.Get,
		List:      p.List,
	}
}
