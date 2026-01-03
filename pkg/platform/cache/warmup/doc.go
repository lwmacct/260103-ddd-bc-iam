// Package warmup 提供缓存预热服务。
//
// TODO: 实现缓存预热功能
//   - setting_warmup.go: Setting 缓存预热
//   - setting_category_warmup.go: SettingCategory 缓存预热
//
// 当前状态：仅定义包结构，实现待补充。
//
// # 设计目标
//
// 应用启动时批量加载热点数据到 Redis，减少首次请求的数据库压力。
//
// # 预热流程（待实现）
//
//  1. 从数据库加载所有数据
//  2. 批量写入缓存
//  3. 失败不阻塞启动，降级为惰性加载
//
// # 依赖注入原则
//
// 预热服务依赖原始仓储（非缓存装饰器），避免循环依赖。
package warmup
