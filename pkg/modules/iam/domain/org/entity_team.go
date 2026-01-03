package org

import "time"

// Team 团队实体
type Team struct {
	ID        uint       `json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`

	// OrgID 所属组织 ID
	OrgID uint `json:"org_id"`

	// Name 团队标识符（slug），组织内唯一。
	// 例如: "engineering", "marketing", "sales"
	Name string `json:"name"`

	// DisplayName 团队显示名称
	DisplayName string `json:"display_name"`

	// Description 团队描述
	Description string `json:"description"`

	// Avatar 团队头像 URL
	Avatar string `json:"avatar"`

	// 聚合关系（仅查询时按需加载）
	Members []TeamMember `json:"members,omitempty"`
}

// BelongsTo 报告团队是否属于指定组织。
func (t *Team) BelongsTo(orgID uint) bool {
	return t.OrgID == orgID
}
