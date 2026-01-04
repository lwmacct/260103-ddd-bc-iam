// Package team 提供团队配置的应用层用例处理器。
//
// 本包实现了团队配置的 CQRS 处理器：
//   - Command: Set, Reset
//   - Query: Get, List（支持三级继承：Team → Org → System）
//
// 团队配置继承链：Team Settings → Org Settings → System Settings（默认值）
package team
