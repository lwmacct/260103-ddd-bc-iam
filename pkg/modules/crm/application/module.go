package application

import (
	"go.uber.org/fx"

	appCompany "github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/crm/application/company"
	appContact "github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/crm/application/contact"
	appLead "github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/crm/application/lead"
	appOpp "github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/crm/application/opportunity"
	crmpersistence "github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/crm/infrastructure/persistence"
)

// --- 用例模块结构体 ---

// ContactUseCases 联系人相关用例处理器
type ContactUseCases struct {
	Create *appContact.CreateHandler
	Update *appContact.UpdateHandler
	Delete *appContact.DeleteHandler
	Get    *appContact.GetHandler
	List   *appContact.ListHandler
}

// CompanyUseCases 公司相关用例处理器
type CompanyUseCases struct {
	Create *appCompany.CreateHandler
	Update *appCompany.UpdateHandler
	Delete *appCompany.DeleteHandler
	Get    *appCompany.GetHandler
	List   *appCompany.ListHandler
}

// LeadUseCases 线索相关用例处理器
type LeadUseCases struct {
	Create  *appLead.CreateHandler
	Update  *appLead.UpdateHandler
	Delete  *appLead.DeleteHandler
	Contact *appLead.ContactHandler
	Qualify *appLead.QualifyHandler
	Convert *appLead.ConvertHandler
	Lose    *appLead.LoseHandler
	Get     *appLead.GetHandler
	List    *appLead.ListHandler
}

// OpportunityUseCases 商机相关用例处理器
type OpportunityUseCases struct {
	Create    *appOpp.CreateHandler
	Update    *appOpp.UpdateHandler
	Delete    *appOpp.DeleteHandler
	Advance   *appOpp.AdvanceHandler
	CloseWon  *appOpp.CloseWonHandler
	CloseLost *appOpp.CloseLostHandler
	Get       *appOpp.GetHandler
	List      *appOpp.ListHandler
}

// --- Fx 模块 ---

// UseCaseModule 提供按领域组织的 CRM 模块用例处理器。
var UseCaseModule = fx.Module("crm.usecase",
	fx.Provide(
		newContactUseCases,
		newCompanyUseCases,
		newLeadUseCases,
		newOpportunityUseCases,
	),
)

// --- 构造函数 ---

func newContactUseCases(repos crmpersistence.ContactRepositories) *ContactUseCases {
	return &ContactUseCases{
		Create: appContact.NewCreateHandler(repos.Command, repos.Query),
		Update: appContact.NewUpdateHandler(repos.Command, repos.Query),
		Delete: appContact.NewDeleteHandler(repos.Command, repos.Query),
		Get:    appContact.NewGetHandler(repos.Query),
		List:   appContact.NewListHandler(repos.Query),
	}
}

func newCompanyUseCases(repos crmpersistence.CompanyRepositories) *CompanyUseCases {
	return &CompanyUseCases{
		Create: appCompany.NewCreateHandler(repos.Command, repos.Query),
		Update: appCompany.NewUpdateHandler(repos.Command, repos.Query),
		Delete: appCompany.NewDeleteHandler(repos.Command, repos.Query),
		Get:    appCompany.NewGetHandler(repos.Query),
		List:   appCompany.NewListHandler(repos.Query),
	}
}

func newLeadUseCases(repos crmpersistence.LeadRepositories) *LeadUseCases {
	return &LeadUseCases{
		Create:  appLead.NewCreateHandler(repos.Command, repos.Query),
		Update:  appLead.NewUpdateHandler(repos.Command, repos.Query),
		Delete:  appLead.NewDeleteHandler(repos.Command, repos.Query),
		Contact: appLead.NewContactHandler(repos.Command, repos.Query),
		Qualify: appLead.NewQualifyHandler(repos.Command, repos.Query),
		Convert: appLead.NewConvertHandler(repos.Command, repos.Query),
		Lose:    appLead.NewLoseHandler(repos.Command, repos.Query),
		Get:     appLead.NewGetHandler(repos.Query),
		List:    appLead.NewListHandler(repos.Query),
	}
}

func newOpportunityUseCases(repos crmpersistence.OpportunityRepositories) *OpportunityUseCases {
	return &OpportunityUseCases{
		Create:    appOpp.NewCreateHandler(repos.Command),
		Update:    appOpp.NewUpdateHandler(repos.Command, repos.Query),
		Delete:    appOpp.NewDeleteHandler(repos.Command, repos.Query),
		Advance:   appOpp.NewAdvanceHandler(repos.Command, repos.Query),
		CloseWon:  appOpp.NewCloseWonHandler(repos.Command, repos.Query),
		CloseLost: appOpp.NewCloseLostHandler(repos.Command, repos.Query),
		Get:       appOpp.NewGetHandler(repos.Query),
		List:      appOpp.NewListHandler(repos.Query),
	}
}
