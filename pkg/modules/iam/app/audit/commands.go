package audit

// CreateCommand 创建审计日志命令
type CreateCommand struct {
	UserID      uint
	Username    string
	Action      string
	Resource    string
	ResourceID  string
	IPAddress   string
	UserAgent   string
	Details     string
	Status      string
	RequestID   string // 请求追踪 ID
	OperationID string // API 操作标识符
}
