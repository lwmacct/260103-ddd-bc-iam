// Package lead 定义线索领域模型和仓储接口。
//
// 本包定义了：
//   - [Lead]: 线索实体（含状态机）
//   - [Status]: 线索状态枚举
//   - [CommandRepository]: 写仓储接口
//   - [QueryRepository]: 读仓储接口
//
// 状态机流程：
//
//	new ──Contact()──→ contacted ──Qualify()──→ qualified ──Convert()──→ converted
//	 │                     │                        │
//	 └─────────────────────┴────────Lose()──────────┴──────────────────→ lost
//
// 依赖倒置：本包仅定义接口，实现位于 infrastructure/persistence。
package lead
