package model

import (
	"time"

	"gorm.io/gorm"
)

// StoreProductInventory 商品库存
type StoreProductInventory struct {
	ID             uint64         `gorm:"primaryKey;column:id" json:"id"`
	StoreID        uint64         `gorm:"column:store_id" json:"store_id"`
	SpuID          uint64         `gorm:"column:spu_id" json:"spu_id"`
	SkuID          uint64         `gorm:"column:sku_id" json:"sku_id"`
	TotalStock     uint64         `gorm:"column:total_stock" json:"total_stock"`
	AvailableStock uint64         `gorm:"column:available_stock" json:"available_stock"`
	LockedStock    uint64         `gorm:"column:locked_stock" json:"locked_stock"`
	WarnStock      uint64         `gorm:"column:warn_stock" json:"warn_stock"`
	CreatedAt      time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt      time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"column:deleted_at;index" json:"deleted_at"`
}
