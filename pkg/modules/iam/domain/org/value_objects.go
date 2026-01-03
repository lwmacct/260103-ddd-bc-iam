package org

// MemberWithUser 包含成员及其用户信息（用于列表查询优化）。
//
// 这是 CQRS 查询端的值对象，用于优化成员列表查询，
// 避免前端多次请求获取用户基本信息。
type MemberWithUser struct {
	Member

	// Username 用户名
	Username string
	// Email 邮箱
	Email string
	// FullName 全名
	FullName string
	// Avatar 头像
	Avatar string
}
