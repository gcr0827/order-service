package model

import (
	"time"

	"gorm.io/gorm"
)

// StoreOrderItem StoreOrder 订单明细
type StoreOrderItem struct {
	ID             uint64         `gorm:"primaryKey;column:id" json:"id"`
	StoreID        uint64         `gorm:"column:store_id" json:"store_id"`
	OrderID        uint64         `gorm:"column:order_id" json:"order_id"`
	SpuID          uint64         `gorm:"column:spu_id" json:"spu_id"`
	SkuID          uint64         `gorm:"column:sku_id" json:"sku_id"`
	SkuCode        string         `gorm:"column:sku_code" json:"sku_code"`
	SkuName        string         `gorm:"column:sku_name" json:"sku_name"`
	SkuImgUrl      string         `gorm:"column:sku_img_url" json:"sku_img_url"`
	Price          int32          `gorm:"column:price" json:"price"`
	Quantity       int32          `gorm:"column:quantity" json:"quantity"`
	TotalAmount    uint64         `gorm:"column:total_amount" json:"total_amount"`
	DiscountAmount uint64         `gorm:"column:discount_amount" json:"discount_amount"`
	PayAmount      uint64         `gorm:"column:pay_amount" json:"pay_amount"`
	Status         int8           `gorm:"column:status" json:"status"`
	CreatedAt      time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt      time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"column:deleted_at;index" json:"deleted_at"`
}
