package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/gcr0827/order-service/internal/model"
	"github.com/gcr0827/order-service/internal/pkg/apperr"
	"gorm.io/gorm"
)

type MysqlProductRepository struct {
	db *gorm.DB
}

func NewMysqlProductRepository(db *gorm.DB) *MysqlProductRepository {
	return &MysqlProductRepository{db: db}
}

// GetSpuByID 根据SPU ID获取SPU信息
func (r *MysqlProductRepository) GetSpuByID(ctx context.Context, id uint64) (*model.StoreProductSpu, error) {
	var spu model.StoreProductSpu
	err := r.db.WithContext(ctx).First(&spu, id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, apperr.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("查询商品SPU失败：%w", err)
	}

	return &spu, nil
}

// GetSkuByID 根据SKU ID获取SKU信息
func (r *MysqlProductRepository) GetSkuByID(ctx context.Context, id uint64) (*model.StoreProductSku, error) {
	var sku model.StoreProductSku
	err := r.db.WithContext(ctx).First(&sku, id).Error

	// 记录不存在，将gorm的哨兵转为项目层的
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, apperr.ErrNotFound
	}

	// 其他错误，如网络问题。数据库连接问题
	if err != nil {
		return nil, fmt.Errorf("查询商品SKU失败：%w", err)
	}
	return &sku, nil
}

// ListSkuBySpuID 根据Spu获取sku列表
func (r *MysqlProductRepository) ListSkuBySpuID(ctx context.Context, storeID, spuID uint64) ([]*model.StoreProductSku, error) {
	var skuList []*model.StoreProductSku
	err := r.db.WithContext(ctx).
		Where("spu_id = ?", spuID).
		Where("store_id = ?", storeID).
		Order("sort desc, id desc").
		Find(&skuList).Error

	// 其他错误，如网络问题。数据库连接问题
	if err != nil {
		return nil, fmt.Errorf("查询商品SKU失败：%w", err)
	}
	return skuList, nil
}
