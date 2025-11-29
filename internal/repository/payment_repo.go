package repository

import (
	"context"

	"github.com/innovationmech/simple-cli/internal/model"
	"gorm.io/gorm"
)

// PaymentRepository 支付数据访问接口
type PaymentRepository interface {
	CreatePayment(ctx context.Context, payment *model.Payment) error
	GetPayment(ctx context.Context, id string) (*model.Payment, error)
	UpdatePayment(ctx context.Context, payment *model.Payment) error
	ListPayments(ctx context.Context, userID, orderID string, offset, limit int) ([]*model.Payment, int64, error)
}

type paymentRepository struct {
	db *gorm.DB
}

// NewPaymentRepository 创建支付仓储实例
// 此函数将作为 fx Provider 使用
func NewPaymentRepository(db *gorm.DB) PaymentRepository {
	return &paymentRepository{db: db}
}

func (r *paymentRepository) CreatePayment(ctx context.Context, payment *model.Payment) error {
	return r.db.WithContext(ctx).Create(payment).Error
}

func (r *paymentRepository) GetPayment(ctx context.Context, id string) (*model.Payment, error) {
	var payment model.Payment
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&payment).Error; err != nil {
		return nil, err
	}
	return &payment, nil
}

func (r *paymentRepository) UpdatePayment(ctx context.Context, payment *model.Payment) error {
	return r.db.WithContext(ctx).Save(payment).Error
}

func (r *paymentRepository) ListPayments(ctx context.Context, userID, orderID string, offset, limit int) ([]*model.Payment, int64, error) {
	var payments []*model.Payment
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Payment{})
	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}
	if orderID != "" {
		query = query.Where("order_id = ?", orderID)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&payments).Error; err != nil {
		return nil, 0, err
	}

	return payments, total, nil
}
