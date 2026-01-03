// Package usersetting 定义用户设置应用层用例处理器。
//
// 本包实现了：
//   - [ListHandler]: 获取用户设置列表（合并 Settings Schema + 用户自定义值）
//   - [GetHandler]: 获取单个用户设置
//   - [UpdateHandler]: 更新用户设置（验证 Schema）
//   - [DeleteHandler]: 删除用户设置（恢复系统默认值）
//
// 跨 BC 依赖：
// 本包注入 Settings BC 的 QueryRepository 用于获取配置 Schema（setting.QueryRepository）。
//
// 设计模式：
// 合并模式（Template Method）：
//  1. 从 Settings BC 获取 Schema（系统定义的配置项）
//  2. 从 UserSetting QueryRepository 获取用户自定义值
//  3. 合并：用户值优先，系统默认值兜底
//  4. 返回统一视图（带 is_custom 标识）
//
// 依赖倒置：本包依赖 Settings BC 的接口（setting.QueryRepository），而非具体实现。
package usersetting
