package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/lwmacct/260101-go-pkg-ddd/pkg/modules/iam/transport/gin/registry"

	"github.com/lwmacct/260101-go-pkg-gin/pkg/response"
)

// OperationHandler 操作列表处理器
type OperationHandler struct{}

// NewOperationHandler 创建操作列表处理器
func NewOperationHandler() *OperationHandler {
	return &OperationHandler{}
}

// ListOperations 获取所有可用操作列表
//
//	@Summary		操作列表
//	@Description	返回所有可用的操作定义，供前端权限配置使用
//	@Tags			System
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	response.DataResponse[[]registry.OperationDefinition]	"操作列表"
//	@Failure		401	{object}	response.ErrorResponse								"未认证"
//	@Failure		403	{object}	response.ErrorResponse								"权限不足"
//	@Router			/api/admin/operations [get]
func (h *OperationHandler) ListOperations(c *gin.Context) {
	ops := registry.AllOperationDefinitions()
	response.OK(c, ops)
}
