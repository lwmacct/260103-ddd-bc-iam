// Package team 定义团队配置领域模型和仓储接口。
//
// 本包定义了：
//   - [TeamSetting]: 团队配置实体
//   - [CommandRepository]: 团队配置写仓储接口
//   - [QueryRepository]: 团队配置读仓储接口
//
// 依赖倒置：本包仅定义接口，实现位于 infrastructure/persistence。
package team
