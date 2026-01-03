package container

import (
	"github.com/gin-gonic/gin"
	"github.com/lwmacct/260101-go-pkg-gin/pkg/routes"

	iamhandler "github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/transport/gin/handler"
	iamroutes "github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/transport/gin/routes"
)

// AllRoutes 聚合所有模块的路由定义。
//
// 职责：
//  1. 从 IAM 模块收集路由定义
//  2. 返回统一的路由列表供注册使用
func AllRoutes(
	// IAM Handlers
	authHandler *iamhandler.AuthHandler,
	twoFAHandler *iamhandler.TwoFAHandler,
	userProfileHandler *iamhandler.UserProfileHandler,
	userOrgHandler *iamhandler.UserOrgHandler,
	patHandler *iamhandler.PATHandler,
	adminUserHandler *iamhandler.AdminUserHandler,
	roleHandler *iamhandler.RoleHandler,
	captchaHandler *iamhandler.CaptchaHandler,
	auditHandler *iamhandler.AuditHandler,
	orgHandler *iamhandler.OrgHandler,
	orgMemberHandler *iamhandler.OrgMemberHandler,
	teamHandler *iamhandler.TeamHandler,
	teamMemberHandler *iamhandler.TeamMemberHandler,
) []routes.Route {
	// IAM 域路由
	return iamroutes.All(
		authHandler,
		twoFAHandler,
		captchaHandler,
		userProfileHandler,
		userOrgHandler,
		adminUserHandler,
		roleHandler,
		patHandler,
		auditHandler,
		orgHandler,
		orgMemberHandler,
		teamHandler,
		teamMemberHandler,
	)
}

// RegisterRoutes 注册路由到 Gin Engine。

// RegisterRoutes 注册路由到 Gin Engine。
//
// 架构特点：
//  1. BC 层定义路由元数据（Operation、Path、Method、Handler）
//  2. 应用层注入中间件（认证、鉴权、审计、日志等）
//  3. 完美解耦 - BC 不依赖具体中间件实现
func RegisterRoutes(engine *gin.Engine, allRoutes []routes.Route, injector *MiddlewareInjector) {
	for _, route := range allRoutes {
		// 应用层注入中间件
		middlewares := injector.InjectMiddlewares(&route)

		// 添加 Handler 作为最后一个处理函数
		middlewares = append(middlewares, route.Handler)

		// 注册到 Gin Engine
		switch route.Method {
		case routes.GET:
			engine.GET(route.Path, middlewares...)
		case routes.POST:
			engine.POST(route.Path, middlewares...)
		case routes.PUT:
			engine.PUT(route.Path, middlewares...)
		case routes.DELETE:
			engine.DELETE(route.Path, middlewares...)
		case routes.PATCH:
			engine.PATCH(route.Path, middlewares...)
		case routes.HEAD:
			engine.HEAD(route.Path, middlewares...)
		case routes.OPTIONS:
			engine.OPTIONS(route.Path, middlewares...)
		}
	}
}
