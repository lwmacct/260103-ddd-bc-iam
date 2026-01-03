package event

// ============================================================================
// 审计事件
// ============================================================================

// AuditableAction 可审计的操作类型
type AuditableAction string

const (
	ActionCreate AuditableAction = "create"
	ActionUpdate AuditableAction = "update"
	ActionDelete AuditableAction = "delete"
	ActionAssign AuditableAction = "assign"
	ActionRevoke AuditableAction = "revoke"
	ActionLogin  AuditableAction = "login"
	ActionLogout AuditableAction = "logout"
)

// CommandExecutedEvent 命令执行完成事件
// 用于记录业务操作的审计日志
type CommandExecutedEvent struct {
	BaseEvent

	// 操作者信息
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`

	// 操作信息
	Action     AuditableAction `json:"action"`      // 操作类型
	Resource   string          `json:"resource"`    // 资源类型（user, role, menu 等）
	ResourceID string          `json:"resource_id"` // 资源ID

	// 执行结果
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`

	// 请求上下文
	IPAddress string `json:"ip_address,omitempty"`
	UserAgent string `json:"user_agent,omitempty"`

	// 额外详情（JSON 格式）
	Details string `json:"details,omitempty"`
}

// NewCommandExecutedEvent 创建命令执行事件
func NewCommandExecutedEvent(
	userID uint,
	username string,
	action AuditableAction,
	resource string,
	resourceID string,
	success bool,
) *CommandExecutedEvent {
	return &CommandExecutedEvent{
		BaseEvent:  NewBaseEvent("audit.command_executed", resource, resourceID),
		UserID:     userID,
		Username:   username,
		Action:     action,
		Resource:   resource,
		ResourceID: resourceID,
		Success:    success,
	}
}

// WithError 设置错误信息
func (e *CommandExecutedEvent) WithError(err error) *CommandExecutedEvent {
	if err != nil {
		e.Error = err.Error()
		e.Success = false
	}
	return e
}

// WithDetails 设置详情
func (e *CommandExecutedEvent) WithDetails(details string) *CommandExecutedEvent {
	e.Details = details
	return e
}

// WithRequestContext 设置请求上下文
func (e *CommandExecutedEvent) WithRequestContext(ipAddress, userAgent string) *CommandExecutedEvent {
	e.IPAddress = ipAddress
	e.UserAgent = userAgent
	return e
}
