package cache

import "go.uber.org/fx"

// CacheModule 缓存 Fx 模块
// TODO: 后续添加缓存服务实现
var CacheModule = fx.Module("settings.cache")
