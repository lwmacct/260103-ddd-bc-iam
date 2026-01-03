// Package contact 定义联系人领域模型和仓储接口。
//
// 本包定义了：
//   - [Contact]: 联系人实体
//   - [CommandRepository]: 写仓储接口
//   - [QueryRepository]: 读仓储接口
//
// 联系人是 CRM 系统的核心实体，可关联公司（Company）。
//
// 依赖倒置：本包仅定义接口，实现位于 infrastructure/persistence。
package contact
