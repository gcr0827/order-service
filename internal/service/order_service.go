package service

import (
	"github.com/gcr0827/order-service/internal/model"
	"github.com/gcr0827/order-service/internal/pkg/apperr"
	"github.com/gcr0827/order-service/internal/repository"
)

type OrderService struct {
	repo repository.OrderRepository
}

// NewOrderService 依赖注入 repository 接口
func NewOrderService(repo repository.OrderRepository) *OrderService {
	return &OrderService{repo: repo}
}

// GetOrder 业务逻辑：错误向上传递时保留类型，由 handler 用 errors.Is 判断
func (s *OrderService) GetOrder(id int64) (*model.Order, error) {
	o, err := s.repo.GetByID(id)
	if err != nil {
		return nil, apperr.Wrap(apperr.CodeNotFound, err)
	}
	return o, nil
}
