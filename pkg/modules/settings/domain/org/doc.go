// Package org 定义组织配置领域模型和仓储接口。
//
// 本包定义了：
//   - [OrgSetting]: 组织配置实体
//   - [CommandRepository]: 组织配置写仓储接口
//   - [QueryRepository]: 组织配置读仓储接口
//
// 依赖倒置：本包仅定义接口，实现位于 infrastructure/persistence。
package org
