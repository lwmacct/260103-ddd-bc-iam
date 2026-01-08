// Package org 定义组织和团队领域模型及仓储接口。
//
// # Overview
//
// 本包是多租户系统的核心，定义了：
//   - [Organization]: 组织实体
//   - [Team]: 团队实体（属于组织）
//   - [Member]: 组织成员实体（用户-组织关联）
//   - [TeamMember]: 团队成员实体（用户-团队关联）
//   - [MemberRole]: 成员角色值对象
//
// 关系模型：
//
//	User <--M:N--> Organization (通过 Member)
//	User <--M:N--> Team (通过 TeamMember)
//	Organization --1:N--> Team
//
// # Usage
//
//	// 创建组织实体
//	org := &org.Organization{
//	    Name:  "Acme Corp",
//	    Slug:  "acme-corp",
//	    Owner: &org.Member{
//	        UserID: 123,
//	        Role:   org.RoleOwner,
//	    },
//	}
//
//	// 添加团队
//	team := org.AddTeam("Engineering", org.RoleOwner, 123)
//
// # Thread Safety
//
// 组织、团队、成员实体都是值类型，是并发安全的。
// Repository 接口的实现需要保证并发安全性（由基础设施层负责）。
//
// # 依赖关系
//
// 本包仅定义接口，实现位于 infrastructure/persistence 包。
package org
