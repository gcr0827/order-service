package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/gcr0827/order-service/internal/model"
	"github.com/gcr0827/order-service/internal/pkg/apperr"
)

type ProductRepository interface {
	// GetSpuByID 根据SPU ID获取SPU信息
	GetSpuByID(ctx context.Context, id uint64) (*model.StoreProductSpu, error)

	// GetSkuByID 根据SKU ID获取SKU信息
	GetSkuByID(ctx context.Context, id uint64) (*model.StoreProductSku, error)

	// ListSkuBySpuID 根据Spu获取sku列表
	ListSkuBySpuID(ctx context.Context, storeID, spuID uint64) ([]*model.StoreProductSku, error)
}

type ProductService struct {
	repo ProductRepository
}

func NewProductService(repo ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

// GetSpu 根据SPU ID获取SPU信息
func (p *ProductService) GetSpu(ctx context.Context, id uint64) (*model.StoreProductSpu, error) {
	data, err := p.repo.GetSpuByID(ctx, id)
	if err != nil {
		// 业务错误：repository 已经给了哨兵 → 直接透传
		if errors.Is(err, apperr.ErrNotFound) {
			return nil, err
		}
		// 系统错误：包上下文 + 保留原因（%w），【不】给它业务码
		return nil, fmt.Errorf("获取商品SPU失败: %w", err)
	}
	return data, nil
}

// GetSku 根据SKU ID获取SKU信息
func (p *ProductService) GetSku(ctx context.Context, id uint64) (*model.StoreProductSku, error) {
	data, err := p.repo.GetSkuByID(ctx, id)
	if err != nil {
		if errors.Is(err, apperr.ErrNotFound) {
			return nil, err
		}
		return nil, fmt.Errorf("获取商品SKU失败: %w", err)
	}

	return data, nil
}

// GetListSkuBySpuID 根据SPU ID获取SKU信息
func (p *ProductService) GetListSkuBySpuID(ctx context.Context, storeID, spuID uint64) ([]*model.StoreProductSku, error) {
	data, err := p.repo.ListSkuBySpuID(ctx, storeID, spuID)
	if err != nil {
		if errors.Is(err, apperr.ErrNotFound) {
			return nil, err
		}
		return nil, fmt.Errorf("获取商品SKU失败: %w", err)
	}

	return data, nil
}
