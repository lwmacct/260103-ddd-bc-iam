// Package task 定义团队任务领域模型和仓储接口。
//
// 本包是团队任务管理的领域层核心，定义了：
//   - [Task]: 任务实体
//   - [Status]: 任务状态值对象
//   - [CommandRepository]: 写仓储接口
//   - [QueryRepository]: 读仓储接口（支持按组织+团队过滤）
//
// 多租户设计：
// 任务通过 OrgID + TeamID 实现组织和团队级别的隔离。
// 所有查询操作都要求指定组织和团队 ID。
//
// 依赖倒置：
// 本包仅定义接口，实现位于 infrastructure/persistence 包。
package task
