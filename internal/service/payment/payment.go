package payment

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/innovationmech/simple-cli/internal/interfaces"
	"github.com/innovationmech/simple-cli/internal/model"
	"github.com/innovationmech/simple-cli/internal/repository"
)

// PaymentSrv 是 PaymentService 接口的别名
type PaymentSrv = interfaces.PaymentService

type paymentService struct {
	paymentRepo repository.PaymentRepository
	orderRepo   repository.OrderRepository
}

// NewPaymentService 创建支付服务实例
// 此函数将作为 fx Provider 使用
func NewPaymentService(
	paymentRepo repository.PaymentRepository,
	orderRepo repository.OrderRepository,
) PaymentSrv {
	return &paymentService{
		paymentRepo: paymentRepo,
		orderRepo:   orderRepo,
	}
}

func (s *paymentService) CreatePayment(ctx context.Context, payment *model.Payment) (string, error) {
	// 验证订单存在
	order, err := s.orderRepo.GetOrder(ctx, payment.OrderID)
	if err != nil {
		return "", errors.New("order not found")
	}

	// 验证订单金额
	if payment.Amount != order.TotalAmount {
		return "", errors.New("payment amount does not match order total")
	}

	// 验证订单状态
	if order.Status != model.OrderStatusPending {
		return "", errors.New("order is not in pending status")
	}

	payment.Status = model.PaymentStatusPending
	if err := s.paymentRepo.CreatePayment(ctx, payment); err != nil {
		return "", err
	}

	// 模拟生成支付链接
	paymentURL := s.generatePaymentURL(payment)
	return paymentURL, nil
}

func (s *paymentService) GetPayment(ctx context.Context, id string) (*model.Payment, error) {
	return s.paymentRepo.GetPayment(ctx, id)
}

func (s *paymentService) ProcessCallback(ctx context.Context, paymentID, transactionID string, success bool) error {
	payment, err := s.paymentRepo.GetPayment(ctx, paymentID)
	if err != nil {
		return err
	}

	if payment.Status != model.PaymentStatusPending {
		return errors.New("payment is not in pending status")
	}

	payment.TransactionID = transactionID
	now := time.Now()

	if success {
		payment.Status = model.PaymentStatusSuccess
		payment.PaidAt = &now

		// 更新订单状态为已支付
		order, err := s.orderRepo.GetOrder(ctx, payment.OrderID)
		if err == nil {
			order.Status = model.OrderStatusPaid
			_ = s.orderRepo.UpdateOrder(ctx, order)
		}
	} else {
		payment.Status = model.PaymentStatusFailed
	}

	return s.paymentRepo.UpdatePayment(ctx, payment)
}

func (s *paymentService) RefundPayment(ctx context.Context, id string, reason string) error {
	payment, err := s.paymentRepo.GetPayment(ctx, id)
	if err != nil {
		return err
	}

	if payment.Status != model.PaymentStatusSuccess {
		return errors.New("only successful payments can be refunded")
	}

	payment.Status = model.PaymentStatusRefunded
	return s.paymentRepo.UpdatePayment(ctx, payment)
}

func (s *paymentService) ListPayments(ctx context.Context, userID, orderID string, page, pageSize int) ([]*model.Payment, int64, error) {
	offset := (page - 1) * pageSize
	if offset < 0 {
		offset = 0
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	return s.paymentRepo.ListPayments(ctx, userID, orderID, offset, pageSize)
}

// generatePaymentURL 模拟生成支付链接
func (s *paymentService) generatePaymentURL(payment *model.Payment) string {
	baseURL := "https://pay.example.com"
	switch payment.Method {
	case model.PaymentMethodAlipay:
		return fmt.Sprintf("%s/alipay?id=%s&amount=%.2f", baseURL, payment.ID, payment.Amount)
	case model.PaymentMethodWechat:
		return fmt.Sprintf("%s/wechat?id=%s&amount=%.2f", baseURL, payment.ID, payment.Amount)
	case model.PaymentMethodCreditCard:
		return fmt.Sprintf("%s/card?id=%s&amount=%.2f", baseURL, payment.ID, payment.Amount)
	default:
		return fmt.Sprintf("%s/pay?id=%s&amount=%.2f", baseURL, payment.ID, payment.Amount)
	}
}
