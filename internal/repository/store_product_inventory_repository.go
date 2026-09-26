package repository

import (
	"context"
	"fmt"

	"github.com/gcr0827/order-service/internal/model"
	"github.com/gcr0827/order-service/internal/pkg/apperr"
	"gorm.io/gorm"
)

type StoreProductInventoryRepository struct {
	db *gorm.DB
}

func NewStoreProductInventoryRepository(db *gorm.DB) *StoreProductInventoryRepository {
	return &StoreProductInventoryRepository{
		db: db,
	}
}

// DeductStock 扣减库存
func (i *StoreProductInventoryRepository) DeductStock(ctx context.Context, storeID, skuID uint64, qty int64) error {

	res := i.db.WithContext(ctx).
		Model(&model.StoreProductInventory{}).
		Where("store_id = ?", storeID).
		Where("sku_id = ?", skuID).
		Where("available_stock >= ?", qty).
		Update("available_stock", gorm.Expr("available_stock - ?", qty))

	if res.Error != nil {
		return fmt.Errorf("库存扣减失败: %w", res.Error)
	}

	if res.RowsAffected == 0 {
		return apperr.ErrInsufficientStock // 没有成功更新
	}

	return nil
}
