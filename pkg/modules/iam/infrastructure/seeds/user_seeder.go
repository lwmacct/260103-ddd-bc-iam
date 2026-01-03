// Package seeds 提供各种领域模型的种子数据
package seeds

import (
	"context"
	"errors"
	"log/slog"

	iampersistence "github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/infrastructure/persistence"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// UserSeeder 用户种子数据
type UserSeeder struct{}

// Seed 执行用户种子数据填充
func (s *UserSeeder) Seed(ctx context.Context, db *gorm.DB) error {
	db = db.WithContext(ctx)

	// 生成密码哈希 (默认密码：admin123)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// 准备种子用户数据
	userData := []struct {
		Username  string
		Email     string
		RealName  string
		Nickname  string
		Phone     string
		Signature string
		Avatar    string
		Bio       string
		Type      string
	}{
		{
			Username:  "admin",
			Email:     "admin@example.com",
			RealName:  "System Administrator",
			Nickname:  "Admin",
			Phone:     "13800138000",
			Signature: "Hello, I am the system administrator.",
			Avatar:    "https://api.dicebear.com/9.x/micah/svg?seed=admin",
			Bio:       "System administrator account",
			Type:      "system",
		},
		{
			Username:  "acme_user",
			Email:     "user@acme.com",
			RealName:  "Acme Employee",
			Nickname:  "AcmeUser",
			Phone:     "13800138101",
			Signature: "Working at Acme",
			Avatar:    "https://api.dicebear.com/9.x/micah/svg?seed=acme_user",
			Bio:       "Acme corporation employee",
			Type:      "human",
		},
		{
			Username:  "globex_user",
			Email:     "user@globex.com",
			RealName:  "Globex Employee",
			Nickname:  "GlobexUser",
			Phone:     "13800138102",
			Signature: "Working at Globex",
			Avatar:    "https://api.dicebear.com/9.x/micah/svg?seed=globex_user",
			Bio:       "Globex corporation employee",
			Type:      "human",
		},
		{
			Username:  "initech_user",
			Email:     "user@initech.com",
			RealName:  "Initech Employee",
			Nickname:  "InitechUser",
			Phone:     "13800138103",
			Signature: "Working at Initech",
			Avatar:    "https://api.dicebear.com/9.x/micah/svg?seed=initech_user",
			Bio:       "Initech employee",
			Type:      "human",
		},
		{
			Username:  "testuser",
			Email:     "test@example.com",
			RealName:  "Test User",
			Nickname:  "Tester",
			Phone:     "13800138001",
			Signature: "Just testing",
			Avatar:    "https://api.dicebear.com/9.x/micah/svg?seed=testuser",
			Bio:       "Test user account",
			Type:      "human",
		},
		{
			Username:  "demo",
			Email:     "demo@example.com",
			RealName:  "Demo User",
			Nickname:  "Demo",
			Phone:     "13800138002",
			Signature: "Welcome to the demo!",
			Avatar:    "https://api.dicebear.com/9.x/micah/svg?seed=demo",
			Bio:       "Demo user account",
			Type:      "human",
		},
	}

	// 转换为 UserModel（Email 和 Phone 使用指针）
	users := make([]iampersistence.UserModel, 0, len(userData))
	for _, u := range userData {
		var emailPtr *string
		if u.Email != "" {
			emailPtr = &u.Email
		}
		var phonePtr *string
		if u.Phone != "" {
			phonePtr = &u.Phone
		}

		users = append(users, iampersistence.UserModel{
			Username:  u.Username,
			Email:     emailPtr,
			Password:  string(hashedPassword),
			RealName:  u.RealName,
			Nickname:  u.Nickname,
			Phone:     phonePtr,
			Signature: u.Signature,
			Avatar:    u.Avatar,
			Bio:       u.Bio,
			Status:    "active",
			Type:      u.Type,
		})
	}

	// 逐个创建或更新用户（不再使用 ON CONFLICT，因为 username 已无唯一约束）
	insertedCount := 0
	for _, userModel := range users {
		var existing iampersistence.UserModel
		switch lookupErr := db.Where("username = ?", userModel.Username).First(&existing).Error; {
		case errors.Is(lookupErr, gorm.ErrRecordNotFound):
			// 用户不存在，创建新用户
			if err := db.Create(&userModel).Error; err != nil {
				return err
			}
			insertedCount++
		case lookupErr == nil:
			// 用户已存在，跳过
			slog.Info("User already exists, skipping", "username", userModel.Username)
		default:
			return lookupErr
		}
	}

	slog.Info("Seeded demo users", "attempted", len(users), "inserted", insertedCount)

	return nil
}

// Name 返回种子器名称。
func (s *UserSeeder) Name() string {
	return "UserSeeder"
}
