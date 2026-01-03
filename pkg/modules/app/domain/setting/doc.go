// Package setting 定义系统配置和用户配置领域模型。
//
// 本包管理应用程序的动态配置项，支持运行时修改而无需重启服务。
//
// # 三表设计
//
// 采用配置定义、分类元数据与用户配置分离的设计：
//   - [Setting]: 配置定义表，存储 Schema 和默认值
//   - [SettingCategory]: 配置分类表，存储 Category 的 UI 元数据
//   - [UserSetting]: 用户配置表，存储用户自定义值（覆盖模式）
//
// # 配置作用域 (Scope)
//
// Setting.Scope 决定配置的作用域和访问方式：
//   - system: 系统设置，全局唯一，管理员直接修改 DefaultValue
//   - user: 用户设置，DefaultValue 作为初始值，用户可在 user_settings 表覆盖
//
// # 配置分类 (Category)
//
// 分类元数据存储在 setting_categories 表，供前端直接渲染：
//   - general: 常规设置
//   - security: 安全设置
//   - notification: 通知设置
//   - backup: 备份设置
//
// # 值类型 (ValueType)
//
// 配置值直接存储为 JSONB 原生类型：
//   - string: 字符串，如 "My Site"
//   - number: 数值，如 30
//   - boolean: 布尔值，如 true
//   - json: JSON 对象/数组
//
// # 仓储接口
//
// 遵循 CQRS 模式，读写分离：
//   - [CommandRepository]: 配置定义写操作
//   - [QueryRepository]: 配置定义读操作
//   - [SettingCategoryQueryRepository]: 配置分类只读操作
//   - [UserSettingCommandRepository]: 用户配置写操作
//   - [UserSettingQueryRepository]: 用户配置读操作
//
// # 与配置文件的区别
//
// 配置文件适合静态配置（数据库连接、端口等），本模块适合需要运行时调整的动态配置。
package setting
