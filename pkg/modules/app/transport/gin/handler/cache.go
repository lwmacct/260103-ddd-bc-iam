package handler

import (
	"errors"

	"github.com/gin-gonic/gin"

	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/app/application/cache"
	cacheDomain "github.com/lwmacct/260101-go-pkg-ddd/pkg/shared/cache"
	"github.com/lwmacct/260101-go-pkg-gin/pkg/response"
)

// CacheHandler 缓存管理 HTTP 处理器（Redis 风格 API）
type CacheHandler struct {
	infoHandler     *cache.InfoHandler
	scanKeysHandler *cache.ScanKeysHandler
	getKeyHandler   *cache.GetKeyHandler
	deleteHandler   *cache.DeleteHandler
}

// NewCacheHandler 创建缓存管理处理器
func NewCacheHandler(
	infoHandler *cache.InfoHandler,
	scanKeysHandler *cache.ScanKeysHandler,
	getKeyHandler *cache.GetKeyHandler,
	deleteHandler *cache.DeleteHandler,
) *CacheHandler {
	return &CacheHandler{
		infoHandler:     infoHandler,
		scanKeysHandler: scanKeysHandler,
		getKeyHandler:   getKeyHandler,
		deleteHandler:   deleteHandler,
	}
}

// Info 获取缓存信息
//
//	@Summary		缓存信息
//	@Description	查看缓存状态信息（类似 redis-cli INFO）
//	@Tags			System
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	response.DataResponse[cache.CacheInfoDTO]	"成功"
//	@Failure		500	{object}	response.ErrorResponse						"服务器错误"
//	@Router			/api/admin/cache/info [get]
func (h *CacheHandler) Info(c *gin.Context) {
	info, err := h.infoHandler.Handle(c.Request.Context())
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, info)
}

// ScanKeysQuery SCAN 查询参数
type ScanKeysQuery struct {
	// Pattern 匹配模式（不含应用前缀，如 "setting:*"）
	Pattern string `form:"pattern"`

	// Cursor 游标（用于分页），首次查询传 "0" 或留空
	Cursor string `form:"cursor"`

	// Limit 每次返回的最大数量
	Limit int `form:"limit" binding:"omitempty,min=1,max=1000"`
}

// ScanKeys 扫描缓存 Keys
//
//	@Summary		扫描缓存键
//	@Description	按 pattern 扫描缓存 keys（类似 redis-cli SCAN）
//	@Tags			System
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			params	query		ScanKeysQuery									false	"查询参数"
//	@Success		200		{object}	response.DataResponse[cache.ScanKeysResultDTO]	"成功"
//	@Failure		400		{object}	response.ErrorResponse							"参数错误"
//	@Failure		500		{object}	response.ErrorResponse							"服务器错误"
//	@Router			/api/admin/cache/keys [get]
func (h *CacheHandler) ScanKeys(c *gin.Context) {
	var q ScanKeysQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	result, err := h.scanKeysHandler.Handle(c.Request.Context(), cache.ScanKeysQuery{
		Pattern: q.Pattern,
		Cursor:  q.Cursor,
		Limit:   q.Limit,
	})
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, result)
}

// GetKeyQuery 获取单个 Key 的查询参数
type GetKeyQuery struct {
	// Key 完整的 key 名称（含前缀）
	Key string `form:"key" binding:"required"`
}

// GetKey 获取单个 Key 的值
//
//	@Summary		获取缓存值
//	@Description	获取指定 key 的完整信息和值（类似 redis-cli GET/JSON.GET）
//	@Tags			System
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			params	query		GetKeyQuery									true	"查询参数"
//	@Success		200		{object}	response.DataResponse[cache.CacheValueDTO]	"成功"
//	@Failure		400		{object}	response.ErrorResponse						"参数错误"
//	@Failure		404		{object}	response.ErrorResponse						"Key 不存在"
//	@Failure		500		{object}	response.ErrorResponse						"服务器错误"
//	@Router			/api/admin/cache/key [get]
func (h *CacheHandler) GetKey(c *gin.Context) {
	var q GetKeyQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	result, err := h.getKeyHandler.Handle(c.Request.Context(), q.Key)
	if err != nil {
		if errors.Is(err, cacheDomain.ErrKeyNotFound) {
			response.NotFoundMessage(c, err.Error())
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, result)
}

// DeleteKeyQuery 删除单个 Key 的查询参数
type DeleteKeyQuery struct {
	// Key 完整的 key 名称（含前缀）
	Key string `form:"key" binding:"required"`
}

// DeleteKey 删除单个 Key
//
//	@Summary		删除缓存键
//	@Description	删除指定的单个 key（类似 redis-cli DEL）
//	@Tags			System
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			params	query		DeleteKeyQuery									true	"查询参数"
//	@Success		200		{object}	response.DataResponse[cache.DeleteResultDTO]	"成功"
//	@Failure		400		{object}	response.ErrorResponse							"参数错误"
//	@Failure		500		{object}	response.ErrorResponse							"服务器错误"
//	@Router			/api/admin/cache/key [delete]
func (h *CacheHandler) DeleteKey(c *gin.Context) {
	var q DeleteKeyQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	result, err := h.deleteHandler.DeleteKey(c.Request.Context(), q.Key)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, result)
}

// DeleteByPatternQuery 按 pattern 删除的查询参数
type DeleteByPatternQuery struct {
	// Pattern 匹配模式（不含应用前缀，如 "setting:*"）
	Pattern string `form:"pattern" binding:"required"`
}

// DeleteByPattern 按 pattern 批量删除 Keys
//
//	@Summary		批量删除缓存
//	@Description	批量删除匹配 pattern 的所有 keys
//	@Tags			System
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			params	query		DeleteByPatternQuery							true	"查询参数"
//	@Success		200		{object}	response.DataResponse[cache.DeleteResultDTO]	"成功"
//	@Failure		400		{object}	response.ErrorResponse							"参数错误"
//	@Failure		500		{object}	response.ErrorResponse							"服务器错误"
//	@Router			/api/admin/cache/keys [delete]
func (h *CacheHandler) DeleteByPattern(c *gin.Context) {
	var q DeleteByPatternQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	result, err := h.deleteHandler.DeleteByPattern(c.Request.Context(), q.Pattern)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, result)
}
