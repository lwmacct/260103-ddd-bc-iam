package persistence

import (
	"time"

	"github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/domain/audit"
	"gorm.io/gorm"
)

// AuditModel 定义审计日志的 GORM 实体
//
//nolint:recvcheck // TableName uses value receiver per GORM convention
type AuditModel struct {
	ID          uint `gorm:"primaryKey"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	UserID      uint           `gorm:"index;not null"`
	Username    string         `gorm:"size:100;not null"`
	Action      string         `gorm:"size:100;not null"`
	Resource    string         `gorm:"size:100;not null"`
	ResourceID  string         `gorm:"size:100"`
	IPAddress   string         `gorm:"size:45"`
	UserAgent   string         `gorm:"size:255"`
	Details     string         `gorm:"type:text"`
	Status      string         `gorm:"size:20;default:'success'"`
	RequestID   string         `gorm:"size:36;index"`  // UUID/Trace ID 长度
	OperationID string         `gorm:"size:100;index"` // Operation ID 索引用于统计
}

// TableName 指定审计日志表名
func (AuditModel) TableName() string {
	return "audit"
}

func newAuditModelFromEntity(entity *audit.Audit) *AuditModel {
	if entity == nil {
		return nil
	}
	model := &AuditModel{
		ID:          entity.ID,
		CreatedAt:   entity.CreatedAt,
		UpdatedAt:   entity.UpdatedAt,
		UserID:      entity.UserID,
		Username:    entity.Username,
		Action:      entity.Action,
		Resource:    entity.Resource,
		ResourceID:  entity.ResourceID,
		IPAddress:   entity.IPAddress,
		UserAgent:   entity.UserAgent,
		Details:     entity.Details,
		Status:      entity.Status,
		RequestID:   entity.RequestID,
		OperationID: entity.OperationID,
	}
	if entity.DeletedAt != nil {
		model.DeletedAt = gorm.DeletedAt{Time: *entity.DeletedAt, Valid: true}
	}
	return model
}

// ToEntity 将 GORM Model 转换为 Domain Entity（实现 Model[E] 接口）
func (m *AuditModel) ToEntity() *audit.Audit {
	if m == nil {
		return nil
	}
	entity := &audit.Audit{
		ID:          m.ID,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
		UserID:      m.UserID,
		Username:    m.Username,
		Action:      m.Action,
		Resource:    m.Resource,
		ResourceID:  m.ResourceID,
		IPAddress:   m.IPAddress,
		UserAgent:   m.UserAgent,
		Details:     m.Details,
		Status:      m.Status,
		RequestID:   m.RequestID,
		OperationID: m.OperationID,
	}
	if m.DeletedAt.Valid {
		t := m.DeletedAt.Time
		entity.DeletedAt = &t
	}
	return entity
}
