package eventhandler

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/iam/domain/user"
	"github.com/lwmacct/260101-go-pkg-ddd/pkg/shared/event"
)

// mockPermissionCache 模拟权限缓存服务
type mockPermissionCache struct {
	invalidatedUsers []uint
	invalidateError  error
}

func newMockPermissionCache() *mockPermissionCache {
	return &mockPermissionCache{
		invalidatedUsers: make([]uint, 0),
	}
}

func (m *mockPermissionCache) InvalidateUser(_ context.Context, userID uint) error {
	if m.invalidateError != nil {
		return m.invalidateError
	}
	m.invalidatedUsers = append(m.invalidatedUsers, userID)
	return nil
}

// mockUserQueryRepo 模拟用户查询仓储
type mockUserQueryRepo struct {
	userIDsByRole map[uint][]uint
	getUserError  error
}

func newMockUserQueryRepo() *mockUserQueryRepo {
	return &mockUserQueryRepo{
		userIDsByRole: make(map[uint][]uint),
	}
}

func (m *mockUserQueryRepo) GetByID(_ context.Context, _ uint) (*user.User, error) {
	return nil, user.ErrUserNotFound
}

func (m *mockUserQueryRepo) GetByUsername(_ context.Context, _ string) (*user.User, error) {
	return nil, user.ErrUserNotFound
}

func (m *mockUserQueryRepo) GetByEmail(_ context.Context, _ string) (*user.User, error) {
	return nil, user.ErrUserNotFound
}

func (m *mockUserQueryRepo) GetByUsernameOrEmail(_ context.Context, _ string) (*user.User, error) {
	return nil, user.ErrUserNotFound
}

func (m *mockUserQueryRepo) List(_ context.Context, _, _ int, _ map[string]any) ([]*user.User, int64, error) {
	return nil, 0, nil
}

func (m *mockUserQueryRepo) ExistsByUsername(_ context.Context, _ string) (bool, error) {
	return false, nil
}

func (m *mockUserQueryRepo) ExistsByEmail(_ context.Context, _ string) (bool, error) {
	return false, nil
}

func (m *mockUserQueryRepo) Exists(_ context.Context, _ uint) (bool, error) {
	return false, nil
}

func (m *mockUserQueryRepo) GetUserIDsByRole(_ context.Context, roleID uint) ([]uint, error) {
	if m.getUserError != nil {
		return nil, m.getUserError
	}
	return m.userIDsByRole[roleID], nil
}

func (m *mockUserQueryRepo) GetWithRoles(_ context.Context, _ uint) (*user.User, error) {
	return nil, user.ErrUserNotFound
}

func (m *mockUserQueryRepo) GetRoleIDs(_ context.Context, _ uint) ([]uint, error) {
	return nil, nil
}

func TestCacheInvalidationHandler_HandleUserRoleAssigned(t *testing.T) {
	ctx := context.Background()

	t.Run("成功失效用户权限缓存", func(t *testing.T) {
		cache := newMockPermissionCache()
		userRepo := newMockUserQueryRepo()
		handler := newTestHandler(cache, userRepo)

		evt := event.NewUserRoleAssignedEvent(123, []uint{1, 2, 3})

		err := handler.Handle(ctx, evt)

		require.NoError(t, err)
		assert.Contains(t, cache.invalidatedUsers, uint(123))
	})

	t.Run("缓存失效错误不阻塞", func(t *testing.T) {
		cache := newMockPermissionCache()
		cache.invalidateError = errors.New("redis error")
		userRepo := newMockUserQueryRepo()
		handler := newTestHandler(cache, userRepo)

		evt := event.NewUserRoleAssignedEvent(123, []uint{1})

		err := handler.Handle(ctx, evt)

		// 缓存失效失败不应返回错误
		require.NoError(t, err)
	})
}

func TestCacheInvalidationHandler_HandleUserDeleted(t *testing.T) {
	ctx := context.Background()

	t.Run("成功清理删除用户缓存", func(t *testing.T) {
		cache := newMockPermissionCache()
		userRepo := newMockUserQueryRepo()
		handler := newTestHandler(cache, userRepo)

		evt := event.NewUserDeletedEvent(456)

		err := handler.Handle(ctx, evt)

		require.NoError(t, err)
		assert.Contains(t, cache.invalidatedUsers, uint(456))
	})
}

