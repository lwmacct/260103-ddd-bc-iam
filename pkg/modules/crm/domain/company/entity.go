package company

import (
	"slices"
	"time"
)

// Company 公司实体。
type Company struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`     // 公司名称（唯一）
	Industry  string    `json:"industry"` // 行业
	Size      string    `json:"size"`     // 规模: small/medium/large/enterprise
	Website   string    `json:"website"`  // 网站
	Address   string    `json:"address"`  // 地址
	OwnerID   uint      `json:"owner_id"` // 负责人（用户 ID）
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ValidSizes 有效的公司规模值。
var ValidSizes = []string{"small", "medium", "large", "enterprise"}

// IsValidSize 检查规模值是否有效。
func IsValidSize(size string) bool {
	return slices.Contains(ValidSizes, size)
}
