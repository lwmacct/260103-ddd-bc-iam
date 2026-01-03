package lead

import (
	"slices"
	"time"
)

// Status 线索状态。
type Status string

// 线索状态常量
const (
	StatusNew       Status = "new"       // 新建
	StatusContacted Status = "contacted" // 已联系
	StatusQualified Status = "qualified" // 已确认
	StatusConverted Status = "converted" // 已转化
	StatusLost      Status = "lost"      // 已丢失
)

// ValidStatuses 有效的状态值。
var ValidStatuses = []Status{StatusNew, StatusContacted, StatusQualified, StatusConverted, StatusLost}

// IsValidStatus 检查状态值是否有效。
func IsValidStatus(s Status) bool {
	return slices.Contains(ValidStatuses, s)
}

// Lead 线索实体。
type Lead struct {
	ID          uint       `json:"id"`
	Title       string     `json:"title"`        // 线索标题
	ContactID   *uint      `json:"contact_id"`   // 可选关联联系人
	CompanyName string     `json:"company_name"` // 公司名（可能还没创建 Company）
	Source      string     `json:"source"`       // 来源: website/referral/campaign/other
	Status      Status     `json:"status"`       // 状态
	Score       int        `json:"score"`        // 线索评分 0-100
	OwnerID     uint       `json:"owner_id"`     // 负责人
	Notes       string     `json:"notes"`        // 备注
	ConvertedAt *time.Time `json:"converted_at"` // 转化时间
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// CanContact 检查是否可以转换到已联系状态。
func (l *Lead) CanContact() bool {
	return l.Status == StatusNew
}

// Contact 将线索状态转换为已联系。
func (l *Lead) Contact() error {
	if !l.CanContact() {
		return ErrCannotContact
	}
	l.Status = StatusContacted
	return nil
}

// CanQualify 检查是否可以转换到已确认状态。
func (l *Lead) CanQualify() bool {
	return l.Status == StatusContacted
}

// Qualify 将线索状态转换为已确认。
func (l *Lead) Qualify() error {
	if !l.CanQualify() {
		return ErrCannotQualify
	}
	l.Status = StatusQualified
	return nil
}

// CanConvert 检查是否可以转化。
func (l *Lead) CanConvert() bool {
	return l.Status == StatusQualified
}

// Convert 将线索转化为商机。
func (l *Lead) Convert() error {
	if !l.CanConvert() {
		return ErrCannotConvert
	}
	l.Status = StatusConverted
	now := time.Now()
	l.ConvertedAt = &now
	return nil
}

// CanLose 检查是否可以标记为丢失。
func (l *Lead) CanLose() bool {
	return l.Status == StatusNew || l.Status == StatusContacted || l.Status == StatusQualified
}

// Lose 将线索标记为丢失。
func (l *Lead) Lose() error {
	if !l.CanLose() {
		return ErrCannotLose
	}
	l.Status = StatusLost
	return nil
}

// IsClosed 检查线索是否已关闭（转化或丢失）。
func (l *Lead) IsClosed() bool {
	return l.Status == StatusConverted || l.Status == StatusLost
}
