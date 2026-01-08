// Package audit 定义审计领域模型和派生逻辑。
//
// # Overview
//
// 本包整合了审计日志实体、审计类型定义和 URN 派生逻辑：
//   - [Audit]: 审计日志实体
//   - [Operation]: 审计操作类型 (create, update, delete, access, authenticate)
//   - [Category]: 审计分类 (auth, user, role, setting 等)
//   - [DeriveCategory], [DeriveOperation], [DeriveAction]: 从 URN 派生审计信息
//   - [FilterOptions]: 日志查询过滤条件
//   - [CommandRepository]: 写仓储接口
//   - [QueryRepository]: 读仓储接口
//
// # 审计类型
//
// [Operation] 定义粗粒度操作分类：
//   - [OpCreate]: 创建操作
//   - [OpUpdate]: 更新操作
//   - [OpDelete]: 删除操作
//   - [OpAccess]: 访问操作
//   - [OpAuthenticate]: 认证操作
//
// [Category] 定义审计分类：
//   - [CatAuth]: 认证相关
//   - [CatUser]: 用户管理
//   - [CatRole]: 角色管理
//   - [CatSetting]: 系统配置
//
// # URN 派生
//
// 从 URN Operation（如 sys:users:create）派生审计信息：
//
//	cat := audit.DeriveCategory("users")     // "user"
//	op := audit.DeriveOperation("create")    // OpCreate
//	action := audit.DeriveAction("users", "create")  // "user.create"
//
// # Usage
//
//	// 创建审计日志实体
//	log := &audit.Audit{
//	    UserID:    123,
//	    Action:    "user.create",
//	    Category:  audit.CatUser,
//	    Operation: audit.OpCreate,
//	    Resource:  "user:456",
//	    Details:   `{"username": "alice"}`,
//	}
//
//	// 从 URN 派生审计信息
//	cat := audit.DeriveCategory("users")
//	op := audit.DeriveOperation("create")
//	action := audit.DeriveAction("users", "create")
//
// # Thread Safety
//
// 审计日志实体和派生函数都是值类型或纯函数，是并发安全的。
// Repository 接口的实现需要保证并发安全性（由基础设施层负责）。
//
// # 依赖关系
//
// 本包仅定义接口，实现位于 infrastructure/persistence 包。
package audit
