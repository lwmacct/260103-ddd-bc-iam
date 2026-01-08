package settings

// 本包的提取辅助函数直接使用 apitest 包实现。
// 使用方式：
//
//	import "github.com/lwmacct/260103-ddd-shared/pkg/shared/apitest"
//
//	ids := apitest.ExtractIDs(settings, func(s SettingDTO) uint { return s.ID })
//	keys := apitest.ExtractStrings(settings, func(s SettingDTO) string { return s.Key })
