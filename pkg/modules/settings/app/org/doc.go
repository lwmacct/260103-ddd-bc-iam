// Package org 提供组织配置的应用层用例处理器。
//
// 本包实现了组织配置的 CQRS 处理器：
//   - Command: Set, Reset
//   - Query: Get, List
//
// 组织配置支持继承链：Org Settings → System Settings（默认值）
package org
