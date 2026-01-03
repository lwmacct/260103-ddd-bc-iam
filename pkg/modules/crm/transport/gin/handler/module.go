package handler

import (
	"go.uber.org/fx"

	crmapplication "github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/crm/application"
)

// HandlersResult 使用 fx.Out 批量返回 CRM 模块的所有 HTTP 处理器。
type HandlersResult struct {
	fx.Out

	Contact     *ContactHandler
	Company     *CompanyHandler
	Lead        *LeadHandler
	Opportunity *OpportunityHandler
}

// HandlerModule 提供 CRM 模块的所有 HTTP 处理器。
var HandlerModule = fx.Module("crm.handler",
	fx.Provide(newAllHandlers),
)

// handlersParams 聚合创建 Handler 所需的依赖。
type handlersParams struct {
	fx.In

	// CRM 模块用例
	Contact     *crmapplication.ContactUseCases
	Company     *crmapplication.CompanyUseCases
	Lead        *crmapplication.LeadUseCases
	Opportunity *crmapplication.OpportunityUseCases
}

func newAllHandlers(p handlersParams) HandlersResult {
	return HandlersResult{
		Contact: NewContactHandler(
			p.Contact.Create,
			p.Contact.Update,
			p.Contact.Delete,
			p.Contact.Get,
			p.Contact.List,
		),
		Company: NewCompanyHandler(
			p.Company.Create,
			p.Company.Update,
			p.Company.Delete,
			p.Company.Get,
			p.Company.List,
		),
		Lead: NewLeadHandler(
			p.Lead.Create,
			p.Lead.Update,
			p.Lead.Delete,
			p.Lead.Contact,
			p.Lead.Qualify,
			p.Lead.Convert,
			p.Lead.Lose,
			p.Lead.Get,
			p.Lead.List,
		),
		Opportunity: NewOpportunityHandler(
			p.Opportunity.Create,
			p.Opportunity.Update,
			p.Opportunity.Delete,
			p.Opportunity.Advance,
			p.Opportunity.CloseWon,
			p.Opportunity.CloseLost,
			p.Opportunity.Get,
			p.Opportunity.List,
		),
	}
}
