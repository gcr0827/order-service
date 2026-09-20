package repository

import (
	"github.com/gcr0827/order-service/internal/model"
	"github.com/gcr0827/order-service/internal/pkg/apperr"
)

// OrderRepository 先定义接口，service 依赖接口而非实现（便于单测 mock）
type OrderRepository interface {
	GetByID(id int64) (*model.Order, error)
}

type orderRepository struct {
	// W1 用内存数据，W2 换成 *gorm.DB
	data map[int64]*model.Order
}

// NewOrderRepository 构造 repository
func NewOrderRepository() OrderRepository {
	return &orderRepository{
		data: map[int64]*model.Order{
			1: {ID: 1, UserID: 100, Amount: 9900, Status: 1, CreateAt: "2026-09-20 16:00:00"},
		},
	}
}

// GetByID 只返回 (data, error)。
// 找不到时返回哨兵错误，不返回 "订单不存在" 这类 message —— 由 handler 决定呈现文案。
func (r *orderRepository) GetByID(id int64) (*model.Order, error) {
	o, ok := r.data[id]
	if !ok {
		return nil, apperr.ErrNotFound
	}
	return o, nil
}
