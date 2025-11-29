package interfaces

import (
	"context"

	"github.com/innovationmech/simple-cli/internal/model"
)

// OrderService 订单服务接口
// 定义订单相关的业务操作
type OrderService interface {
	CreateOrder(ctx context.Context, order *model.Order) error
	GetOrder(ctx context.Context, id string) (*model.Order, error)
	UpdateOrderStatus(ctx context.Context, id string, status model.OrderStatus) error
	CancelOrder(ctx context.Context, id string) error
	ListOrdersByUser(ctx context.Context, userID string, page, pageSize int) ([]*model.Order, int64, error)
}
