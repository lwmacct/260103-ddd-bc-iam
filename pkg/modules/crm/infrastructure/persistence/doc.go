// Package persistence 实现 CRM 域的 Repository 接口。
//
// 本包提供 CRM 域所有实体的持久化实现：
//   - Lead: 线索
//   - Opportunity: 商机
//   - Contact: 联系人
//   - Company: 公司
//
// 依赖关系：
//   - 本包只依赖 CRM Domain 层接口
//   - 不依赖其他域的 Persistence 实现
//   - 通过逻辑外键关联 User、Org 等实体
package persistence
