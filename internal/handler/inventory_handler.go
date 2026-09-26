package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gcr0827/order-service/internal/pkg/apperr"
	"github.com/gcr0827/order-service/internal/pkg/response"
	"github.com/gcr0827/order-service/internal/service"
	"github.com/gin-gonic/gin"
)

type InventoryHandler struct {
	svc *service.InventoryService
}

func NewInventoryHandler(svc *service.InventoryService) *InventoryHandler {
	return &InventoryHandler{svc: svc}
}

// DeductStock 库存扣减
func (p *InventoryHandler) DeductStock(c *gin.Context) {
	// 接收参数
	storeID, err := strconv.ParseInt(c.PostForm("store_id"), 10, 64)
	if err != nil || storeID <= 0 {
		response.Fail(c, http.StatusBadRequest, apperr.CodeInvalidParam, "店铺ID不合法")
		return
	}

	skuID, err := strconv.ParseInt(c.PostForm("sku_id"), 10, 64)
	if err != nil || skuID <= 0 {
		response.Fail(c, http.StatusBadRequest, apperr.CodeInvalidParam, "SKU ID不合法")
		return
	}

	qty, err := strconv.ParseInt(c.PostForm("qty"), 10, 64)
	if err != nil || qty <= 0 {
		response.Fail(c, http.StatusBadRequest, apperr.CodeInvalidParam, "扣减不合法")
		return
	}

	err = p.svc.DeductStock(c, uint64(storeID), uint64(skuID), qty)
	if err != nil {
		switch {
		case errors.Is(err, apperr.ErrInsufficientStock):
			response.Fail(c, http.StatusInternalServerError, apperr.CodeInternalError, "库存无法扣减")
		default:
			response.Fail(c, http.StatusInternalServerError, apperr.CodeInternalError, "服务器内部错误")

		}
		return
	}
	response.OK(c, nil)
}
