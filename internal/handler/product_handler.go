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

type ProductHandler struct {
	svc *service.ProductService
}

func NewProductHandler(svc *service.ProductService) *ProductHandler {
	return &ProductHandler{svc: svc}
}

// GetSpuByID 获取指定SPU的信息
func (p *ProductHandler) GetSpuByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		response.Fail(c, http.StatusBadRequest, apperr.CodeInvalidParam, "SPU ID不合法")
		return
	}
	spu, err := p.svc.GetSpu(c, uint64(id))
	if err != nil {
		switch {
		case errors.Is(err, apperr.ErrNotFound):
			response.Fail(c, http.StatusNotFound, apperr.CodeNotFound, "SPU ID 不存在")
		default:
			response.Fail(c, http.StatusInternalServerError, apperr.CodeInternalError, "服务器内部错误")

		}
		return
	}
	response.OK(c, spu)
}

// GetSkuByID 获取指定的Sku
func (p *ProductHandler) GetSkuByID(c *gin.Context) {
	// 接收入参
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		response.Fail(c, http.StatusBadRequest, apperr.CodeInvalidParam, "SKU ID 不合法")
		return
	}
	sku, err := p.svc.GetSku(c, uint64(id))
	if err != nil {
		switch {
		case errors.Is(err, apperr.ErrNotFound):
			response.Fail(c, http.StatusNotFound, apperr.CodeNotFound, "SPU ID 不存在")
		default:
			response.Fail(c, http.StatusInternalServerError, apperr.CodeInternalError, "服务器内部错误")

		}
		return
	}
	response.OK(c, sku)
}

// GetListSkuBySpuID 获取指定的Sku
func (p *ProductHandler) GetListSkuBySpuID(c *gin.Context) {
	// 接收入参
	spuID, err := strconv.ParseInt(c.Query("spu_id"), 10, 64)
	if err != nil || spuID <= 0 {
		response.Fail(c, http.StatusBadRequest, apperr.CodeInvalidParam, "SPU ID 不合法")
		return
	}

	storeID, err := strconv.ParseInt(c.Query("store_id"), 10, 64)
	if err != nil || storeID <= 0 {
		response.Fail(c, http.StatusBadRequest, apperr.CodeInvalidParam, "店铺 ID 不合法")
		return
	}

	sku, err := p.svc.GetListSkuBySpuID(c, uint64(storeID), uint64(spuID))
	if err != nil {
		switch {
		case errors.Is(err, apperr.ErrNotFound):
			response.Fail(c, http.StatusNotFound, apperr.CodeNotFound, "SPU 不存在")
		default:
			response.Fail(c, http.StatusInternalServerError, apperr.CodeInternalError, "服务器内部错误")

		}
		return
	}
	response.OK(c, sku)
}
