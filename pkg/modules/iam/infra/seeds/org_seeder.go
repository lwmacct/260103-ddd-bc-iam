package seeds

import (
	"context"
	"log/slog"
	"time"

	persistence "github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/infra/persistence"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// OrganizationSeeder 创建示例组织和团队。
// 依赖 RBACSeeder 和 UserSeeder 已创建 admin 用户。
//
// 创建内容：
//   - 3 个组织（acme, globex, initech）
//   - 每个组织 2 个团队（engineering, product）
//   - admin 用户作为每个组织的 owner
type OrganizationSeeder struct{}

// orgConfig 组织配置
type orgConfig struct {
	name        string
	displayName string
	description string
	memberUser  string // 组织专属成员用户名
}

// teamConfig 团队配置
type teamConfig struct {
	name        string
	displayName string
	description string
}

// Seed implements database.Seeder interface.
func (s *OrganizationSeeder) Seed(ctx context.Context, db *gorm.DB) error {
	// 组织配置
	orgs := []orgConfig{
		{
			name:        "acme",
			displayName: "Acme Corporation",
			description: "Acme 示例组织",
			memberUser:  "acme_user",
		},
		{
			name:        "globex",
			displayName: "Globex Corporation",
			description: "Globex 示例组织",
			memberUser:  "globex_user",
		},
		{
			name:        "initech",
			displayName: "Initech Inc",
			description: "Initech 示例组织",
			memberUser:  "initech_user",
		},
	}

	// 团队配置（每个组织都会创建这些团队）
	teams := []teamConfig{
		{
			name:        "engineering",
			displayName: "Engineering Team",
			description: "工程团队",
		},
		{
			name:        "product",
			displayName: "Product Team",
			description: "产品团队",
		},
	}

	return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. 查找 admin 用户
		var admin persistence.UserModel
		if err := tx.Where("username = ?", "admin").First(&admin).Error; err != nil {
			slog.Warn("admin user not found, skipping member seeding", "error", err.Error())
			return nil
		}

		// 2. 创建组织和团队
		for _, orgCfg := range orgs {
			org := &persistence.OrgModel{
				Name:        orgCfg.name,
				DisplayName: orgCfg.displayName,
				Description: orgCfg.description,
				Avatar:      "https://api.dicebear.com/9.x/identicon/svg?seed=" + orgCfg.name,
				Status:      "active",
			}
			if err := tx.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "name"}},
				DoNothing: true,
			}).Create(org).Error; err != nil {
				return err
			}

			// 如果是冲突跳过，需要重新查询获取 ID
			if org.ID == 0 {
				if err := tx.Where("name = ?", orgCfg.name).First(org).Error; err != nil {
					return err
				}
			}
			slog.Info("seeded organization", "name", org.Name, "id", org.ID)

			// 3. 查找组织专属成员用户
			var orgMember persistence.UserModel
			if err := tx.Where("username = ?", orgCfg.memberUser).First(&orgMember).Error; err != nil {
				slog.Warn("org member user not found", "user", orgCfg.memberUser, "error", err.Error())
				// 继续处理，不中断
			}

			// 4. 添加组织成员：admin 作为 owner
			members := []struct {
				userID uint
				role   string
			}{
				{admin.ID, "owner"},
			}
			if orgMember.ID > 0 {
				members = append(members, struct {
					userID uint
					role   string
				}{orgMember.ID, "member"})
			}

			for _, m := range members {
				member := &persistence.OrgMemberModel{
					OrgID:    org.ID,
					UserID:   m.userID,
					Role:     m.role,
					JoinedAt: time.Now(),
				}
				if err := tx.Clauses(clause.OnConflict{
					Columns:   []clause.Column{{Name: "org_id"}, {Name: "user_id"}},
					DoNothing: true,
				}).Create(member).Error; err != nil {
					return err
				}
				username := admin.Username
				if m.userID == orgMember.ID {
					username = orgMember.Username
				}
				slog.Info("seeded organization member", "org", org.Name, "user", username, "role", member.Role)
			}

			// 5. 为每个组织创建 2 个团队
			for _, teamCfg := range teams {
				team := &persistence.TeamModel{
					OrgID:       org.ID,
					Name:        teamCfg.name,
					DisplayName: teamCfg.displayName,
					Description: teamCfg.description,
				}
				if err := tx.Clauses(clause.OnConflict{
					Columns:   []clause.Column{{Name: "org_id"}, {Name: "name"}},
					DoNothing: true,
				}).Create(team).Error; err != nil {
					return err
				}

				// 如果是冲突跳过，需要重新查询获取 ID
				if team.ID == 0 {
					if err := tx.Where("org_id = ? AND name = ?", org.ID, teamCfg.name).First(team).Error; err != nil {
						return err
					}
				}
				slog.Info("seeded team", "name", team.Name, "org", org.Name, "id", team.ID)

				// 6. 将用户加入团队
				// admin: engineering 团队为 lead，product 为 member
				// orgMember: 两个团队都是 member
				teamMembers := []struct {
					userID uint
					role   string
				}{
					{admin.ID, "member"},
				}
				if teamCfg.name == "engineering" {
					teamMembers[0].role = "lead"
				}
				if orgMember.ID > 0 {
					teamMembers = append(teamMembers, struct {
						userID uint
						role   string
					}{orgMember.ID, "member"})
				}

				for _, tm := range teamMembers {
					teamMember := &persistence.TeamMemberModel{
						TeamID:   team.ID,
						UserID:   tm.userID,
						Role:     tm.role,
						JoinedAt: time.Now(),
					}
					if err := tx.Clauses(clause.OnConflict{
						Columns:   []clause.Column{{Name: "team_id"}, {Name: "user_id"}},
						DoNothing: true,
					}).Create(teamMember).Error; err != nil {
						return err
					}
					username := admin.Username
					if tm.userID == orgMember.ID {
						username = orgMember.Username
					}
					slog.Info("seeded team member", "team", team.Name, "user", username, "role", teamMember.Role)
				}
			}
		}

		return nil
	})
}

// Name implements database.Seeder interface.
func (s *OrganizationSeeder) Name() string {
	return "OrganizationSeeder"
}
