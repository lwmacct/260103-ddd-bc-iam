// Package routes 定义 HTTP 路由配置。
//
// 本包是适配器层的路由配置中心，集中管理所有 API 操作的 HTTP 元数据：
//   - HTTP 路由（Method, Path）
//   - Swagger 注解（Tags, Summary, Description）
//   - 审计开关（Audit bool，详情从 Operation 派生）
//
// # 单一数据源
//
// 路由配置定义在 http.AllRoutes()，启动时由 [BuildRegistryFromRoutes] 构建 Registry。
// 其他模块通过函数获取：
//   - [Method]: 获取 HTTP 方法
//   - [Path]: 获取路由路径
//   - [NeedsAudit]: 判断是否需要审计
//   - [AuditAction]: 获取审计操作标识（派生自 Operation）
//   - [AllOperationDefinitions]: 权限列表（供前端）
//
// # 审计信息派生
//
// 审计详情从 URN Operation 自动派生，无需手动配置：
//   - Category: 从 Operation.Type() 映射（使用 [audit.DeriveCategory]）
//   - Action: 格式为 {category}.{identifier}（使用 [audit.DeriveAction]）
//   - Operation: 从 Identifier 映射（使用 [audit.DeriveOperation]）
//
// # 依赖关系
//
// 本包依赖领域层类型：
//   - [permission.Operation]: URN 操作标识符
//   - [audit.Category]: 审计分类
//   - [audit.Operation]: 审计操作类型
//
// # 设计原则
//
// HTTP 配置属于适配器层，与领域概念分离：
//   - 领域概念（Operation, Resource, Matcher）→ domain/permission
//   - 审计派生（Category, Operation, Action）→ domain/audit
//   - HTTP 配置（Method, Path, Swagger, Audit 开关）→ adapters/http/routes
package registry
