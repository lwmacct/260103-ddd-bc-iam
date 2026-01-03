// Package opportunity 提供商机管理用例。
//
// 本包实现商机的完整生命周期管理，包括：
//   - CRUD 操作：Create, Update, Delete, Get, List
//   - 阶段推进：AdvanceHandler（prospecting → proposal → negotiation）
//   - 成交/丢单：CloseWonHandler, CloseLostHandler
//
// 遵循 CQRS 模式，命令和查询分离处理。
package opportunity
