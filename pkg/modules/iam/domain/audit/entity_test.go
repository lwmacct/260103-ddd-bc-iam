package audit

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAuditLog_StatusChecks(t *testing.T) {
	tests := []struct {
		name      string
		status    string
		isSuccess bool
		isFailed  bool
		isPending bool
	}{
		{"success status", StatusSuccess, true, false, false},
		{"failed status", StatusFailed, false, true, false},
		{"pending status", StatusPending, false, false, true},
		{"unknown status", "unknown", false, false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Audit{Status: tt.status}
			assert.Equal(t, tt.isSuccess, a.IsSuccess())
			assert.Equal(t, tt.isFailed, a.IsFailed())
			assert.Equal(t, tt.isPending, a.IsPending())
		})
	}
}

func TestAuditLog_GetResourceIdentifier(t *testing.T) {
	tests := []struct {
		name       string
		resource   string
		resourceID string
		want       string
	}{
		{"with resource id", "user", "123", "user:123"},
		{"without resource id", "settings", "", "settings"},
		{"empty resource with id", "", "456", ":456"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Audit{Resource: tt.resource, ResourceID: tt.resourceID}
			assert.Equal(t, tt.want, a.GetResourceIdentifier())
		})
	}
}

func TestAuditLog_IsRecentlyCreated(t *testing.T) {
	t.Run("recently created", func(t *testing.T) {
		a := &Audit{CreatedAt: time.Now().Add(-5 * time.Minute)}
		assert.True(t, a.IsRecentlyCreated(10*time.Minute))
	})

	t.Run("not recently created", func(t *testing.T) {
		a := &Audit{CreatedAt: time.Now().Add(-1 * time.Hour)}
		assert.False(t, a.IsRecentlyCreated(10*time.Minute))
	})

	t.Run("exactly at boundary", func(t *testing.T) {
		a := &Audit{CreatedAt: time.Now()}
		// With 1 second tolerance, recently created should pass
		assert.True(t, a.IsRecentlyCreated(1*time.Second))
	})
}

func TestAuditLog_HasDetails(t *testing.T) {
	tests := []struct {
		name    string
		details string
		want    bool
	}{
		{"with details", "some details", true},
		{"empty details", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Audit{Details: tt.details}
			assert.Equal(t, tt.want, a.HasDetails())
		})
	}
}

func TestAuditLog_ActionTypes(t *testing.T) {
	t.Run("user action", func(t *testing.T) {
		a := &Audit{UserID: 1}
		assert.True(t, a.IsUserAction())
		assert.False(t, a.IsSystemAction())
	})

	t.Run("system action", func(t *testing.T) {
		a := &Audit{UserID: 0}
		assert.False(t, a.IsUserAction())
		assert.True(t, a.IsSystemAction())
	})
}

func TestAuditLog_MatchesFilter(t *testing.T) {
	now := time.Now()
	userID := uint(1)

	baseLog := &Audit{
		UserID:    1,
		Action:    "create",
		Resource:  "user",
		Status:    StatusSuccess,
		CreatedAt: now,
	}

	tests := []struct {
		name   string
		filter FilterOptions
		want   bool
	}{
		{
			name:   "empty filter matches all",
			filter: FilterOptions{},
			want:   true,
		},
		{
			name:   "matching user id",
			filter: FilterOptions{UserID: &userID},
			want:   true,
		},
		{
			name:   "non-matching user id",
			filter: FilterOptions{UserID: func() *uint { id := uint(999); return &id }()},
			want:   false,
		},
		{
			name:   "matching action",
			filter: FilterOptions{Action: "create"},
			want:   true,
		},
		{
			name:   "non-matching action",
			filter: FilterOptions{Action: "delete"},
			want:   false,
		},
		{
			name:   "matching resource",
			filter: FilterOptions{Resource: "user"},
			want:   true,
		},
		{
			name:   "non-matching resource",
			filter: FilterOptions{Resource: "role"},
			want:   false,
		},
		{
			name:   "matching status",
			filter: FilterOptions{Status: StatusSuccess},
			want:   true,
		},
		{
			name:   "non-matching status",
			filter: FilterOptions{Status: StatusFailed},
			want:   false,
		},
		{
			name:   "within date range",
			filter: FilterOptions{StartDate: func() *time.Time { t := now.Add(-1 * time.Hour); return &t }(), EndDate: func() *time.Time { t := now.Add(1 * time.Hour); return &t }()},
			want:   true,
		},
		{
			name:   "before start date",
			filter: FilterOptions{StartDate: func() *time.Time { t := now.Add(1 * time.Hour); return &t }()},
			want:   false,
		},
		{
			name:   "after end date",
			filter: FilterOptions{EndDate: func() *time.Time { t := now.Add(-1 * time.Hour); return &t }()},
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, baseLog.MatchesFilter(tt.filter))
		})
	}
}

func TestFilterOptions_IsValidFilter(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name   string
		filter FilterOptions
		want   bool
	}{
		{
			name:   "empty filter is valid",
			filter: FilterOptions{},
			want:   true,
		},
		{
			name:   "valid pagination",
			filter: FilterOptions{Page: 1, Limit: 20},
			want:   true,
		},
		{
			name:   "negative page",
			filter: FilterOptions{Page: -1},
			want:   false,
		},
		{
			name:   "negative limit",
			filter: FilterOptions{Limit: -1},
			want:   false,
		},
		{
			name:   "valid date range",
			filter: FilterOptions{StartDate: func() *time.Time { t := now.Add(-1 * time.Hour); return &t }(), EndDate: &now},
			want:   true,
		},
		{
			name:   "invalid date range",
			filter: FilterOptions{StartDate: &now, EndDate: func() *time.Time { t := now.Add(-1 * time.Hour); return &t }()},
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.filter.IsValidFilter())
		})
	}
}

func TestFilterOptions_SetDefaults(t *testing.T) {
	t.Run("sets default page and limit", func(t *testing.T) {
		f := &FilterOptions{}
		f.SetDefaults()
		assert.Equal(t, 1, f.Page)
		assert.Equal(t, 20, f.Limit)
	})

	t.Run("caps limit at 100", func(t *testing.T) {
		f := &FilterOptions{Limit: 200}
		f.SetDefaults()
		assert.Equal(t, 100, f.Limit)
	})

	t.Run("preserves valid values", func(t *testing.T) {
		f := &FilterOptions{Page: 5, Limit: 50}
		f.SetDefaults()
		assert.Equal(t, 5, f.Page)
		assert.Equal(t, 50, f.Limit)
	})

	t.Run("fixes zero page", func(t *testing.T) {
		f := &FilterOptions{Page: 0, Limit: 50}
		f.SetDefaults()
		assert.Equal(t, 1, f.Page)
		assert.Equal(t, 50, f.Limit)
	})
}
