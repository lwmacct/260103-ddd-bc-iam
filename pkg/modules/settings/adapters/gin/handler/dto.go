package handler

// SetRequest 设置配置请求
type SetRequest struct {
	Value any `json:"value" binding:"required"`
}

// BatchSetRequest 批量设置配置请求
type BatchSetRequest struct {
	Settings []BatchSetItem `json:"settings" binding:"required,min=1"`
}

// BatchSetItem 批量设置配置项
type BatchSetItem struct {
	Key   string `json:"key" binding:"required"`
	Value any    `json:"value" binding:"required"`
}
