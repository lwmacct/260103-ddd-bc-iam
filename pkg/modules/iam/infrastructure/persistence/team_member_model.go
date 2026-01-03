package persistence

import (
	"time"

	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/iam/domain/org"
)

// TeamMemberModel 定义团队成员的 GORM 持久化模型
//
//nolint:recvcheck // TableName uses value receiver per GORM convention
type TeamMemberModel struct {
	ID       uint      `gorm:"primaryKey"`
	TeamID   uint      `gorm:"index:idx_team_member,unique,priority:1;index;not null"`
	UserID   uint      `gorm:"index:idx_team_member,unique,priority:2;index;not null"`
	Role     string    `gorm:"size:20;default:'member';not null"`
	JoinedAt time.Time `gorm:"not null"`

	CreatedAt time.Time
}

// TableName 指定团队成员表名
func (TeamMemberModel) TableName() string {
	return "team_members"
}

func newTeamMemberModelFromEntity(entity *org.TeamMember) *TeamMemberModel {
	if entity == nil {
		return nil
	}

	return &TeamMemberModel{
		ID:        entity.ID,
		TeamID:    entity.TeamID,
		UserID:    entity.UserID,
		Role:      string(entity.Role),
		JoinedAt:  entity.JoinedAt,
		CreatedAt: entity.CreatedAt,
	}
}

// ToEntity 将 GORM Model 转换为 Domain Entity
func (m *TeamMemberModel) ToEntity() *org.TeamMember {
	if m == nil {
		return nil
	}

	return &org.TeamMember{
		ID:        m.ID,
		TeamID:    m.TeamID,
		UserID:    m.UserID,
		Role:      org.TeamMemberRole(m.Role),
		JoinedAt:  m.JoinedAt,
		CreatedAt: m.CreatedAt,
	}
}

func mapTeamMemberModelsToEntities(models []TeamMemberModel) []org.TeamMember {
	if len(models) == 0 {
		return nil
	}

	members := make([]org.TeamMember, 0, len(models))
	for i := range models {
		if entity := models[i].ToEntity(); entity != nil {
			members = append(members, *entity)
		}
	}
	return members
}

func mapTeamMemberModelPtrsToEntities(models []*TeamMemberModel) []*org.TeamMember {
	if len(models) == 0 {
		return nil
	}

	members := make([]*org.TeamMember, 0, len(models))
	for _, m := range models {
		if entity := m.ToEntity(); entity != nil {
			members = append(members, entity)
		}
	}
	return members
}
