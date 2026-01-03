package persistence

import (
	"time"

	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/iam/domain/org"
)

// OrgMemberModel 定义组织成员的 GORM 持久化模型
//
//nolint:recvcheck // TableName uses value receiver per GORM convention
type OrgMemberModel struct {
	ID       uint      `gorm:"primaryKey"`
	OrgID    uint      `gorm:"index:idx_org_member,unique,priority:1;not null"`
	UserID   uint      `gorm:"index:idx_org_member,unique,priority:2;index;not null"`
	Role     string    `gorm:"size:20;default:'member';not null"`
	JoinedAt time.Time `gorm:"not null"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

// TableName 指定组织成员表名
func (OrgMemberModel) TableName() string {
	return "org_members"
}

func newOrgMemberModelFromEntity(entity *org.Member) *OrgMemberModel {
	if entity == nil {
		return nil
	}

	return &OrgMemberModel{
		ID:        entity.ID,
		OrgID:     entity.OrgID,
		UserID:    entity.UserID,
		Role:      string(entity.Role),
		JoinedAt:  entity.JoinedAt,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
	}
}

// ToEntity 将 GORM Model 转换为 Domain Entity
func (m *OrgMemberModel) ToEntity() *org.Member {
	if m == nil {
		return nil
	}

	return &org.Member{
		ID:        m.ID,
		OrgID:     m.OrgID,
		UserID:    m.UserID,
		Role:      org.MemberRole(m.Role),
		JoinedAt:  m.JoinedAt,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func mapOrgMemberModelsToEntities(models []OrgMemberModel) []org.Member {
	if len(models) == 0 {
		return nil
	}

	members := make([]org.Member, 0, len(models))
	for i := range models {
		if entity := models[i].ToEntity(); entity != nil {
			members = append(members, *entity)
		}
	}
	return members
}

func mapOrgMemberModelPtrsToEntities(models []*OrgMemberModel) []*org.Member {
	if len(models) == 0 {
		return nil
	}

	members := make([]*org.Member, 0, len(models))
	for _, m := range models {
		if entity := m.ToEntity(); entity != nil {
			members = append(members, entity)
		}
	}
	return members
}
