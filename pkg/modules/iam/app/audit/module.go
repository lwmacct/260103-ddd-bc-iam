package audit

import (
	"go.uber.org/fx"

	"github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/infra/persistence"
)

// AuditUseCases 审计日志用例处理器聚合
type AuditUseCases struct {
	CreateLog *CreateHandler
	Get       *GetHandler
	List      *ListHandler
}

// Module 注册 Audit 子模块依赖
var Module = fx.Module("iam.audit",
	fx.Provide(newAuditUseCases),
)

type auditUseCasesParams struct {
	fx.In

	Repos persistence.AuditRepositories
}

func newAuditUseCases(p auditUseCasesParams) *AuditUseCases {
	return &AuditUseCases{
		CreateLog: NewCreateHandler(p.Repos.Command),
		Get:       NewGetHandler(p.Repos.Query),
		List:      NewListHandler(p.Repos.Query),
	}
}
