package team

import "go.uber.org/fx"

// UseCaseModule 团队配置用例 Fx 模块
// Note: Handlers are constructed by app.newUseCases, not provided directly
var UseCaseModule = fx.Module("settings.team.usecase")
