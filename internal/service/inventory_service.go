package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/gcr0827/order-service/internal/pkg/apperr"
)

type InventoryRepository interface {
	// DeductStock 库存扣减
	DeductStock(ctx context.Context, storeID, skuId uint64, qty int64) error
}

type InventoryService struct {
	repo InventoryRepository
}

func NewInventoryService(repo InventoryRepository) *InventoryService {
	return &InventoryService{repo: repo}
}

func (i *InventoryService) DeductStock(ctx context.Context, storeID, skuId uint64, qty int64) error {
	err := i.repo.DeductStock(ctx, storeID, skuId, qty)
	if err != nil {
		// 库存不足，直接透传
		if errors.Is(err, apperr.ErrInsufficientStock) {
			return err
		}

		// 其他情况
		return fmt.Errorf("扣减库存失败: %w", err)
	}

	return nil
}
