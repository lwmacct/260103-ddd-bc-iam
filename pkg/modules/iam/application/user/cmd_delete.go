package user

import (
	"context"
	"fmt"

	"github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/domain/user"
	"github.com/lwmacct/260103-ddd-bc-iam/pkg/shared/event"
)

// DeleteHandler 删除用户命令处理器
type DeleteHandler struct {
	userCommandRepo user.CommandRepository
	userQueryRepo   user.QueryRepository
	eventBus        event.EventBus
}

// NewDeleteHandler 创建删除用户命令处理器
func NewDeleteHandler(
	userCommandRepo user.CommandRepository,
	userQueryRepo user.QueryRepository,
	eventBus event.EventBus,
) *DeleteHandler {
	return &DeleteHandler{
		userCommandRepo: userCommandRepo,
		userQueryRepo:   userQueryRepo,
		eventBus:        eventBus,
	}
}

// Handle 处理删除用户命令
func (h *DeleteHandler) Handle(ctx context.Context, cmd DeleteCommand) error {
	// 1. 获取用户（需要检查是否为系统用户）
	u, err := h.userQueryRepo.GetByID(ctx, cmd.UserID)
	if err != nil {
		return user.ErrUserNotFound
	}

	// 2. 系统用户保护：系统用户不可删除
	if !u.CanBeDeleted() {
		return user.ErrCannotDeleteSystemUser
	}

	// 3. 执行删除（软删除）
	if err := h.userCommandRepo.Delete(ctx, cmd.UserID); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	// 4. 发布用户删除事件，触发缓存清理
	evt := event.NewUserDeletedEvent(cmd.UserID)
	if h.eventBus != nil {
		_ = h.eventBus.Publish(ctx, evt) // 缓存清理失败不阻塞业务
	}

	return nil
}
