package org

import "context"

// TeamCommandRepository 团队命令仓储接口
type TeamCommandRepository interface {
	// Create 创建团队
	Create(ctx context.Context, team *Team) error

	// Update 更新团队
	Update(ctx context.Context, team *Team) error

	// Delete 删除团队（软删除）
	Delete(ctx context.Context, id uint) error
}
