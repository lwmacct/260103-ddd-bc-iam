package persistence

import (
	"time"

	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/iam/domain/pat"
	"gorm.io/gorm"
)

// PersonalAccessTokenModel 定义 PAT 的 GORM 实体
//
//nolint:recvcheck // TableName uses value receiver per GORM convention
type PersonalAccessTokenModel struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	UserID      uint   `gorm:"not null;index"`
	Name        string `gorm:"size:100;not null"`
	Token       string `gorm:"size:255;uniqueIndex;not null"`
	TokenPrefix string `gorm:"size:20;not null;index"`

	Scopes pat.StringList `gorm:"type:jsonb;not null;default:'[\"full\"]'"` // 权限范围

	ExpiresAt  *time.Time
	LastUsedAt *time.Time
	Status     string `gorm:"size:20;not null;default:'active';index"`

	IPWhitelist pat.StringList `gorm:"type:jsonb"`
	Description string         `gorm:"type:text"`
}

// TableName 指定 PAT 表名
func (PersonalAccessTokenModel) TableName() string {
	return "personal_access_tokens"
}

func newPATModelFromEntity(entity *pat.PersonalAccessToken) *PersonalAccessTokenModel {
	if entity == nil {
		return nil
	}

	model := &PersonalAccessTokenModel{
		ID:          entity.ID,
		CreatedAt:   entity.CreatedAt,
		UpdatedAt:   entity.UpdatedAt,
		UserID:      entity.UserID,
		Name:        entity.Name,
		Token:       entity.Token,
		TokenPrefix: entity.TokenPrefix,
		Scopes:      entity.Scopes,
		ExpiresAt:   entity.ExpiresAt,
		LastUsedAt:  entity.LastUsedAt,
		Status:      entity.Status,
		IPWhitelist: entity.IPWhitelist,
		Description: entity.Description,
	}

	if entity.DeletedAt != nil {
		model.DeletedAt = gorm.DeletedAt{Time: *entity.DeletedAt, Valid: true}
	}

	return model
}

// ToEntity 将 GORM Model 转换为 Domain Entity（实现 Model[E] 接口）
func (m *PersonalAccessTokenModel) ToEntity() *pat.PersonalAccessToken {
	if m == nil {
		return nil
	}

	entity := &pat.PersonalAccessToken{
		ID:          m.ID,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
		UserID:      m.UserID,
		Name:        m.Name,
		Token:       m.Token,
		TokenPrefix: m.TokenPrefix,
		Scopes:      m.Scopes,
		ExpiresAt:   m.ExpiresAt,
		LastUsedAt:  m.LastUsedAt,
		Status:      m.Status,
		IPWhitelist: m.IPWhitelist,
		Description: m.Description,
	}

	if m.DeletedAt.Valid {
		t := m.DeletedAt.Time
		entity.DeletedAt = &t
	}

	return entity
}

func mapPATModelsToEntities(models []PersonalAccessTokenModel) []*pat.PersonalAccessToken {
	if len(models) == 0 {
		return nil
	}

	entities := make([]*pat.PersonalAccessToken, 0, len(models))
	for i := range models {
		if entity := models[i].ToEntity(); entity != nil {
			entities = append(entities, entity)
		}
	}
	return entities
}
