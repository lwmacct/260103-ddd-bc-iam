package pat

import (
	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/iam/domain/pat"
)

// ToTokenDTO 将领域模型 PersonalAccessToken 转换为应用层 TokenDTO
func ToTokenDTO(token *pat.PersonalAccessToken) *TokenDTO {
	if token == nil {
		return nil
	}

	return &TokenDTO{
		ID:          token.ID,
		UserID:      token.UserID,
		Name:        token.Name,
		TokenPrefix: token.TokenPrefix,
		Scopes:      token.Scopes,
		IPWhitelist: token.IPWhitelist,
		Status:      token.Status,
		ExpiresAt:   token.ExpiresAt,
		LastUsedAt:  token.LastUsedAt,
		CreatedAt:   token.CreatedAt,
		UpdatedAt:   token.UpdatedAt,
	}
}

// ToCreateResultDTO 将领域模型 PersonalAccessToken 转换为创建响应 DTO（携带一次性明文 token）
func ToCreateResultDTO(token *pat.PersonalAccessToken, plainToken string) *CreateResultDTO {
	if token == nil {
		return nil
	}

	return &CreateResultDTO{
		PlainToken: plainToken,
		Token:      ToTokenDTO(token),
	}
}

// ToTokenListDTO 将领域模型 TokenListItem 数组转换为应用层 TokenListDTO
func ToTokenListDTO(items []*pat.TokenListItem) *TokenListDTO {
	responses := make([]*TokenDTO, len(items))
	for i, item := range items {
		responses[i] = &TokenDTO{
			ID:          item.ID,
			Name:        item.Name,
			TokenPrefix: item.TokenPrefix,
			Scopes:      item.Scopes,
			Status:      item.Status,
			ExpiresAt:   item.ExpiresAt,
			LastUsedAt:  item.LastUsedAt,
			CreatedAt:   item.CreatedAt,
			UpdatedAt:   item.CreatedAt,
		}
	}

	return &TokenListDTO{
		Tokens: responses,
		Total:  int64(len(responses)),
	}
}

// ToTokenInfoDTO 将领域实体转换为 TokenInfoDTO（不包含 token）
func ToTokenInfoDTO(token *pat.PersonalAccessToken) *TokenInfoDTO {
	if token == nil {
		return nil
	}

	return &TokenInfoDTO{
		ID:          token.ID,
		UserID:      token.UserID,
		Name:        token.Name,
		TokenPrefix: token.TokenPrefix,
		Scopes:      token.Scopes,
		IPWhitelist: token.IPWhitelist,
		Status:      token.Status,
		ExpiresAt:   token.ExpiresAt,
		LastUsedAt:  token.LastUsedAt,
		CreatedAt:   token.CreatedAt,
		UpdatedAt:   token.UpdatedAt,
	}
}
