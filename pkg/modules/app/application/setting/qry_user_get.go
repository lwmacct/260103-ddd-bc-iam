package setting

import (
	"context"
	"errors"
	"fmt"

	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/app/domain/setting"
)

// UserGetHandler 获取用户配置查询处理器
type UserGetHandler struct {
	settingQueryRepo setting.QueryRepository
	queryRepo        setting.UserSettingQueryRepository
}

// NewUserGetHandler 创建 UserGetHandler 实例
func NewUserGetHandler(
	settingQueryRepo setting.QueryRepository,
	queryRepo setting.UserSettingQueryRepository,
) *UserGetHandler {
	return &UserGetHandler{
		settingQueryRepo: settingQueryRepo,
		queryRepo:        queryRepo,
	}
}

// Handle 处理获取用户配置查询
// 返回合并后的配置：如果用户有自定义值则使用用户值，否则使用默认值
func (h *UserGetHandler) Handle(ctx context.Context, query UserGetQuery) (*UserSettingDTO, error) {
	// 1. 查找配置定义
	def, err := h.settingQueryRepo.FindByKey(ctx, query.Key)
	if err != nil {
		return nil, fmt.Errorf("failed to find setting definition: %w", err)
	}
	if def == nil {
		return nil, errors.New("setting key does not exist")
	}

	// 2. 查找用户自定义配置
	us, err := h.queryRepo.FindByUserAndKey(ctx, query.UserID, query.Key)
	if err != nil {
		return nil, fmt.Errorf("failed to find user setting: %w", err)
	}

	// 3. 合并返回
	return ToUserSettingDTO(def, us), nil
}
