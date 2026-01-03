package usersetting

import (
	"context"
	"encoding/json"
	"fmt"

	usersettingdomain "github.com/lwmacct/260103-ddd-bc-iam/pkg/modules/iam/domain/usersetting"
	setting "github.com/lwmacct/260103-ddd-bc-settings/pkg/modules/settings/domain/setting"
)

// ListHandler 获取用户设置列表命令处理器
// 核心逻辑：合并 Settings Schema + 用户自定义值
type ListHandler struct {
	userSettingQueryRepo usersettingdomain.QueryRepository
	settingQueryRepo     setting.QueryRepository // 注入 Settings BC 的 QueryRepository
}

// NewListHandler 创建获取用户设置列表命令处理器
func NewListHandler(
	userSettingQueryRepo usersettingdomain.QueryRepository,
	settingQueryRepo setting.QueryRepository,
) *ListHandler {
	return &ListHandler{
		userSettingQueryRepo: userSettingQueryRepo,
		settingQueryRepo:     settingQueryRepo,
	}
}

// Handle 处理获取用户设置列表查询
// 返回：系统 Schema + 用户自定义值的合并视图
func (h *ListHandler) Handle(ctx context.Context, userID uint, category string) (*UserSettingListDTO, error) {
	// 1. 获取用户可见的 Schema（scope=user + scope=system 且 public=true）
	schemas, err := h.settingQueryRepo.FindVisibleToUser(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch settings schema: %w", err)
	}

	// 2. 获取用户自定义值
	userValues, err := h.userSettingQueryRepo.FindByUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user settings: %w", err)
	}

	// 3. 构建用户值的 Map（用于快速查找）
	userValueMap := make(map[string]*usersettingdomain.UserSetting)
	for _, uv := range userValues {
		userValueMap[uv.Key] = uv
	}

	// 4. 合并 Schema 和用户值
	result := make([]*UserSettingDTO, 0, len(schemas))
	for _, schema := range schemas {
		// 分类过滤（按 CategoryID）
		// TODO: 实现分类过滤功能（需要关联查询 SettingCategory 表）
		// 当前 category 参数暂不使用，所有用户可见设置都返回

		// 将 DefaultValue 转换为 JSON 字符串
		defaultValueJSON, err := json.Marshal(schema.DefaultValue)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal default value: %w", err)
		}

		dto := &UserSettingDTO{
			Key:        schema.Key,
			CategoryID: schema.CategoryID,
			ValueType:  schema.ValueType,
			Label:      schema.Label,
		}

		// 如果用户有自定义值，使用用户值；否则使用系统默认值
		if userValue, exists := userValueMap[schema.Key]; exists {
			dto.Value = userValue.Value
			dto.IsCustom = true
		} else {
			dto.Value = string(defaultValueJSON)
			dto.IsCustom = false
		}

		result = append(result, dto)
	}

	return &UserSettingListDTO{
		Settings: result,
		Total:    int64(len(result)),
	}, nil
}
