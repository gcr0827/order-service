package model

import (
	"time"

	"gorm.io/gorm"
)

// StoreProductSpu 商品SPU
type StoreProductSpu struct {
	ID         uint64         `gorm:"primaryKey;column:id" json:"id"`
	StoreID    uint64         `gorm:"column:store_id" json:"store_id"`
	CategoryID uint64         `gorm:"column:category_id" json:"category_id"`
	Name       string         `gorm:"column:name" json:"name"`
	Code       string         `gorm:"column:code" json:"code"`
	ImgUrl     string         `gorm:"column:img_url" json:"img_url"`
	Status     int8           `gorm:"column:status" json:"status"`
	Sort       uint32         `gorm:"column:sort" json:"sort"`
	CreatedAt  time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt  time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"column:deleted_at;index" json:"-"`
}
