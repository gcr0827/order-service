package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/gcr0827/order-service/internal/pkg/apperr"
	"github.com/gcr0827/order-service/internal/pkg/response"
	"github.com/gcr0827/order-service/internal/service"
)

type OrderHandler struct {
	svc *service.OrderService
}

func NewOrderHandler(svc *service.OrderService) *OrderHandler {
	return &OrderHandler{svc: svc}
}

// GetOrder 薄层：只做「参数校验 → 调 service → 写响应」。
// 业务逻辑一律不下沉到 handler。
func (h *OrderHandler) GetOrder(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		response.Fail(c, http.StatusBadRequest, apperr.CodeInvalidParam, "订单 ID 不合法")
		return
	}

	o, err := h.svc.GetOrder(id)
	if err != nil {
		fail(c, err)
		return
	}
	response.OK(c, o)
}
