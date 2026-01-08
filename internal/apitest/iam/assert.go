package iam

// 本包的提取辅助函数直接使用 apitest 包实现。
// 使用方式：
//
//	import "github.com/lwmacct/260103-ddd-shared/pkg/shared/apitest"
//
//	ids := apitest.ExtractIDs(users, func(u UserDTO) uint { return u.ID })
//	names := apitest.ExtractStrings(users, func(u UserDTO) string { return u.Username })
