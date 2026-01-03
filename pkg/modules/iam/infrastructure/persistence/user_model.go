package persistence

import (
	"encoding/json"
	"time"

	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/iam/domain/user"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// UserModel 定义用户的 GORM 持久化模型
//
//nolint:recvcheck // TableName uses value receiver per GORM convention
type UserModel struct {
	ID        uint    `gorm:"primaryKey"`
	Username  string  `gorm:"size:50;not null;index"` // 移除 uniqueIndex，改为部分唯一索引
	Email     *string `gorm:"size:100;index"`         // nullable，移除 uniqueIndex
	Password  string  `gorm:"size:255;not null"`
	RealName  string  `gorm:"size:100"`      // 真实姓名
	Nickname  string  `gorm:"size:50"`       // 昵称
	Phone     *string `gorm:"size:20;index"` // nullable，移除 uniqueIndex
	Signature string  `gorm:"size:255"`      // 个性签名
	Avatar    string  `gorm:"size:255"`
	Bio       string  `gorm:"type:text"` // 个人简介
	Status    string  `gorm:"size:20;default:'active'"`

	// Type 用户类型：human（人类用户）、service（服务账户）、system（系统用户）
	Type string `gorm:"size:20;default:'human';not null;index"`

	// Extra 扩展数据（JSONB 存储），留空备用
	Extra datatypes.JSON `gorm:"type:jsonb;default:'{}'"`

	Roles     []RoleModel `gorm:"many2many:user_roles;joinForeignKey:UserID;joinReferences:RoleID;foreignKey:ID;references:ID"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// TableName 指定用户表名
func (UserModel) TableName() string {
	return "users"
}

func newUserModelFromEntity(entity *user.User) *UserModel {
	if entity == nil {
		return nil
	}

	model := &UserModel{
		ID:        entity.ID,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
		Username:  entity.Username,
		Password:  entity.Password,
		RealName:  entity.RealName,
		Nickname:  entity.Nickname,
		Signature: entity.Signature,
		Avatar:    entity.Avatar,
		Bio:       entity.Bio,
		Status:    entity.Status,
		Type:      string(entity.Type),
		// Roles 不在这里映射，通过 user_roles 关联表管理
	}

	// 处理指针类型字段
	if entity.Email != nil {
		model.Email = entity.Email
	}
	if entity.Phone != nil {
		model.Phone = entity.Phone
	}

	// 处理 Extra 字段
	if len(entity.Extra) > 0 {
		extraJSON, err := json.Marshal(entity.Extra)
		if err != nil {
			// 序列化失败时忽略 Extra 字段
			extraJSON = []byte("{}")
		}
		model.Extra = datatypes.JSON(extraJSON)
	}

	if entity.DeletedAt != nil {
		model.DeletedAt = gorm.DeletedAt{Time: *entity.DeletedAt, Valid: true}
	}

	return model
}

// ToEntity 将 GORM Model 转换为 Domain Entity（实现 Model[E] 接口）
func (m *UserModel) ToEntity() *user.User {
	if m == nil {
		return nil
	}

	entity := &user.User{
		ID:        m.ID,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
		Username:  m.Username,
		Password:  m.Password,
		RealName:  m.RealName,
		Nickname:  m.Nickname,
		Signature: m.Signature,
		Avatar:    m.Avatar,
		Bio:       m.Bio,
		Status:    m.Status,
		Type:      user.UserType(m.Type),
		Roles:     mapRoleModelsToEntities(m.Roles),
	}

	// 处理指针类型字段
	if m.Email != nil {
		entity.Email = m.Email
	}
	if m.Phone != nil {
		entity.Phone = m.Phone
	}

	// 处理 Extra 字段
	if len(m.Extra) > 0 {
		_ = json.Unmarshal(m.Extra, &entity.Extra)
	}

	if m.DeletedAt.Valid {
		t := m.DeletedAt.Time
		entity.DeletedAt = &t
	}

	return entity
}

func mapUserModelsToEntities(models []UserModel) []*user.User {
	if len(models) == 0 {
		return nil
	}

	users := make([]*user.User, 0, len(models))
	for i := range models {
		if entity := models[i].ToEntity(); entity != nil {
			users = append(users, entity)
		}
	}
	return users
}
