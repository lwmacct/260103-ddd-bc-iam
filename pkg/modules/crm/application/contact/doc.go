// Package contact 提供联系人管理的应用层用例。
//
// 本包实现了 CQRS 模式的 Command/Query Handler：
//   - [CreateHandler]: 创建联系人
//   - [UpdateHandler]: 更新联系人
//   - [DeleteHandler]: 删除联系人
//   - [GetHandler]: 获取联系人详情
//   - [ListHandler]: 联系人列表查询
//
// 所有 Handler 依赖 Domain 层接口，通过 DI 注入实现。
package contact
