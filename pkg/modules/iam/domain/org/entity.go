package org

import "time"

// Org 组织实体
type Org struct {
	ID        uint       `json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`

	// Name 组织标识符（slug），全局唯一，用于 URL 和 API。
	// 例如: "acme", "github", "google"
	Name string `json:"name"`

	// DisplayName 组织显示名称
	DisplayName string `json:"display_name"`

	// Description 组织描述
	Description string `json:"description"`

	// Avatar 组织头像 URL
	Avatar string `json:"avatar"`

	// Status 组织状态: active, suspended
	Status string `json:"status"`

	// 聚合关系（仅查询时按需加载）
	Teams   []Team   `json:"teams,omitempty"`
	Members []Member `json:"members,omitempty"`
}

// IsActive 报告组织是否处于活跃状态。
func (o *Org) IsActive() bool {
	return o.Status == StatusActive
}

// IsSuspended 报告组织是否被暂停。
func (o *Org) IsSuspended() bool {
	return o.Status == StatusSuspended
}

// Activate 激活组织。
func (o *Org) Activate() {
	o.Status = StatusActive
}

// Suspend 暂停组织。
func (o *Org) Suspend() {
	o.Status = StatusSuspended
}

// 组织状态常量
const (
	StatusActive    = "active"
	StatusSuspended = "suspended"
)