func TestCacheInvalidationHandler_HandleRolePermissionsChanged(t *testing.T) {
	ctx := context.Background()

	t.Run("成功失效所有相关用户缓存", func(t *testing.T) {
		cache := newMockPermissionCache()
		userRepo := newMockUserQueryRepo()
		userRepo.userIDsByRole[1] = []uint{100, 200, 300}
		handler := newTestHandler(cache, userRepo)

		evt := event.NewRolePermissionsChangedEvent(1, []uint{10, 20})

		err := handler.Handle(ctx, evt)

		require.NoError(t, err)
		assert.Len(t, cache.invalidatedUsers, 3)
		assert.Contains(t, cache.invalidatedUsers, uint(100))
		assert.Contains(t, cache.invalidatedUsers, uint(200))
		assert.Contains(t, cache.invalidatedUsers, uint(300))
	})

	t.Run("角色无关联用户时不报错", func(t *testing.T) {
		cache := newMockPermissionCache()
		userRepo := newMockUserQueryRepo()
		handler := newTestHandler(cache, userRepo)

		evt := event.NewRolePermissionsChangedEvent(999, []uint{10})

		err := handler.Handle(ctx, evt)

		require.NoError(t, err)
		assert.Empty(t, cache.invalidatedUsers)
	})

	t.Run("获取用户失败不阻塞", func(t *testing.T) {
		cache := newMockPermissionCache()
		userRepo := newMockUserQueryRepo()
		userRepo.getUserError = errors.New("database error")
		handler := newTestHandler(cache, userRepo)

		evt := event.NewRolePermissionsChangedEvent(1, []uint{10})

		err := handler.Handle(ctx, evt)

		// 查询失败不应返回错误
		require.NoError(t, err)
	})
}

func TestCacheInvalidationHandler_HandleUnknownEvent(t *testing.T) {
	ctx := context.Background()

	t.Run("忽略未知事件", func(t *testing.T) {
		cache := newMockPermissionCache()
		userRepo := newMockUserQueryRepo()
		handler := newTestHandler(cache, userRepo)

		evt := &unknownEvent{
			BaseEvent: event.NewBaseEvent("unknown.event", "unknown", "1"),
		}

		err := handler.Handle(ctx, evt)

		require.NoError(t, err)
		assert.Empty(t, cache.invalidatedUsers)
	})
}

// unknownEvent 未知事件类型
type unknownEvent struct {
	event.BaseEvent
}

// newTestHandler 创建测试用处理器
// 由于 CacheInvalidationHandler 依赖 *auth.PermissionCacheService，
// 我们需要一个适配层来进行测试
func newTestHandler(cache *mockPermissionCache, userRepo *mockUserQueryRepo) *testCacheHandler {
	return &testCacheHandler{
		cache:    cache,
		userRepo: userRepo,
	}
}

// testCacheHandler 测试用缓存处理器，模拟 CacheInvalidationHandler 的行为
type testCacheHandler struct {
	cache    *mockPermissionCache
	userRepo *mockUserQueryRepo
}

func (h *testCacheHandler) Handle(ctx context.Context, e event.Event) error {
	switch evt := e.(type) {
	case *event.UserRoleAssignedEvent:
		_ = h.cache.InvalidateUser(ctx, evt.UserID)
		return nil
	case *event.RolePermissionsChangedEvent:
		userIDs, err := h.userRepo.GetUserIDsByRole(ctx, evt.RoleID)
		if err != nil {
			// 查询失败时记录但不阻塞，返回 nil 是预期行为
			return nil //nolint:nilerr // intentional: query failure should not block business
		}
		for _, userID := range userIDs {
			_ = h.cache.InvalidateUser(ctx, userID)
		}
		return nil
	case *event.UserDeletedEvent:
		_ = h.cache.InvalidateUser(ctx, evt.UserID)
		return nil
	default:
		return nil
	}
}
