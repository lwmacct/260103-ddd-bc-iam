// Package opportunity 定义商机领域模型和仓储接口。
//
// 本包定义了：
//   - [Opportunity]: 商机实体（带阶段管理）
//   - [Stage]: 商机阶段值对象
//   - [CommandRepository]: 写仓储接口
//   - [QueryRepository]: 读仓储接口
//
// 阶段转换图：
//
//	prospecting → proposal → negotiation ─┬─ CloseWon() → closed_won
//	                                       └─ CloseLost() → closed_lost
//
// 依赖倒置：本包仅定义接口，实现位于 infrastructure/persistence。
package opportunity
