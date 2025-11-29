package interfaces

import (
	"context"

	"github.com/innovationmech/simple-cli/internal/model"
)

// PaymentService 支付服务接口
// 定义支付相关的业务操作
type PaymentService interface {
	CreatePayment(ctx context.Context, payment *model.Payment) (paymentURL string, err error)
	GetPayment(ctx context.Context, id string) (*model.Payment, error)
	ProcessCallback(ctx context.Context, paymentID, transactionID string, success bool) error
	RefundPayment(ctx context.Context, id string, reason string) error
	ListPayments(ctx context.Context, userID, orderID string, page, pageSize int) ([]*model.Payment, int64, error)
}
