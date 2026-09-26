package model

import (
	"time"

	"gorm.io/gorm"
)

// StoreProductSku 商品SKU
type StoreProductSku struct {
	ID        uint64         `gorm:"primaryKey;column:id" json:"id"`
	StoreID   uint64         `gorm:"column:store_id" json:"store_id"`
	SpuID     uint64         `gorm:"column:spu_id" json:"spu_id"`
	Name      string         `gorm:"column:name" json:"name"`
	Code      string         `gorm:"column:code" json:"code"`
	Price     int64          `gorm:"column:price" json:"price"`
	CostPrice int64          `gorm:"column:cost_price" json:"cost_price"`
	ImgUrl    string         `gorm:"column:img_url" json:"img_url"`
	Status    int8           `gorm:"column:status" json:"status"`
	Sort      uint32         `gorm:"column:sort" json:"sort"`
	CreatedAt time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at" json:"-"`
}
