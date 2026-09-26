package model

import (
	"time"

	"gorm.io/gorm"
)

// StoreOrder 订单主表
type StoreOrder struct {
	ID             uint64         `gorm:"primaryKey;column:id" json:"id"`
	StoreID        uint64         `gorm:"column:store_id" json:"store_id"`
	OrderNo        string         `gorm:"column:order_no" json:"order_no"`
	UserID         uint64         `gorm:"column:user_id" json:"user_id"`
	OrderType      int8           `gorm:"column:order_type" json:"order_type"`
	Status         int8           `gorm:"column:status" json:"status"`
	PayStatus      int8           `gorm:"column:pay_status" json:"pay_status"`
	PayTime        time.Time      `gorm:"column:pay_time" json:"pay_time"`
	TotalAmount    uint64         `gorm:"column:total_amount" json:"total_amount"`
	DiscountAmount uint64         `gorm:"column:discount_amount" json:"discount_amount"`
	PayAmount      uint64         `gorm:"column:pay_amount" json:"pay_amount"`
	Remark         string         `gorm:"column:remark" json:"remark"`
	CreatedAt      time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt      time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"column:deleted_at;index" json:"deleted_at"`
}
