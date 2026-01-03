package setting

import (
	"context"
	"fmt"

	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/app/domain/setting"
)

// UserListHandler 获取用户配置列表查询处理器
type UserListHandler struct {
	settingQueryRepo setting.QueryRepository
	queryRepo        setting.UserSettingQueryRepository
}

// NewUserListHandler 创建 UserListHandler 实例
func NewUserListHandler(
	settingQueryRepo setting.QueryRepository,
	queryRepo setting.UserSettingQueryRepository,
) *UserListHandler {
	return &UserListHandler{
		settingQueryRepo: settingQueryRepo,
		queryRepo:        queryRepo,
	}
}

// Handle 处理获取用户配置列表查询
// 返回合并后的配置列表：如果用户有自定义值则使用用户值，否则使用默认值
func (h *UserListHandler) Handle(ctx context.Context, query UserListQuery) ([]*UserSettingDTO, error) {
	// 1. 查找所有配置定义
	var defs []*setting.Setting
	var err error

	if query.CategoryID != 0 {
		defs, err = h.settingQueryRepo.FindByCategoryID(ctx, query.CategoryID)
	} else {
		defs, err = h.settingQueryRepo.FindAll(ctx)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find setting definitions: %w", err)
	}

	// 2. 查找用户的所有自定义配置
	userSettings, err := h.queryRepo.FindByUser(ctx, query.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to find user settings: %w", err)
	}

	// 3. 构建用户配置映射
	userSettingMap := make(map[string]*setting.UserSetting, len(userSettings))
	for _, us := range userSettings {
		userSettingMap[us.SettingKey] = us
	}

	// 4. 合并返回
	result := make([]*UserSettingDTO, 0, len(defs))
	for _, def := range defs {
		us := userSettingMap[def.Key]
		result = append(result, ToUserSettingDTO(def, us))
	}

	return result, nil
}
