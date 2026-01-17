package cache

import "go.uber.org/fx"

// CacheModule 缓存 Fx 模块
//
// 注意：Settings 缓存服务由容器层提供（来自 260103-ddd-settings-bc 外部模块）。
// 本模块为占位符，保持模块架构一致性。
var CacheModule = fx.Module("settings.cache")
