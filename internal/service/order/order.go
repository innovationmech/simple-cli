package order

import (
	"context"
	"errors"

	"github.com/innovationmech/simple-cli/internal/interfaces"
	"github.com/innovationmech/simple-cli/internal/model"
	"github.com/innovationmech/simple-cli/internal/repository"
)

// OrderSrv 是 OrderService 接口的别名
type OrderSrv = interfaces.OrderService

type orderService struct {
	orderRepo   repository.OrderRepository
	productRepo repository.ProductRepository
}

// NewOrderService 创建订单服务实例
// 此函数将作为 Wire Provider 使用
func NewOrderService(
	orderRepo repository.OrderRepository,
	productRepo repository.ProductRepository,
) OrderSrv {
	return &orderService{
		orderRepo:   orderRepo,
		productRepo: productRepo,
	}
}

func (s *orderService) CreateOrder(ctx context.Context, order *model.Order) error {
	// 获取商品信息计算总价
	product, err := s.productRepo.GetProduct(ctx, order.ProductID)
	if err != nil {
		return errors.New("product not found")
	}

	// 检查库存
	if product.Stock < order.Quantity {
		return errors.New("insufficient stock")
	}

	// 计算总金额
	order.TotalAmount = product.Price * float64(order.Quantity)
	order.Status = model.OrderStatusPending

	return s.orderRepo.CreateOrder(ctx, order)
}

func (s *orderService) GetOrder(ctx context.Context, id string) (*model.Order, error) {
	return s.orderRepo.GetOrder(ctx, id)
}

func (s *orderService) UpdateOrderStatus(ctx context.Context, id string, status model.OrderStatus) error {
	order, err := s.orderRepo.GetOrder(ctx, id)
	if err != nil {
		return err
	}

	// 状态流转验证
	if !isValidStatusTransition(order.Status, status) {
		return errors.New("invalid status transition")
	}

	order.Status = status
	return s.orderRepo.UpdateOrder(ctx, order)
}

func (s *orderService) CancelOrder(ctx context.Context, id string) error {
	order, err := s.orderRepo.GetOrder(ctx, id)
	if err != nil {
		return err
	}

	// 只有待支付状态可以取消
	if order.Status != model.OrderStatusPending {
		return errors.New("only pending orders can be cancelled")
	}

	order.Status = model.OrderStatusCancelled
	return s.orderRepo.UpdateOrder(ctx, order)
}

func (s *orderService) ListOrdersByUser(ctx context.Context, userID string, page, pageSize int) ([]*model.Order, int64, error) {
	offset := (page - 1) * pageSize
	if offset < 0 {
		offset = 0
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	return s.orderRepo.ListOrdersByUser(ctx, userID, offset, pageSize)
}

// isValidStatusTransition 验证状态流转是否合法
func isValidStatusTransition(from, to model.OrderStatus) bool {
	validTransitions := map[model.OrderStatus][]model.OrderStatus{
		model.OrderStatusPending:   {model.OrderStatusPaid, model.OrderStatusCancelled},
		model.OrderStatusPaid:      {model.OrderStatusShipped},
		model.OrderStatusShipped:   {model.OrderStatusCompleted},
		model.OrderStatusCompleted: {},
		model.OrderStatusCancelled: {},
	}

	allowed, ok := validTransitions[from]
	if !ok {
		return false
	}

	for _, s := range allowed {
		if s == to {
			return true
		}
	}
	return false
}
