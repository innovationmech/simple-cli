package repository

import (
	"context"

	"github.com/innovationmech/simple-cli/internal/model"
	"gorm.io/gorm"
)

// OrderRepository 订单数据访问接口
type OrderRepository interface {
	CreateOrder(ctx context.Context, order *model.Order) error
	GetOrder(ctx context.Context, id string) (*model.Order, error)
	UpdateOrder(ctx context.Context, order *model.Order) error
	ListOrdersByUser(ctx context.Context, userID string, offset, limit int) ([]*model.Order, int64, error)
}

type orderRepository struct {
	db *gorm.DB
}

// NewOrderRepository 创建订单仓储实例
// 此函数将作为 Wire Provider 使用
func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) CreateOrder(ctx context.Context, order *model.Order) error {
	return r.db.WithContext(ctx).Create(order).Error
}

func (r *orderRepository) GetOrder(ctx context.Context, id string) (*model.Order, error) {
	var order model.Order
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&order).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *orderRepository) UpdateOrder(ctx context.Context, order *model.Order) error {
	return r.db.WithContext(ctx).Save(order).Error
}

func (r *orderRepository) ListOrdersByUser(ctx context.Context, userID string, offset, limit int) ([]*model.Order, int64, error) {
	var orders []*model.Order
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Order{})
	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&orders).Error; err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}
