package usersetting

// UserSettingDTO 用户设置响应 DTO
type UserSettingDTO struct {
	Key        string `json:"key"`
	Value      string `json:"value"`
	CategoryID uint   `json:"category_id"`
	ValueType  string `json:"value_type"` // 从 Schema 继承（string, boolean, number, json）
	Label      string `json:"label"`      // 从 Schema 继承
	IsCustom   bool   `json:"is_custom"`  // true=用户自定义，false=系统默认值
}

// UserSettingListDTO 用户设置列表响应 DTO
type UserSettingListDTO struct {
	Settings []*UserSettingDTO `json:"settings"`
	Total    int64             `json:"total"`
}

// UpdateDTO 更新用户设置请求 DTO
type UpdateDTO struct {
	Value string `json:"value" binding:"required" example:"true"`
}
