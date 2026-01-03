package manualtest

// ExtractIDs 从结构体切片中提取 ID。
// 使用泛型实现，通过 getter 函数获取每个元素的 ID。
//
// 示例:
//
//	ids := ExtractIDs(settings, func(s UserSettingDTO) uint { return s.ID })
//	assert.Contains(t, ids, expectedID)
func ExtractIDs[T any](items []T, getter func(T) uint) []uint {
	ids := make([]uint, len(items))
	for i, item := range items {
		ids[i] = getter(item)
	}
	return ids
}

// ExtractStrings 从结构体切片中提取字符串字段。
// 使用泛型实现，通过 getter 函数获取每个元素的字符串值。
//
// 示例:
//
//	keys := ExtractStrings(settings, func(s UserSettingDTO) string { return s.Key })
//	assert.Contains(t, keys, expectedKey)
func ExtractStrings[T any](items []T, getter func(T) string) []string {
	strs := make([]string, len(items))
	for i, item := range items {
		strs[i] = getter(item)
	}
	return strs
}
