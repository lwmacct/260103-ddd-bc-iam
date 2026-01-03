// Package company 定义公司领域模型和仓储接口。
//
// 本包定义了：
//   - [Company]: 公司实体
//   - [CommandRepository]: 写仓储接口
//   - [QueryRepository]: 读仓储接口
//
// 公司与联系人为 1:N 关系，一个公司可关联多个联系人。
//
// 依赖倒置：本包仅定义接口，实现位于 infrastructure/persistence。
package company
