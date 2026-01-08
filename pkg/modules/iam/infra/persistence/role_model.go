package persistence

import (
	"encoding/json"
	"time"

	"github.com/lwmacct/260103-ddd-iam-bc/pkg/modules/iam/domain/role"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// ============================================================================
// Role Model
// ============================================================================

// RoleModel 定义角色的 GORM 持久化模型
//
// Permissions 字段使用 JSONB 存储权限列表，替代原来的 role_permissions 关联表。
// 结构示例：[{"operation_pattern":"sys:*:*","resource_pattern":"*"}]
//
//nolint:recvcheck // TableName uses value receiver per GORM convention
type RoleModel struct {
	ID          uint `gorm:"primaryKey"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	Name        string         `gorm:"size:50;uniqueIndex;not null"`
	DisplayName string         `gorm:"size:100;not null"`
	Description string         `gorm:"size:255"`
	IsSystem    bool           `gorm:"default:false;not null"`
	Permissions datatypes.JSON `gorm:"type:jsonb;default:'[]';not null"`
}

// TableName 指定角色表名
func (RoleModel) TableName() string {
	return "roles"
}

// permissionJSON 用于 JSONB 序列化/反序列化的内部结构
type permissionJSON struct {
	OperationPattern string `json:"operation_pattern"`
	ResourcePattern  string `json:"resource_pattern"`
}

func newRoleModelFromEntity(entity *role.Role) *RoleModel {
	if entity == nil {
		return nil
	}

	model := &RoleModel{
		ID:          entity.ID,
		CreatedAt:   entity.CreatedAt,
		UpdatedAt:   entity.UpdatedAt,
		Name:        entity.Name,
		DisplayName: entity.DisplayName,
		Description: entity.Description,
		IsSystem:    entity.IsSystem,
		Permissions: marshalPermissions(entity.Permissions),
	}

	if entity.DeletedAt != nil {
		model.DeletedAt = gorm.DeletedAt{Time: *entity.DeletedAt, Valid: true}
	}

	return model
}

// ToEntity 将 GORM Model 转换为 Domain Entity（含权限）
func (m *RoleModel) ToEntity() *role.Role {
	if m == nil {
		return nil
	}

	entity := &role.Role{
		ID:          m.ID,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
		Name:        m.Name,
		DisplayName: m.DisplayName,
		Description: m.Description,
		IsSystem:    m.IsSystem,
		Permissions: unmarshalPermissions(m.Permissions),
	}

	if m.DeletedAt.Valid {
		t := m.DeletedAt.Time
		entity.DeletedAt = &t
	}

	return entity
}

func mapRoleModelsToEntities(models []RoleModel) []role.Role {
	if len(models) == 0 {
		return nil
	}

	roles := make([]role.Role, 0, len(models))
	for i := range models {
		if entity := models[i].ToEntity(); entity != nil {
			roles = append(roles, *entity)
		}
	}
	return roles
}

// ============================================================================
// Permission JSONB 序列化/反序列化
// ============================================================================

// marshalPermissions 将权限列表序列化为 JSONB
func marshalPermissions(permissions []role.Permission) datatypes.JSON {
	if len(permissions) == 0 {
		return datatypes.JSON("[]")
	}

	items := make([]permissionJSON, len(permissions))
	for i, p := range permissions {
		resPattern := p.ResourcePattern
		if resPattern == "" {
			resPattern = "*"
		}
		items[i] = permissionJSON{
			OperationPattern: p.OperationPattern,
			ResourcePattern:  resPattern,
		}
	}

	data, err := json.Marshal(items)
	if err != nil {
		return datatypes.JSON("[]")
	}
	return data
}

// unmarshalPermissions 从 JSONB 反序列化权限列表
func unmarshalPermissions(data datatypes.JSON) []role.Permission {
	if len(data) == 0 {
		return nil
	}

	var items []permissionJSON
	if err := json.Unmarshal(data, &items); err != nil {
		return nil
	}

	if len(items) == 0 {
		return nil
	}

	permissions := make([]role.Permission, len(items))
	for i, item := range items {
		permissions[i] = role.Permission{
			OperationPattern: item.OperationPattern,
			ResourcePattern:  item.ResourcePattern,
		}
	}
	return permissions
}
