package opportunity

import (
	"time"

	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/crm/domain/opportunity"
)

// CreateCommand 创建商机命令。
type CreateCommand struct {
	Name          string
	ContactID     uint
	CompanyID     *uint
	LeadID        *uint
	Amount        float64
	Probability   int
	ExpectedClose *time.Time
	OwnerID       uint
	Notes         string
}

// UpdateCommand 更新商机命令。
type UpdateCommand struct {
	ID            uint
	Name          *string
	ContactID     *uint
	CompanyID     *uint
	Amount        *float64
	Probability   *int
	ExpectedClose *time.Time
	Notes         *string
}

// DeleteCommand 删除商机命令。
type DeleteCommand struct {
	ID uint
}

// AdvanceCommand 推进商机阶段命令。
type AdvanceCommand struct {
	ID    uint
	Stage opportunity.Stage
}

// CloseWonCommand 成交命令。
type CloseWonCommand struct {
	ID uint
}

// CloseLostCommand 丢单命令。
type CloseLostCommand struct {
	ID uint
}
