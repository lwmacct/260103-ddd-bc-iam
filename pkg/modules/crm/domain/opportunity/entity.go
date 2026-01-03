package opportunity

import (
	"time"
)

// Stage 商机阶段。
type Stage string

// 商机阶段常量。
const (
	StageProspecting Stage = "prospecting" // 初步接触
	StageProposal    Stage = "proposal"    // 提案阶段
	StageNegotiation Stage = "negotiation" // 谈判阶段
	StageClosedWon   Stage = "closed_won"  // 成交
	StageClosedLost  Stage = "closed_lost" // 丢单
)

// Opportunity 商机实体。
type Opportunity struct {
	ID            uint       `json:"id"`
	Name          string     `json:"name"`           // 商机名称
	ContactID     uint       `json:"contact_id"`     // 必须关联联系人
	CompanyID     *uint      `json:"company_id"`     // 可选关联公司
	LeadID        *uint      `json:"lead_id"`        // 来源线索（如果从 Lead 转化）
	Stage         Stage      `json:"stage"`          // 当前阶段
	Amount        float64    `json:"amount"`         // 预计金额
	Probability   int        `json:"probability"`    // 成交概率 0-100
	ExpectedClose *time.Time `json:"expected_close"` // 预计成交日期
	OwnerID       uint       `json:"owner_id"`       // 负责人
	Notes         string     `json:"notes"`          // 备注
	ClosedAt      *time.Time `json:"closed_at"`      // 实际关闭时间
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

// WeightedAmount 返回加权金额（Amount × Probability / 100）。
func (o *Opportunity) WeightedAmount() float64 {
	return o.Amount * float64(o.Probability) / 100
}

// IsClosed 报告商机是否已关闭。
func (o *Opportunity) IsClosed() bool {
	return o.Stage == StageClosedWon || o.Stage == StageClosedLost
}

// CanAdvanceTo 报告是否可以推进到指定阶段。
func (o *Opportunity) CanAdvanceTo(stage Stage) bool {
	if o.IsClosed() {
		return false
	}
	// 只能按顺序推进：prospecting → proposal → negotiation
	switch o.Stage {
	case StageProspecting:
		return stage == StageProposal
	case StageProposal:
		return stage == StageNegotiation
	case StageNegotiation:
		return false // 不能通过 Advance 到 closed_*, 必须用 CloseWon/CloseLost
	case StageClosedWon, StageClosedLost:
		return false // 已关闭的商机不能推进
	}
	return false
}

// AdvanceTo 将商机推进到指定阶段。
func (o *Opportunity) AdvanceTo(stage Stage) error {
	if !o.CanAdvanceTo(stage) {
		return ErrInvalidStageTransition
	}
	o.Stage = stage
	return nil
}

// CanClose 报告是否可以关闭商机。
func (o *Opportunity) CanClose() bool {
	// 只有在 negotiation 阶段才能关闭
	return o.Stage == StageNegotiation
}

// CloseWon 将商机标记为成交。
func (o *Opportunity) CloseWon() error {
	if !o.CanClose() {
		return ErrCannotClose
	}
	now := time.Now()
	o.Stage = StageClosedWon
	o.ClosedAt = &now
	o.Probability = 100
	return nil
}

// CloseLost 将商机标记为丢单。
func (o *Opportunity) CloseLost() error {
	if !o.CanClose() {
		return ErrCannotClose
	}
	now := time.Now()
	o.Stage = StageClosedLost
	o.ClosedAt = &now
	o.Probability = 0
	return nil
}
