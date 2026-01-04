package user

import "go.uber.org/fx"

// UseCaseModule 用户配置用例 Fx 模块
// Note: Handlers are constructed by app.newUseCases, not provided directly
var UseCaseModule = fx.Module("settings.user.usecase")
