package model

import (
	"time"

	"gorm.io/gorm"
)

// StoreProductInventoryLog 商品库存流水
type StoreProductInventoryLog struct {
	ID          uint64         `gorm:"primaryKey;column:id" json:"id"`
	StoreID     uint64         `gorm:"column:store_id" json:"store_id"`
	SpuID       uint64         `gorm:"column:spu_id" json:"spu_id"`
	SkuID       uint64         `gorm:"column:sku_id" json:"sku_id"`
	Quantity    uint64         `gorm:"column:quantity" json:"quantity"`
	BeforeStock uint64         `gorm:"column:before_stock" json:"before_stock"`
	AfterStock  uint64         `gorm:"column:after_stock" json:"after_stock"`
	BizType     int8           `gorm:"column:biz_type" json:"biz_type"`
	BizNo       string         `gorm:"column:biz_no" json:"biz_no"`
	CreatedAt   time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at;index" json:"deleted_at"`
}
